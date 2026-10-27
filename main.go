package main

import (

	"fmt"
	"log"

	"os"

	

	"github.com/spf13/cobra"
)





func main(){

	var port int

	var rootCmd = &cobra.Command{
		Use: "bytecache",
		Short: "ByteCache is a high-performance in-memory key-value database",
	}

	var startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start the ByteCache TCP server",
		Run: func(cmd *cobra.Command, args []string) {
			cache := NewCache()
			addr := fmt.Sprintf(":%d", port)
			server := NewServer(addr, cache)
			log.Printf("Starting ByteCache on port %d...", port)
			if err := server.Start(); err != nil {
				log.Fatalf("Server failed: %v", err)
			}
		},
	}

	// Bind --port / -p flag to start command, default to 6379 (Redis default port!)

	startCmd.Flags().IntVarP(&port, "port", "p", 6379, "Port to listen on")

	rootCmd.AddCommand(startCmd)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}


}
