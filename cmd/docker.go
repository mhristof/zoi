package cmd

import (
	"fmt"

	"github.com/mhristof/zoi/docker"
	"github.com/spf13/cobra"
)

var (
	dockerfile string
)

var (
	dockerCmd = &cobra.Command{
		Use:   "docker",
		Short: "Update a dockerfile",
		Run: func(cmd *cobra.Command, args []string) {
			doc := docker.New(dockerfile)
			fmt.Println(doc.Render())
		},
	}
)

func init() {
	dockerCmd.Flags().StringVarP(&dockerfile, "file", "f", "Dockerfile", "Name of the Dockerfile")

	rootCmd.AddCommand(dockerCmd)
}
