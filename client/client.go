package main

import (
	"context"
	"fmt"
	"time"

	pb "ruspro/api"
	"ruspro/internal/logging"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	logging.InitLogger()

	var configFile string
	pflag.StringVar(&configFile, "config", "config.yaml", "Path to configuration file")

	pflag.Parse()

	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		logging.Log.Fatalf("Failed to read config file: %v", err)
	}

	grpcServerAddress := "localhost" + viper.GetString("grpc_server_address")

	conn, err := grpc.Dial(grpcServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logging.Log.Fatalf("Не удалось установить соединение с сервером: %v", err)
	}
	defer conn.Close()

	client := pb.NewRusProfileServiceClient(conn)

	inn := "7736207543"
	req := &pb.CompanyRequest{
		Inn: inn,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.GetCompanyInfo(ctx, req)
	if err != nil {
		logging.Log.Fatalf("Ошибка при вызове метода GetCompanyInfo: %v", err)
	}

	fmt.Printf("ИНН: %s\n", resp.GetInn())
	fmt.Printf("КПП: %s\n", resp.GetKpp())
	fmt.Printf("Название: %s\n", resp.GetName())
	fmt.Printf("ФИО руководителя: %s\n", resp.GetCeo())
}
