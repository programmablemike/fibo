package cmd

import (
	"fmt"
	"net/http"

	"github.com/programmablemike/fibo/internal/router"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	serverCmd.PersistentFlags().String("host", "localhost", "HTTP server hostname to bind (default: localhost)")
	serverCmd.PersistentFlags().Int("port", 8080, "HTTP server port to bind (default: 8080)")
	serverCmd.PersistentFlags().String("pguser", "fibo", "Postgres database user (default: fibo)")
	serverCmd.PersistentFlags().String("pgpassword", "", "Postgres database password (default: \"\")")
	serverCmd.PersistentFlags().String("pghost", "localhost", "Postgres database hostname (default: localhost)")
	serverCmd.PersistentFlags().Int("pgport", 5432, "Postgres database port (default: 5432)")
	serverCmd.PersistentFlags().String("pgdb", "fibo", "Postgres database name (default: fibo)")
	viper.BindPFlag("host", serverCmd.PersistentFlags().Lookup("host"))
	viper.BindPFlag("port", serverCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("pguser", serverCmd.PersistentFlags().Lookup("pguser"))
	viper.BindPFlag("pgpassword", serverCmd.PersistentFlags().Lookup("pgpassword"))
	viper.BindPFlag("pghost", serverCmd.PersistentFlags().Lookup("pghost"))
	viper.BindPFlag("pgport", serverCmd.PersistentFlags().Lookup("pgport"))
	viper.BindPFlag("pgdb", serverCmd.PersistentFlags().Lookup("pgdb"))
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the API server for memoized Fibonacci generation",
	Long:  `Run the API server for memoized Fibonacci generation`,
	Run: func(cmd *cobra.Command, args []string) {
		// Print the command line options for debugging purposes
		log.Debugf("host: %s", viper.GetString("host"))
		log.Debugf("port: %d", viper.GetInt("port"))
		log.Debugf("pguser: %s", viper.GetString("pguser"))
		log.Debugf("pghost: %s", viper.GetString("pghost"))
		log.Debugf("pgport: %s", viper.GetString("pgport"))
		log.Debugf("pgdb: %s", viper.GetString("pgdb"))

		r := router.NewRouter()
		addr := fmt.Sprintf("%s:%d", viper.GetString("host"), viper.GetInt("port"))
		log.Info("Started server at ", addr)
		http.ListenAndServe(addr, r)
	},
}
