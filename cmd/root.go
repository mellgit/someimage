package cmd

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "someimage",
		Short: "Some image downloader",
		Long:  `Saves images`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				log.WithFields(log.Fields{
					"component": "rootCmd",
					"err":       err,
				}).Fatal("CLI error")
				os.Exit(1)
			}
		},
	}
	configPath string
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	rootCmd.PersistentFlags().StringVarP(&configPath,
		"config",
		"c",
		"./config.yml",
		"Path to config file")
	cobra.OnInitialize(loadConfig)

}

func loadConfig() {
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		log.WithFields(log.Fields{
			"function":  "loadConfig",
			"component": "viper",
			"err":       err,
		}).Fatal("Failed to load config file")
	}
	// Loading specific settings
	if err := setUpLogger(); err != nil {
		log.WithFields(log.Fields{
			"function":  "loadConfig",
			"component": "setUpLogger",
			"err":       err,
		}).Fatal("Failed to set up logger")
	}
}

func setUpLogger() error {

	level, err := log.ParseLevel(viper.GetString("logging.level"))
	if err != nil {
		level = log.InfoLevel
	}
	log.SetLevel(level)
	if level != log.InfoLevel {
		log.SetReportCaller(true) // display filename and line number
	}

	// Setting log formatter
	formatter := viper.GetString("logging.formatter")
	switch formatter {
	case "text":
		log.SetFormatter(&log.TextFormatter{
			DisableLevelTruncation: true,
			FullTimestamp:          true,
			DisableColors:          true,
		})
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	default:
		return fmt.Errorf("unsupported log formatter: %s", formatter)
	}

	// Setting log handler
	var handler io.Writer
	switch viper.GetString("logging.handler") {
	case "file":
		handler = &lumberjack.Logger{
			Filename: filepath.Join(viper.GetString("logging.path"), "entry.log"),
			MaxSize:  128,
		}
	case "console":
		handler = os.Stdout
	default:
		return fmt.Errorf("unsupported log handler type %s", viper.GetString("logging.handler"))
	}
	log.SetOutput(handler)

	return nil

}
