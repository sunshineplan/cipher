package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sunshineplan/cipher"
)

func run() {
	if logPath != "" {
		f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
		if err != nil {
			log.Fatalln("Failed to open log file:", err)
		}
		gin.DefaultWriter = f
		gin.DefaultErrorWriter = f
		log.SetOutput(f)
	}

	router := gin.Default()
	server.Handler = router
	router.StaticFS("/static", http.Dir(filepath.Join(filepath.Dir(self), "static")))
	router.LoadHTMLGlob(filepath.Join(filepath.Dir(self), "templates/*"))
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	router.POST("/do", func(c *gin.Context) {
		mode := c.PostForm("mode")
		key := c.PostForm("key")
		content := c.PostForm("content")
		switch mode {
		case "encrypt":
			c.JSON(200, gin.H{"result": cipher.Encrypt(key, content)})
		case "decrypt":
			result, err := cipher.Decrypt(key, strings.TrimSpace(content))
			if err != nil {
				c.JSON(200, gin.H{"result": nil})
			} else {
				c.JSON(200, gin.H{"result": result})
			}
		default:
			c.String(400, "")
		}
	})

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
