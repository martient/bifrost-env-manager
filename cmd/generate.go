package cmd

import (
	"fmt"
	"io"
	"os"

	environmentmanager "github.com/martient/bifrost-env-manager/pkg/environment_manager"
	"github.com/spf13/cobra"
)

var newEnvFilePath string
var readOnlyEnvFilesPath string

// generateCmd represents the load command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a new version of the env file",
	Long:  `Generate a new version of the environement file in function of the config given`,
	Run: func(cmd *cobra.Command, args []string) {
		if !disableUpdateCheck {
			doConfirmAndSelfUpdate()
		}
		jsonFile, err := os.Open(jsonConfigFile)
		if err != nil {
			fmt.Println(err)
		}
		defer jsonFile.Close()
		byteValue, _ := io.ReadAll(jsonFile)
		environmentmanager.GenerateEnvFile(byteValue, newEnvFilePath, readOnlyEnvFilesPath)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	rootCmd.PersistentFlags().StringVar(&newEnvFilePath, "path", "", "Path for the new file folder, ex: /home/ubuntu/code/")
	rootCmd.PersistentFlags().StringVar(&readOnlyEnvFilesPath, "read-only-env", "", "Path for read-only environement config, ex: \".api.env;.redis.env\"")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
