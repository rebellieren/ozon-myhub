package main

import (
	"log"
	"myhub/db"
	"myhub/internal/config"
	storage "myhub/internal/storage"
	"myhub/server"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ .env файл не найден, использую системные переменные")
	}
	config.LoadEnv()
	dbInstance := db.ConnectDB()
	db.AutoMigrate(dbInstance)
	log.Println("USE_POSTGRES:", os.Getenv("USE_POSTGRES"))
	var dataStore storage.DataStore
	if os.Getenv("USE_POSTGRES") == "true" {
		log.Println("В качестве основного хранилища был выбран PostgreSql")
		database := db.ConnectDB()
		dataStore = storage.NewPostgresDataStore(database)
		db.CreateTestUsersIfNotExist(database)
	} else {
		log.Println("В качестве основного хранилища было выбрано локальное")
		inMemoryStore := storage.NewInMemoryDataStore()
		inMemoryStore.CreateTestUsersIfNotExist()
		dataStore = inMemoryStore
	}

	log.Println("Starting GraphQL server...")
	server.StartServer(dataStore)
}
