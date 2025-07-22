package main
// Constantes globales

import "image"

const (
    MaxFileSize       = 6 * 1024 * 1024 // 5MB (Mantener, ya que es el límite final)
    // MinQuality y InitialQuality ya no son usados si se usa una calidad fija.
    // Si quieres reintroducir la heurística de tamaño de entrada, podrías definirlos aquí
    // o hardcodearlos en la heurística.
    // TargetSizeKB ya no es usado con la eliminación del bucle iterativo.
    ExpirationTimeSec = 60 * 60         // 1 hora
)

// RESIZE_LIMIT actualizado para coincidir con el original de Node.js (555x555)
var ResizeLimit = image.Rect(0, 0, 555, 555) // Ancho y alto máximos del original

var AllowedMimeTypes = []string{
    "image/jpeg",
    "image/png",
    // "image/webp", // Añadir si también procesas webp como se ve en s3_services.go
}

var AllowedExtensions = []string{
    "jpeg",
    "jpg",
    "png",
    // "webp", // Añadir si se soporta webp
}