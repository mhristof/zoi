package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/mhristof/zoi/gh"
	"github.com/mhristof/zoi/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "zoi",
	Short: "Ze Ongoing Improvement",
	Long:  `Update libraries and packages with the latest versions`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("File not provided")
		}

		for _, arg := range args {
			if _, err := os.Stat(arg); os.IsNotExist(err) {
				return errors.New("File not found")
			}
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		Verbose(cmd)

		lines, err := ioutil.ReadFile(args[0])
		if err != nil {
			log.WithFields(log.Fields{
				"err":     err,
				"args[0]": args[0],
			}).Error("Could not read file")
		}

		for _, line := range strings.Split(string(lines), "\n") {
			fmt.Println(gh.Release(line))
		}

	},
}

// Verbose Increase verbosity
func Verbose(cmd *cobra.Command) {
	verbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		log.Panic(err)
	}

	if verbose {
		log.SetLevel(log.DebugLevel)
	}
}
func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Increase verbosity")
}

// Execute The main function for the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
