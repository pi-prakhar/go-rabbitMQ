package main

import (
	"go-rabbitmq-server2/internal/models"
	"log"
)

const (
	CONFIG_FILE = "config/config.yml"
)

func main() {
	ConfigData := &models.ConfigData{}
	if err := ConfigData.LoadConfig(CONFIG_FILE); err != nil {
		log.Fatal("Error: Loading config file", err)
	}

	config := ConfigData.GetConfig()

	log.Println(config.Port)

	// srv := &http.Server{
	// 	Addr:         config.Port,
	// 	Handler:      api.NewRouter(),
	// 	ReadTimeout:  time.Second * time.Duration(config.ReadTimeout),
	// 	WriteTimeout: time.Second * time.Duration(config.WriteTimeout),
	// }

	// if err := srv.ListenAndServe(); err != nil {
	// 	log.Fatalf("Failed to start server at port %s", config.Port)
	// }
}
