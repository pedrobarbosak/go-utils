package captcha

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Request struct {
	CaptchaToken string `form:"captchaResponse" json:"captchaResponse" validate:"required"`
}

func Middleware(service Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request Request
		if err := ctx.ShouldBindBodyWith(&request, binding.JSON); err != nil {
			log.Println("[service Middleware] failed to bind request to captcha model:", err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid captcha"})
			ctx.Abort()
			return
		}

		if err := service.Verify(request.CaptchaToken); err != nil {
			log.Println("[service Middleware] failed to verify captcha token:", err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid captcha"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
