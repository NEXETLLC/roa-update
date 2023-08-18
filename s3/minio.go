package s3

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"os"
	"time"
	"update-roa/env"
)

var Client *minio.Client

type File struct {
	Name    string
	Options minio.Options
	Path    string
}

//Default set context time out 10s.

var ctx context.Context

func Init() {
	var err error
	Client, err = minio.New(
		env.MINIO_ENDPOINT, &minio.Options{
			Creds:  credentials.NewStaticV4(env.MINIO_ACCESS_KEY, env.MINIO_SECRET_KEY, ""),
			Secure: true,
		})
	if err != nil {
		panic(err)
	}
	//test connection
	ctx, _ = context.WithTimeout(context.Background(), time.Second*5)
	//list bucket content
	obj := Client.ListObjects(ctx, env.MINIO_BUCKET_NAME, minio.ListObjectsOptions{})
	for object := range obj {
		if object.Err != nil {
			panic(object.Err)
		}
	}
}

func (b *File) IfFileExist() bool {
	//set ctx deadline
	ctx, _ = context.WithTimeout(context.Background(), time.Second*5)
	var exist minio.ObjectInfo

	exist, err := Client.StatObject(ctx, env.MINIO_BUCKET_NAME, b.Name, minio.StatObjectOptions{})
	if err != nil {
		return false
	} else if exist.Err != nil {
		return false
	}
	return exist.Size > 0
}

func (f *File) Upload() {
	ctx, _ = context.WithTimeout(context.Background(), time.Second*5)
	//func (c *Client) PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64,
	//	opts PutObjectOptions,
	//) (info UploadInfo, err error)
	//openfile
	file, err := os.Open(f.Path)
	if err != nil {
		fmt.Printf("open file error: %v\n", err)
		return
	}
	defer file.Close()
	//calculate file size
	fileInfo, err := file.Stat()

	if err != nil {
		fmt.Printf("get file info error: %v\n", err)
		return
	}
	info, err := Client.PutObject(ctx, env.MINIO_BUCKET_NAME, f.Name, file, fileInfo.Size(), minio.PutObjectOptions{
		ContentType: "text/plain",
	})
	if err != nil {
		fmt.Printf("upload file error: %v\n", err)
		return
	}
	fmt.Printf("Successfully uploaded %s of size %d\n", f.Name, info.Size)
}
