package oss

import (
	"ice-chat/config"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinioClient(ossConfig config.OssConfig) *minio.Client {
	minioClient, err := minio.New(ossConfig.EndPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(ossConfig.AccessKeyID, ossConfig.SecretAccessKey, ""),
		Secure: ossConfig.Secure,
	})
	if err != nil {
		log.Fatalln(err)
	}

	return minioClient
}
