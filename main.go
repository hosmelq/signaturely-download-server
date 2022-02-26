package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/download", func(c *gin.Context) {
		file := c.PostForm("file")
		buf := new(bytes.Buffer)
		reader := base64.NewDecoder(
			base64.StdEncoding,
			strings.NewReader(strings.TrimPrefix(file, "data:image/png;base64,")),
		)
		i, _, err := image.Decode(reader)

		if err != nil {
			fmt.Println("Failed to decode file:", err)
			c.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}

		err = png.Encode(buf, i)

		if err != nil {
			fmt.Println("Failed to create buffer:", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		data := buf.Bytes()

		c.Writer.WriteHeader(http.StatusOK)
		c.Header("Content-Disposition", "attachment; filename=signature.png")
		c.Header("Content-Length", fmt.Sprintf("%d", len(data)))
		c.Header("Content-Type", "image/png")

		_, err = c.Writer.Write(data)

		if err != nil {
			fmt.Println("Failed to write data to response:", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	})

	r.Run()
}
