package middlewares

import (
	"halalcloud/auth"
	"halalcloud/constants"
	"net/http"

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
		serv, err := auth.NewAuthService(constants.AppID, constants.AppVersion, constants.AppSecret, refreshToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			c.Abort()
			return
		}
		c.Set("authService", serv)
		//client := pubUserFile.NewPubUserFileClient(serv.GetGrpcConnection())
		//c.Set("userClient", client)

		c.Next()
	}
}
