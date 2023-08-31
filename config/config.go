package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config interface {
	Get(key string) string
	GetWithDefault(key, defaultVal string) string
	GetOrPanic(key string) string
}

type configImpl struct {
}

func (config *configImpl) GetOrPanic(key string) string {
	var str string

	if str = os.Getenv(key); str == "" {
		log.Fatalf("Given key '%s' doesn't exist", key)
	}

	return str
}

func (config *configImpl) Get(key string) string {
	return os.Getenv(key)
}

func (config *configImpl) GetWithDefault(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}

func (config *configImpl) GetInt(key string, defaultVal int) int {
	var valueString string
	var value int
	if valueString := os.Getenv(key); valueString == "" {
		return defaultVal
	}

	value, err := strconv.Atoi(valueString)
	if err != nil {
		return defaultVal
	}

	return value
}

func New(filenames ...string) Config {
	err := godotenv.Load(filenames...)
	if err != nil {
		panic(err)
	}
	return &configImpl{}
}
