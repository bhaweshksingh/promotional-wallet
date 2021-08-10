package main

import (
	"account/pkg/app"
	"account/pkg/repository"
	"fmt"
	"log"
)

const (
	httpServeCommand = "http-serve"
	migrateCommand   = "migrate"
	rollbackCommand  = "rollback"
)

func commands() map[string]func(configFile string) {
	return map[string]func(configFile string){
		httpServeCommand: app.StartHTTPServer,
		migrateCommand:   repository.RunMigrations,
		rollbackCommand:  repository.RollBackMigrations,
	}
}

func execute(cmd string, configFile string) {
	fmt.Println("config : " + configFile)
	run, ok := commands()[cmd]
	if !ok {
		log.Fatal("invalid command")
	}

	run(configFile)
}
