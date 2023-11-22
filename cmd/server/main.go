package main

import (
	"context"
	"net"
	"net/http"
	pb "ruspro/api"
	"ruspro/internal/logging"
	"ruspro/server"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func startGRPCServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logging.Log.Fatalf("Failed to listen for gRPC server: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterRusProfileServiceServer(grpcServer, &server.Server{})
	logging.Log.Info("gRPC server started on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		logging.Log.Fatalf("failed to serve gRPC server: %v", err)
	}
}

func startHTTPServer() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	httpMux := http.NewServeMux()
	httpMux.Handle("/", mux)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := pb.RegisterRusProfileServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts)
	if err != nil {
		logging.Log.Fatalf("failed to register gateway: %v", err)
	}

	httpMux.HandleFunc("/swagger.json", serveSwagger)

	httpMux.Handle("/swaggerui/", http.StripPrefix("/swaggerui/", http.FileServer(http.Dir("swaggerui"))))

	gwServer := &http.Server{
		Addr:    ":8080",
		Handler: httpMux,
	}

	logging.Log.Printf("HTTP server started on :8080")
	logging.Log.Fatal(gwServer.ListenAndServe())
}

func serveSwagger(w http.ResponseWriter, r *http.Request) {
	logging.Log.Info("Serving Swagger JSON")
	http.ServeFile(w, r, "api/rusprofile.swagger.json")
}

func main() {
	logging.InitLogger()
	go startGRPCServer()
	startHTTPServer()
}
