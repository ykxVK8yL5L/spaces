package handlers

import (
	"context"
	"fmt"
	"halalcloud/auth"
	"net/http"
	"time"

	pbDavConfig "github.com/city404/v6-public-rpc-proto/go/v6/webdavconfig"
	"github.com/gin-gonic/gin"
)

func EnableWebdav(c *gin.Context) {
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

	var webdavRequestBody pbDavConfig.DavConfig
	// 从请求体中解析 JSON 到 requestBody 结构体
	if err := c.ShouldBindJSON(&webdavRequestBody); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据",
		})
		return
	}
	client := pbDavConfig.NewPubDavConfigClient(fserv.GetGrpcConnection())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	result, err := client.Enable(ctx, &webdavRequestBody)

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

func DisableWebdav(c *gin.Context) {
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

	var webdavRequestBody pbDavConfig.DavConfig
	// 从请求体中解析 JSON 到 requestBody 结构体
	if err := c.ShouldBindJSON(&webdavRequestBody); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据",
		})
		return
	}
	client := pbDavConfig.NewPubDavConfigClient(fserv.GetGrpcConnection())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	result, err := client.Disable(ctx, &webdavRequestBody)

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

func GetWebdav(c *gin.Context) {
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

	var webdavRequestBody pbDavConfig.DavConfig
	// 从请求体中解析 JSON 到 requestBody 结构体
	if err := c.ShouldBindJSON(&webdavRequestBody); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据",
		})
		return
	}
	client := pbDavConfig.NewPubDavConfigClient(fserv.GetGrpcConnection())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	result, err := client.Get(ctx, &webdavRequestBody)

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

func DeleteWebdav(c *gin.Context) {
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

	var webdavRequestBody pbDavConfig.DavConfig
	// 从请求体中解析 JSON 到 requestBody 结构体
	if err := c.ShouldBindJSON(&webdavRequestBody); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据",
		})
		return
	}
	client := pbDavConfig.NewPubDavConfigClient(fserv.GetGrpcConnection())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	result, err := client.Delete(ctx, &webdavRequestBody)

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

func CreateWebdav(c *gin.Context) {
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

	var webdavRequestBody pbDavConfig.DavConfig
	// 从请求体中解析 JSON 到 requestBody 结构体
	if err := c.ShouldBindJSON(&webdavRequestBody); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据",
		})
		return
	}
	client := pbDavConfig.NewPubDavConfigClient(fserv.GetGrpcConnection())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	result, err := client.Create(ctx, &webdavRequestBody)

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

func ListAllWebdav(c *gin.Context) {
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

	var webdavRequestBody pbDavConfig.DavConfig
	// 从请求体中解析 JSON 到 requestBody 结构体
	if err := c.ShouldBindJSON(&webdavRequestBody); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据",
		})
		return
	}
	client := pbDavConfig.NewPubDavConfigClient(fserv.GetGrpcConnection())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	result, err := client.ListAll(ctx, &webdavRequestBody)

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

func UpdateWebdav(c *gin.Context) {
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

	var webdavRequestBody pbDavConfig.DavConfig
	// 从请求体中解析 JSON 到 requestBody 结构体
	if err := c.ShouldBindJSON(&webdavRequestBody); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据",
		})
		return
	}
	client := pbDavConfig.NewPubDavConfigClient(fserv.GetGrpcConnection())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	result, err := client.Update(ctx, &webdavRequestBody)

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

func ValidateUserNameWebdav(c *gin.Context) {
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

	var webdavRequestBody pbDavConfig.DavConfig
	// 从请求体中解析 JSON 到 requestBody 结构体
	if err := c.ShouldBindJSON(&webdavRequestBody); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据",
		})
		return
	}
	client := pbDavConfig.NewPubDavConfigClient(fserv.GetGrpcConnection())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	result, err := client.ValidateUserName(ctx, &webdavRequestBody)

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
