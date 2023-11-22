package main

import (
	"context"
	"log"
	"net"
	"net/http"
	pb "ruspro/api"
	"ruspro/server"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func startGRPCServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterRusProfileServiceServer(grpcServer, &server.Server{})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func startHTTPServer() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	httpMux := http.NewServeMux()
	httpMux.Handle("/", mux) // grpc-gateway
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := pb.RegisterRusProfileServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts)
	if err != nil {
		log.Fatalf("failed to register gateway: %v", err)
	}

	httpMux.HandleFunc("/swagger.json", serveSwagger)

	// Обслуживание статических файлов Swagger UI
	httpMux.Handle("/swaggerui/", http.StripPrefix("/swaggerui/", http.FileServer(http.Dir("swaggerui"))))

	// Запуск HTTP сервера
	gwServer := &http.Server{
		Addr:    ":8080",
		Handler: httpMux,
	}

	log.Println("Сервер запущен на :8080")
	log.Fatal(gwServer.ListenAndServe())
}

func serveSwagger(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "api/rusprofile.swagger.json")
}

func main() {
	go startGRPCServer()
	startHTTPServer()
}
