package main

import (
	"github.com/spf13/cobra"
	"yumiko_kawaii.com/yine/applications/orchestrator/pkg/logger"
	"yumiko_kawaii.com/yine/applications/orchestrator/serve"
)

func main() {
	cmd := &cobra.Command{
		Use: "rpc-runtime",
	}

	cmd.AddCommand(&cobra.Command{
		Use: "server",
		Run: serve.Serve,
	})

	if err := cmd.Execute(); err != nil {
		logger.Fatalf("failed to execute command: %v", err)
	}
}
