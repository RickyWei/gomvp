package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"golang.org/x/text/language"

	"github.com/rickywei/sparrow/project/conf"
	"github.com/rickywei/sparrow/project/logger"
)

var (
	bundle = i18n.NewBundle(language.English)
	langs  = []string{"en", "zh"}
)

func init() {
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	for _, lang := range langs {
		fn := fmt.Sprintf("active.%s.toml", lang)
		if _, err := bundle.LoadMessageFile(fn); err != nil {
			logger.L().Fatal(fmt.Sprintf("load %s failed", fn), zap.Error(err))
		}
	}
}

func I18N() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accept := ctx.GetHeader("Accept-Language")
		localizer := i18n.NewLocalizer(bundle, accept)
		ctx.Set("localizer", localizer)

		ctx.Next()
	}
}

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		wb := &bodyWriter{
			body:           &bytes.Buffer{},
			ResponseWriter: ctx.Writer,
		}
		ctx.Writer = wb

		ctx.Next()

		bs := wb.body.Bytes()
		r := make(map[string]any)
		_ = json.Unmarshal(bs, &r)

		logger.L().Info(
			"HTTP",
			zap.String("clientIP", ctx.ClientIP()),
			zap.String("method", ctx.Request.Method),
			zap.String("url", ctx.Request.URL.String()),
			zap.Duration("duration", time.Since(start)),
			zap.Int("status", ctx.Writer.Status()),
			zap.Int("code", cast.ToInt(r["code"])),
			zap.String("msg", cast.ToString(r["msg"])),
			zap.Any("data", r["data"]),
		)

		wb.ResponseWriter.Write(bs)
	}
}

type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

func Recover() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logger.L().Debug("recover", zap.Any("panic", r))
			}
		}()

		ctx.Next()
	}
}

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error
		defer func() {
			if err != nil {
				ctx.AbortWithError(http.StatusUnauthorized, err)
			}
		}()

		tokenStr, err := ctx.Cookie("token")
		if err != nil {
			return
		}

		claims := &jwt.MapClaims{}
		if _, err := jwt.ParseWithClaims(
			tokenStr,
			claims,
			func(t *jwt.Token) (interface{}, error) { return conf.String("jwt.secret"), nil },
			jwt.WithIssuedAt(),
		); err != nil {
			return
		}

		sub, err := claims.GetSubject()
		if err != nil {
			return
		}
		uid, err := strconv.ParseInt(sub, 10, 64)
		if err != nil {
			return
		}

		ctx.Set("uid", uid)

		ctx.Next()
	}
}

type ctxKey string

func GinContextToContext() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c := context.WithValue(ctx.Request.Context(), ctxKey("GinContextKey"), ctx)
		ctx.Request = ctx.Request.WithContext(c)
		ctx.Next()
	}
}

func GinContextFromContext(ctx context.Context) (gc *gin.Context, err error) {
	ginContext := ctx.Value(ctxKey("GinContextKey"))
	if ginContext == nil {
		err = fmt.Errorf("could not retrieve gin.Context")
		return
	}

	gc, ok := ginContext.(*gin.Context)
	if !ok {
		err = fmt.Errorf("gin.Context has wrong type")
		return
	}

	return
}
