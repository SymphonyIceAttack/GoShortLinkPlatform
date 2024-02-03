package main

import (
	databaseUtil "GoShortLinkPlatform/DataBase"
	linkurl "GoShortLinkPlatform/LinkUrl"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {

	godotenv.Load(".env.local")

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
