package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"listes_back/src/database"
	"listes_back/src/users"
	"log"
	"net/http"
	"os"
)

func main() {

	fmt.Println("\nStarting ...\n")

	configPath := flag.String("config", "", "The path of the configuration file")
	flag.Parse()

	if len(*configPath) == 0 {
		log.Fatal("No configuration file path provided !")
	}

	fileContent, err := os.ReadFile(*configPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("file %q not found", *configPath)
		}
		log.Fatal(err)
	}
	var config Config
	err = yaml.Unmarshal(fileContent, &config)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Loaded config:")
	fmt.Println(config)

	db, err := database.New(config.Database)
	if err != nil {
		log.Fatal(err)
	}
	database.SetCommonDb(db)

	err = users.InitAvatarsDir(config.Server.AvatarsDirectory)
	if err != nil {
		log.Fatal("failed to init avatars directory: ", err)
	}

	fmt.Println("Registering routes ...")
	router := initRoutes()

	fmt.Println("Starting server on port", config.Server.Port, "...")

	log.Fatal(http.ListenAndServe(config.Server.GetStringAddress(), router))
}
