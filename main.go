package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Body struct {
	Image string `json:"image"`
}

func main() {
	r := gin.Default()

	r.Use(cors.Default())

	r.POST("/signatures", func(c *gin.Context) {
		body := Body{}

		if err := c.BindJSON(&body); err != nil {
			fmt.Println("Failed to bind JSON:", err)
			c.AbortWithStatus(http.StatusUnprocessableEntity)
		}

		buf := new(bytes.Buffer)
		reader := base64.NewDecoder(
			base64.StdEncoding,
			strings.NewReader(
				strings.TrimPrefix(body.Image, "data:image/png;base64,"),
			),
		)
		i, _, err := image.Decode(reader)

		if err != nil {
			fmt.Println("Failed to decode file:", err)
			c.AbortWithStatus(http.StatusUnprocessableEntity)
		}

		err = png.Encode(buf, i)

		if err != nil {
			fmt.Println("Failed to create buffer:", err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		uuid := uuid.NewString()

		os.WriteFile(fmt.Sprintf("%s.png", uuid), buf.Bytes(), 0644)

		c.IndentedJSON(http.StatusOK, gin.H{
			"uuid": uuid,
		})
	})

	r.GET("/signatures/:uuid/download", func(c *gin.Context) {
		file := fmt.Sprintf("%s.png", c.Param("uuid"))
		fi, err := os.Stat(file)

		if os.IsNotExist(err) {
			fmt.Println("File not found:", err)
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.FileAttachment(fi.Name(), "signature.png")
		os.Remove(file)
	})

	r.Run()
}
