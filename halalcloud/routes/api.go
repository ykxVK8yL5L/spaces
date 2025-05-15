package routes

import (
	"halalcloud/handlers"
	"halalcloud/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 用来配置路由
func SetupRoutes(router *gin.Engine) {
	// API 路由分组
	api := router.Group("/api")
	api.Use(middlewares.AuthMiddleware())

	file := api.Group("/file")
	{

		file.POST("/listtrash", handlers.ListTrashFiles)
		file.POST("/get", handlers.GetFile)
		file.POST("/search", handlers.SearchFiles)
		file.POST("/list", handlers.ListFiles)
		file.POST("/create", handlers.CreateFolderFile)
		file.POST("/move", handlers.MoveFiles)
		file.POST("/copy", handlers.CopyFiles)
		file.POST("/trash", handlers.TrashFiles)
		file.POST("/rename", handlers.RenameFiles)
		file.POST("/batchrename", handlers.BatchRenameFiles)
		file.POST("/recover", handlers.RecoverFiles)
		file.POST("/parsefileslice", handlers.ParseFileSliceFiles)
		file.POST("/getslicedownloadaddress", handlers.GetSliceDownloadAddressFiles)
		file.DELETE("/delete", handlers.DeleteFiles)
		file.DELETE("/deletetrash", handlers.DeleteTrashFiles)
	}

	user := api.Group("/user")
	{
		user.POST("/get", handlers.GetUser)
		user.POST("/quoto", handlers.QuotoUser)
		user.POST("/refresh", handlers.RefreshUserAccessToken)
		user.POST("/logoff", handlers.LogoffUser)
		user.POST("/update", handlers.UpdateUser)
		user.POST("/changepassword", handlers.LogoffUser)
	}

	offline := api.Group("/offline")
	{
		offline.POST("/list", handlers.ListOfflines)
		offline.POST("/add", handlers.AddOffline)
		offline.POST("/parse", handlers.ParseOffline)
		offline.DELETE("/delete", handlers.DeleteOffline)
	}

	share := api.Group("/share")
	{
		share.POST("/list", handlers.ListShares)
		share.POST("/get", handlers.GetShare)
		share.POST("/create", handlers.CreateShare)
		share.POST("/save", handlers.SaveShare)
		share.DELETE("/delete", handlers.DeleteShare)
	}

	webdav := api.Group("/webdav")
	{
		webdav.POST("/get", handlers.GetWebdav)
		webdav.POST("/listall", handlers.ListAllWebdav)
		webdav.DELETE("/delete", handlers.DeleteWebdav)
		webdav.POST("/create", handlers.CreateWebdav)
		webdav.POST("/enable", handlers.EnableWebdav)
		webdav.POST("/disable", handlers.DisableWebdav)
		webdav.POST("/update", handlers.UpdateWebdav)
		webdav.POST("/validateusername", handlers.ValidateUserNameWebdav)
	}

	ftp := api.Group("/ftp")
	{
		ftp.POST("/get", handlers.GetFtp)
		ftp.POST("/enable", handlers.EnableFtp)
		ftp.POST("/disable", handlers.DisableFtp)
		ftp.POST("/update", handlers.UpdateFtp)
		ftp.POST("/updatekeys", handlers.UpdateKeysFtp)
		ftp.POST("/validateusername", handlers.ValidateUserNameFtp)
	}

}
