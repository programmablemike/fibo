package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	agentCmd.PersistentFlags().String("host", "localhost", "HTTP server hostname to bind (default: localhost)")
	agentCmd.PersistentFlags().Int("port", 8080, "HTTP server port to bind (default: 8080)")
	viper.BindPFlag("host", agentCmd.PersistentFlags().Lookup("host"))
	viper.BindPFlag("port", agentCmd.PersistentFlags().Lookup("port"))
	rootCmd.AddCommand(agentCmd)
}

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "CLI agent to call the Fibo API server",
	Long:  `CLI agent used to call the Fibo API server`,
	Run: func(cmd *cobra.Command, args []string) {
		// Print the command line options for debugging purposes
		log.Debugf("host: %s", viper.GetString("host"))
		log.Debugf("port: %d", viper.GetInt("port"))
	},
}
