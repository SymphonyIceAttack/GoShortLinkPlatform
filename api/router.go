package api

import (
	databaseUtil "GoShortLinkPlatform/DataBase"
	linkurl "GoShortLinkPlatform/LinkUrl"
	"GoShortLinkPlatform/handler"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var router *gin.Engine

//	func init() {
//		Init()
//	}
func init() {

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
	router.GET("/:shortLinkUrl", wrapDB(linkurl.ParseShortLink, db))

}
func wrapDB(fn linkurl.BeWrapDbFnType, db *gorm.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		fn(ctx, db)
	}
}

func Listen(w http.ResponseWriter, r *http.Request) {
	router.ServeHTTP(w, r)
}
