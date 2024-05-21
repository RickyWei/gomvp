package middleware

import (
	"fmt"

	"github.com/rickywei/sparrow/project/logger"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pelletier/go-toml/v2"
	"go.uber.org/zap"
	"golang.org/x/text/language"
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
