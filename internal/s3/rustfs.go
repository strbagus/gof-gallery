package s3

import (
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var S3Client *minio.Client
var BucketName string

// InitS3 initializes the S3-compatible client for RustFS using environment variables
func InitS3() {
	s3Domain := os.Getenv("S3_DOMAIN")
	if s3Domain == "" {
		log.Fatal("S3_DOMAIN tidak diatur di file .env")
	}

	accessKey := os.Getenv("S3_ACCESS_KEY")
	secretKey := os.Getenv("S3_SECRET_KEY")
	BucketName = os.Getenv("S3_BUCKET")

	var endpoint string
	var useSSL bool

	// Parse URL schema to support domains with http/https prefix
	if strings.Contains(s3Domain, "://") {
		u, err := url.Parse(s3Domain)
		if err != nil {
			log.Fatalf("Gagal memproses S3_DOMAIN: %v", err)
		}
		endpoint = u.Host
		useSSL = u.Scheme == "https"
	} else {
		endpoint = s3Domain
		useSSL = false
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalf("Gagal menginisialisasi client RustFS/S3: %v", err)
	}

	S3Client = client
	log.Printf("Client RustFS (S3) berhasil diinisialisasi ke %s (Secure: %t)!", endpoint, useSSL)
}
