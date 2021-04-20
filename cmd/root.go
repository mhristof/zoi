package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/mhristof/zoi/docker"
	"github.com/mhristof/zoi/gh"
	"github.com/mhristof/zoi/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "zoi",
	Short: "Ze Ongoing Improvement",
	Long: heredoc.Doc(`
		Update libraries and packages with the latest versions.

		To pin versions for Docker files, run 'docker build' with
			zoi -- docker build -t foo .
		where 'docker build -t foo .' would be the command to build your
		docker container.

		To update a file containing supported versions, feed it in as
			zoi file.txt
		and updated version of the file will be output to stdout.
	`),
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 2 && args[0] == "docker" && args[1] == "build" {
			return nil
		}

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

		if args[0] == "docker" && args[1] == "build" {
			docker.Build(args)
			return
		}

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
