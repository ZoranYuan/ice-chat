package oss

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"ice-chat/config"
	"log"
)

func NewMinioClient(ossConfig config.OssConfig) *minio.Client {

	minioClient, err := minio.New(ossConfig.EndPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(ossConfig.AccessKey, ossConfig.SecretKey, ""),
		Secure: ossConfig.Secure,
	})

	if err != nil {
		log.Fatalln(err)
	}

	return minioClient
}
