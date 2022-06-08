package controllers

import (
	"Toy_Cryptocurrency/configs"
	"Toy_Cryptocurrency/functions"
	"Toy_Cryptocurrency/models"
	"Toy_Cryptocurrency/responses"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/gomail.v2"
)

var userCollection = configs.GetCollection(configs.DB, "Users")
var validateUser = validator.New()

func SendSecurityCodeLogin() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var user models.User
		defer cancel()

		// Validar que el body está en formato JSON
		if err := json.NewDecoder(request.Body).Decode(&user); err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{
				Status:  http.StatusBadRequest,
				Message: "Formato del contenido de la solicitud no válido",
				Data:    err.Error(),
			}
			_ = json.NewEncoder(writer).Encode(response)
			return
		}

		// Se usa la librería para validar los campos del body
		if validationErr := validateUser.Struct(&user); validationErr != nil {
			writer.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{
				Status:  http.StatusBadRequest,
				Message: "Campos del contenido de la solicitud no válidos",
				Data:    validationErr.Error(),
			}
			_ = json.NewEncoder(writer).Encode(response)
			return
		}

		// Verificar que el usuario ya esté registrado
		var tempUser models.User
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&tempUser)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "Usuario no registrado",
				Data:    "Usuario no registrado",
			}
			_ = json.NewEncoder(writer).Encode(response)
			return
		}

		// Verificar que la contraseña sea correcta
		if tempUser.Password != user.Password {
			writer.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "Contraseña incorrecta",
				Data:    "Contraseña incorrecta",
			}
			_ = json.NewEncoder(writer).Encode(response)
			return
		}

		// Crear imagen y txt con el código de seguridad
		errorCode, securityCodeUserImageRoute := functions.CreateSecurityCode(user)
		if errorCode != 0 {
			writer.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "Error al generar el código de seguridad",
				Data:    fmt.Sprintf("Error al generar el código de seguridad %d", errorCode),
			}
			_ = json.NewEncoder(writer).Encode(response)
			return
		}

		// Enviar correo con el código de verificación
		message := gomail.NewMessage()
		message.SetHeader("From", message.FormatAddress("toy.cryptocurrency@gmail.com", "Toy Cryptocurrency"))
		message.SetHeader("To", user.Email)
		message.SetHeader("Subject", "Código de autorización Toy Cryptocurrency")
		emailContent := "<p>Hola <b>" + tempUser.Email + "</b>,</p><p>Tu código de verificación se encuentra adjunto como imagen en este correo.</p><p>No compartas este código con nadie.</p><p><b>Equipo de Toy Cryptocurrency</b></p>"
		message.SetBody("text/html", emailContent)
		message.Attach(securityCodeUserImageRoute)
		dialer := gomail.NewDialer("smtp.gmail.com", 587, "toy.cryptocurrency@gmail.com", "aihulcizgjoxwqwr")
		if err := dialer.DialAndSend(message); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "Error al generar el código de seguridad 4",
				Data:    "Error al generar el código de seguridad",
			}
			_ = json.NewEncoder(writer).Encode(response)
			return
		}

		writer.WriteHeader(http.StatusOK)
		response := responses.UserResponse{
			Status:  http.StatusOK,
			Message: "Código de seguridad enviado al correo electrónico",
			Data:    "Código de seguridad enviado al correo electrónico",
		}
		_ = json.NewEncoder(writer).Encode(response)
		fmt.Printf("Código de verificación del usuario %s ha sido enviado\n", tempUser.Email)
	}
}

func VerifySecurityCodeLogin() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		params := mux.Vars(request)
		securityCode := params["securityCode"]
		var user models.User
		defer cancel()

		// Validar que el body está en formato JSON
		if err := json.NewDecoder(request.Body).Decode(&user); err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{
				Status:  http.StatusBadRequest,
				Message: "Formato del contenido de la solicitud no válido",
				Data:    err.Error(),
			}
			_ = json.NewEncoder(writer).Encode(response)
			return
		}

		// Leer código de archivo txt
		securityCodeTxtRoute := ""
		if runtime.GOOS == "windows" {
			securityCodeTxtRoute = "security_codes/users_codes_texts/" + user.Email + ".txt"
		} else {
			securityCodeTxtRoute = "/home/piero/Encrypted_Instant_Messaging_Backend/security_codes/users_codes_texts/" + user.Email + ".txt"
		}
		securityCodeFile, err := ioutil.ReadFile(securityCodeTxtRoute)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "No se ha solicitado código de seguridad",
				Data:    "No se ha solicitado código de seguridad",
			}
			_ = json.NewEncoder(writer).Encode(response)
			return
		}

		// Verificar código de seguridad
		if securityCode != string(securityCodeFile) {
			writer.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "Código de seguridad incorrecto",
				Data:    "Código de seguridad incorrecto",
			}
			_ = json.NewEncoder(writer).Encode(response)
			return
		}

		// Eliminar archivos con códigos de seguridad
		securityCodeImageRoute := ""
		if runtime.GOOS == "windows" {
			securityCodeImageRoute = "security_codes/users_codes_images/" + user.Email + ".png"
		} else {
			securityCodeImageRoute = "/home/piero/Encrypted_Instant_Messaging_Backend/security_codes/users_codes_images/" + user.Email + ".png"
		}
		err = os.Remove(securityCodeTxtRoute)
		err = os.Remove(securityCodeImageRoute)

		// Si no están los archivos, significa que no se ha solicitado un código de seguridad
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "No se ha solicitado código de seguridad",
				Data:    "No se ha solicitado código de seguridad",
			}
			_ = json.NewEncoder(writer).Encode(response)
			return
		}

		// Extraer datos del usuario de la base de datos
		var dbUser models.User
		_ = userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&dbUser)

		// Se retorna los campos del usuario autenticado
		writer.WriteHeader(http.StatusOK)
		response := responses.UserResponse{
			Status:  http.StatusOK,
			Message: "Usuario logeado con éxito",
			Data: models.User{
				Id:         dbUser.Id,
				FirstName:  dbUser.FirstName,
				LastName:   dbUser.LastName,
				Country:    dbUser.Country,
				Email:      dbUser.Email,
				Password:   "",
				PublicKey:  dbUser.PublicKey,
				PrivateKey: dbUser.PrivateKey,
			},
		}
		_ = json.NewEncoder(writer).Encode(response)
		fmt.Printf("Código verificado, usuario %s logeado con éxito\n", user.Email)
	}
}

func SendSecurityCodeRegister() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var user models.User
		defer cancel()

		// Validar que el body está en formato JSON
		if err := json.NewDecoder(request.Body).Decode(&user); err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{
				Status:  http.StatusBadRequest,
				Message: "Formato del contenido de la solicitud no válido",
				Data:    err.Error(),
			}
			_ = json.NewEncoder(writer).Encode(response)
			return
		}

		// Se usa librería para validar los campos del body
		if validationErr := validateUser.Struct(&user); validationErr != nil {
			writer.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{
				Status:  http.StatusBadRequest,
				Message: "Campos del contenido de la solicitud no válidos",
				Data:    validationErr.Error(),
			}
			_ = json.NewEncoder(writer).Encode(response)
			return
		}

		// Verificar que el usuario no esté registrado en la base de datos
		var tempUser models.User
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&tempUser)
		if err == nil {
			writer.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "Usuario ya existente",
				Data:    "Usuario ya existente",
			}
			_ = json.NewEncoder(writer).Encode(response)
			return
		}

		// Crear imagen y txt con el código de seguridad
		errorCode, securityCodeUserImageRoute := functions.CreateSecurityCode(user)
		if errorCode != 0 {
			writer.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "Error al generar el código de seguridad",
				Data:    fmt.Sprintf("Error al generar el código de seguridad %d", errorCode),
			}
			_ = json.NewEncoder(writer).Encode(response)
			return
		}

		// Enviar correo con el código de verificación
		message := gomail.NewMessage()
		message.SetHeader("From", message.FormatAddress("toy.cryptocurrency@gmail.com", "Toy Cryptocurrency"))
		message.SetHeader("To", user.Email)
		message.SetHeader("Subject", "Código de autorización Toy Cryptocurrency")
		emailContent := "<p>Hola <b>" + user.Email + "</b>,</p><p>Tu código de verificación se encuentra adjunto como imagen en este correo.</p><p>No compartas este código con nadie.</p><p><b>Equipo de Toy Cryptocurrency</b></p>"
		message.SetBody("text/html", emailContent)
		message.Attach(securityCodeUserImageRoute)
		dialer := gomail.NewDialer("smtp.gmail.com", 587, "toy.cryptocurrency@gmail.com", "aihulcizgjoxwqwr")
		if err := dialer.DialAndSend(message); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "Error al generar el código de seguridad 4",
				Data:    "Error al generar el código de seguridad",
			}
			_ = json.NewEncoder(writer).Encode(response)
			return
		}

		writer.WriteHeader(http.StatusOK)
		response := responses.UserResponse{
			Status:  http.StatusOK,
			Message: "Código de seguridad enviado al correo electrónico",
			Data:    "Código de seguridad enviado al correo electrónico",
		}
		_ = json.NewEncoder(writer).Encode(response)
		fmt.Printf("Código de verificación del usuario %s ha sido enviado\n", user.Email)
	}
}

func VerifySecurityCodeRegister() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		params := mux.Vars(request)
		securityCode := params["securityCode"]
		var user models.User
		defer cancel()

		// Validar que el body está en formato JSON
		if err := json.NewDecoder(request.Body).Decode(&user); err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{
				Status:  http.StatusBadRequest,
				Message: "Formato del contenido de la solicitud no válido",
				Data:    err.Error(),
			}
			_ = json.NewEncoder(writer).Encode(response)
			return
		}

		// Leer código de archivo txt
		securityCodeTxtRoute := ""
		if runtime.GOOS == "windows" {
			securityCodeTxtRoute = "security_codes/users_codes_texts/" + user.Email + ".txt"
		} else {
			securityCodeTxtRoute = "/home/piero/Encrypted_Instant_Messaging_Backend/security_codes/users_codes_texts/" + user.Email + ".txt"
		}
		securityCodeFile, err := ioutil.ReadFile(securityCodeTxtRoute)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "No se ha solicitado código de seguridad",
				Data:    "No se ha solicitado código de seguridad",
			}
			_ = json.NewEncoder(writer).Encode(response)
			return
		}

		// Verificar código de seguridad
		if securityCode != string(securityCodeFile) {
			writer.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "Código de seguridad incorrecto",
				Data:    "Código de seguridad incorrecto",
			}
			_ = json.NewEncoder(writer).Encode(response)
			return
		}

		// Eliminar archivos con códigos de seguridad
		securityCodeImageRoute := ""
		if runtime.GOOS == "windows" {
			securityCodeImageRoute = "security_codes/users_codes_images/" + user.Email + ".png"
		} else {
			securityCodeImageRoute = "/home/piero/Encrypted_Instant_Messaging_Backend/security_codes/users_codes_images/" + user.Email + ".png"
		}
		err = os.Remove(securityCodeTxtRoute)
		err = os.Remove(securityCodeImageRoute)

		// Si no están los archivos, significa que no se ha solicitado un código de seguridad
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "No se ha solicitado código de seguridad",
				Data:    "No se ha solicitado código de seguridad",
			}
			_ = json.NewEncoder(writer).Encode(response)
			return
		}

		// Crear llaves pública y privada
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "Error al generar la llave privada",
				Data:    err.Error(),
			}
			_ = json.NewEncoder(writer).Encode(response)
			return
		}
		publicKey := &privateKey.PublicKey

		// Convertir llaves a string
		privateKeyString := base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PrivateKey(privateKey))
		publicKeyString := base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PublicKey(publicKey))

		// Crear modelo de usuario con sus campos completos
		newUser := models.User{
			Id:         primitive.NewObjectID(),
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			Country:    user.Country,
			Email:      user.Email,
			Password:   user.Password,
			PublicKey:  publicKeyString,
			PrivateKey: privateKeyString,
		}
		_, err = userCollection.InsertOne(ctx, newUser)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "Error al registrar el nuevo usuario",
				Data:    err.Error(),
			}
			_ = json.NewEncoder(writer).Encode(response)
			return
		}

		var dbUser models.User
		_ = userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&dbUser)

		writer.WriteHeader(http.StatusCreated)
		response := responses.UserResponse{
			Status:  http.StatusCreated,
			Message: "Usuario registrado con éxito",
			Data:    dbUser,
		}
		_ = json.NewEncoder(writer).Encode(response)
		fmt.Printf("Código verificado, usuario %s registrado con éxito\n", user.Email)
	}
}
