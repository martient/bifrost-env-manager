package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var jsonConfigFile string

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
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&jsonConfigFile, "config", "", "config file for this software environement")
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().BoolP("disable-update-check", "", false, "Disable auto update checking before execution")

	// rootCmd.PersistentFlags().CallbackVarP(runner.GetUpdateCallback(), "update", "up", "update notify to latest version")
	// rootCmd.PersistentFlags().BoolVarP(&options.DisableUpdateCheck, "disable-update-check", "duc", false, "disable automatic notify update check")
}

func initConfig() {
	if jsonConfigFile != "" {
		viper.SetConfigFile(jsonConfigFile)
	} else {
		os.Exit(1)
	}

	fmt.Println(rootCmd.Flag("disable-update-check").Value)

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
