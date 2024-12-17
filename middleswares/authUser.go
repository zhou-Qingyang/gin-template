package middleswares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"tz-gin/global"
	"tz-gin/models/common/response"
	"tz-gin/utils"
)

func AuthUserCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("x-token")
		if token == "" {
			response.FailWithErrorCode(3, "token is empty", c)
			c.Abort()
			return
		}
		j := utils.NewJWT()
		claims, err := j.ParseToken(token)
		if err != nil {
			response.FailWithErrorCode(4, "token is invalid", c)
			c.Abort()
			return
		}

		if c.Request.Method == "DELETE" && c.Request.URL.Path == "/api/user" {
			global.LocalCache.Delete(fmt.Sprintf("token_%s", token))
			c.Next()
			return
		}

		// 判断登录是否过期
		_, found := global.LocalCache.Get(fmt.Sprintf("token_%s", token))
		if found {
			c.Set("claims", claims)
			c.Next()
		} else {
			response.FailWithErrorCode(3, "登录过期 请重新登录", c)
			c.Abort()
			return
		}
	}
}
