package main

import (
	"performance-analyzer/handlers"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func main() {
	debug.SetGCPercent(1)
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(CORSMiddleware())

	r.POST("/performance-analyzer/signals/parse", handlers.ParseSignalsHandler)
	r.POST("/performance-analyzer/signals/endresponse", handlers.EndpointResponseHandler)
	r.GET("/performance-analyzer/signals/analyzedata/:TsInterval", handlers.GetAnalyzedDataHandler)
	r.POST("/performance-analyzer/signals/analyzetelegrams/:TsInterval", handlers.AnalyzeCapMqttDataHandler)

	err := r.Run(":4300")
	if err != nil {
		panic(err)
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, accept")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
