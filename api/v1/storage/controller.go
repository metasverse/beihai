package storage

import (
	"fmt"
	"github.com/eatmoreapple/regia"
	"github.com/google/uuid"
	"lihood/g"
	"lihood/pkg/storage"
	"net/http"
	"time"
)

func newController() *controller {
	return &controller{}
}

type controller struct{}

func (c controller) productUpload() regia.HandleFunc {
	return g.Wrapper(func(context *regia.Context) error {
		file, _, err := context.Request.FormFile("file")
		if err != nil {
			return g.BadRequest(context, "请上传文件")
		}
		defer file.Close()
		var data = make([]byte, 512)
		if _, err = file.Read(data); err != nil {
			return g.ServerError(context)
		}
		contentType := http.DetectContentType(data)
		if _, err = file.Seek(0, 0); err != nil {
			return g.ServerError(context)
		}
		if contentType != "image/jpeg" && contentType != "image/png" {
			return g.BadRequest(context, "png或jpg格式的图片")
		}
		uw, err := uuid.NewUUID()
		if err != nil {
			return err
		}
		filename := fmt.Sprintf("products/%s/%s.%s", time.Now().Format("20060102"), uw.String(),
			contentType[6:])
		sto := storage.NewOssFileStorage(filename, contentType)
		if path, err := sto.Save(file); err != nil {
			return err
		} else {
			return g.OK(context, path)
		}
	})
}

func (c controller) avatarUpload() regia.HandleFunc {
	return g.Wrapper(func(context *regia.Context) error {
		file, _, err := context.Request.FormFile("file")
		if err != nil {
			return g.BadRequest(context, "请上传文件")
		}
		defer file.Close()
		var data = make([]byte, 512)
		if _, err = file.Read(data); err != nil {
			return g.ServerError(context)
		}
		contentType := http.DetectContentType(data)
		if _, err = file.Seek(0, 0); err != nil {
			return g.ServerError(context)
		}
		if contentType != "image/jpeg" && contentType != "image/png" {
			return g.BadRequest(context, "png或jpg格式的图片")
		}
		uw, err := uuid.NewUUID()
		if err != nil {
			return err
		}
		filename := fmt.Sprintf("avatars/%s/%s.%s", time.Now().Format("20060102"), uw.String(),
			contentType[6:])
		sto := storage.NewOssFileStorage(filename, contentType)
		if path, err := sto.Save(file); err != nil {
			return err
		} else {
			return g.OK(context, path)
		}
	})
}

func (c controller) idcardImageUpload() regia.HandleFunc {
	return g.Wrapper(func(context *regia.Context) error {
		file, _, err := context.Request.FormFile("file")
		if err != nil {
			return g.BadRequest(context, "请上传文件")
		}
		defer file.Close()
		var data = make([]byte, 512)
		if _, err = file.Read(data); err != nil {
			return g.ServerError(context)
		}
		contentType := http.DetectContentType(data)
		if _, err = file.Seek(0, 0); err != nil {
			return g.ServerError(context)
		}
		if contentType != "image/jpeg" && contentType != "image/png" {
			return g.BadRequest(context, "png或jpg格式的图片")
		}
		uw, err := uuid.NewUUID()
		if err != nil {
			return err
		}
		filename := fmt.Sprintf("idcard/%s/%s.%s", time.Now().Format("20060102"), uw.String(),
			contentType[6:])
		sto := storage.NewOssFileStorage(filename, contentType)
		if path, err := sto.Save(file); err != nil {
			return err
		} else {
			return g.OK(context, path)
		}
	})
}
