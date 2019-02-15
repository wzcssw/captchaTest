package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"time"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
)

const (
	// Default number of digits in captcha solution.
	DefaultLen = 6
	// The number of captchas created that triggers garbage collection used
	// by default store.
	CollectNum = 100
	// Expiration time of captchas used by default store.
	Expiration = 10 * time.Minute
	// Standard width and height of a captcha image.
	StdWidth  = 210
	StdHeight = 60
)

var (
	ErrNotFound = errors.New("captcha: id not found")
)

func main() {
	r := gin.Default()
	r.Static("/assets", "./static")
	// r.GET("/captcha_img", func(c *gin.Context) {
	// 	cap := captcha.New()
	// 	writer := c.Writer
	// 	captcha.WriteImage(writer, cap, StdWidth, StdHeight)
	// })

	r.GET("/captcha", func(c *gin.Context) {
		captchaID := captcha.New()
		base64 := GetCaptchaImageBase64(captchaID)
		c.JSON(200, gin.H{
			"captchaID": captchaID,
			"img":       base64,
		})
	})

	r.POST("/verify", func(c *gin.Context) {
		captchaID := c.PostForm("captcha_id")
		result := c.PostForm("result")
		success := captcha.VerifyString(captchaID, result)
		c.JSON(200, gin.H{"success": success})
	})
	r.Run(":8666")

}

// GetCaptchaImageBase64 GetCaptchaImageBase64
func GetCaptchaImageBase64(captchaID string) string {

	bytesContainer := make([]byte, 0)
	byteBuffer := bytes.NewBuffer(bytesContainer)
	bufioWriter := bufio.NewWriter(byteBuffer)

	captcha.WriteImage(bufioWriter, captchaID, StdWidth, StdHeight)
	bufioWriter.Flush()
	return base64.StdEncoding.EncodeToString(byteBuffer.Bytes())
}
