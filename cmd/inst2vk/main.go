package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/s-yakubovskiy/inst2any/pkg/config"
	"github.com/s-yakubovskiy/inst2any/pkg/daemon"
	"github.com/s-yakubovskiy/inst2any/pkg/db"
	"github.com/s-yakubovskiy/inst2any/pkg/instagram"
	"github.com/s-yakubovskiy/inst2any/pkg/storage"
	"github.com/s-yakubovskiy/inst2any/pkg/vk"
)

func main() {
	// parse flags
	configFile := flag.String("config", "./configs/config.yaml", "Configuration file path")
	flag.Parse()

	// load config
	cfg, err := config.Load(*configFile)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	fmt.Printf("[config] %+v\n", cfg)

	// Setup the database
	database, err := db.SetupDB(cfg.Database.DSN)
	if err != nil {
		log.Fatalf("Failed to setup database: %v", err)
	}
	defer database.Close()

	// Setup GCS
	gcsClient, err := storage.NewGCS(cfg.GCS.BucketName, cfg.GCS.CredentialsFilePath)
	if err != nil {
		log.Fatalf("Failed to initialize GCS client: %v", err)
	}

	// Setup instagram metaClient
	metaClient := instagram.NewClient(cfg.Instagram)

	// Setup vk apiClient
	vkClient := vk.NewClient(cfg.VK)

	// Create the Daemon
	mediaWorker := daemon.NewMediaWorker(cfg, database, gcsClient, metaClient, vkClient)
	storyWorker := daemon.NewStoryWorker(cfg, database, gcsClient, metaClient, vkClient)

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // ensure all paths cancel the context to avoid context leak

	// Create Daemon and start
	daemon := daemon.NewDaemon(mediaWorker, storyWorker)
	daemon.Start(ctx)

	// Handle SIGINT and SIGTERM.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-quit

	log.Println("Shutting down...")
}
