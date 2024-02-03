package main

import (
	databaseUtil "GoShortLinkPlatform/DataBase"
	linkurl "GoShortLinkPlatform/LinkUrl"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	//部署时无需这个env文件直接在vercel中设置中设置环境变量
	godotenv.Overload("./.env.local")

	db, err := databaseUtil.LoadDataBase()
	// db.AutoMigrate(&databaseUtil.LinkObject{})

	if err != nil {
		log.Fatal(nil)
	}

	// 设置定时器清除任务
	ticker := time.Tick(3600 * time.Second) // 每3600秒执行一次任务

	go func() {
		for range ticker {
			// 这里执行定时任务的具体逻辑
			result := db.Where("created_at < ?", time.Now().Add(-24*time.Hour)).Delete(&databaseUtil.LinkObject{})
			if result.Error != nil {
				fmt.Println("Error deleting records:", result.Error)
				return
			}
			// 输出受影响的行数
			fmt.Println("Rows affected:", result.RowsAffected)
		}
	}()
	router := gin.Default()
	router.Use(cors.Default())
	// Handling routing errors
	router.NoRoute(func(c *gin.Context) {
		sb := &strings.Builder{}
		sb.WriteString("routing err: no route, try this:\n")
		for _, v := range router.Routes() {
			sb.WriteString(fmt.Sprintf("%s %s\n", v.Method, v.Path))
		}
		c.String(http.StatusBadRequest, sb.String())
	})
	router.POST("/generateLink", wrapDB(linkurl.GenerateLink, db))
	router.GET("/:s", wrapDB(linkurl.ParseShortLink, db))
	router.Run(fmt.Sprintf(":%s", port))
}

func wrapDB(fn linkurl.BeWrapDbFnType, db *gorm.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		fn(ctx, db)
	}
}
