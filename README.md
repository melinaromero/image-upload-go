🚀 Image Upload & Optimization API - Go
Este proyecto implementa una API RESTful en Go para la subida y optimización automática de imágenes, almacenándolas eficientemente en Amazon S3. Diseñado para alto rendimiento y escalabilidad en entornos serverless con AWS Lambda.

✨ Características Principales
Recepción de Archivos: Acepta imágenes a través de solicitudes multipart/form-data.

Optimización Inteligente:

Redimensionamiento a 555x555 píxeles: Las imágenes se ajustan a una resolución máxima para optimizar el tamaño.

Calidad JPEG Fija: Utiliza una calidad JPEG fija (70) para un equilibrio óptimo entre tamaño y fidelidad visual, eliminando el bucle adaptativo de compresión para un procesamiento más rápido y predecible.

Soporte de Formatos: Procesa imágenes JPEG y PNG.

Almacenamiento en S3: Las imágenes optimizadas se suben automáticamente a un bucket S3.

Validaciones Robustas: Incluye validaciones de tipo MIME y tamaño de archivo para asegurar la integridad de los datos.

Logs Detallados: Genera logs exhaustivos con tiempos de optimización y tamaño de imagen para facilitar el monitoreo en CloudWatch.

🛠️ Tecnologías Utilizadas
Go: Lenguaje de programación (versión 1.20+).

AWS SDK v2 para Go: Para la interacción fluida con servicios de AWS como S3.

github.com/disintegration/imaging: Librería de procesamiento de imágenes para Go, utilizada para redimensionamiento y codificación.

github.com/aws/aws-lambda-go/lambda: Para la construcción de la función Lambda.

CloudWatch Logs Insights: Herramienta de AWS para el análisis y monitoreo de logs.

📁 Estructura del Proyecto
image-upload-go/
├── go.mod                     # Módulos y dependencias de Go
├── go.sum
├── main.go                    # Punto de entrada de la función Lambda (maneja el `init` y `lambda.Start`)
├── handlers.go                # (Si aplica) Lógica para manejar eventos HTTP/API Gateway
├── s3_service.go              # Funciones para la interacción con AWS S3 (subida, inicialización)
├── image_processor.go         # Contiene la lógica principal de optimización de imágenes
├── constants.go               # Definiciones centralizadas de constantes globales (límites de tamaño, calidad, etc.)
└── utils.go                   # (Opcional) Funciones auxiliares como `generateImageID`

🚀 Cómo Ejecutar y Probar Localmente (Para desarrollo)
Clonar el repositorio:

Bash

git clone https://github.com/melinaromero/image-upload-go.git
cd image-upload-go
Configurar Variables de Entorno (AWS):
Asegúrate de tener tus credenciales de AWS configuradas en tu entorno. Puedes definirlas directamente o usar el archivo ~/.aws/credentials.

Bash

export AWS_REGION=us-east-1
export AWS_ACCESS_KEY_ID=TU_ACCESS_KEY_ID
export AWS_SECRET_ACCESS_KEY=TU_SECRET_ACCESS_KEY
export S3_BUCKET=nombre-de-tu-bucket # El nombre de tu bucket S3 donde se subirán las imágenes

Instalar Dependencias y Ejecutar:

Bash

go mod tidy
go run main.go
La API quedará escuchando por defecto en http://localhost:8080 si tienes un servidor HTTP montado en main.go para pruebas locales.

📦 Despliegue en AWS Lambda
Este proyecto está diseñado para ser desplegado como una función AWS Lambda con un runtime personalizado (Amazon Linux 2023 - provided.al2023).

Limpiar archivos de compilación previos (opcional):
Si ya tienes binarios o ZIPs de compilaciones anteriores, es buena práctica eliminarlos:

PowerShell

Remove-Item -Path .\function.zip -ErrorAction SilentlyContinue
Remove-Item -Path .\bootstrap -ErrorAction SilentlyContinue
go clean -modcache

Compilar el código Go para Lambda:
Compila tu aplicación Go, asegurándote de que el binario resultante se llame bootstrap (este es el nombre esperado por el runtime personalizado de Lambda):

PowerShell

$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o bootstrap
Después de ejecutar este comando, NO DEBES VER NINGÚN ERROR en la salida. Solo se mostrarán las descargas de paquetes si es la primera vez o si se limpió la caché.

Verificar el binario:
Confirma que el archivo bootstrap fue creado en tu directorio raíz:

PowerShell

dir bootstrap
Este comando DEBE mostrarte el archivo bootstrap ahora.

Comprimir el binario:
Crea un archivo function.zip que contenga el binario bootstrap:

PowerShell

Compress-Archive -Path .\bootstrap -DestinationPath .\function.zip

Subir function.zip a tu función Lambda en AWS:

Ve a la Consola de AWS Lambda.

Navega a tu función Lambda (ej. image-upload-go).

Haz clic en la pestaña "Código".

Debajo del editor de código, busca el botón "Cargar desde".

Selecciona la opción "Archivo .zip".

Haz clic en "Cargar" y selecciona el function.zip recién creado.

Haz clic en "Implementar" (o "Deploy") si aparece.

Confirmar la configuración del "Controlador" (Handler):

En la página de tu función Lambda, haz clic en la pestaña "Configuración".

En el menú lateral izquierdo, haz clic en "Configuración general".

Haz clic en "Editar" y asegúrate de que el campo "Controlador" (Handler) esté establecido exactamente a bootstrap.

Asegúrate también de que el Runtime esté configurado como Custom runtime on Amazon Linux 2 o Custom runtime on Amazon Linux 2023.

Ajustar la memoria de la Lambda (Performance Tuning):
Para optimizar el rendimiento y el costo, se recomienda 1024MB de memoria para esta función, ya que Go se beneficia de más CPU y tu Max Memory Used (117MB) está bien dentro de este límite. Si buscas un rendimiento aún mayor, puedes experimentar con 2048MB o 3008MB utilizando herramientas como AWS Lambda Power Tuning.

En la pestaña "Configuración", ve a "General configuration".

Haz clic en "Edit" y ajusta el Memory (MB) a 1024 o el valor deseado.

Probar tu función Lambda:
Invoca tu función Lambda desde tu cliente (ej. Postman) a través del API Gateway que tengas configurado.

📤 Endpoint Principal
POST /upload

Request:

Content-Type: multipart/form-data

Campo: file (el archivo de imagen, e.g., .jpg, .jpeg, .png)

Ejemplo de Request Body (Postman/cURL):
{ "file": <tu archivo de imagen> }

Response (Ejemplo de éxito):

JSON

{
"url": "https://your-bucket-name.s3.amazonaws.com/uploads/unique_image_id.jpg"
}
📊 Logs Útiles para Métricas (CloudWatch Logs Insights)
Puedes usar CloudWatch Logs Insights para monitorear y analizar el rendimiento de tu función. Busca los siguientes patrones en los logs:

Tiempo de inicialización del cliente S3 (parte del cold start):
"✅ Cliente S3 inicializado en"

Tamaño original de la imagen recibida:
"📏 Tamaño original de imagen:"

Tiempo de optimización de la imagen:
"🛠 Optimización completada en"

Tiempo de subida a S3:
"☁️ Subida a S3 completada en"

Tiempo total de procesamiento de tu función (excluyendo cold start):
"✅ Imagen procesada y subida correctamente en"

Duración total de la invocación de Lambda (incluye overhead):
REPORT RequestId: ... Duration:

🧠 Optimización y Consideraciones Avanzadas
Redimensionamiento a 555x555 píxeles: Esta medida es la más efectiva para reducir drásticamente los tiempos de procesamiento y el tamaño de los archivos finales. Si se requieren tiempos de respuesta por debajo de 100ms, considera revisar si la resolución 555x555 puede ser aún menor, o si la subida a S3 puede desacoplarse.

Límites de API Gateway (HTTP API): El límite de payload para HTTP API es de 10MB. Si subes imágenes que, tras la codificación Base64 y el overhead de multipart/form-data, superan este tamaño (como tu ejemplo de 4.8MB), recibirás un error 413 Request Entity Too Large antes de que la solicitud llegue a tu Lambda. Para archivos más grandes, la solución recomendada es implementar S3 Pre-Signed URLs para la subida directa desde el cliente a S3.

Rendimiento Extremo (Menos de 100ms):

Desacoplar subida a S3: Si el objetivo de 100ms es el tiempo de respuesta al cliente, considera que la Lambda devuelva una respuesta rápida y luego use una invocación asíncrona de otra Lambda o un evento S3 para realizar la subida y post-procesamiento.

Reducción de Resolución/Calidad: Experimentar con ResizeLimit (ej. 300x300) y quality JPEG (ej. 60) puede reducir aún más los tiempos de procesamiento.

Librerías de bajo nivel: Para un rendimiento extremo, bibliotecas como go-libvips (un binding a la librería C++ libvips) son más rápidas, pero su implementación en Lambda es considerablemente más compleja debido a las dependencias nativas.

📌 Posibles Mejoras Futuras
Soporte de Formatos Adicionales: Implementar decodificación/codificación para formatos como WebP (si se convierte en un requisito de negocio).

Manejo de Errores: Reintentos automáticos ante fallas transitorias de S3 o dependencias externas.

Pruebas: Implementar pruebas unitarias y de integración para las funciones clave.

Contenedorización: Dockerizar la aplicación para entornos de desarrollo y despliegue más consistentes.

Interfaz Web de Prueba: Desarrollar una pequeña UI para facilitar las pruebas manuales.

Configuración Externa: Mover s3BucketName y otras configuraciones a variables de entorno de Lambda o AWS Secrets Manager para mayor flexibilidad.