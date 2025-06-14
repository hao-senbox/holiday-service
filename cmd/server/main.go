package main

import (
	"context"
	"log"
	"os"
	"time"
	"worktime-service/config"
	"worktime-service/internal/attendance"
	"worktime-service/internal/leave"
	"worktime-service/internal/user"
	"worktime-service/pkg/consul"
	"worktime-service/pkg/zap"

	// cronV3 "github.com/robfig/cron/v3"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {

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

	mongoClient, err := connectToMongoDB(cfg.MongoURI)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()

	consulConn := consul.NewConsulConn(logger, cfg)
	consulClient := consulConn.Connect()
	defer consulConn.Deregister()

	settingCollection := mongoClient.Database(cfg.MongoDB).Collection("settings")
	leaveRequestCollection := mongoClient.Database(cfg.MongoDB).Collection("leave_requests")
	dailyLeaveSlotsCollection := mongoClient.Database(cfg.MongoDB).Collection("daily_leave_slots")
	leaveBalanceCollection := mongoClient.Database(cfg.MongoDB).Collection("leave_balance")
	attendanceCollection := mongoClient.Database(cfg.MongoDB).Collection("attendance_logs")
	attendanceDailyCollection := mongoClient.Database(cfg.MongoDB).Collection("attendances_daily")
	
	userService := user.NewUserService(consulClient)
	attendanceRepository := attendance.NewAttendanceRepository(attendanceCollection, attendanceDailyCollection)
	attendanceService := attendance.NewAttendanceService(attendanceRepository, userService)
	attendanceHandler := attendance.NewAttendanceHandler(attendanceService)

	leaveRepository := leave.NewLeaveRepository(leaveRequestCollection, settingCollection, dailyLeaveSlotsCollection, leaveBalanceCollection)
	leaveService := leave.NewLeaveService(leaveRepository, userService)
	leaveHandler := leave.NewLeaveHandler(leaveService)

	r := gin.Default()

	// r.LoadHTMLGlob("web/*")

	// r.GET("/", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "index.html", gin.H{})
	// })

	leave.RegisterRoutes(r, leaveHandler)
	attendance.RegisterRoutes(r, attendanceHandler)

	// c := cronV3.New(cronV3.WithSeconds(), cronV3.WithLogger(cronV3.DefaultLogger))

	// _, err = c.AddFunc("*/30 * * * * *", func() {

	// 	log.Println("Bắt đầu chạy job...")
	// 	ctx := context.Background()

	// 	if err := leaveService.AddCronLeavesBalance(ctx); err != nil {
	// 		log.Printf("Lỗi khi thêm ngày nghỉ cho user %v", err)
	// 	}

	// })

	// if err != nil {
	// 	log.Printf("Lỗi khi thiết lập cron job: %v", err)
	// } else {
	// 	c.Start()
	// 	log.Println("Cron job đã được khởi động")
	// }

	port := os.Getenv("PORT")
	if port == "" {
		port = "8008"
	}

	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)

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
		log.Println("Failed to ping to MongoDB")
		return nil, err
	}

	log.Println("Successfully connected to MongoDB")
	return client, nil
}


