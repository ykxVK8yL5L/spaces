package handlers

import (
	"context"
	"fmt"
	"halalcloud/auth"
	"net/http"
	"time"

	pubUserFile "github.com/city404/v6-public-rpc-proto/go/v6/userfile"
	"github.com/gin-gonic/gin"
)

func ListFiles(c *gin.Context) {
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
	client := pubUserFile.NewPubUserFileClient(fserv.GetGrpcConnection())
	var fileRequestBody pubUserFile.FileListRequest

	// 从请求体中解析 JSON 到 requestBody 结构体
	if err := c.ShouldBindJSON(&fileRequestBody); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	result, err := client.List(ctx, &fileRequestBody)

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

func ListTrashFiles(c *gin.Context) {
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
	client := pubUserFile.NewPubUserFileClient(fserv.GetGrpcConnection())

	var fileRequestBody pubUserFile.FileListRequest

	// 从请求体中解析 JSON 到 requestBody 结构体
	if err := c.ShouldBindJSON(&fileRequestBody); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	result, err := client.ListTrash(ctx, &fileRequestBody)
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

func GetFile(c *gin.Context) {
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
	client := pubUserFile.NewPubUserFileClient(fserv.GetGrpcConnection())

	var fileRequestBody pubUserFile.File

	// 从请求体中解析 JSON 到 requestBody 结构体
	if err := c.ShouldBindJSON(&fileRequestBody); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	result, err := client.Get(ctx, &fileRequestBody)
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

func SearchFiles(c *gin.Context) {
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
	client := pubUserFile.NewPubUserFileClient(fserv.GetGrpcConnection())

	var searchFileRequestBody pubUserFile.SearchRequest

	// 从请求体中解析 JSON 到 requestBody 结构体
	if err := c.ShouldBindJSON(&searchFileRequestBody); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	result, err := client.Search(ctx, &searchFileRequestBody)
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

func CreateFolderFile(c *gin.Context) {
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
	client := pubUserFile.NewPubUserFileClient(fserv.GetGrpcConnection())

	var fileRequestBody pubUserFile.File

	// 从请求体中解析 JSON 到 requestBody 结构体
	if err := c.ShouldBindJSON(&fileRequestBody); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	result, err := client.Create(ctx, &fileRequestBody)
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

func MoveFiles(c *gin.Context) {
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
	client := pubUserFile.NewPubUserFileClient(fserv.GetGrpcConnection())

	var batchOperationRequestBody pubUserFile.BatchOperationRequest

	// 从请求体中解析 JSON 到 requestBody 结构体
	if err := c.ShouldBindJSON(&batchOperationRequestBody); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	result, err := client.Move(ctx, &batchOperationRequestBody)
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

func CopyFiles(c *gin.Context) {
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
	client := pubUserFile.NewPubUserFileClient(fserv.GetGrpcConnection())

	var batchOperationRequestBody pubUserFile.BatchOperationRequest

	// 从请求体中解析 JSON 到 requestBody 结构体
	if err := c.ShouldBindJSON(&batchOperationRequestBody); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	result, err := client.Copy(ctx, &batchOperationRequestBody)
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

func TrashFiles(c *gin.Context) {
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
	client := pubUserFile.NewPubUserFileClient(fserv.GetGrpcConnection())

	var batchOperationRequestBody pubUserFile.BatchOperationRequest

	// 从请求体中解析 JSON 到 requestBody 结构体
	if err := c.ShouldBindJSON(&batchOperationRequestBody); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	result, err := client.Trash(ctx, &batchOperationRequestBody)
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

func DeleteFiles(c *gin.Context) {
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
	client := pubUserFile.NewPubUserFileClient(fserv.GetGrpcConnection())

	var batchOperationRequestBody pubUserFile.BatchOperationRequest

	// 从请求体中解析 JSON 到 requestBody 结构体
	if err := c.ShouldBindJSON(&batchOperationRequestBody); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	result, err := client.Delete(ctx, &batchOperationRequestBody)
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

func RenameFiles(c *gin.Context) {
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
	client := pubUserFile.NewPubUserFileClient(fserv.GetGrpcConnection())

	var fileRequestBody pubUserFile.File

	// 从请求体中解析 JSON 到 requestBody 结构体
	if err := c.ShouldBindJSON(&fileRequestBody); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	result, err := client.Rename(ctx, &fileRequestBody)
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

func BatchRenameFiles(c *gin.Context) {
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
	client := pubUserFile.NewPubUserFileClient(fserv.GetGrpcConnection())

	var batchOperationRequestBody pubUserFile.BatchOperationRequest

	// 从请求体中解析 JSON 到 requestBody 结构体
	if err := c.ShouldBindJSON(&batchOperationRequestBody); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	result, err := client.BatchRename(ctx, &batchOperationRequestBody)
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

func RecoverFiles(c *gin.Context) {
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
	client := pubUserFile.NewPubUserFileClient(fserv.GetGrpcConnection())

	var batchOperationRequestBody pubUserFile.BatchOperationRequest

	// 从请求体中解析 JSON 到 requestBody 结构体
	if err := c.ShouldBindJSON(&batchOperationRequestBody); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	result, err := client.Recover(ctx, &batchOperationRequestBody)
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

func DeleteTrashFiles(c *gin.Context) {
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
	client := pubUserFile.NewPubUserFileClient(fserv.GetGrpcConnection())

	var batchOperationRequestBody pubUserFile.BatchOperationRequest

	// 从请求体中解析 JSON 到 requestBody 结构体
	if err := c.ShouldBindJSON(&batchOperationRequestBody); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	result, err := client.DeleteTrash(ctx, &batchOperationRequestBody)
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

func ParseFileSliceFiles(c *gin.Context) {
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
	client := pubUserFile.NewPubUserFileClient(fserv.GetGrpcConnection())

	var fileRequestBody pubUserFile.File

	// 从请求体中解析 JSON 到 requestBody 结构体
	if err := c.ShouldBindJSON(&fileRequestBody); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	result, err := client.ParseFileSlice(ctx, &fileRequestBody)
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

func GetSliceDownloadAddressFiles(c *gin.Context) {
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
	client := pubUserFile.NewPubUserFileClient(fserv.GetGrpcConnection())

	var sliceRequestBody pubUserFile.SliceDownloadAddressRequest

	// 从请求体中解析 JSON 到 requestBody 结构体
	if err := c.ShouldBindJSON(&sliceRequestBody); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	result, err := client.GetSliceDownloadAddress(ctx, &sliceRequestBody)
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
