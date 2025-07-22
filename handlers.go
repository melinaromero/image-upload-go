package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	//"fmt"
	"io"
	"log"
	"mime/multipart"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

type UploadedFile struct {
	Content     []byte
	Filename    string
	ContentType string
}

func LambdaHandler(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	start := time.Now()
	userName := "Not Yet Identified"

	// Validar Content-Type del request
	contentType := event.Headers["content-type"]
	if contentType == "" {
		contentType = event.Headers["Content-Type"]
	}
	if contentType == "" {
		return response(400, "Falta el Content-Type")
	}

	if !event.IsBase64Encoded {
		return response(400, "El cuerpo debe estar en base64")
	}

	bodyBytes, err := base64.StdEncoding.DecodeString(event.Body)
	if err != nil {
		log.Printf("❌ Error al decodificar base64: %v", err)
		return response(400, "Base64 inválido")
	}

	boundary := extractBoundary(contentType)
	if boundary == "" {
		return response(400, "No se encontró el boundary")
	}

	reader := multipart.NewReader(bytes.NewReader(bodyBytes), boundary)
	form, err := reader.ReadForm(MaxFileSize)
	if err != nil {
		log.Printf("❌ Error al parsear el form: %v", err)
		return response(400, "No se pudo parsear el formulario")
	}

	files := form.File["image"]
	if len(files) == 0 {
		return response(400, "Se requiere el campo 'image'")
	}

	fileHeader := files[0]
	file, err := fileHeader.Open()
	if err != nil {
		log.Printf("❌ No se pudo abrir el archivo: %v", err)
		return response(500, "Error abriendo archivo")
	}
	defer file.Close()

	var buf bytes.Buffer
	// Limita la lectura a 5MB + 1 byte para detectar exceso
	n, err := io.CopyN(&buf, file, MaxFileSize+1)
	if err != nil && err != io.EOF {
		log.Printf("❌ Error al leer archivo: %v", err)
		return response(500, "Error leyendo archivo")
	}
	if n > MaxFileSize {
		return response(400, "El archivo excede el tamaño máximo permitido (5 MB)")
	}

	uploadedFile := &UploadedFile{
		Content:     buf.Bytes(),
		Filename:    fileHeader.Filename,
		ContentType: fileHeader.Header.Get("Content-Type"),
	}

	url, err := ProcessAndUploadImage(uploadedFile, userName)
	if err != nil {
		log.Printf("❌ Error al procesar y subir imagen: %v", err)
		return response(500, "Error al procesar imagen")
	}

	log.Printf("✅ Imagen procesada y subida en %v\n", time.Since(start))

	responseBody, _ := json.Marshal(map[string]string{
		"imageUrl": url,
	})
	return events.APIGatewayProxyResponse{
		StatusCode:      200,
		Headers:         map[string]string{"Content-Type": "application/json"},
		Body:            string(responseBody),
		IsBase64Encoded: false,
	}, nil
}

// Extra helpers

func extractBoundary(contentType string) string {
	const boundaryPrefix = "boundary="
	index := strings.Index(contentType, boundaryPrefix)
	if index == -1 {
		return ""
	}
	return contentType[index+len(boundaryPrefix):]
}

func response(status int, message string) (events.APIGatewayProxyResponse, error) {
	body, _ := json.Marshal(map[string]string{"error": message})
	return events.APIGatewayProxyResponse{
		StatusCode:      status,
		Headers:         map[string]string{"Content-Type": "application/json"},
		Body:            string(body),
		IsBase64Encoded: false,
	}, nil
}
