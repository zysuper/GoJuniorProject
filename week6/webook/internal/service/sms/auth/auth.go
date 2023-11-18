package auth

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/service/sms"
	"github.com/golang-jwt/jwt/v5"
)

type AuthSmsService struct {
	svc sms.Service
	key []byte
}

type SmsClaims struct {
	jwt.RegisteredClaims
	tpl string
}

func (a *AuthSmsService) Send(ctx context.Context, tplToken string, args []string, numbers ...string) error {
	var claims SmsClaims
	_, err := jwt.ParseWithClaims(tplToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return a.key, nil
	})

	if err != nil {
		return err
	}

	return a.svc.Send(ctx, claims.tpl, args, numbers...)
}
