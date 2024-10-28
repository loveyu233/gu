package client

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-redis/redis/v8"
	"github.com/loveyu233/gu/public"
	"github.com/minio/minio-go/v7"
	"github.com/olivere/elastic/v7"
	"gorm.io/gorm"
)

var (
	PgSqlClient func(cfg ...public.UseConfig) *gorm.DB

	MySqlClient func(cfg ...public.UseConfig) *gorm.DB

	EsClient *elastic.Client

	RedisClient *redis.Client

	MinioClient *minio.Client

	S3Client *s3.S3
)
