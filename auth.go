package gu

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"strings"
)

type FiberAuthConfig struct {
	Secret          string
	GetTokenStrFunc func(c *fiber.Ctx) string
}

var (
	defaultConfig = &FiberAuthConfig{
		Secret: JWTSECRET,
		GetTokenStrFunc: func(c *fiber.Ctx) string {
			token := c.Cookies("token")
			if token == "" {
				token = strings.TrimSpace(strings.TrimPrefix(c.Get("Authorization"), "Bearer"))
			}
			return token
		},
	}
)

type jwtDecodeFunc func(token string) (map[string]any, error)

func configVerify(configs ...*FiberAuthConfig) (jwtDecodeFunc, *FiberAuthConfig) {
	var (
		cfg *FiberAuthConfig
	)

	if len(configs) > 0 && configs[0] != nil {
		cfg = configs[0]
	} else {
		cfg = &FiberAuthConfig{}
	}

	if cfg.GetTokenStrFunc == nil {
		cfg.GetTokenStrFunc = defaultConfig.GetTokenStrFunc
	}

	jwtDecode := func(token string) (map[string]any, error) {
		var (
			pt     *jwt.Token
			claims jwt.MapClaims
			fe     error
			ok     bool
		)

		if pt, fe = jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("加密方法错误: %v", t.Header["alg"])
			}

			return []byte(cfg.Secret), nil
		}); fe != nil {
			logrus.Debugf("token解析失败: %v", fe)
			return nil, fe
		}

		if !pt.Valid {
			return nil, errors.New("解析token无效")
		}

		if claims, ok = pt.Claims.(jwt.MapClaims); !ok {
			return nil, errors.New("转化jwt claims失败")
		}

		return claims, nil
	}
	return jwtDecode, cfg
}

func FiberAuth(user any, configs ...*FiberAuthConfig) fiber.Handler {
	jwtDecode, cfg := configVerify(configs...)
	return func(c *fiber.Ctx) error {
		var (
			err     error
			dataMap = map[string]any{}
			token   = cfg.GetTokenStrFunc(c)
		)

		if token == "" {
			return Resp401(c, "token为空")
		}

		if dataMap, err = jwtDecode(token); err != nil {
			return Resp401(c, err.Error())
		}

		bytes, _ := json.Marshal(dataMap["user"])
		if err = json.Unmarshal(bytes, &user); err != nil {
			logrus.Fatalf("序列化失败: %v", err)
		}

		c.Locals("user", user)
		c.Locals("token", token)

		return c.Next()
	}
}
