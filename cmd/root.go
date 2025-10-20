/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"cogmoteHub/internal/db"
	"cogmoteHub/internal/devices"
	"cogmoteHub/internal/logger"
	"log/slog"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cogmoteHub",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		Serve()
	},
}

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

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cogmoteHub.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func Serve() {
	logger.Init()

	host := os.Getenv("POSTGRES_SERVER")
	if host == "" {
		host = "localhost"
	}

	user := os.Getenv("POSTGRES_USER")
	dbName := os.Getenv("POSTGRES_USER")
	password := loadSecret("POSTGRES_PASSWORD", "POSTGRES_PASSWORD_FILE")
	if password == "" {
		slog.Warn("database password not provided")
	}

	db.Init(host, user, password, dbName)

	r := gin.Default()
	r.UseH2C = true

	api := r.Group("/api")

	devices.RegisterRoutes(api)

	r.Run(":9013")
}

func loadSecret(envKey, fileKey string) string {
	if value := os.Getenv(envKey); value != "" {
		return value
	}

	filePath := os.Getenv(fileKey)
	if filePath == "" {
		return ""
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		slog.Error("unable to read secret file", "error", err, "path", filePath)
		return ""
	}

	return strings.TrimSpace(string(data))
}
