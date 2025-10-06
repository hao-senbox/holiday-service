package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"worktime-service/config"
	"worktime-service/internal/attendance"
	"worktime-service/internal/leave"
	"worktime-service/internal/user"
	"worktime-service/pkg/consul"
	"worktime-service/pkg/zap"

	"github.com/gin-gonic/gin"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	// Load env
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Printf("Warning: Error loading .env file: %v", err)
		} else {
			log.Println("Successfully loaded .env file")
		}
	} else {
		log.Println("No .env file found, using environment variables")
	}

	cfg := config.LoadConfig()

	logger, err := zap.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	consulConn := consul.NewConsulConn(logger, cfg)
	consulClient := consulConn.Connect()
	defer consulConn.Deregister()
	
	mongoClient, err := connectToMongoDB(cfg.MongoURI)
	if err != nil {
		panic(err)
	}

	if err := waitPassing(consulClient, "go-main-service", 60*time.Second); err != nil {
		logger.Fatalf("Dependency not ready: %v", err)
	}

	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()

	// Handle OS signal để deregister
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Shutting down server... De-registering from Consul...")
		consulConn.Deregister()
		os.Exit(0)
	}()

	settingCollection := mongoClient.Database(cfg.MongoDB).Collection("settings")
	leaveRequestCollection := mongoClient.Database(cfg.MongoDB).Collection("leave_requests")
	dailyLeaveSlotsCollection := mongoClient.Database(cfg.MongoDB).Collection("daily_leave_slots")
	leaveBalanceCollection := mongoClient.Database(cfg.MongoDB).Collection("leave_balance")
	attendanceCollection := mongoClient.Database(cfg.MongoDB).Collection("attendance_logs")
	attendanceDailyCollection := mongoClient.Database(cfg.MongoDB).Collection("attendances_daily")
	attendanceDailyStudentCollection := mongoClient.Database(cfg.MongoDB).Collection("attendances_daily_students")
	userService := user.NewUserService(consulClient)
	attendanceRepository := attendance.NewAttendanceRepository(attendanceCollection, attendanceDailyCollection, attendanceDailyStudentCollection)
	attendanceService := attendance.NewAttendanceService(attendanceRepository, userService)
	attendanceHandler := attendance.NewAttendanceHandler(attendanceService)

	leaveRepository := leave.NewLeaveRepository(leaveRequestCollection, settingCollection, dailyLeaveSlotsCollection, leaveBalanceCollection)
	leaveService := leave.NewLeaveService(leaveRepository, userService)
	leaveHandler := leave.NewLeaveHandler(leaveService)

	r := gin.Default()

	leave.RegisterRoutes(r, leaveHandler)
	attendance.RegisterRoutes(r, attendanceHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8008"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server stopped with error: %v", err)
	}
}

func connectToMongoDB(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Println("Failed to connect to MongoDB")
		return nil, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Println("Failed to ping MongoDB")
		return nil, err
	}

	log.Println("Successfully connected to MongoDB")
	return client, nil
}

func waitPassing(cli *consulapi.Client, name string, timeout time.Duration) error {
	dl := time.Now().Add(timeout)
	for time.Now().Before(dl) {
		entries, _, err := cli.Health().Service(name, "", true, nil)
		if err == nil && len(entries) > 0 {
			return nil // đã sẵn sàng
		}
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("%s not ready in consul", name)
}
