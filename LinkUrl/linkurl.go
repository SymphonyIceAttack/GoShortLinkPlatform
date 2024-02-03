package linkurl

import (
	databaseUtil "GoShortLinkPlatform/DataBase"
	"crypto/md5"
	"encoding/base64"

	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BeWrapDbFnType func(ctx *gin.Context, db *gorm.DB)

type LinkUrlBody struct {
	LinkUrl string `json:"url" binding:"required"`
}

var GenerateLink = BeWrapDbFnType(func(c *gin.Context, db *gorm.DB) {

	newLinkUrlbody := new(LinkUrlBody)

	if err := c.BindJSON(newLinkUrlbody); err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}

	if !isValidURL(newLinkUrlbody.LinkUrl) {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "Url UnValidURL"})
		return
	}
	// if !isURLAccessible(newLinkUrlbody.LinkUrl) {
	// 	c.JSON(http.StatusNotAcceptable, gin.H{"error": "Url UnValidURL"})
	// 	return
	// }

	shortLinkUrl := shortLink(newLinkUrlbody.LinkUrl)
	LinkObject := databaseUtil.LinkObject{ShortUrl: shortLinkUrl, WholeUrl: newLinkUrlbody.LinkUrl}
	result := db.Create(&LinkObject)
	if result.Error != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"url": shortLinkUrl})

})

func shortLink(LinkUrl string) string {
	url := []byte(LinkUrl)

	// 计算MD5哈希值
	hash := md5.Sum(url)

	var newHash [8]byte

	for i := 0; i < 8; i++ {
		newHash[i] = hash[2*i] + hash[2*i+1]
	}

	// 使用Base64编码
	hashString := base64.StdEncoding.EncodeToString(newHash[:])

	return hashString

}

var ParseShortLink = BeWrapDbFnType(func(ctx *gin.Context, db *gorm.DB) {
	shortLinkUrl := ctx.Query("s")
	QuerylinkurlBody := databaseUtil.LinkObject{}

	result := db.Where(&databaseUtil.LinkObject{ShortUrl: shortLinkUrl}).First(&QuerylinkurlBody)
	if result.Error != nil {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"error": result.Error.Error()})
		return
	}

	ctx.Redirect(http.StatusMovedPermanently, QuerylinkurlBody.WholeUrl)
})

func isValidURL(inputURL string) bool {
	parsedURL, err := url.ParseRequestURI(inputURL)
	if err != nil {
		return false
	}

	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return false
	}

	return true
}

// func isURLAccessible(url string) bool {
// 	response, err := http.Head(url)
// 	if err != nil {
// 		return false
// 	}

// 	defer response.Body.Close()

// 	return response.StatusCode == http.StatusOK
// }
