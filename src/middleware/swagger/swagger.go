package swagger

import (
	"backend-crowdfunding/config"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
)

func New(cfg config.Config) fiber.Handler {
	return swagger.New(initConfig(cfg))
}

func initConfig(cfg config.Config) swagger.Config {
	return swagger.Config{
		BasePath: "/", //swagger ui base path
		FilePath: cfg.GetWithDefault("SWAGGER_FILEPATH", "./docs/swagger.json"),
	}
}
