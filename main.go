package main

import (
	"net/http"
	"os"
	"pricesapi/internal/api"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Error("error of .env file loading")
		return
	}

	port := os.Getenv("PORT")

	if port == "" {
		log.Error("port wasn't specified")
		return
	}

	r := gin.Default()

	r.GET("/price/coingecko", func(c *gin.Context) {
		currency := c.Query("currency")
		in := c.Query("in")

		if currency == "" || in == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "wrong query parameters"})
			return
		}

		price, err := api.GetPriceFromCoingecko(currency, in)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"base_currency":   currency,
			"quoted_currency": in,
			"price":           price,
		})
	})

	r.GET("/price/binance", func(c *gin.Context) {
		currency := c.Query("currency")
		in := c.Query("in")

		if currency == "" || in == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "wrong query parameters"})
			return
		}

		price, err := api.GetPriceFromBinance(currency, in)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"base_currency":   currency,
			"quoted_currency": in,
			"price":           price,
		})
	})

	r.Run(":" + port)
}
