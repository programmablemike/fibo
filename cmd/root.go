package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile = ""
)

var rootCmd = &cobra.Command{
	Use:   "fibo",
	Short: "Fibo is an API server and CLI client for generating Fibonacci sequences",
	Long: `Fibo is an API server and CLI client for generating Fibonacci sequences
				 It uses dynamic programming techniques (memoization) to speed up processing.
				`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.fibo.yaml)")
	rootCmd.PersistentFlags().Bool("debug", false, "Turns on debugging mode")
	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
}

func initConfig() {
	// Config file name (without extension - valid extensions are .yaml, .toml, .json, etc.)
	viper.SetConfigName(".fiborc")
	// Default configuration type when extension is missing
	viper.SetConfigType("toml")
	// Prefix for environment variables (ex. FIBO_POSTGRES_USER)
	viper.SetEnvPrefix("fibo")
	// Automatically source environment variables
	viper.AutomaticEnv()

	// Don't forget to read config either from cfgFile or from home directory!
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// Find the current working directory
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search the home directory
		viper.AddConfigPath(home)
		// Search the current working directory
		viper.AddConfigPath(cwd)
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}

	// Turn on debug is toggled
	if viper.GetBool("debug") {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}
