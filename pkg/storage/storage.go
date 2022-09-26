package storage

import (
	"context"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io"
	"net/http"
	"net/url"
)

type Storage interface {
	Save(reader io.Reader) (string, error)
}

// NewOssFileStorage 创建一个新的OssFileStorage
func NewOssFileStorage(filename string, contentType string) Storage {
	return OssFileStorage{
		bucketName:  "lihood-1306623008",
		region:      "https://lihood-1306623008.cos.ap-nanjing.myqcloud.com",
		accessKey:   "AKIDOcZFgqiJfkY5MM64mDNGsEJ6cstd6IDr",
		secretKey:   "OWP1LOBhwtZy6xXrviEUTBbkdFbfSyDB",
		contentType: contentType,
		filename:    filename,
	}
}

// OssFileStorage implements Storage interface
// It uses Tencent Cloud OSS as the storage backend
// 将文件存储到腾讯云对象存储服务
type OssFileStorage struct {
	bucketName string
	region     string
	accessKey  string
	secretKey  string
	//functionName string
	//zipRegion    string
	//endpoint     string
	contentType string
	filename    string
}

func (o OssFileStorage) Save(reader io.Reader) (string, error) {
	// 获取随机文件名
	return o.SaveFile(reader, o.filename)
}

func (o OssFileStorage) NewClient() (*cos.Client, error) {
	path, err := url.Parse(o.region)
	if err != nil {
		return nil, err
	}
	baseURL := &cos.BaseURL{
		BucketURL: path,
	}
	return cos.NewClient(baseURL, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  o.accessKey,
			SecretKey: o.secretKey,
		},
	}), nil
}

func (o OssFileStorage) SaveFile(file io.Reader, filename string) (string, error) {
	client, err := o.NewClient()
	if err != nil {
		return "", err
	}
	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType: o.contentType,
		},
	}
	_, err = client.Object.Put(context.Background(), filename, file, opt)
	if err != nil {
		return "", err
	}
	return o.region + "/" + filename, nil
}
