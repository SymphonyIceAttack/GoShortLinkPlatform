package main

import (
	databaseUtil "GoShortLinkPlatform/DataBase"
	linkurl "GoShortLinkPlatform/LinkUrl"
	"GoShortLinkPlatform/handler"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

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
	db.AutoMigrate(&databaseUtil.LinkObject{})

	if err != nil {
		log.Fatal(nil)
	}

	router := gin.Default()
	// Handling routing errors
	router.NoRoute(func(c *gin.Context) {
		sb := &strings.Builder{}
		sb.WriteString("routing err: no route, try this:\n")
		for _, v := range router.Routes() {
			sb.WriteString(fmt.Sprintf("%s %s\n", v.Method, v.Path))
		}
		c.String(http.StatusBadRequest, sb.String())
	})
	router.POST("/generateLink", handler.Cors, wrapDB(linkurl.GenerateLink, db))
	router.GET("/s", handler.Cors, wrapDB(linkurl.ParseShortLink, db))
	router.Run(fmt.Sprintf(":%s", port))
}

func wrapDB(fn linkurl.BeWrapDbFnType, db *gorm.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		fn(ctx, db)
	}
}
