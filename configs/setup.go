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
	serverOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(EnvMongoURI()).SetServerAPIOptions(serverOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
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
