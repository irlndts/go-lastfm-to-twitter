package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "go-lastfm-to-twitter",
	Short: "Get your last.fm chart and publish it on twitter",
	Long:  `Publish your current last.fm chart into your twitter account`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Config ...
type Config struct {
	Twitter struct {
		Consumer    string
		Secret      string
		Token       string
		TokenSecret string
	}
	LastFM struct {
		Key string
	}
}

var cfg = &Config{}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfg.Twitter.Consumer, "twitter-consumer", "", "twitter consumer key")
	rootCmd.PersistentFlags().StringVar(&cfg.Twitter.Token, "twitter-token", "", "twitter token")
	rootCmd.PersistentFlags().StringVar(&cfg.Twitter.TokenSecret, "twitter-token-secret", "", "twitter token secret")
	rootCmd.PersistentFlags().StringVar(&cfg.Twitter.Secret, "twitter-secret", "", "twitter secret")
	rootCmd.PersistentFlags().StringVar(&cfg.LastFM.Key, "lastfm-key", "", "LastFM application key")

	rootCmd.AddCommand(
		newUserTopArtistsCommand(),
		newUserTopArtistsServerCommand(),
	)
}
