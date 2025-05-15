package handlers

import (
	"halalcloud/auth"
	"halalcloud/constants"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	login_url, err := auth.NewAuthServiceWithOauth(constants.AppID, constants.AppVersion, constants.AppSecret)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "发生错误！"})
	}
	c.JSON(http.StatusOK, gin.H{"url": login_url})
	// htmlContent := fmt.Sprintf(`
	// 		<!DOCTYPE html>
	// 		<html lang="en">
	// 		<head>
	// 		    <meta charset="UTF-8">
	// 		    <meta name="viewport" content="width=device-width, initial-scale=1.0">
	// 		    <title>信息</title>
	// 		</head>
	// 		<body>
	// 		    <h1>请打开以下链接，登陆并记下设置Refresh Token</h1>
	// 			<a %s>%s</a>
	// 		</body>
	// 		</html>
	// 	`, login_url, login_url)

	// // 设置响应头 Content-Type 为 text/html
	// c.Header("Content-Type", "text/html")
	// // 返回HTML字符串
	// c.String(http.StatusOK, htmlContent)

}
