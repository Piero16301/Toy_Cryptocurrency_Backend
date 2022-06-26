package configs

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(EnvMongoURI()))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Hacer un ping a la BD
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Conectado a Base de Datos Blockchain de MongoDB")
	return client
}

// DB Instancia de Cliente
var DB = ConnectDB()

// GetCollection Obtener una colleción de la BD
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("Cryptocurrency").Collection(collectionName)
	return collection
}
