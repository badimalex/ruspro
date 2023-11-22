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

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func startGRPCServer(rusprofileAPIURL, grpcServerAddress string) {
	lis, err := net.Listen("tcp", grpcServerAddress)
	if err != nil {
		logging.Log.Fatalf("Failed to listen for gRPC server: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterRusProfileServiceServer(grpcServer, &server.Server{ApiUrl: rusprofileAPIURL})
	logging.Log.Infof("gRPC server started on %s", grpcServerAddress)
	if err := grpcServer.Serve(lis); err != nil {
		logging.Log.Fatalf("failed to serve gRPC server: %v", err)
	}
}

func startHTTPServer(grpcServerAddress, httpServerAddress string) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	httpMux := http.NewServeMux()
	httpMux.Handle("/", mux)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := pb.RegisterRusProfileServiceHandlerFromEndpoint(ctx, mux, "localhost"+grpcServerAddress, opts)
	if err != nil {
		logging.Log.Fatalf("failed to register gateway: %v", err)
	}

	httpMux.HandleFunc("/swagger.json", serveSwagger)

	httpMux.Handle("/swaggerui/", http.StripPrefix("/swaggerui/", http.FileServer(http.Dir("swaggerui"))))

	gwServer := &http.Server{
		Addr:    httpServerAddress,
		Handler: httpMux,
	}

	logging.Log.Printf("HTTP server started on :%s", httpServerAddress)
	logging.Log.Fatal(gwServer.ListenAndServe())
}

func serveSwagger(w http.ResponseWriter, r *http.Request) {
	logging.Log.Info("Serving Swagger JSON")
	http.ServeFile(w, r, "api/rusprofile.swagger.json")
}

func main() {
	var configFile string
	pflag.StringVar(&configFile, "config", "config.yaml", "Path to configuration file")

	pflag.Parse()

	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		logging.Log.Fatalf("Failed to read config file: %v", err)
	}

	// Чтение настроек из viper
	grpcServerAddress := viper.GetString("grpc_server_address")
	httpServerAddress := viper.GetString("http_server_address")
	rusprofileAPIURL := viper.GetString("rusprofile_api_url")

	logging.InitLogger()
	go startGRPCServer(rusprofileAPIURL, grpcServerAddress)
	startHTTPServer(grpcServerAddress, httpServerAddress)
}
