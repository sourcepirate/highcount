package main

import (
	"fmt"
	"os"

	"github.com/sourcepirate/highcount/pkg/gitcount"
	"github.com/spf13/cobra"
)

var gitLabUsername string
var gitLabPassword string

var rootCommand = &cobra.Command{
	Use:   "highcount",
	Short: "High count issue with issue spend",
	Long:  "time  count",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		project := args[0]
		label := args[1]
		client := gitcount.New(gitLabUsername, gitLabPassword, project)
		foundProject, err := client.GetProject()
		if err != nil {
			fmt.Println("Cannot find the project")
			os.Exit(1)
		}
		client.PrintStats(foundProject, label)
	},
}

func main() {
	rootCommand.PersistentFlags().StringVar(&gitLabUsername, "username", "", "User username")
	rootCommand.PersistentFlags().StringVar(&gitLabPassword, "password", "", "Password")
	if err := rootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
