package env

import (
	"github.com/joho/godotenv"
	"os"
)

var (
	MINIO_ENDPOINT    string
	MINIO_ACCESS_KEY  string
	MINIO_SECRET_KEY  string
	MINIO_BUCKET_NAME string
)

func Init() {
	godotenv.Load()
	MINIO_ENDPOINT = os.Getenv("MINIO_ENDPOINT")
	MINIO_ACCESS_KEY = os.Getenv("MINIO_ACCESS_KEY")
	MINIO_SECRET_KEY = os.Getenv("MINIO_SECRET_KEY")
	MINIO_BUCKET_NAME = os.Getenv("MINIO_BUCKET_NAME")
}
