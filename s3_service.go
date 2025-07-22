package main

import (
    "bytes"
    "context"
    "fmt"
    "log"
    "time"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
    s3Client     *s3.Client
    s3BucketName = "image-upload-go-melina"
    // maxImageSize eliminado, se usa MaxFileSize de constants.go
    // allowedTypes eliminado, se usa AllowedMimeTypes de constants.go
)

func InitS3Client() error {
    start := time.Now()
    cfg, err := config.LoadDefaultConfig(context.Background())
    if err != nil {
       return fmt.Errorf("error al cargar configuración de AWS: %w", err)
    }
    s3Client = s3.NewFromConfig(cfg)
    log.Printf("✅ Cliente S3 inicializado en %v\n", time.Since(start))
    return nil
}

func ProcessAndUploadImage(uploaded *UploadedFile, userName string) (string, error) {
    start := time.Now() // ⏱ Inicio del proceso
    log.Println("🚀 Inicio del procesamiento de imagen")
    if uploaded == nil || len(uploaded.Content) == 0 {
       return "", fmt.Errorf("no se recibió contenido de imagen")
    }

    // Validar tipo MIME usando AllowedMimeTypes de constants.go
    foundMime := false
    for _, mime := range AllowedMimeTypes {
        if uploaded.ContentType == mime {
            foundMime = true
            break
        }
    }
    if !foundMime {
       log.Printf("❌ Tipo de archivo no permitido: %s\n", uploaded.ContentType)
       return "", fmt.Errorf("tipo de archivo no permitido: %s", uploaded.ContentType)
    }

    // Validar tamaño inicial usando MaxFileSize de constants.go
    imageSize := len(uploaded.Content)
    log.Printf("📏 Tamaño original de imagen: %d bytes", imageSize)
    if imageSize > MaxFileSize { // Usar MaxFileSize de constants.go
       log.Printf("❌ Imagen demasiado grande: %d bytes\n", imageSize)
       return "", fmt.Errorf("la imagen supera el tamaño máximo permitido (5 MB)")
    }

    // Optimizar imagen
    optStart := time.Now()
    optimized, contentType, err := OptimizeImage(uploaded.Content, uploaded.ContentType) // OptimizeImage ya hace su propia validación y lógica.
    if err != nil {
       return "", fmt.Errorf("error al optimizar la imagen: %w", err)
    }
    log.Printf("🛠 Optimización completada en %v", time.Since(optStart))

    // Usar versión más pequeña
    // Esta lógica de "usar la original si es más pequeña" es interesante,
    // pero podría ser contraproducente si el objetivo es un TargetSizeKB bajo.
    // Si tu objetivo es MaxFileSize (5MB) y la optimizada es más grande que la original,
    // la original ya habrá pasado la validación de MaxFileSize.
    finalBuffer := optimized
    // if len(uploaded.Content) < len(optimized) { // Descomentar si quieres esta lógica específica
    //    finalBuffer = uploaded.Content
    //    contentType = uploaded.ContentType
    // }

    // Generar clave S3
    imageID := generateImageID(uploaded.Filename)
    s3Key := fmt.Sprintf("uploads/%s", imageID)

    // 🟢 Tiempo de subida a S3
    s3Start := time.Now()
    // Preparar subida
    input := &s3.PutObjectInput{
       Bucket:      aws.String(s3BucketName),
       Key:         aws.String(s3Key),
       Body:        bytes.NewReader(finalBuffer),
       ContentType: aws.String(contentType),
       Metadata: map[string]string{
          "user": userName,
       },
    }

    // Subir a S3
    _, err = s3Client.PutObject(context.Background(), input)
    if err != nil {
       return "", fmt.Errorf("error al subir a S3: %w", err)
    }
    log.Printf("☁️ Subida a S3 completada en %v", time.Since(s3Start))

    // URL pública
    url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s3BucketName, s3Key)
    log.Printf("✅ Imagen procesada y subida correctamente en %v\n", time.Since(start)) // ⏱ Fin del proceso
    return url, nil
}