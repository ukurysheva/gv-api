package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// "golang/go/pkg/mod/github.com/sirupsen/logrus@v1.8.1"

// _ "github.com/lib/pq"

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := InitConfig(); err != nil {
		logrus.Fatalf("error initializaing configs: %s", err.Error())
		fmt.Println(err)
	}
	return
}

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
