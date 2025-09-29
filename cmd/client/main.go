package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/tempizhere/vaultfactory/internal/client/commands"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	// Создаём контекст с обработкой сигналов
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Обрабатываем сигналы прерывания
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nReceived interrupt signal, shutting down gracefully...")
		cancel()
	}()

	rootCmd := &cobra.Command{
		Use:   "vaultfactory",
		Short: "VaultFactory - Secure data storage and synchronization",
		Long:  "VaultFactory is a secure data storage system with client-server synchronization",
	}

	rootCmd.AddCommand(commands.NewVersionCommand(buildVersion, buildDate, buildCommit))
	rootCmd.AddCommand(commands.NewAuthCommands())
	rootCmd.AddCommand(commands.NewDataCommands())

	// Устанавливаем контекст для команды
	rootCmd.SetContext(ctx)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func getBuildInfo(value string) string {
	if value == "" {
		return "N/A"
	}
	return value
}
