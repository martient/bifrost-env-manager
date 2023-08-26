package cmd

import (
	"github.com/martient/golang-utils/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var jsonConfigFile string
var BEMversion string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "Bifrost-env-manager",
	Short: "Env files manager",
	Long:  `Software environement files manager`,
	// examples and usage of using your application. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(versionFormated string, version string) {
	rootCmd.Version = versionFormated
	BEMversion = version
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SetVersionTemplate(`{{with .Name}}{{printf   .}}{{end}} {{printf  .Version}}`)
	rootCmd.PersistentFlags().StringVar(&jsonConfigFile, "config", "config.json", "config file for this software environement")
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().BoolP("yes", "y", false, "Auto accept manual question y/n")
	rootCmd.PersistentFlags().BoolP("disable-update-check", "", false, "Disable auto update checking before execution")
}

func initConfig() {
	if jsonConfigFile != "" {
		viper.SetConfigFile(jsonConfigFile)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		utils.LogInfo("Using config file: %s", viper.ConfigFileUsed(), "CLI")
	}
}
