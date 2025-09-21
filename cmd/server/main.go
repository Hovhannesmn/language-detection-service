package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"language-detection-service/internal/language_detection/application"
	"language-detection-service/internal/language_detection/domain"
	"language-detection-service/internal/language_detection/infrastructure/adapters"
	"language-detection-service/internal/language_detection/infrastructure/config"
	"language-detection-service/internal/language_detection/infrastructure/grpc"
)

// createSignalContext creates a context that gets cancelled on SIGINT or SIGTERM
func createSignalContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Received shutdown signal, cancelling context...")
		cancel()
	}()

	return ctx, cancel
}

func main() {
	// Load configuration
	configProvider := config.NewConfigProvider()

	// Validate configuration
	if err := configProvider.ValidateConfig(); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	cfg := configProvider.GetConfig()
	log.Printf("Starting Language Detection Service with configuration:")
	log.Printf("  Server Address: %s:%d", cfg.ServerAddress, cfg.ServerPort)
	log.Printf("  AWS Comprehend: %v", cfg.UseAWSComprehend)
	log.Printf("  AWS Region: %s", cfg.AWSRegion)
	log.Printf("  Max Text Length: %d", cfg.MaxTextLength)
	log.Printf("  Min Confidence: %.2f", cfg.MinConfidenceThreshold)
	log.Printf("  Supported Languages: %v", cfg.SupportedLanguages)

	// Create language detector based on configuration
	var detector domain.LanguageDetector
	var err error

	if cfg.UseAWSComprehend {
		detector, err = adapters.NewAWSComprehendAdapter(cfg.AWSRegion, 3)
		if err != nil {
			log.Printf("Warning: Failed to create AWS Comprehend adapter: %v", err)
			log.Printf("Falling back to pattern-based detection")
			detector = adapters.NewFallbackAdapter()
		} else {
			log.Println("Using AWS Comprehend for language detection")
		}
	} else {
		detector = adapters.NewFallbackAdapter()
		log.Println("Using fallback pattern-based language detection")
	}

	// Create application service
	service := application.NewLanguageDetectionService(detector, configProvider)

	// Create gRPC server
	grpcServer := grpc.NewServer(service)

	// Create server address
	address := fmt.Sprintf("%s:%d", cfg.ServerAddress, cfg.ServerPort)

	// Create context that will be cancelled on signal
	ctx, cancel := createSignalContext()
	defer cancel()

	// Start gRPC server in a goroutine with context support
	serverErr := make(chan error, 1)
	go func() {
		log.Printf("Starting gRPC server on %s", address)
		if err := grpcServer.StartWithContext(ctx, address); err != nil {
			serverErr <- fmt.Errorf("server failed to start: %w", err)
		}
	}()

	// Wait for server to start
	time.Sleep(2 * time.Second)

	// Wait for context cancellation (signal) or server error
	select {
	case err := <-serverErr:
		log.Fatalf("Server error: %v", err)
	case <-ctx.Done():
		log.Printf("Shutdown signal received, shutting down gracefully...")

		// Graceful shutdown with timeout
		_, shutdownCancel := context.WithTimeout(context.Background(),
			time.Duration(cfg.ShutdownTimeoutSeconds)*time.Second)
		defer shutdownCancel()

		// Stop the server
		if err := grpcServer.Stop(); err != nil {
			log.Printf("Error during server shutdown: %v", err)
		}

		log.Println("Language Detection Service stopped")
	}
}
