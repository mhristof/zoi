package cmd

import (
	"fmt"
	"os"

	"github.com/mhristof/zoi/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "zoi",
	Short: "Ze Ongoing Improvement",
	Long:  `Update libraries and packages with the latest versions`,
	Run: func(cmd *cobra.Command, args []string) {
		Verbose(cmd)
	},
}

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

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
