package handlers

import (
	"context"
	"fmt"
	"halalcloud/auth"
	"net/http"
	"time"

	pbPublicUser "github.com/city404/v6-public-rpc-proto/go/v6/user"
	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) {

	serv, exists := c.Get("authService")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "authService not found"})
		return
	}
	fserv, ok := serv.(*auth.AuthService)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "认证服务出错！！！",
		})
		return
	}

	var getUserRequestBody pbPublicUser.User
	// 从请求体中解析 JSON 到 requestBody 结构体
	if err := c.ShouldBindJSON(&getUserRequestBody); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据",
		})
		return
	}
	client := pbPublicUser.NewPubUserClient(fserv.GetGrpcConnection())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	result, err := client.Get(ctx, &getUserRequestBody)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result":    result,
		"auth_info": fserv.GetAuth(),
	})
}

func UpdateUser(c *gin.Context) {

	serv, exists := c.Get("authService")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "authService not found"})
		return
	}
	fserv, ok := serv.(*auth.AuthService)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "认证服务出错！！！",
		})
		return
	}

	var getUserRequestBody pbPublicUser.User
	// 从请求体中解析 JSON 到 requestBody 结构体
	if err := c.ShouldBindJSON(&getUserRequestBody); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据",
		})
		return
	}
	client := pbPublicUser.NewPubUserClient(fserv.GetGrpcConnection())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	result, err := client.Update(ctx, &getUserRequestBody)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result":    result,
		"auth_info": fserv.GetAuth(),
	})
}

func ChangePassword(c *gin.Context) {
	serv, exists := c.Get("authService")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "authService not found"})
		return
	}
	fserv, ok := serv.(*auth.AuthService)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "认证服务出错！！！",
		})
		return
	}

	var chagePasswordRequestBody pbPublicUser.ChangePasswordRequest
	// 从请求体中解析 JSON 到 requestBody 结构体
	if err := c.ShouldBindJSON(&chagePasswordRequestBody); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据",
		})
		return
	}
	client := pbPublicUser.NewPubUserClient(fserv.GetGrpcConnection())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	result, err := client.ChangePassword(ctx, &chagePasswordRequestBody)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result":    result,
		"auth_info": fserv.GetAuth(),
	})
}

func RefreshUserAccessToken(c *gin.Context) {
	serv, exists := c.Get("authService")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "authService not found"})
		return
	}
	fserv, ok := serv.(*auth.AuthService)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "认证服务出错！！！",
		})
		return
	}

	var tokenRequestBody pbPublicUser.Token
	// 从请求体中解析 JSON 到 requestBody 结构体
	if err := c.ShouldBindJSON(&tokenRequestBody); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据",
		})
		return
	}
	client := pbPublicUser.NewPubUserClient(fserv.GetGrpcConnection())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	result, err := client.Refresh(ctx, &tokenRequestBody)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result":    result,
		"auth_info": fserv.GetAuth(),
	})
}

func QuotoUser(c *gin.Context) {
	serv, exists := c.Get("authService")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "authService not found"})
		return
	}
	fserv, ok := serv.(*auth.AuthService)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "认证服务出错！！！",
		})
		return
	}

	var quotoUserRequestBody pbPublicUser.User
	// 从请求体中解析 JSON 到 requestBody 结构体
	if err := c.ShouldBindJSON(&quotoUserRequestBody); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据",
		})
		return
	}
	client := pbPublicUser.NewPubUserClient(fserv.GetGrpcConnection())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	result, err := client.GetStatisticsAndQuota(ctx, &quotoUserRequestBody)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result":    result,
		"auth_info": fserv.GetAuth(),
	})
}
func LogoffUser(c *gin.Context) {
	serv, exists := c.Get("authService")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "authService not found"})
		return
	}
	fserv, ok := serv.(*auth.AuthService)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "认证服务出错！！！",
		})
		return
	}

	var tokenRequestBody pbPublicUser.Token
	// 从请求体中解析 JSON 到 requestBody 结构体
	if err := c.ShouldBindJSON(&tokenRequestBody); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据",
		})
		return
	}
	client := pbPublicUser.NewPubUserClient(fserv.GetGrpcConnection())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	result, err := client.Logoff(ctx, &tokenRequestBody)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result":    result,
		"auth_info": fserv.GetAuth(),
	})
}
