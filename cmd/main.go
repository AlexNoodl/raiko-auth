package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github/alexnoodl/raiko-auth/docs"
	"github/alexnoodl/raiko-auth/internal/config"
	"github/alexnoodl/raiko-auth/internal/handler"
	"github/alexnoodl/raiko-auth/internal/services"
	"github/alexnoodl/raiko-auth/pkg/database"
	pb "github/alexnoodl/raiko-auth/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

// @title API Авторизации
// @version 1.0
// @description Сервер авторизации.
// @host localhost:8080
// @BasePath /api/v1
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	if cfg.Logger == nil {
		log.Fatal("Logger is nil after config initialization")
	}

	db, err := database.InitMongoDB(cfg)
	if err != nil {
		cfg.Logger.Fatal("Failed to connect to MongoDB: ", err)
	}
	cfg.Logger.Info("Connected to MongoDB successfully")

	router := gin.Default()

	router.Use(gin.Recovery())

	authService := services.NewAuthService(db, cfg.Logger, []byte(cfg.JWTKey))
	authHandler := handler.NewAuthHandler(authService, cfg.Logger)

	{
		v1 := router.Group("/api/v1")
		v1.POST("/register", authHandler.Register)
		v1.POST("/login", authHandler.Login)
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	go func() {
		cfg.Logger.Info("Server starting on port: ", cfg.Port)
		if err := router.Run(":" + cfg.Port); err != nil {
			cfg.Logger.Fatal("HTTP server failed: ", err)
		}
	}()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		cfg.Logger.Fatal("Failed to listen for gRPC: ", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, services.NewAuthGrpcServer(authService, cfg.Logger))
	reflection.Register(grpcServer)

	cfg.Logger.Info("Starting gRPC server on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		cfg.Logger.Fatal("gRPC server failed: ", err)
	}
}
