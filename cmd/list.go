package cmd

import (
	"fmt"
	"os"

	lastfm "github.com/irlndts/go-lastfm"
	twitter "github.com/irlndts/go-twitter"
	"github.com/spf13/cobra"
)

type userTopArtistsOptions struct {
	user    string
	period  string
	limit   int
	offset  int
	publish bool
}

func newUserTopArtistsCommand() *cobra.Command {
	var opts userTopArtistsOptions
	c := &cobra.Command{
		Use:   "list",
		Short: "Show a top by list without publishing",
		Run: func(cmd *cobra.Command, args []string) {
			if opts.user == "" {
				fmt.Fprintln(os.Stderr, "Username doesn't available")
				os.Exit(0)
			}
			if err := userTopArtists(cmd, args, &opts); err != nil {
				fmt.Fprintln(os.Stderr, "Error: ", err)
				os.Exit(0)
			}
		},
	}
	c.PersistentFlags().StringVar(&opts.user, "user", "", "The user name to fetch top artists for.")
	c.PersistentFlags().StringVar(&opts.period, "period", "week", "The time period over which to retrieve top artists for (overall|week|month|quartal|halfyear|year).")
	c.PersistentFlags().IntVar(&opts.limit, "limit", 10, "The number of results to fetch per page.")
	c.PersistentFlags().IntVar(&opts.offset, "offset", 0, "The page number to fetch.")
	c.PersistentFlags().BoolVar(&opts.publish, "publish", false, "Publish the data to twitter")
	return c
}

const (
	msgStart = "This week Last.fm top:\n"
	msgEnd   = "\nMade by github.com/irlndts/go-lastfm-to-twitter"
)

// PeriodType ...
type PeriodType string

// Periods
const (
	PeriodOverall  PeriodType = "overall"
	PeriodWeek     PeriodType = "7day"
	PeriodMonth    PeriodType = "1month"
	PeriodQuartal  PeriodType = "3month"
	PeriodHalfYear PeriodType = "6month"
	PeriodYear     PeriodType = "12month"
)

func period(p string) PeriodType {
	return map[string]PeriodType{
		"overall":  PeriodOverall,
		"week":     PeriodWeek,
		"month":    PeriodMonth,
		"quartal":  PeriodQuartal,
		"halfyear": PeriodHalfYear,
		"year":     PeriodYear,
	}[p]
}

func userTopArtists(cmd *cobra.Command, args []string, opts *userTopArtistsOptions) error {
	lastfm, err := lastfm.New(cfg.LastFM.Key)
	if err != nil {
		return err
	}
	top, err := lastfm.User.TopArtists(opts.user, string(period(opts.period)), opts.limit, opts.offset)
	if err != nil {
		return err
	}
	if top.Total == 0 {
		fmt.Println("Nothing was listened")
		return nil
	}

	for _, a := range top.Artists {
		fmt.Printf("%d\t%s\n", a.Playcount, a.Name)
	}

	if opts.publish {
		if cfg.Twitter.Consumer == "" || cfg.Twitter.Secret == "" ||
			cfg.Twitter.Token == "" || cfg.Twitter.TokenSecret == "" {
			return fmt.Errorf("Bad twitter arguments")
		}
		twitter, err := twitter.NewClient(cfg.Twitter.Consumer, cfg.Twitter.Secret)
		if err != nil {
			return err
		}
		twitter.Token(cfg.Twitter.Token, cfg.Twitter.TokenSecret)

		msg := msgStart
		for _, a := range top.Artists {
			line := fmt.Sprintf("(%d) %s\n", a.Playcount, a.Name)
			if len(line)+len(msg)+len(msgEnd) > 280 {
				break
			}
			msg += line
		}
		msg += msgEnd
		if err := twitter.Update(msg); err != nil {
			return err
		}
		fmt.Println("message published")
	}
	return nil
}
