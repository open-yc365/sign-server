package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prediction-platform/sign-server/utils"
)

func main() {
	_ = godotenv.Load()

	// 读取 config.json
	type Config struct {
		RootPath   string   `json:"rootpath"`
		Port       int      `json:"port"`
		AllowedIPs []string `json:"allowed_ips"`
	}
	var config Config
	configData, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Failed to read config.json: %v", err)
	}
	if err := json.Unmarshal(configData, &config); err != nil {
		log.Fatalf("Failed to parse config.json: %v", err)
	}

	r := gin.Default()
	//r.Use(IPWhitelistMiddleware(config.AllowedIPs))
	root := r.Group(config.RootPath)

	// 生成地址接口
	root.GET("/address", func(c *gin.Context) {
		indexStr := c.Query("index")
		if indexStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "index is required"})
			return
		}
		index, err := strconv.ParseInt(indexStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "index must be an integer"})
			return
		}
		mnemonic := os.Getenv("MNEMONIC")
		_, address, _, err := utils.CreateAddress(mnemonic, index)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"address": address, "index": index})
	})

	// 签名接口
	root.POST("/sign", func(c *gin.Context) {
		var req struct {
			Index      int64  `json:"index"`
			DigestHash string `json:"digestHash"` // 16进制字符串
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.Index < 1 || len(req.DigestHash) < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "params is error"})
			return
		}
		cleanHex := strings.TrimPrefix(req.DigestHash, "0x")
		hashBytes, err := hex.DecodeString(cleanHex)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hash"})
			return
		}
		if len(hashBytes) != 32 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hash length, expected 32 bytes"})
			return
		}
		mnemonic := os.Getenv("MNEMONIC")
		signature, err := utils.SignTransaction(mnemonic, req.Index, hashBytes)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"signature": hexutil.Encode(signature)})
	})

	// 优雅关闭
	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(config.Port),
		Handler: r,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("Server forced to shutdown: %v", err)
		}
		log.Println("Server exiting")
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}

// IP 白名单中间件
func IPWhitelistMiddleware(allowedIPs []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		for _, ip := range allowedIPs {
			if clientIP == ip {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
	}
}
