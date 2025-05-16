package handlers

import (
	"context"
	"fmt"
	"halalcloud/auth"
	"net/http"
	"time"

	pubUserOffline "github.com/city404/v6-public-rpc-proto/go/v6/offline"
	"github.com/gin-gonic/gin"
)

func ListOfflines(c *gin.Context) {

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

	var offlineRequestBody pubUserOffline.OfflineTaskListRequest
	// 从请求体中解析 JSON 到 requestBody 结构体
	if err := c.ShouldBindJSON(&offlineRequestBody); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据",
		})
		return
	}
	client := pubUserOffline.NewPubOfflineTaskClient(fserv.GetGrpcConnection())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	result, err := client.List(ctx, &offlineRequestBody)

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

func AddOffline(c *gin.Context) {

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

	var addOfflineRequestBody pubUserOffline.UserTask

	// 从请求体中解析 JSON 到 requestBody 结构体
	if err := c.ShouldBindJSON(&addOfflineRequestBody); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据",
		})
		return
	}

	client := pubUserOffline.NewPubOfflineTaskClient(fserv.GetGrpcConnection())
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	result, err := client.Add(ctx, &addOfflineRequestBody)

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

func ParseOffline(c *gin.Context) {

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

	var parseOfflineRequestBody pubUserOffline.TaskParseRequest

	// 从请求体中解析 JSON 到 requestBody 结构体
	if err := c.ShouldBindJSON(&parseOfflineRequestBody); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据",
		})
		return
	}

	client := pubUserOffline.NewPubOfflineTaskClient(fserv.GetGrpcConnection())
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	result, err := client.Parse(ctx, &parseOfflineRequestBody)

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

func DeleteOffline(c *gin.Context) {

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

	var deleteOfflineRequestBody pubUserOffline.OfflineTaskDeleteRequest

	// 从请求体中解析 JSON 到 requestBody 结构体
	if err := c.ShouldBindJSON(&deleteOfflineRequestBody); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据",
		})
		return
	}

	client := pubUserOffline.NewPubOfflineTaskClient(fserv.GetGrpcConnection())
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	result, err := client.Delete(ctx, &deleteOfflineRequestBody)

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
