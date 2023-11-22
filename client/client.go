package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "ruspro/api"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Не удалось установить соединение с сервером: %v", err)
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
		log.Fatalf("Ошибка при вызове метода GetCompanyInfo: %v", err)
	}

	fmt.Printf("ИНН: %s\n", resp.GetInn())
	fmt.Printf("КПП: %s\n", resp.GetKpp())
	fmt.Printf("Название: %s\n", resp.GetName())
	fmt.Printf("ФИО руководителя: %s\n", resp.GetCeo())
}
