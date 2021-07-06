package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"syscall"

	"github.com/MakeNowJust/heredoc"
	"github.com/mhristof/zoi/gh"
	"github.com/mhristof/zoi/log"
	"github.com/mhristof/zoi/precommit"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	version string
)

var rootCmd = &cobra.Command{
	Use:     "zoi",
	Short:   "Ze Ongoing Improvement",
	Version: version,
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
				return errors.Wrap(err, "File not found")
			}
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		Verbose(cmd)

		byteLines, err := ioutil.ReadFile(args[0])
		if err != nil {
			log.WithFields(log.Fields{
				"err":     err,
				"args[0]": args[0],
			}).Error("Could not read file")
		}

		out := os.Stdout
		inplace, err := cmd.Flags().GetBool("inplace")
		if err != nil {
			panic(err)
		}

		if inplace {
			out, err = os.Create(args[0])
			if err != nil {
				panic(err)
			}

			defer out.Close()
		}

		ghToken := getGithubToken()
		precommitContents, err := precommit.Update(byteLines, ghToken)
		if err == nil {
			fmt.Fprintf(out, "%s", precommitContents)

			return
		}

		log.WithFields(log.Fields{
			"err": err,
		}).Debug("Handling liny by line")

		// lines ends up having one extra line at the end. Im sure there is a
		// better fix, but meh.
		llines := strings.Split(string(byteLines), "\n")
		for _, line := range llines[0 : len(llines)-1] {
			fmt.Fprintf(out, "%s\n", gh.Release(line, ghToken))
		}
	},
}

func getGithubToken() string {
	ghToken := os.Getenv("GITHUB_READONLY_TOKEN")
	if ghToken != "" {
		return ghToken
	}

	fmt.Print("Enter github token: ")

	byteToken, err := terminal.ReadPassword(syscall.Stdin)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Cannot read github token")
	}

	return string(byteToken)
}

// Verbose Increase verbosity.
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
	rootCmd.PersistentFlags().BoolP("inplace", "i", false, "Inplace replacement of the target file")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Increase verbosity")
}

// Execute The main function for the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
