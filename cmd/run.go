package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	lastfm "github.com/irlndts/go-lastfm"
	twitter "github.com/irlndts/go-twitter"
	"github.com/spf13/cobra"
)

func newUserTopArtistsServerCommand() *cobra.Command {
	var opts userTopArtistsOptions
	c := &cobra.Command{
		Use:   "run",
		Short: "Run a daemon to publish top once a week",
		Run: func(cmd *cobra.Command, args []string) {
			if opts.user == "" {
				fmt.Fprintln(os.Stderr, "Username doesn't available")
				os.Exit(0)
			}
			if err := userTopArtistsRun(cmd, args, &opts); err != nil {
				fmt.Fprintln(os.Stderr, "Error: ", err)
				os.Exit(0)
			}
		},
	}
	c.PersistentFlags().StringVar(&opts.user, "user", "", "The user name to fetch top artists for.")
	return c
}

func userTopArtistsRun(cmd *cobra.Command, args []string, opts *userTopArtistsOptions) error {
	if opts.user == "" {
		return fmt.Errorf("user is unknown")
	}
	lastfm, err := lastfm.New(cfg.LastFM.Key)
	if err != nil {
		return err
	}

	if cfg.Twitter.Consumer == "" || cfg.Twitter.Secret == "" ||
		cfg.Twitter.Token == "" || cfg.Twitter.TokenSecret == "" {
		return fmt.Errorf("Bad twitter arguments")
	}
	twitter, err := twitter.NewClient(cfg.Twitter.Consumer, cfg.Twitter.Secret)
	if err != nil {
		return err
	}
	twitter.Token(cfg.Twitter.Token, cfg.Twitter.TokenSecret)

	server := newServer(opts.user, lastfm, twitter)
	go server.Run()
	fmt.Println("server started...")

	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	sig := <-term
	fmt.Printf("received %s, exiting gracefully...\n", sig)

	stopped := make(chan struct{})
	go func() {
		// Web server must first to stop accepting requests
		server.Shutdown()

		stopped <- struct{}{}
	}()

	select {
	case <-time.After(5 * time.Second):
		fmt.Println("shutdown timeout expired, server stopped")
	case <-stopped:
		fmt.Println("server gracefully stopped")
	}

	return nil

}

type l2tServer struct {
	user string

	lastFM  *lastfm.Client
	twitter *twitter.Twitter

	quitCh chan struct{}
}

func newServer(user string, lastFM *lastfm.Client, twitter *twitter.Twitter) *l2tServer {
	return &l2tServer{
		user:    user,
		lastFM:  lastFM,
		twitter: twitter,

		quitCh: make(chan struct{}),
	}
}

func (s *l2tServer) Run() {
	var now time.Time
	for {
		select {
		case now = <-time.Tick(time.Hour):
			if now.Weekday() == time.Monday && now.Hour() > 9 && now.Hour() < 10 {
				top, err := s.lastFM.User.TopArtists(s.user, lastfm.PeriodWeek, 10, 0)
				if err != nil {
					fmt.Printf("failed to get top artists: %s\n", err)
					continue
				}
				if top.Total == 0 {
					continue
				}

				msg := fmt.Sprintf(msgStart, "week")
				for _, a := range top.Artists {
					line := fmt.Sprintf("(%d) %s\n", a.Playcount, a.Name)
					if len(line)+len(msg)+len(msgEnd) > 280 {
						break
					}
					msg += line
				}
				msg += msgEnd
				if err := s.twitter.Update(msg); err != nil {
					fmt.Printf("failed to publish twit: %s\n", err)
					continue
				}
			}

		case <-s.quitCh:
			return
		}
	}
}

func (s *l2tServer) Shutdown() {
	close(s.quitCh)
	// Wait until all jobs (from queue) will be processed.
	time.Sleep(time.Second)
}
