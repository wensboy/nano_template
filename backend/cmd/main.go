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

// @BasePath /api/v1
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
