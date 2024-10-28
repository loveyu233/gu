package gu

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
	"path/filepath"
)

type S3Config struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Region    string
}

var DefaultS3Config = S3Config{
	Endpoint:  "http://127.0.0.1:9000",
	AccessKey: "minio-minio",
	SecretKey: "minio-minio",
	Region:    "us-east-1",
}

func MustInitS3Client(s3Config ...S3Config) *s3.S3 {
	var config S3Config
	if len(s3Config) >= 1 {
		config = s3Config[0]
	} else {
		config = DefaultS3Config
	}

	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(config.Region),
		Endpoint:         aws.String(config.Endpoint),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      credentials.NewStaticCredentials(config.AccessKey, config.SecretKey, ""),
	})
	if err != nil {
		log.Fatalf("failed to create session: %v\n", err)
	}

	logrus.Info("Successfully created S3 client")
	S3Client = s3.New(sess)

	return S3Client
}

func S3CreateBucket(bucket string, isPublic ...bool) error {
	_, err := S3Client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return err
	}

	logrus.Printf("Bucket %s created successfully.", bucket)

	if len(isPublic) == 0 || !isPublic[0] {
		return nil
	}

	// 设置 Bucket 策略以公开访问
	policy := `{
        "Version": "2012-10-17",
        "Statement": [
            {
                "Effect": "Allow",
                "Principal": "*",
                "Action": "s3:GetObject",
                "Resource": "arn:aws:s3:::` + bucket + `/*"
            }
        ]
    }`

	_, err = S3Client.PutBucketPolicy(&s3.PutBucketPolicyInput{
		Bucket: aws.String(bucket),
		Policy: aws.String(policy),
	})
	if err != nil {
		return err
	}

	logrus.Printf("Bucket %s is now public.", bucket)
	return nil
}

func S3DeleteObj(bucket, key string) error {
	_, err := S3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}
	logrus.Printf("Deleted object %s from bucket %s", key, bucket)
	return nil
}

func S3ClearBucket(bucket string, isDeleteBucket ...bool) error {
	objects, err := S3Client.ListObjects(&s3.ListObjectsInput{Bucket: aws.String(bucket)})
	if err != nil {
		return err
	}

	for _, obj := range objects.Contents {
		err = S3DeleteObj(bucket, *obj.Key)
		if err != nil {
			return err
		}
	}

	if len(isDeleteBucket) == 0 || !isDeleteBucket[0] {
		return nil
	}

	_, err = S3Client.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return err
	}

	logrus.Printf("Deleting bucket %s from s3 successfully.", bucket)

	return nil
}

func S3PutObj(bucket, key, contentType string, size int64, data io.ReadSeeker) error {
	_, err := S3Client.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(key),
		Body:          data,
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(size),
	})

	if err != nil {
		return err
	}

	logrus.Printf("Put object %s to bucket %s", key, bucket)
	return nil
}

func S3PutFile(bucket string, file *os.File, key ...string) error {
	contentType, err := GetFileContentType(file)
	if err != nil {
		return err
	}
	var name string
	if len(key) > 0 {
		name = key[0]
	} else {
		name = filepath.Base(file.Name())
	}
	_, err = S3Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(name),
		Body:        file,
		ContentType: aws.String(contentType),
	})

	if err != nil {
		return err
	}

	logrus.Printf("Put object %s to bucket %s", name, bucket)
	return nil
}

// S3PutDir 例如：./dirPath/ 最后要加上一个 ‘/‘
func S3PutDir(bucket, dirPath string, isNested ...bool) error {
	dir, err := os.ReadDir(dirPath)
	if err != nil {
		panic(err)
	}

	for _, file := range dir {
		if file.IsDir() {
			if len(isNested) == 0 || !isNested[0] {
				continue
			}
			if err = S3PutDir(bucket, dirPath+file.Name()+"/"); err != nil {
				return err
			}
		}

		open, err := os.Open(dirPath + file.Name())
		if err != nil {
			return err
		}

		if err = S3PutFile(bucket, open); err != nil {
			return err
		}
	}

	return nil
}
