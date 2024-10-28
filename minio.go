package gu

import (
	"bytes"
	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
)

type MinioConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Secure    bool
}

var DefaultMinioConfig = MinioConfig{
	Endpoint:  "127.0.0.1:9000",
	AccessKey: "minio-minio",
	SecretKey: "minio-minio",
	Secure:    false,
}

func MustInitMinioClient(cfg ...MinioConfig) *minio.Client {
	var dMinioConfig MinioConfig
	if len(cfg) > 0 {
		dMinioConfig = cfg[0]
	} else {
		dMinioConfig = DefaultMinioConfig
	}
	minioClient, err := minio.New(dMinioConfig.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(dMinioConfig.AccessKey, dMinioConfig.SecretKey, ""),
		Secure: dMinioConfig.Secure,
	})
	if err != nil {
		logrus.Panicf("failed to init minio client: %v\n", err)
	}

	logrus.Info("successfully init minio client")
	MinioClient = minioClient

	return minioClient
}

func CreateBucket(bucketName string, location ...string) error {
	var l string
	if len(location) > 0 {
		l = location[0]
	} else {
		l = "us-east-1"
	}
	err := MinioClient.MakeBucket(Timeout(), bucketName, minio.MakeBucketOptions{Region: l})
	if err != nil {
		exists, errBucketExists := MinioClient.BucketExists(Timeout(), bucketName)
		if errBucketExists == nil && exists {
			logrus.Infof("bucket %s 已存在", bucketName)
		} else {
			return errBucketExists
		}
	} else {
		logrus.Infof("bucket %s 创建成功", bucketName)
	}
	return nil
}

func UploadDataByte(data []byte, bucketName, objectName string) error {
	contentType := http.DetectContentType(data)
	object, err := MinioClient.PutObject(Timeout(), bucketName, objectName, bytes.NewBuffer(data), int64(len(data)), minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return err
	}
	logrus.Infof("minio 上传 %s 成功 大小: %+v", object.Key, object.Size)
	return nil
}

func UploadDataReader(data io.Reader, dataSize int64, bucketName, objectName, contentType string) error {
	object, err := MinioClient.PutObject(Timeout(), bucketName, objectName, data, dataSize, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return err
	}
	logrus.Infof("minio 上传 %s 成功 大小: %+v", object.Key, object.Size)
	return nil
}

type Download struct {
	Data        []byte
	Name        string
	ContextType string
	Size        int64
}

func DownloadData(bucketName, objectName string) (*Download, error) {
	object, err := MinioClient.GetObject(Timeout(), bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer object.Close()
	data, err := io.ReadAll(object)
	if err != nil {
		return nil, err
	}
	stat, err := object.Stat()
	if err != nil {
		return nil, err
	}
	logrus.Infof("minio 获取文件 %s 成功 size: %d", stat.Key, stat.Size)
	return &Download{data, stat.Key, stat.ContentType, stat.Size}, nil
}

func StartS3FiberHandler(prefix string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		path := c.Path()
		path = strings.TrimPrefix(path, prefix)
		arr := strings.Split(path, "/")
		if len(arr) < 2 {
			return Resp400(c, nil, "参数错误")
		}
		bucket := arr[0]
		object := arr[1]
		data, err := DownloadData(bucket, object)
		if err != nil {
			return Resp500(c, err, "文件获取失败")
		}
		c.Response().Header.Set("content-type", data.ContextType)
		return c.Send(data.Data)
	}
}
