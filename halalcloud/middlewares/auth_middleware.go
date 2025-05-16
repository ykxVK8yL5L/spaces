package middlewares

import (
	"halalcloud/auth"
	"halalcloud/constants"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Middleware to create serv from Refresh-Token and store them in the context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		refreshToken := c.GetHeader("Refresh-Token")
		if refreshToken == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Refresh token is missing",
			})
			c.Abort()
			return
		}

		accessToken := ""
		var accessTokenExpiredAt int64 = 0
		accessToken = c.GetHeader("Access-Token")

		expiredAtStr := c.GetHeader("Access-Token-Expired-At")
		if expiredAtStr != "" {
			var err error
			accessTokenExpiredAt, err = strconv.ParseInt(expiredAtStr, 10, 64)
			if err != nil {
				// 如果转换失败，返回错误响应
				c.JSON(400, gin.H{"error": "Invalid Access-Token-Expired-At"})
				return
			}
		}

		serv, err := auth.NewAuthService(constants.AppID, constants.AppVersion, constants.AppSecret, refreshToken, accessToken, accessTokenExpiredAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			c.Abort()
			return
		}
		c.Set("authService", serv)
		// println(serv.GetAuth())

		// userJson, err := json.Marshal(serv.GetAuth())
		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{
		// 		"message": "无法获取用户信息",
		// 	})
		// 	c.Abort()
		// 	return
		// }

		//client := pubUserFile.NewPubUserFileClient(serv.GetGrpcConnection())
		//c.Set("userClient", client)

		c.Next()
	}
}
