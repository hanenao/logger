package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hanenao/logger/log"
)

type LogEntry struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	engine := gin.Default()
	engine.Use(gin.Recovery())

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	engine.Use(log.SetGCPLogger(projectID, "logger"))

	engine.GET("/", func(c *gin.Context) {
		ctx := c.Request.Context()
		log.Infof(ctx, "info hogehogehoge")
		log.Warningf(ctx, "warn piyopiyo")
		log.Errorf(ctx, "err hogehogehoge")
		log.DebugObj(ctx, "hoge", LogEntry{
			Name: "Akahane",
			Age:  99,
		})
		log.ErrorObj(ctx, "hoge", LogEntry{
			Name: "Hogehoge",
			Age:  40,
		})
		log.ErrorObj(ctx, "hoge", LogEntry{
			Name: "Fugafuga",
			Age:  80,
		})

		log.Errorf(ctx, "err fugafugafuga")
		c.JSON(http.StatusOK, gin.H{
			"message": "hello world",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)
	engine.Run(addr)
}
