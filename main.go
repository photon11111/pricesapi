package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"pricesapi/internal/api"
	"pricesapi/internal/db"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
		log.Warn("%s not set, using default: %s", key, defaultValue)
	}

	return value
}

func updatePrices(currencies []string, quotedCurrency string,) {
	for _, currency := range currencies {
		price, err := api.GetPriceFromBinance(currency, quotedCurrency)
		if err != nil {
			log.WithFields(log.Fields{
				"currency": currency,
				"error":    err,
			}).Error("Failed to get price from Binance")
			continue
		}

		err = db.SaveInstrumentInfoToDB(db.InstrumentInfo{
			UpdateTime:       time.Now(),
			Instrument:       currency,
			Price:            price,
			QuotedInstrument: quotedCurrency,
		})

		if err != nil {
			log.WithFields(log.Fields{
				"currency": currency,
				"error":    err,
			}).Error("Failed to save currency information to database")
		}
	}
}

func handlePriceRequest(c *gin.Context) {
	currency := c.Query("currency")
	if currency == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong query parameters"})
		return
	}

	curInfo, err := db.GetInstrumentInfoFromDB(currency, "USDT")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"base_currency":   curInfo.Instrument,
		"quoted_currency": curInfo.QuotedInstrument,
		"price":           curInfo.Price,
		"update_time":     curInfo.UpdateTime,
	})
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.WithError(err).Warn("Failed to load .env file")
	}

	quotedCurrency := getEnv("QUOTED_CURRENCY", "USDT")

	port := getEnv("PORT", "8080")

	currencies := []string{"BTC", "ETH", "BNB", "SOL", "TON", "TRX", "DOGE", "ADA", "AVAX"}

	if err := db.OpenDB(); err != nil {
		log.WithError(err).Fatal("Failed to open database")
	}
	defer db.CloseDB()

	r := gin.Default()
	r.GET("/price/binance", handlePriceRequest)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func (ctx context.Context, interval time.Duration, currencies []string, quotedCurrency string) {
		updatePrices(currencies, quotedCurrency)
	
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
	
		for {
			select {
			case <-ctx.Done():
				log.Info("Price update stopped")
				return
			case <-ticker.C:
				updatePrices(currencies, quotedCurrency)
			}
		}
	}(ctx, 1*time.Hour, currencies, quotedCurrency)


	startedCh := make(chan bool, 1)
	go func() {
		startedCh <- true

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("Server encountered an error")
		}
	}()

	<-startedCh
	log.Info("Listening and serve at :", port)

	<-quit
	log.Info("Shutting down server...")

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.WithError(err).Fatal("Server forced to shutdown")
	}

	log.Info("Server shut down gracefully")
}