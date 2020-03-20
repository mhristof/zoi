package cmd

import (
	"fmt"
	"os"

	"github.com/mhristof/zoi/ansible"
	"github.com/mhristof/zoi/log"
	"github.com/spf13/cobra"
)

var (
	updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update the requirements file to the latest versions",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			Verbose(cmd)
			reqFile, err := cmd.Flags().GetString("requirements")
			if err != nil {
				log.Panic(err)
			}

			if _, err := os.Stat(reqFile); os.IsNotExist(err) {
				fmt.Println(fmt.Sprintf("Error, file %s not found", reqFile))
				os.Exit(1)
			}

			reqs := ansible.Requirements{}
			reqs.LoadFromFile(reqFile)
			out, err := cmd.Flags().GetString("output")
			if err != nil {
				log.Panic(err)
			}

			reqs.Update().SaveToFile(out)
		},
	}
)

func init() {
	updateCmd.PersistentFlags().StringP("requirements", "r", "requirements.yml", "Requirements file to update")
	updateCmd.PersistentFlags().StringP("output", "o", "latest.yml", "Output file for the latest and greatest")

	rootCmd.AddCommand(updateCmd)
}
