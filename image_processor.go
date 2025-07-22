package main

import (
    "bytes"
    "errors"
    "fmt"
    "image"
    "image/jpeg"
    "image/png"

    "github.com/disintegration/imaging"
)

// OptimizeImage optimiza una imagen, redimensiona y ajusta calidad si es necesario.
func OptimizeImage(imageData []byte, contentType string) ([]byte, string, error) {
    var img image.Image
    var format string
    reader := bytes.NewReader(imageData)

    // Validar tipo MIME permitido ANTES de decodificar
    foundMime := false
    for _, mime := range AllowedMimeTypes { // Usar AllowedMimeTypes de constants.go
        if contentType == mime {
            foundMime = true
            break
        }
    }
    if !foundMime {
        return nil, "", errors.New("tipo de imagen no soportado o no permitido")
    }

    switch contentType {
    case "image/jpeg", "image/jpg":
       decoded, err := jpeg.Decode(reader)
       if err != nil {
          return nil, "", fmt.Errorf("no se pudo decodificar la imagen JPEG: %w", err)
       }
       img = decoded
       format = "jpeg"
    case "image/png":
       decoded, err := png.Decode(reader)
       if err != nil {
          return nil, "", fmt.Errorf("no se pudo decodificar la imagen PNG: %w", err)
       }
       img = decoded
       format = "png"

    default:
       // Esto ya fue validado arriba, pero por si acaso.
return nil, "", fmt.Errorf("tipo de imagen no soportado internamente: %s", contentType)}
    // Redimensionar si es necesario
    bounds := img.Bounds()
    // Utiliza ResizeLimit de constants.go
    if bounds.Dx() > ResizeLimit.Dx() || bounds.Dy() > ResizeLimit.Dy() {
       img = imaging.Fit(img, ResizeLimit.Dx(), ResizeLimit.Dy(), imaging.Linear)
    }

    // Buffer de salida
    var optimizedBuffer bytes.Buffer
    var encodeErr error

    if format == "jpeg" {
       // Calidad fija. Se puede usar 70, o implementar la heurística de tamaño de entrada aquí.
       quality := 70 // Esta calidad es la que ha demostrado buen rendimiento en tus pruebas.
       encodeErr = jpeg.Encode(&optimizedBuffer, img, &jpeg.Options{Quality: quality})
       if encodeErr != nil {
          return nil, "", fmt.Errorf("error al codificar JPEG: %w", encodeErr)
       }
    } else if format == "png" {
       encodeErr = png.Encode(&optimizedBuffer, img)
       if encodeErr != nil {
          return nil, "", fmt.Errorf("error al codificar PNG: %w", encodeErr)
       }
    }

    // Validación final contra MaxFileSize de constants.go
    if optimizedBuffer.Len() > MaxFileSize {
       return nil, "", fmt.Errorf("la imagen optimizada (%d bytes) sigue siendo demasiado grande (max: %d bytes)", optimizedBuffer.Len(), MaxFileSize)
    }

    return optimizedBuffer.Bytes(), fmt.Sprintf("image/%s", format), nil
}