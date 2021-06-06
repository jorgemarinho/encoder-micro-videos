package main

import (
	"encoder/application/services"
	"encoder/framework/database"
	"encoder/framework/queue"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"log"
	"os"
	"strconv"
)

var db database.Database

func init() {

	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	automMigrateDb, err := strconv.ParseBool(os.Getenv("AUTO_MIGRATE_DB"))

	if err != nil {
		log.Fatalf("Error parsing boolean env var")
	}

	debug, err := strconv.ParseBool(os.Getenv("DEBUG"))

	if err != nil {
		log.Fatalf("Error parsing boolean env debug var")
	}

	db.AutomMigrateDb = automMigrateDb
	db.Debug = debug
	db.DsnTest = os.Getenv("DSN_TEST")
	db.Dsn = os.Getenv("DSN")
	db.DbTypeTest = os.Getenv("DB_TYPE_TEST")
	db.DbType = os.Getenv("DB_TYPE")
	db.Env = os.Getenv("ENV")
}

func main() {

	messageChannel := make(chan amqp.Delivery)
	jobReturnChanel := make(chan services.JobWorkerResult)

	dbConnection, err := db.Connect()

	if err != nil {
		log.Fatalf("Error connecting to DB")
	}

	defer dbConnection.Close()

	rabbitMQ := queue.NewRabbitMQ()
	ch := rabbitMQ.Connect()
	defer ch.Close()

	rabbitMQ.Consume(messageChannel)

	jobManager := services.NewJobManager(dbConnection, rabbitMQ, jobReturnChanel, messageChannel)
	jobManager.Start(ch)
}
