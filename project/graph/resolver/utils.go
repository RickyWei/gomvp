package resolver

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"github.com/rickywei/sparrow/project/conf"
)

func setJwtCookie(ctx *gin.Context, uid string) error {
	exp := conf.Int("jwt.exp")
	now := time.Now()
	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Subject: uid,
			ExpiresAt: &jwt.NumericDate{
				Time: now.Add(time.Second * time.Duration(exp)),
			},
			NotBefore: &jwt.NumericDate{
				Time: now,
			},
			IssuedAt: &jwt.NumericDate{
				Time: now,
			},
		},
	)
	token, err := t.SignedString([]byte(conf.String("jwt.secret")))
	if err != nil {
		return errors.Wrap(err, "jwt sign failed")
	}

	ctx.SetCookie("token", token, exp, "/", "localhost", false, true)

	return nil
}
