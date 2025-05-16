package handlers

import (
	"context"
	"fmt"
	"halalcloud/auth"
	"net/http"
	"time"

	pubFileShare "github.com/city404/v6-public-rpc-proto/go/v6/fileshare"
	"github.com/gin-gonic/gin"
)

func ListShares(c *gin.Context) {

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

	var shareRequestBody pubFileShare.FileShareListRequest
	// 从请求体中解析 JSON 到 requestBody 结构体
	if err := c.ShouldBindJSON(&shareRequestBody); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据",
		})
		return
	}
	client := pubFileShare.NewPubFileShareClient(fserv.GetGrpcConnection())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	result, err := client.List(ctx, &shareRequestBody)

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

func CreateShare(c *gin.Context) {

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

	var createShareRequestBody pubFileShare.FileShare

	// 从请求体中解析 JSON 到 requestBody 结构体
	if err := c.ShouldBindJSON(&createShareRequestBody); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据",
		})
		return
	}

	client := pubFileShare.NewPubFileShareClient(fserv.GetGrpcConnection())
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	result, err := client.Create(ctx, &createShareRequestBody)

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

func DeleteShare(c *gin.Context) {
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

	var deleteShareRequestBody pubFileShare.FileShareDeleteRequest

	// 从请求体中解析 JSON 到 requestBody 结构体
	if err := c.ShouldBindJSON(&deleteShareRequestBody); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据",
		})
		return
	}

	client := pubFileShare.NewPubFileShareClient(fserv.GetGrpcConnection())
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	result, err := client.Delete(ctx, &deleteShareRequestBody)

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

func GetShare(c *gin.Context) {
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

	var getShareRequestBody pubFileShare.FileShare

	// 从请求体中解析 JSON 到 requestBody 结构体
	if err := c.ShouldBindJSON(&getShareRequestBody); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据",
		})
		return
	}

	client := pubFileShare.NewPubFileShareClient(fserv.GetGrpcConnection())
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	result, err := client.Get(ctx, &getShareRequestBody)

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

func SaveShare(c *gin.Context) {
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

	var saveShareRequestBody pubFileShare.FileShare

	// 从请求体中解析 JSON 到 requestBody 结构体
	if err := c.ShouldBindJSON(&saveShareRequestBody); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据",
		})
		return
	}

	client := pubFileShare.NewPubFileShareClient(fserv.GetGrpcConnection())
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	result, err := client.Get(ctx, &saveShareRequestBody)

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
