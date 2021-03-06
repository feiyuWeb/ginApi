package main

import (
	"ginApi/controls"
	_ "ginApi/docs" // 注意这个一定要引入自己的docs
	"ginApi/middleware"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Golang Gin API
// @version 2.0
// @description An example of gin
// @termsOfService 运行地址：http://localhost:8081/swagger/index.html
// @license.name MIT //localhost:8081
func main() {
	f, _ := os.Create("gin.log") // 创建gin.log日志文件
	// gin.DefaultWriter = io.MultiReader(f)
	gin.DefaultErrorWriter = io.MultiWriter(f) // 错误信息写入log日志文件
	/*
		cors.Default() 默认允许所有源跨域
		跨域要写在路由前面，需要先执行
	*/
	router := gin.Default()
	router.Use(cors.Default())

	url := ginSwagger.URL(":8081/swagger/doc.json") // The url pointing to API definition
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	// 测试api
	router.GET("/api", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello go !",
		})
	})

	v2 := router.Group("/api/v2")
	{
		// 音乐接口
		v2.GET("/music", controls.MusicList)
		v2.POST("/music", controls.MusicCreate)
		v2.PUT("/music/:id", controls.MusicUpdate)
		v2.DELETE("/music/:id", controls.MusicDelete)

		// 电影接口
		v2.GET("/film", controls.FilmList)
		v2.POST("/film", controls.FilmCreate)
		v2.PUT("/film/:id", controls.FilmUpdate)
		v2.DELETE("/film/:id", controls.FilmDelete)

		// 书籍列表
		v2.GET("/book", controls.BookList)
		v2.POST("/book", controls.BookCreate)
		v2.PUT("/book/:id", controls.BookUpdate)
		v2.DELETE("/book/:id", controls.BookDelete)

		// 文件上传
		v2.POST("/upload", controls.UploadFile)

		// 注册及登录
		v2.POST("/register", controls.RegisterUser)
		v2.POST("/login", controls.Login)

		v2.GET("/userList",middleware.JWTAuth(), controls.UserList)
		v2.DELETE("/userList/:id", middleware.JWTAuth(), controls.UserDelete)
	}

	router.Run(":8081")
	reLaunch()
}

// 自动部署
func reLaunch() {
	cmd := exec.Command("sh", "./deploy.sh")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.Wait()
}
