# back-proyecto
🚀 Instalación de Go
✅ Requisitos mínimos

    Sistema operativo: Linux, Windows o macOS

    Espacio en disco: al menos 200 MB

    Acceso a terminal o consola

🐧 En Linux (Ubuntu/Debian)

# Descargar el último archivo tar.gz de Go
wget https://go.dev/dl/go1.22.2.linux-amd64.tar.gz

# Eliminar cualquier instalación anterior
sudo rm -rf /usr/local/go

# Extraer Go en /usr/local
sudo tar -C /usr/local -xzf go1.22.2.linux-amd64.tar.gz

# Configurar variables de entorno (agrega esto en ~/.bashrc o ~/.zshrc)
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Verificar instalación
go version

🪟 En Windows

    Ve al sitio oficial: https://go.dev/dl/

    Descarga el instalador de Windows (.msi) — por ejemplo: go1.22.2.windows-amd64.msi

    Ejecuta el instalador y sigue los pasos

    Asegúrate de que Go se agregue al PATH del sistema (el instalador lo hace por defecto)

    Verifica en terminal:

go version

📦 Dependencias del Proyecto

Después de instalar Go, instala las dependencias del proyecto:

cd go-xlsx-client
go mod tidy

📂 Estructura del Proyecto (Resumen)

go-xlsx-client/
├── cmd/                  # Punto de entrada
├── internal/             # Lógica de negocio, puertos, adaptadores
├── go.mod / go.sum       # Dependencias
└── clientes.xlsx         # Archivo de prueba

🚀 Ejecutar la aplicación

go run cmd/server/main.go

📦 back-proyecto

Este proyecto expone un API REST con arquitectura hexagonal (ports & adapters) desarrollada en Go. Permite:

    Obtener los nombres de las hojas (clientes) en un archivo .xlsx.

    Consultar los datos de un cliente específico desde su hoja correspondiente.

⚙️ Requisitos

    Go 1.20 o superior (recomendado 1.22)

    Git

    Un archivo .xlsx con los datos de los clientes (cada hoja representa un cliente)

🖥️ Instalación
En Linux

# Instalar Go (Ubuntu/Debian)
sudo apt update
sudo apt install golang -y

# Verificar instalación
go version

En Windows

    Descarga el instalador desde: https://golang.org/dl/

    Ejecuta el instalador y sigue los pasos.

    Verifica en terminal:

go version

📁 Clonar el repositorio

git clone https://github.com/richardsuan/back-proyecto.git
cd back-proyecto

🛠️ Configuración inicial

Inicializa el módulo Go (solo la primera vez):

go mod tidy

Esto descargará automáticamente las dependencias como:

    Gin: framework web

    tealeg/xlsx: lectura de archivos Excel

▶️ Ejecución

go run main.go

El servidor quedará corriendo en:

http://localhost:8080

📌 Endpoints disponibles
🔹 Obtener nombres de clientes (nombres de hojas del Excel)

GET /clientes

Respuesta:

{
  "clientes": ["Cliente1", "Cliente2", "Cliente3"]
}

🔹 Obtener datos de un cliente específico

GET /clientes/{nombre}

Ejemplo:

GET /clientes/Cliente1

Respuesta:

{
  "cliente": "Cliente1",
  "datos": [
    ["Fecha", "Presion", "Volumen"],
    ["2023-01-01", 15.5, 200],
    ...
  ]
}

🧪 Probar con curl

curl http://localhost:8080/clientes
curl http://localhost:8080/clientes/Cliente1

📂 Estructura del proyecto (Hexagonal)

back-proyecto/
│
├── cmd/                    # Punto de entrada principal
│   └── main.go
│
├── internal/
│   ├── domain/             # Entidades del dominio
│   ├── ports/              # Interfaces
│   ├── services/           # Lógica de negocio
│   └── adapters/           # Adaptadores (Excel, HTTP)
│
├── go.mod
└── go.sum

📄 Licencia

MIT © Richard Suan

¿Quieres que también incluya una plantilla inicial del código base del proyecto?














## librerias necesarias
go get github.com/tealeg/xlsx
