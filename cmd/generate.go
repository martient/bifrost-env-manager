/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"

	environmentmanager "github.com/martient/Bifrost-env-manager/pkg/environment_manager"
	"github.com/spf13/cobra"
)

var newEnvFilePath string

// generateCmd represents the load command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a new version of the env file",
	Long:  `Generate a new version of the environement file in function of the config given`,
	Run: func(cmd *cobra.Command, args []string) {
		jsonFile, err := os.Open(jsonConfigFile)
		if err != nil {
			fmt.Println(err)
		}
		defer jsonFile.Close()
		byteValue, _ := io.ReadAll(jsonFile)
		environmentmanager.GenerateEnvFile(byteValue, newEnvFilePath)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	rootCmd.PersistentFlags().StringVar(&newEnvFilePath, "path", "", "Path for the new file folder, ex: /home/ubuntu/code/")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
