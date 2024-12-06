package gu

import (
	"errors"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	trzh "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"strings"
)

var Verify *verify

type VerifyErrFunc func(fe validator.FieldError) error

type verify struct {
	validate *validator.Validate
	trans    ut.Translator
	errMap   map[string]VerifyErrFunc
}

func init() {
	zhTranslator := zh.New()
	uni := ut.New(zhTranslator)
	trans, _ := uni.GetTranslator("zh")
	validate := validator.New()
	if err := trzh.RegisterDefaultTranslations(validate, trans); err != nil {
		panic(err)
	}

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		jsonTag := fld.Tag.Get("json")
		guTag := fld.Tag.Get("gu")
		if guTag == "-" {
			return ""
		}
		if guTag != "" {
			return guTag
		}
		return jsonTag
	})

	Verify = &verify{
		validate: validate,
		trans:    trans,
		errMap:   make(map[string]VerifyErrFunc),
	}
}

func (v *verify) Register(tag string, fn validator.Func, err VerifyErrFunc) error {
	v.errMap[tag] = err
	return v.validate.RegisterValidation(tag, fn)
}

func (v *verify) Struct(s interface{}) error {
	err := v.validate.Struct(s)
	if err == nil {
		return nil
	}

	var verr validator.ValidationErrors
	if !errors.As(err, &verr) {
		return err
	}

	errs := make([]string, len(verr))
	for i := range verr {
		tag := verr[i].Tag()
		if fun, ok := v.errMap[tag]; ok {
			errs[i] = fun(verr[i]).Error()
			continue
		}
		errs[i] = verr[i].Translate(v.trans)
	}

	return errors.New(strings.Join(errs, ","))
}
