package api

import (
	databaseUtil "GoShortLinkPlatform/DataBase"
	linkurl "GoShortLinkPlatform/LinkUrl"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var router *gin.Engine

func init() {
	Init()
}
func Init() {

	db, err := databaseUtil.LoadDataBase()
	// db.AutoMigrate(&databaseUtil.LinkObject{})

	if err != nil {
		log.Fatal(nil)
	}

	router := gin.Default()

	router.POST("/generateLink", wrapDB(linkurl.GenerateLink, db))
	router.GET("/:shortLinkUrl", wrapDB(linkurl.ParseShortLink, db))

	router.Run("localhost:8080")
}
func wrapDB(fn linkurl.BeWrapDbFnType, db *gorm.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		fn(ctx, db)
	}
}

func Listen(w http.ResponseWriter, r *http.Request) {
	router.ServeHTTP(w, r)
}
