package controllers

import (
	"Toy_Cryptocurrency/configs"
	"Toy_Cryptocurrency/functions"
	"Toy_Cryptocurrency/models"
	"Toy_Cryptocurrency/responses"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strings"
	"time"
)

var blockCollection = configs.GetCollection(configs.DB, "Blockchain")
var validateTransaction = validator.New()

func NewTransaction() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		var transaction models.Transaction
		defer cancel()

		// Validar que el body está en formato JSON
		if err := json.NewDecoder(request.Body).Decode(&transaction); err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			response := responses.BlockResponse{
				Status:  http.StatusBadRequest,
				Message: "Formato del contenido de la solicitud no válido",
				Data:    err.Error(),
			}
			_ = json.NewEncoder(writer).Encode(response)
			return
		}

		// Se usa la librería para validar los campos del body
		if validationErr := validateTransaction.Struct(&transaction); validationErr != nil {
			writer.WriteHeader(http.StatusBadRequest)
			response := responses.BlockResponse{
				Status:  http.StatusBadRequest,
				Message: "Campos del contenido de la solicitud no válidos",
				Data:    validationErr.Error(),
			}
			_ = json.NewEncoder(writer).Encode(response)
			return
		}

		// Obtener cadena de bloques de la base de datos
		var blocks []models.Block

		results, err := blockCollection.Find(ctx, bson.M{})

		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			response := responses.BlockResponse{
				Status:  http.StatusInternalServerError,
				Message: "Error al leer la base de datos",
				Data:    err.Error(),
			}
			_ = json.NewEncoder(writer).Encode(response)
			return
		}

		// Lectura de resultados de bloques
		defer func(results *mongo.Cursor, ctx context.Context) {
			_ = results.Close(ctx)
		}(results, ctx)

		for results.Next(ctx) {
			var singleBlock models.Block
			if err = results.Decode(&singleBlock); err != nil {
				writer.WriteHeader(http.StatusInternalServerError)
				response := responses.BlockResponse{
					Status:  http.StatusInternalServerError,
					Message: "Error al leer la base de datos",
					Data:    err.Error(),
				}
				_ = json.NewEncoder(writer).Encode(response)
			}
			blocks = append(blocks, singleBlock)
		}

		// Configuración de zona horaria
		timeZone, err := time.LoadLocation("America/Lima")

		// Se realiza el proof of work (minado) SHA256 (64 bytes)
		previousBlock := blocks[len(blocks)-1]
		proofOfWork := functions.GetProofOfWork(previousBlock.Proof)
		hashPreviousBlock := functions.EncryptSHA256Block(previousBlock)

		// Se crea un nuevo bloque
		newBlock := models.Block{
			Id:           primitive.NewObjectID(),
			Index:        len(blocks) + 1,
			PreviousHash: hashPreviousBlock,
			Proof:        proofOfWork,
			Timestamp:    time.Now().In(timeZone),
			Miner:        functions.EncryptSHA256String(functions.GetMacAddress()),
			Transaction: models.Transaction{
				From:   transaction.From,
				To:     transaction.To,
				Amount: transaction.Amount,
				Fee:    transaction.Fee,
			},
		}

		// Se inserta nuevo bloque en la base de datos
		_, err = blockCollection.InsertOne(ctx, newBlock)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			response := responses.BlockResponse{
				Status:  http.StatusInternalServerError,
				Message: "Error al insertar nuevo bloque en la cadena",
				Data:    err.Error(),
			}
			_ = json.NewEncoder(writer).Encode(response)
			return
		}

		// Devolver response con el nuevo bloque
		writer.WriteHeader(http.StatusOK)
		response := responses.BlockResponse{
			Status:  http.StatusOK,
			Message: "Nuevo bloque insertado con éxito",
			Data:    newBlock,
		}
		_ = json.NewEncoder(writer).Encode(response)
		fmt.Println("Nuevo bloque insertado con éxito")
	}
}

func GetBlockchain() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var blocks []models.Block
		defer cancel()

		results, err := blockCollection.Find(ctx, bson.M{})

		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			response := responses.BlockResponse{
				Status:  http.StatusInternalServerError,
				Message: "Error al leer la base de datos",
				Data:    err.Error(),
			}
			_ = json.NewEncoder(writer).Encode(response)
			return
		}

		// Lectura de resultados de bloques
		defer func(results *mongo.Cursor, ctx context.Context) {
			_ = results.Close(ctx)
		}(results, ctx)

		for results.Next(ctx) {
			var singleBlock models.Block
			if err = results.Decode(&singleBlock); err != nil {
				writer.WriteHeader(http.StatusInternalServerError)
				response := responses.BlockResponse{
					Status:  http.StatusInternalServerError,
					Message: "Error al leer la base de datos",
					Data:    err.Error(),
				}
				_ = json.NewEncoder(writer).Encode(response)
			}
			blocks = append(blocks, singleBlock)
		}

		writer.WriteHeader(http.StatusOK)
		response := responses.BlockResponse{
			Status:  http.StatusOK,
			Message: "Bloques leídos con éxito",
			Data:    blocks,
		}
		_ = json.NewEncoder(writer).Encode(response)
		fmt.Println("Bloques leídos con éxito")
	}
}

func ValidateBlockchain() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var blocks []models.Block
		defer cancel()

		results, err := blockCollection.Find(ctx, bson.M{})

		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			response := responses.BlockResponse{
				Status:  http.StatusInternalServerError,
				Message: "Error al leer la base de datos",
				Data:    err.Error(),
			}
			_ = json.NewEncoder(writer).Encode(response)
			return
		}

		// Lectura de resultados de bloques
		defer func(results *mongo.Cursor, ctx context.Context) {
			_ = results.Close(ctx)
		}(results, ctx)

		for results.Next(ctx) {
			var singleBlock models.Block
			if err = results.Decode(&singleBlock); err != nil {
				writer.WriteHeader(http.StatusInternalServerError)
				response := responses.BlockResponse{
					Status:  http.StatusInternalServerError,
					Message: "Error al leer la base de datos",
					Data:    err.Error(),
				}
				_ = json.NewEncoder(writer).Encode(response)
			}
			blocks = append(blocks, singleBlock)
		}

		previousBlock := blocks[0]
		blockIndex := 1

		for blockIndex < len(blocks) {
			currentBlock := blocks[blockIndex]
			if currentBlock.PreviousHash != functions.EncryptSHA256Block(previousBlock) {
				writer.WriteHeader(http.StatusInternalServerError)
				response := responses.BlockResponse{
					Status:  http.StatusInternalServerError,
					Message: "El blockchain no es válido",
					Data:    nil,
				}
				_ = json.NewEncoder(writer).Encode(response)
				return
			}
			var numberOfZeros = 6
			previousProof := previousBlock.Proof
			currentProof := currentBlock.Proof
			hash := functions.EncryptSHA256Int(currentProof*currentProof - previousProof*previousProof)
			if hash[:numberOfZeros] != strings.Repeat("0", numberOfZeros) {
				writer.WriteHeader(http.StatusInternalServerError)
				response := responses.BlockResponse{
					Status:  http.StatusInternalServerError,
					Message: "El blockchain no es válido",
					Data:    nil,
				}
				_ = json.NewEncoder(writer).Encode(response)
				return
			}
			previousBlock = currentBlock
			blockIndex += 1
		}

		writer.WriteHeader(http.StatusOK)
		response := responses.BlockResponse{
			Status:  http.StatusOK,
			Message: "El blockchain sí es válido",
			Data:    nil,
		}
		_ = json.NewEncoder(writer).Encode(response)
		fmt.Println("Bloques validados con éxito")
	}
}
