package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type store struct {
	Database  database  `mapstructure:"db"`
	Metrics   metrics   `mapstructure:"metrics"`
	API       api       `mapstructure:"api"`
	Addresses addresses `mapstructure:"addresses"`
}

var Store store

func Load() {
	err := viper.Unmarshal(&Store)
	if err != nil {
		logrus.Fatal(err)
	}
}
