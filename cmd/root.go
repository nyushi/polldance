package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/nyushi/polldance"
	"github.com/spf13/cobra"
)

var (
	inFile        string
	outHTTP       string
	outHTTPMethod string
	debug         bool
	filter        string
)

func init() {
	rootCmd.PersistentFlags().StringVar(&inFile, "in-file", "", "path of input file")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "debug")
	rootCmd.PersistentFlags().StringVar(&outHTTP, "out-http", "", "output url, {{.Filename}} and {{.Path}} are provided")
	rootCmd.PersistentFlags().StringVar(&outHTTPMethod, "out-http-method", "POST", "")
	rootCmd.PersistentFlags().StringVar(&filter, "filter", "", "filter command")
}

var rootCmd = &cobra.Command{
	Use:   "polldance",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		c := &polldance.PollConfig{
			InputFilePaths:   []string{inFile},
			OutputHTTPURL:    outHTTP,
			OutputHTTPMethod: outHTTPMethod,
			FilterCommand:    filter,
			Debug:            debug,
		}
		if err := polldance.Poll(c); err != nil {
			log.Fatal(err)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
