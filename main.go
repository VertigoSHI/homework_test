package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code  int         `json:"code"`
	Error string      `json:"error"`
	Data  interface{} `json:"data"`
}

func main() {
	r := gin.Default()
	parser := NewParser("https://cloudflare-eth.com")
	r.GET("/current_block", func(c *gin.Context) {
		block := parser.GetCurrentBlock()
		c.JSON(200, Response{
			http.StatusOK,
			"",
			block,
		})
	})
	r.POST("/:address/subscribe", func(c *gin.Context) {
		address, _ := c.Params.Get("address")
		c.JSON(200, Response{
			http.StatusOK,
			"",
			parser.Subscribe(address),
		})
	})
	r.GET("/:address/transaction", func(c *gin.Context) {
		address, _ := c.Params.Get("address")
		result, err := parser.GetTransactions(address)
		if err != nil {
			c.JSON(500, Response{http.StatusInternalServerError, err.Error(), nil})
		} else {
			c.JSON(200, Response{
				http.StatusOK,
				"",
				result,
			})
		}
	})
	r.Run(":8080")
}
