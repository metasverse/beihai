package account

import (
	"bytes"
	"encoding/base64"
	"github.com/eatmoreapple/regia"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	"lihood/g"
	"lihood/internal/enum"
	"lihood/internal/repository"
	"lihood/internal/requests"
	"lihood/internal/services"
	"lihood/pkg/storage"
)

func newSmsLoginController() *smsLoginController {
	return &smsLoginController{}
}

type smsLoginController struct{}

// 发送登录短信验证码
func (s smsLoginController) sendLoginCode() regia.HandleFunc {
	return g.Wrapper(func(context *regia.Context) error {
		var req requests.SMSSendRequest
		if err := context.Data(&req); err != nil {
			return g.BadRequest(context, err.Error())
		}
		if err := services.NewPhoneLoginService().Send(req.Phone); err != nil {
			return err
		}
		return g.OK(context, nil)
	})
}

// 短信验证码登录
func (s smsLoginController) login() regia.HandleFunc {
	return g.Wrapper(func(context *regia.Context) error {
		var req requests.SMSLoginRequest
		if err := context.Data(&req); err != nil {
			return g.BadRequest(context, err.Error())
		}
		token, err := services.NewPhoneLoginService().LoginWithPhone(req.Phone, req.Code, req.Invitation)
		if err != nil {
			return err
		}
		return g.OK(context, token)
	})
}

func newAccountController() *accountController {
	return &accountController{}
}

type accountController struct{}

// 更新手机号
func (a accountController) updatePhone() regia.HandleFunc {
	return g.Wrapper(func(context *regia.Context) error {
		var req requests.UpdatePhoneRequest
		if err := context.Data(&req); err != nil {
			return g.BadRequest(context, err.Error())
		}
		uid := g.CurrentUserID(context)
		if err := services.NewAccountService(g.DB).UpdatePhone(uid, req.Phone, req.Code); err != nil {
			return err
		}
		return g.OK(context, nil)
	})
}

// 获取更新手机号码验证码
func (a accountController) updatePhoneCode() regia.HandleFunc {
	return g.Wrapper(func(context *regia.Context) error {
		var req requests.UpdatePhoneCodeRequest
		if err := context.Data(&req); err != nil {
			return g.BadRequest(context, err.Error())
		}
		uid := g.CurrentUserID(context)
		if err := services.NewAccountService(g.DB).SendUpdatePhoneCode(uid, req.Phone); err != nil {
			return err
		}
		return g.OK(context, nil)
	})
}

// 获取当前用户详情
func (a accountController) info() regia.HandleFunc {
	return g.Wrapper(func(context *regia.Context) error {
		uid := g.CurrentUserID(context)
		service := services.NewAccountService(g.DB)
		user, err := service.GetByID(uid)
		if err != nil {
			return err
		}
		return g.OK(context, user)
	})
}

// 获取用户详情
func (a accountController) userInfo() regia.HandleFunc {
	return g.Wrapper(func(context *regia.Context) error {
		pk, err := context.Params.Get("id").Int64()
		if err != nil {
			context.Engine.NotFoundHandle(context)
			return nil
		}
		service := services.NewAccountService(g.DB)
		user, err := service.GetByID(pk)
		if err != nil {
			return err
		}
		user.Amount = 0
		user.IdCardPositiveImageUrl = ""
		user.IdCardNegativeImageUrl = ""
		user.IDCardNum = ""
		user.Name = ""
		user.Phone = ""
		return g.OK(context, user)
	})
}

// 用户认证
func (a accountController) authentication() regia.HandleFunc {
	return g.Wrapper(func(context *regia.Context) error {
		var req requests.AuthenticationRequest
		if err := context.Data(&req); err != nil {
			return g.BadRequest(context, err.Error())
		}
		uid := g.CurrentUserID(context)
		if err := services.NewAccountService(g.DB).Authentication(
			uid, req.Name, req.IDCard, req.PositiveImage, req.NegativeImage); err != nil {
			return err
		}
		return g.OK(context, nil)
	})
}

func (a accountController) updateAccountInfo() regia.HandleFunc {
	return g.Wrapper(func(context *regia.Context) error {
		var req requests.UpdateAccountInfoRequest
		if err := context.Data(&req); err != nil {
			return g.BadRequest(context, err.Error())
		}
		uid := g.CurrentUserID(context)
		if err := services.NewAccountService(g.DB).UpdateAccountInfo(uid, req.Nickname, req.Avatar, req.Desc); err != nil {
			return err
		}
		return g.OK(context, nil)
	})
}

// 我的交易历史记录
func (a accountController) tradeHistory() regia.HandleFunc {
	return g.Wrapper(func(context *regia.Context) error {
		t, _ := context.QueryValue("type").Int()
		uid := g.CurrentUserID(context)
		service := services.NewIncomeService(g.DB)
		paging := g.NewQueryPagination(context)
		result, err := service.QueryByType(uid, enum.IncomeType(t), paging.Page(), paging.PageSize())
		if err != nil {
			return err
		}
		count, err := service.CountByType(uid, enum.IncomeType(t))
		if err != nil {
			return err
		}
		return g.Many(context, result, count)
	})
}

func (a accountController) qrcode() regia.HandleFunc {
	return g.Wrapper(func(context *regia.Context) error {
		uid := g.CurrentUserID(context)
		service := services.NewAccountService(g.DB)
		user, err := service.GetByID(uid)
		if err != nil {
			return err
		}
		qr, err := qrcode.New(user.BsnAddress, qrcode.Medium)
		if err != nil {
			return err
		}
		context.ResponseWriter.Header().Set("Content-Type", "image/png")
		context.ResponseWriter.Header().Set("Expires", "86400")
		context.ResponseWriter.Header().Set("Cache-Control", "public, max-age=86400")
		var buffer = new(bytes.Buffer)
		if err = qr.Write(400, buffer); err != nil {
			return err
		}
		img := base64.StdEncoding.EncodeToString(buffer.Bytes())
		return g.OK(context, img)
	})
}

type public struct{}

func (a public) createUser() regia.HandleFunc {
	type request struct {
		Phone string `json:"phone" validate:"phone(m=请输入正确的手机号)"`
		Code  string `json:"code" validate:"required(m=请输入code)"`
	}
	return g.Wrapper(func(context *regia.Context) error {
		var req request
		if err := context.Data(&req); err != nil {
			return g.BadRequest(context, err.Error())
		}
		if req.Code != "dJY1a4bXm+!C[zqrKvwE" {
			return g.BadRequest(context, "code错误")
		}
		// 先去account账户中查询是否存在该用户
		service := services.NewPhoneLoginService()
		_, err := service.PhoneLogin(req.Phone, "")
		if err != nil {
			return err
		}
		return g.OK(context, nil)
	})
}

func (a public) authentication() regia.HandleFunc {
	type request struct {
		Phone         string `json:"phone" validate:"phone(m=请输入正确的手机号)"`
		Code          string `json:"code" validate:"required(m=请输入code)"`
		Name          string `json:"name" validate:"required(m=用户名不能为空)"`
		IDCard        string `json:"id_card" validate:"required(m=身份证号不能为空)"`
		PositiveImage string `json:"positive_image" validate:"required(m=身份证正面照不能为空)"`
		NegativeImage string `json:"negative_image" validate:"required(m=身份证反面照不能为空)"`
	}
	return g.Wrapper(func(context *regia.Context) error {
		var req request
		if err := context.Data(&req); err != nil {
			return g.BadRequest(context, err.Error())
		}
		if req.Code != "dJY1a4bXm+!C[zqrKvwE" {
			return g.BadRequest(context, "code错误")
		}
		user, err := services.NewAccountService(g.DB).GetByPhone(req.Phone)
		if err != nil {
			return err
		}
		if user == nil {
			return g.Fail(context, "该用户不存在")
		}
		// 上传base64图片
		positiveImageData, err := base64.StdEncoding.DecodeString(req.PositiveImage)
		if err != nil {
			return g.Fail(context, "身份证正面照base64解码失败")
		}
		negativeImageData, err := base64.StdEncoding.DecodeString(req.NegativeImage)
		if err != nil {
			return g.Fail(context, "身份证反面照base64解码失败")
		}
		// 上传身份证正面照
		positiveImageURL, err := storage.NewOssFileStorage(uuid.New().String()+".png", "image/png").Save(bytes.NewReader(positiveImageData))
		if err != nil {
			return err
		}
		// 上传身份证反面照
		negativeImageURL, err := storage.NewOssFileStorage(uuid.New().String()+".png", "image/png").Save(bytes.NewReader(negativeImageData))
		if err != nil {
			return err
		}
		if err := services.NewAccountService(g.DB).Authentication(
			user.ID, req.Name, req.IDCard, positiveImageURL, negativeImageURL); err != nil {
			return err
		}
		return g.OK(context, nil)
	})
}

func artistList() regia.HandleFunc {
	return g.Wrapper(func(context *regia.Context) error {
		repo := repository.NewAccountRepository(g.DB)
		paging := g.NewQueryPagination(context)
		list, err := repo.AuthorList(paging.Limit(), paging.Offset())
		if err != nil {
			return err
		}
		count, err := repo.AuthorCount()
		if err != nil {
			return err
		}
		return g.Many(context, list, count)
	})
}
