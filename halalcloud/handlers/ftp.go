package handlers

import (
	"context"
	"fmt"
	"halalcloud/auth"
	"net/http"
	"time"

	pbSftpConfig "github.com/city404/v6-public-rpc-proto/go/v6/sftpconfig"
	"github.com/gin-gonic/gin"
)

func EnableFtp(c *gin.Context) {
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

	var ftpRequestBody pbSftpConfig.SftpConfig
	// 从请求体中解析 JSON 到 requestBody 结构体
	if err := c.ShouldBindJSON(&ftpRequestBody); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据",
		})
		return
	}
	client := pbSftpConfig.NewPubSftpConfigClient(fserv.GetGrpcConnection())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	result, err := client.Enable(ctx, &ftpRequestBody)

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

func DisableFtp(c *gin.Context) {
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

	var ftpRequestBody pbSftpConfig.SftpConfig
	// 从请求体中解析 JSON 到 requestBody 结构体
	if err := c.ShouldBindJSON(&ftpRequestBody); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据",
		})
		return
	}
	client := pbSftpConfig.NewPubSftpConfigClient(fserv.GetGrpcConnection())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	result, err := client.Disable(ctx, &ftpRequestBody)

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

func GetFtp(c *gin.Context) {
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

	var ftpRequestBody pbSftpConfig.SftpConfig
	// 从请求体中解析 JSON 到 requestBody 结构体
	if err := c.ShouldBindJSON(&ftpRequestBody); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据",
		})
		return
	}
	client := pbSftpConfig.NewPubSftpConfigClient(fserv.GetGrpcConnection())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	result, err := client.Get(ctx, &ftpRequestBody)

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

func UpdateFtp(c *gin.Context) {
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

	var ftpRequestBody pbSftpConfig.SftpConfig
	// 从请求体中解析 JSON 到 requestBody 结构体
	if err := c.ShouldBindJSON(&ftpRequestBody); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据",
		})
		return
	}
	client := pbSftpConfig.NewPubSftpConfigClient(fserv.GetGrpcConnection())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	result, err := client.Update(ctx, &ftpRequestBody)

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

func UpdateKeysFtp(c *gin.Context) {
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

	var ftpRequestBody pbSftpConfig.SftpConfig
	// 从请求体中解析 JSON 到 requestBody 结构体
	if err := c.ShouldBindJSON(&ftpRequestBody); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据",
		})
		return
	}
	client := pbSftpConfig.NewPubSftpConfigClient(fserv.GetGrpcConnection())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	result, err := client.UpdateKeys(ctx, &ftpRequestBody)

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

func ValidateUserNameFtp(c *gin.Context) {
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

	var ftpRequestBody pbSftpConfig.SftpConfig
	// 从请求体中解析 JSON 到 requestBody 结构体
	if err := c.ShouldBindJSON(&ftpRequestBody); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据",
		})
		return
	}
	client := pbSftpConfig.NewPubSftpConfigClient(fserv.GetGrpcConnection())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	result, err := client.ValidateUserName(ctx, &ftpRequestBody)

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
