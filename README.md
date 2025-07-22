# Image Upload & Optimization API - Go

Este proyecto es una API REST escrita en Go que permite subir imágenes, optimizarlas automáticamente y almacenarlas en un bucket de Amazon S3. Fue desarrollada con foco en eficiencia y escalabilidad.

## ✨ Características

- ✅ Recepción de archivos mediante `multipart/form-data`.
- 📦 Optimización automática de imágenes JPEG hasta un tamaño objetivo en KB.
- 🚀 Subida a S3 una vez optimizadas.
- 🛡️ Control de calidad mínima para no perder fidelidad visual.
- 📊 Logs detallados con tiempos de optimización y peso de las imágenes.
- 🔁 Bucle adaptativo de compresión hasta alcanzar tamaño deseado.

## 🛠️ Tecnologías

- **Go** 1.20+
- **AWS SDK v2 para S3**
- **JPEG image encoding**
- **CloudWatch Logs Insights (para métricas y seguimiento)**

## 📁 Estructura del proyecto

image-upload-go/
├── go.mod # Dependencias del proyecto
├── go.sum
├── main.go # Punto de entrada y servidor HTTP
├── handlers.go # Endpoints y lógica HTTP
├── s3_service.go # Interacción con AWS S3
├── image_processor.go # Lógica de optimización de imágenes
├── constants.go # Constantes del sistema (quality, size, etc.)


## 🚀 Cómo ejecutar el proyecto

### 1. Clonar el repositorio

```bash
git clone https://github.com/melinaromero/image-upload-go.git
cd image-upload-go

2. Configurar variables de entorno (AWS)
export AWS_REGION=us-east-1
export AWS_ACCESS_KEY_ID=TU_ACCESS_KEY
export AWS_SECRET_ACCESS_KEY=TU_SECRET_KEY
export S3_BUCKET=nombre-de-tu-bucket


3. Instalar dependencias y ejecutar
```bash
go mod tidy
go run main.go
La API quedará escuchando por defecto en http://localhost:8080.

📤 Endpoint principal
POST /upload
Request:

Content-Type: multipart/form-data

Campo: file (imagen .jpg)

Response:
{
  "url": "https://bucket.s3.amazonaws.com/imagen_optimizada.jpg"
}

📊 Logs útiles para métricas
Tamaño original de imagen:
Tamaño original de imagen: 1532 KB

Tiempo de optimización:
Optimización completada en 57ms

Tiempo total de procesamiento:
Imagen procesada y subida correctamente en 115ms

Se pueden visualizar con AWS CloudWatch Logs Insights usando patrones de búsqueda.

🧠 Optimización y mejoras
Actualmente, se usa un bucle que reduce la calidad JPEG de 90 a 40 con decremento de 5 hasta alcanzar el tamaño objetivo.

Si el rendimiento no es suficiente, podés usar una librería externa más optimizada como libvips vía bindings o servicio externo.

📌 Posibles mejoras futuras
Agregar soporte para PNG y WebP.

Reintentos automáticos ante fallas de S3.

Test unitarios para funciones clave.

Dockerización.

Interfaz web de prueba.

Optimización de codigo