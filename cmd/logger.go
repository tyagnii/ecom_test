/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/tyagnii/ecom_test/logger"
)

var (
	logLevel    string
	logFormat   string
	logOutput   string
	enableCaller bool
)

// loggerCmd represents the logger command
var loggerCmd = &cobra.Command{
	Use:   "logger",
	Short: "Logger configuration and testing",
	Long:  `Configure and test the logging system.`,
}

// loggerTestCmd represents the logger test command
var loggerTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Test logger with different levels",
	Long:  `Test the logger with different log levels and formats.`,
	Run: func(cmd *cobra.Command, args []string) {
		testLogger()
	},
}

// loggerConfigCmd represents the logger config command
var loggerConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Show logger configuration",
	Long:  `Show current logger configuration and settings.`,
	Run: func(cmd *cobra.Command, args []string) {
		showLoggerConfig()
	},
}

func init() {
	rootCmd.AddCommand(loggerCmd)
	loggerCmd.AddCommand(loggerTestCmd)
	loggerCmd.AddCommand(loggerConfigCmd)
	
	// Logger configuration flags
	loggerCmd.PersistentFlags().StringVar(&logLevel, "level", "INFO", "Log level (DEBUG, INFO, WARN, ERROR, FATAL)")
	loggerCmd.PersistentFlags().StringVar(&logFormat, "format", "structured", "Log format (structured, simple)")
	loggerCmd.PersistentFlags().StringVar(&logOutput, "output", "stdout", "Log output (stdout, stderr)")
	loggerCmd.PersistentFlags().BoolVar(&enableCaller, "caller", false, "Enable caller information in logs")
}

func testLogger() {
	fmt.Println("Testing Logger System")
	fmt.Println("====================")
	
	// Create logger based on configuration
	var testLogger logger.Logger
	
	switch logFormat {
	case "structured":
		structuredLogger := logger.NewStructuredLogger(parseLogLevel(logLevel), getOutput())
		if enableCaller {
			structuredLogger.EnableCaller()
		}
		testLogger = structuredLogger
	case "simple":
		simpleLogger := logger.NewSimpleLogger(parseLogLevel(logLevel))
		testLogger = simpleLogger
	default:
		log.Fatalf("Invalid log format: %s", logFormat)
	}
	
	// Test different log levels
	fmt.Println("\nTesting different log levels:")
	fmt.Println("----------------------------")
	
	testLogger.Debug("This is a debug message", 
		logger.NewField("test_id", 1),
		logger.NewField("component", "logger_test"))
	
	testLogger.Info("This is an info message", 
		logger.NewField("test_id", 2),
		logger.NewField("component", "logger_test"))
	
	testLogger.Warn("This is a warning message", 
		logger.NewField("test_id", 3),
		logger.NewField("component", "logger_test"))
	
	testLogger.Error("This is an error message", 
		logger.NewField("test_id", 4),
		logger.NewField("component", "logger_test"))
	
	// Test with context
	fmt.Println("\nTesting with context:")
	fmt.Println("-------------------")
	
	contextLogger := testLogger.WithContext("banner_service")
	contextLogger.Info("Banner operation", 
		logger.NewField("operation", "create_banner"),
		logger.NewField("banner_name", "Test Banner"))
	
	// Test with error
	fmt.Println("\nTesting with error:")
	fmt.Println("------------------")
	
	errorLogger := testLogger.WithError(fmt.Errorf("database connection failed"))
	errorLogger.Error("Failed to connect to database")
	
	// Test with fields
	fmt.Println("\nTesting with fields:")
	fmt.Println("------------------")
	
	fieldsLogger := testLogger.WithFields(
		logger.NewField("user_id", 123),
		logger.NewField("session_id", "abc123"),
		logger.NewField("request_id", "req-456"))
	
	fieldsLogger.Info("User action performed")
	
	fmt.Println("\nLogger test completed!")
}

func showLoggerConfig() {
	fmt.Println("Logger Configuration")
	fmt.Println("===================")
	fmt.Printf("Level: %s\n", logLevel)
	fmt.Printf("Format: %s\n", logFormat)
	fmt.Printf("Output: %s\n", logOutput)
	fmt.Printf("Caller: %t\n", enableCaller)
	
	// Show available levels
	fmt.Println("\nAvailable Log Levels:")
	fmt.Println("  DEBUG - Detailed information for debugging")
	fmt.Println("  INFO  - General information about program execution")
	fmt.Println("  WARN  - Warning messages for potential issues")
	fmt.Println("  ERROR - Error messages for recoverable errors")
	fmt.Println("  FATAL - Fatal errors that cause program termination")
	
	// Show available formats
	fmt.Println("\nAvailable Formats:")
	fmt.Println("  structured - JSON structured logging with fields")
	fmt.Println("  simple     - Simple text logging with timestamps")
	
	// Show available outputs
	fmt.Println("\nAvailable Outputs:")
	fmt.Println("  stdout - Standard output")
	fmt.Println("  stderr - Standard error")
}

func parseLogLevel(level string) logger.LogLevel {
	switch level {
	case "DEBUG":
		return logger.DEBUG
	case "INFO":
		return logger.INFO
	case "WARN":
		return logger.WARN
	case "ERROR":
		return logger.ERROR
	case "FATAL":
		return logger.FATAL
	default:
		return logger.INFO
	}
}

func getOutput() io.Writer {
	switch logOutput {
	case "stderr":
		return os.Stderr
	case "stdout":
		return os.Stdout
	default:
		return os.Stdout
	}
}

// Example usage function
func ExampleLoggerUsage() {
	fmt.Println("Logger Usage Examples")
	fmt.Println("====================")
	
	// Create a development logger
	devLogger := logger.NewDevelopmentLogger()
	devLogger.Info("Application started", 
		logger.NewField("version", "1.0.0"),
		logger.NewField("environment", "development"))
	
	// Create a production logger
	prodLogger := logger.NewProductionLogger()
	prodLogger.Info("Application started", 
		logger.NewField("version", "1.0.0"),
		logger.NewField("environment", "production"))
	
	// Use global logger
	logger.SetGlobalLogger(devLogger)
	logger.Info("Using global logger")
	
	// Create logger with context
	serviceLogger := devLogger.WithContext("banner_service")
	serviceLogger.Info("Banner created", 
		logger.NewField("banner_id", 1),
		logger.NewField("banner_name", "Summer Sale"))
	
	// Create logger with error
	errorLogger := devLogger.WithError(fmt.Errorf("database timeout"))
	errorLogger.Error("Database operation failed")
	
	// Create logger with multiple fields
	fieldsLogger := devLogger.WithFields(
		logger.NewField("user_id", 123),
		logger.NewField("request_id", "req-456"))
	fieldsLogger.Info("User action completed")
}
