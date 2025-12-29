package oss

import (
	"ice-chat/config"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func InitMinio(ossConfig config.OssConfig) {
	minioClient, err := minio.New(ossConfig.EndPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(ossConfig.AccessKeyID, ossConfig.SecretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}
}
