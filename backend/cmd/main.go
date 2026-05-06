package main

import (
	"flag"

	"example.com/nano_template/pkg/config"
	"example.com/nano_template/pkg/util"
	"github.com/joho/godotenv"

	docs "example.com/nano_template/cmd/docs"
)

func init() {
	docs.SwaggerInfo.BasePath = "/api/v1"
}

// @title           Nano Template Api
// @version         1.0
// @description     swagger for nano template
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:3000
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description jwt format： Bearer {token}

func main() {

	// Load environment variables from .env in local/dev; ignore if file does not exist.
	_ = godotenv.Load("./.env")

	util.InitLogger(false, "")

	cfg, err := config.LoadConfig("./env.yaml")
	if err != nil {
		panic(err)
	}
	flag.Parse()
	srv := NewServer(&cfg.ServerConfig)
	srv.Start(cfg)
}
