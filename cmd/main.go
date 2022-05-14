package main

import (
	"fmt"
	"log"
	"rag/internal"
	"rag/internal/repository"
	"rag/internal/server"
	"rag/internal/service"
)

func main() {
	// init settings
	settings, err := internal.NewSettings()
	if err != nil {
		log.Panic(fmt.Errorf("internal.NewSettings(), err: %w", err))
	}
	// init logger
	logger, err := internal.NewLogger()
	if err != nil {
		log.Panic(fmt.Errorf("internal.NewLogger(), err: %w", err))
	}
	// init repository
	db, err := repository.NewRepository(settings, logger)
	if err != nil {
		log.Panic(fmt.Errorf("repository.NewRepository(), err: %w", err))
	}
	// init service
	userService, err := service.NewUserService(db, settings, logger)
	if err != nil {
		log.Panic(fmt.Errorf("service.NewUserService(), err: %w", err))
	}
	//init server
	srv := server.NewServer(userService, logger, settings)
	if err != nil {
		log.Panic(fmt.Errorf("server.NewServer(), err: %w", err))
	}
	// run server
	log.Printf("Listening and serving HTTP on %s\n", settings.Database.PORT)
	log.Panic(srv.ListenAndServe())
}
