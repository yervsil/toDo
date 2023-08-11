package main

import (
	"log"
	"time"

	"github.com/yervsil/toDo-microservice/config"
	handler "github.com/yervsil/toDo-microservice/internal/delivery/http"
	"github.com/yervsil/toDo-microservice/internal/repository"
	"github.com/yervsil/toDo-microservice/internal/server"
	"github.com/yervsil/toDo-microservice/internal/service"
	"github.com/yervsil/toDo-microservice/pkg/database/mongodb"
	"github.com/yervsil/toDo-microservice/pkg/logger"
)

// @title Todo App API
// @version 1.0
// @description API Server for TodoList Application

// @host localhost:8000
// @BasePath /

func main() {
	cfg, err := config.InitConfig()
	time.Sleep(5 * time.Minute)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	l := logger.New(cfg.Env)

	client, err := mongodb.NewClient(cfg.Mongo.URI, cfg.Mongo.User, cfg.Mongo.Password)
	if err != nil {
		l.Fatal(err)
	}

	db := client.Database(cfg.Mongo.Name)

	repository := repository.NewRepository(db)
	service := service.NewService(repository)
	handler := handler.NewHandler(service, l)

	srv := server.NewServer(cfg, handler.InitRoutes())
	if err = srv.Run(); err != nil {
		l.Fatal(err)
	}
}