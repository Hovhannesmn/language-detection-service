package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	pb "github.com/hovman/ld-proto/pb"
	"language-detection-service/internal/language_detection/domain"
)

// Server represents the gRPC server for language detection
type Server struct {
	pb.UnimplementedLanguageDetectionServiceServer
	service         domain.LanguageDetectionService
	healthServer    *health.Server
	server          *grpc.Server
	shutdownTimeout time.Duration
}

// NewServer creates a new gRPC server
func NewServer(service domain.LanguageDetectionService, opts ...grpc.ServerOption) *Server {
	// Add default options
	opts = append(opts,
		grpc.ConnectionTimeout(30*time.Second),
		grpc.MaxRecvMsgSize(4*1024*1024), // 4MB
		grpc.MaxSendMsgSize(4*1024*1024), // 4MB
	)

	server := grpc.NewServer(opts...)
	healthServer := health.NewServer()

	// Register services
	pb.RegisterLanguageDetectionServiceServer(server, &Server{
		service:         service,
		healthServer:    healthServer,
		server:          server,
		shutdownTimeout: 30 * time.Second,
	})

	// Register health service
	grpc_health_v1.RegisterHealthServer(server, healthServer)

	// Enable reflection for debugging
	reflection.Register(server)

	// Set health status
	healthServer.SetServingStatus("language_detection.LanguageDetectionService", grpc_health_v1.HealthCheckResponse_SERVING)

	return &Server{
		service:         service,
		healthServer:    healthServer,
		server:          server,
		shutdownTimeout: 30 * time.Second,
	}
}

// StartWithContext starts the gRPC server with context support
func (s *Server) StartWithContext(ctx context.Context, address string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", address, err)
	}

	log.Printf("Starting gRPC Language Detection Service on %s", address)

	// Set health status to serving
	s.healthServer.SetServingStatus("language_detection.LanguageDetectionService", grpc_health_v1.HealthCheckResponse_SERVING)

	// Start server in a goroutine to allow context cancellation
	serverErr := make(chan error, 1)
	go func() {
		serverErr <- s.server.Serve(lis)
	}()

	// Wait for context cancellation or server error
	select {
	case <-ctx.Done():
		log.Println("Context cancelled, stopping gRPC server...")
		s.server.Stop()
		return ctx.Err()
	case err := <-serverErr:
		return err
	}
}

// Stop gracefully stops the gRPC server
func (s *Server) Stop() error {
	log.Println("Shutting down gRPC Language Detection Service...")

	// Set health status to not serving
	s.healthServer.SetServingStatus("language_detection.LanguageDetectionService", grpc_health_v1.HealthCheckResponse_NOT_SERVING)

	// Graceful shutdown
	stopped := make(chan struct{})
	go func() {
		s.server.GracefulStop()
		close(stopped)
	}()

	// Wait for graceful shutdown or timeout
	select {
	case <-stopped:
		log.Println("gRPC server stopped gracefully")
		return nil
	case <-time.After(s.shutdownTimeout):
		log.Println("gRPC server shutdown timeout, forcing stop")
		s.server.Stop()
		return fmt.Errorf("server shutdown timeout")
	}
}

// DetectLanguage implements the DetectLanguage gRPC method
func (s *Server) DetectLanguage(
	ctx context.Context,
	req *pb.DetectLanguageRequest,
) (*pb.DetectLanguageResponse, error) {
	// Convert protobuf request to domain request
	domainReq := &domain.LanguageDetectionRequest{
		Text:       domain.Text(req.Text),
		DocumentID: req.DocumentId,
		Metadata:   req.Metadata,
	}

	// Call the application service
	domainResp, err := s.service.DetectLanguage(ctx, domainReq)
	if err != nil {
		return nil, fmt.Errorf("language detection failed: %w", err)
	}

	// Convert domain response to protobuf response
	return s.convertToProtobufResponse(domainResp), nil
}

// convertToProtobufResponse converts domain response to protobuf response
func (s *Server) convertToProtobufResponse(resp *domain.LanguageDetectionResponse) *pb.DetectLanguageResponse {
	var alternatives []*pb.LanguageAlternative
	for _, alt := range resp.Alternatives {
		alternatives = append(alternatives, &pb.LanguageAlternative{
			LanguageCode: string(alt.LanguageCode),
			Confidence:   float32(alt.Confidence),
		})
	}

	return &pb.DetectLanguageResponse{
		LanguageCode: string(resp.LanguageCode),
		Confidence:   float32(resp.Confidence),
		Alternatives: alternatives,
		DocumentId:   resp.DocumentID,
		Metadata: &pb.ProcessingMetadata{
			ProcessingTimeMs: resp.Metadata.ProcessingTimeMs,
			ServiceVersion:   resp.Metadata.ServiceVersion,
			ModelVersion:     resp.Metadata.ModelVersion,
			Provider:         resp.Metadata.Provider,
		},
	}
}
