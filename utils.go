package main

import (
	"path/filepath"         // Mejor que strings.Split para extraer extensiones
	"strings"
	"github.com/google/uuid"
)

// generateImageID genera un nombre único para la imagen, conservando su extensión (si es válida).
func generateImageID(filename string) string {
	id := uuid.New().String()

	// Extraemos la extensión con seguridad, y la forzamos a minúscula
	ext := strings.ToLower(filepath.Ext(filename))

	// Si no hay extensión o es muy rara, usar .jpg como fallback
	if ext == "" || len(ext) > 10 {
		ext = ".jpg"
	}

	return id + ext
}
