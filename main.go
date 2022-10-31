package main

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

const __idletters string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var __idlength int = 10
var __dirname string = "/home/user/workspace/projects/shortak/database/"
var __port string = "3200"

type igeneratebody struct {
	Url string `form:"url" json:"url"  binding:"required"`
}

func getShort(key string) string {
	if hasShort(key) {
		str, _ := os.ReadFile(__dirname + key)
		return string(str)
	}

	return ""
}

func setShort(key string, value string) {
	if hasShort(key) == false {
		os.WriteFile(__dirname+key, []byte(value), 0644)
	}
}

func hasShort(key string) bool {
	if _, err := os.Stat(__dirname + key); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}

func uuid() string {
	b := make([]byte, __idlength)
	for i := range b {
		b[i] = __idletters[rand.Int63()%int64(len(__idletters))]
	}
	return string(b)
}

func setup() {
	_port := os.Getenv("PORT")
	_dirname := os.Getenv("DBDIR")
	_length := os.Getenv("LENGTH")

	if len(_port) != 0 {
		__port = _port
	}

	if len(_dirname) != 0 {
		__dirname = _dirname
	}

	if hasShort("") == false {
		os.Mkdir(__dirname, 0755)
	}

	if len(_length) != 0 {
		__length, _ := strconv.Atoi(_length)
		__idlength = __length
	}

	fmt.Println("=== SETUP Start ===")
	fmt.Printf("PORT: %s\n", __port)
	fmt.Printf("Database Path: %s\n", __dirname)
	fmt.Printf("ID Length: %v\n", __idlength)
	fmt.Println("=== SETUP End ===")
}

func useAPI() {
	gin.SetMode(gin.ReleaseMode)
	server := gin.Default()

	server.POST("/api/v1/short", func(c *gin.Context) {
		var body igeneratebody

		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(400, gin.H{
				"status":  false,
				"code":    400,
				"message": "bad data",
			})
			return
		}

		if len(body.Url) == 0 {
			c.JSON(400, gin.H{
				"status":  false,
				"code":    400,
				"message": "url needed in body",
			})
		} else {
			id := uuid()
			setShort(id, body.Url)

			c.JSON(200, gin.H{
				"status": true,
				"code":   200,
				"data": gin.H{
					"id":  id,
					"url": body.Url,
				},
			})
		}
	})

	server.GET("/:key", func(c *gin.Context) {
		key := c.Param("key")

		if len(key) != __idlength {
			c.JSON(400, gin.H{
				"status":  false,
				"code":    400,
				"message": "Bad data",
			})
		} else if hasShort(key) == false {
			c.JSON(404, gin.H{
				"status":  false,
				"code":    404,
				"message": "Not found",
			})
		} else {
			c.Redirect(http.StatusMovedPermanently, getShort(key))
		}
	})

	fmt.Printf("Server starting at port %s\n", __port)
	server.Run(":" + __port)

}

func main() {
	setup()
	useAPI()
}