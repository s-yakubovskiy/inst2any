package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/s-yakubovskiy/inst2any/pkg/config"
	"github.com/s-yakubovskiy/inst2any/pkg/server"
	"github.com/s-yakubovskiy/inst2any/pkg/vk"
)

func main() {
	cfg, err := config.Load("./configs/config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	vkClient := vk.NewClient(cfg.VK)
	vkService := &vk.VideoService{Client: vkClient}

	mux := server.NewServer(vkService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Starting server at :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
