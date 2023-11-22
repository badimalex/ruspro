package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "ruspro/api" // Импортируйте сгенерированный клиентский код

	"google.golang.org/grpc"
)

func main() {
	// Устанавливаем соединение с gRPC сервером на порту 50051
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Не удалось установить соединение с сервером: %v", err)
	}
	defer conn.Close()

	// Создаем клиент
	client := pb.NewRusProfileServiceClient(conn)

	// Вызываем метод GetCompanyInfo с тестовыми данными (замените их на реальные)
	inn := "7736207543" // Замените на нужный ИНН
	req := &pb.CompanyRequest{
		Inn: inn,
	}

	// Устанавливаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Вызываем метод сервера
	resp, err := client.GetCompanyInfo(ctx, req)
	if err != nil {
		log.Fatalf("Ошибка при вызове метода GetCompanyInfo: %v", err)
	}

	// Выводим результат
	fmt.Printf("ИНН: %s\n", resp.GetInn())
	fmt.Printf("КПП: %s\n", resp.GetKpp())
	fmt.Printf("Название: %s\n", resp.GetName())
	fmt.Printf("ФИО руководителя: %s\n", resp.GetCeo())
}
