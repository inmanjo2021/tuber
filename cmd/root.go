package cmd

import (
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"strings"
)

var (
	rootCmd = &cobra.Command{
		Use:   "tuber",
		Short: "",
	}
)

func init() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	rootCmd.PersistentFlags().BoolP("debug", "d", false, "debug")
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
}

func createLogger() (logger *zap.Logger) {
	if viper.GetBool("debug") {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}
	return
}

// Execute executes
func Execute() error {
	err := godotenv.Load()
	if err == nil {
		log.Println(".env file loaded")
	}

	return rootCmd.Execute()
}
