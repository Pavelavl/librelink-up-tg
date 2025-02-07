package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"librelink-up-tg/config"
	"librelink-up-tg/internal/clients/libre"
	"librelink-up-tg/internal/clients/tg"
)

func main() {
	// Чтение конфигурации
	config, err := config.Read("config/config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	if !libre.IsValidRegion(config.LinkUpRegion) {
		log.Fatal("Region is not valid")
	}

	client, err := libre.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Login(); err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}

	// Создаем контекст для graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Обработка сигналов для graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Запускаем джобу в отдельной горутине
	go func() {
		ticker := time.NewTicker(time.Duration(config.LinkUpTimeInterval) * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("Shutting down job...")
				return
			case <-ticker.C:
				// Выполнение задачи
				data, err := client.GetGlucoseData()
				if err != nil {
					log.Printf("Failed to get glucose data: %v", err)
					continue
				}

				if data.IsBullshit() {
					tgClient := tg.NewClient(config)

					if err := tgClient.SendToFriends(data); err != nil {
						log.Printf("Failed to send data to friends: %v", err)
					}
				} else {
					log.Printf("Data is not bullshit, notification skipped")
				}
			}
		}
	}()

	// Ожидание сигнала для завершения
	<-signalChan
	log.Println("Received shutdown signal, exiting...")
	cancel()
}
