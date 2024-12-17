package middleswares

import (
	"github.com/gin-gonic/gin"
	"tz-gin/models/common/response"
	"tz-gin/utils"
)

func AuthAdminCheck() gin.HandlerFunc {
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
		if claims.IsAdmin != 1 {
			response.FailWithErrorCode(5, "you are not admin", c)
			c.Abort()
			return
		}
		// 判断权限
		c.Set("claims", claims)
		c.Next()
	}
}
