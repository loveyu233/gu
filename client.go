package gu

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-redis/redis/v8"
	"github.com/minio/minio-go/v7"
	"github.com/olivere/elastic/v7"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

var (
	PgSqlClient func(cfg ...UseConfig) *gorm.DB

	MySqlClient func(cfg ...UseConfig) *gorm.DB

	EsClient *elastic.Client

	RedisClient *redis.Client

	MinioClient *minio.Client

	S3Client *s3.S3

	RabbitClient *amqp.Connection
)
