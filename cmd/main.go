package main

import (
	"fmt"
	"os"

	"github.com/go-redis/redis/v7"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	gvapi "github.com/ukurysheva/gv-api"
	"github.com/ukurysheva/gv-api/pkg/handler"
	"github.com/ukurysheva/gv-api/pkg/repository"
	"github.com/ukurysheva/gv-api/pkg/service"
)

// "golang/go/pkg/mod/github.com/sirupsen/logrus@v1.8.1"

// _ "github.com/lib/pq"

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := InitConfig(); err != nil {
		logrus.Fatalf("error initializaing configs: %s", err.Error())
		fmt.Println(err)
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}
	fmt.Println(viper.GetString("db.port"))
	// db, err := repository.NewPostgresDB(repository.Config{})
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	redis, err := initReddis()
	if err != nil {
		logrus.Fatalf("failed to initialize reddis: %s", err.Error())
	}
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services, redis)

	srv := new(gvapi.Server)
	// go func() {
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error occured while running http server: %s", err.Error())
	}
	// }()

	logrus.Print("TodoApp Started")

	return
}

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func initReddis() (*redis.Client, error) {
	var client *redis.Client
	//Initializing redis
	// dsn := os.Getenv("REDIS_PORT")
	// if len(dsn) == 0 {
	// 	dsn = "localhost:6379"
	// }

	client = redis.NewClient(&redis.Options{
		Addr: "redis:6379", //redis port
	})
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
