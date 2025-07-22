package main

import (
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

func init() {
	log.Println("⏳ Inicializando cliente S3...")
	err := InitS3Client()
	if err != nil {
		log.Fatalf("❌ No se pudo inicializar el cliente S3: %v", err)
	}
	log.Println("✅ Cliente S3 inicializado correctamente")
}

func main() {
	lambda.Start(LambdaHandler)
}
