package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/programmablemike/fibo/internal/router"
	"github.com/programmablemike/fibo/internal/tracing"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile = ""
)

var calculateCmd = &cobra.Command{
	Use:   "calculate N",
	Short: "Calculates the Fibonacci number for the given ordinal N",
	Long:  `Calculates the Fibonacci number for the given ordinal N`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		host := viper.GetString("host")
		port := viper.GetInt("port")

		closer := tracing.SetupTracing("fibo-client")
		defer closer.Close()
		tracer := opentracing.GlobalTracer()
		span := tracer.StartSpan("calculate")
		defer span.Finish()

		client := &http.Client{}
		uri := fmt.Sprintf("http://%s:%d/fibo/calculate/%s", host, port, args[0])
		req, _ := http.NewRequest("GET", uri, nil)

		// inject tracing headers to match up client requests with server responses
		tracer.Inject(span.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(req.Header))

		res, err := client.Do(req)
		if err != nil {
			log.Fatalf("error: %s\n", err)
		}
		defer res.Body.Close()

		v := router.GenericResponse{}
		err = json.NewDecoder(res.Body).Decode(&v)
		if err != nil {
			log.Fatalf("error: failed to decode res.Body, %s\n", err)
		}
		fmt.Printf("Fibonacci number: %s\n", v.Value)
	},
}

var countCmd = &cobra.Command{
	Use:   "count NUM",
	Short: "Counts the number of ordinals in the Fibonacci value range (0, NUM)",
	Long:  `Counts the number of ordinals in the Fibonacci value range (0, NUM)`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		host := viper.GetString("host")
		port := viper.GetInt("port")

		closer := tracing.SetupTracing("fibo-client")
		defer closer.Close()
		tracer := opentracing.GlobalTracer()
		span := tracer.StartSpan("count")
		defer span.Finish()

		client := &http.Client{}
		uri := fmt.Sprintf("http://%s:%d/fibo/count/%s", host, port, args[0])
		req, _ := http.NewRequest("GET", uri, nil)

		// inject tracing headers to match up client requests with server responses
		tracer.Inject(span.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(req.Header))

		res, err := client.Do(req)
		if err != nil {
			log.Fatalf("error: %s\n", err)
		}
		defer res.Body.Close()

		v := router.GenericResponse{}
		err = json.NewDecoder(res.Body).Decode(&v)
		if err != nil {
			log.Fatalf("error: failed to decode res.Body, %s\n", err)
		}
		fmt.Printf("Ordinals in this range: %s\n", v.Value)
	},
}

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clears the memoizer cache",
	Long:  `Clears the memoizer cache`,
	Run: func(cmd *cobra.Command, args []string) {
		host := viper.GetString("host")
		port := viper.GetInt("port")

		closer := tracing.SetupTracing("fibo-client")
		defer closer.Close()
		tracer := opentracing.GlobalTracer()
		span := tracer.StartSpan("clear")
		defer span.Finish()

		client := &http.Client{}
		uri := fmt.Sprintf("http://%s:%d/fibo/cache", host, port)
		req, _ := http.NewRequest("DELETE", uri, nil)

		// inject tracing headers to match up client requests with server responses
		tracer.Inject(span.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(req.Header))

		res, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
			return
		}
		defer res.Body.Close()

		v := router.GenericResponse{}
		err = json.NewDecoder(res.Body).Decode(&v)
		if err != nil {
			log.Fatalf("error: failed to decode res.Body, %s\n", err)
		}
		fmt.Println("Successfully cleared cache")
	},
}

var rootCmd = &cobra.Command{
	Use:   "fibo",
	Short: "Fibo is an API server and CLI client for generating Fibonacci sequences",
	Long: `Fibo is an API server and CLI client for generating Fibonacci sequences
It uses dynamic programming techniques (memoization) to speed up processing.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Print the command line options for debugging purposes
		log.Debugf("host: %s", viper.GetString("host"))
		log.Debugf("port: %d", viper.GetInt("port"))
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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.fiborc)")
	rootCmd.PersistentFlags().Bool("debug", false, "Turns on debugging mode")
	rootCmd.PersistentFlags().String("host", "localhost", "HTTP server hostname to bind (default: localhost)")
	rootCmd.PersistentFlags().Int("port", 8080, "HTTP server port to bind (default: 8080)")
	rootCmd.AddCommand(calculateCmd, countCmd, clearCmd)
	viper.BindPFlag("host", rootCmd.PersistentFlags().Lookup("host"))
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
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
		fmt.Println("Can't find config:", err)
		fmt.Println("Falling back to command-line defaults.")
	}

	// Turn on debug is toggled
	if viper.GetBool("debug") {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}
