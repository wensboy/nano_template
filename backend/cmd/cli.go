package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"example.com/nano_template/pkg/config"
	"example.com/nano_template/pkg/rpc"
	"example.com/nano_template/pkg/util"
	"example.com/nano_template/pkg/web"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v3"
	"go.uber.org/zap"
)

const (
	VERSION_TEMPLATE = "Nano Template Cli [current: %s, last-modify: %s]\n"
)

func MountCommands(cmd *cli.Command) {
	// Load environment variables from .env in local/dev; ignore if file does not exist.
	_ = godotenv.Load("./.env")

	util.InitLogger(false, "")

	cfg, err := config.LoadConfig("./env.yaml")
	if err != nil {
		util.Error(err.Error())
		os.Exit(1)
	}

	cli.VersionPrinter = showVersion

	cmd.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "server",
			Value:       "web",
			Destination: &cfg.FlagConfig.Server,
			Usage:       `start server with web/rpc`,
		},
	}

	cmd.Action = Excute(cfg)
}

func Excute(cfg *config.Config) cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		var server Server
		switch cmd.String("server") {
		case "web":
			server = web.NewServer()
		case "rpc":
			server = rpc.NewRpcServer()
		default:
			util.Warn("[cli]", zap.String("server-start", "failed"))
			return errors.New("invalid server option")
		}
		runner := NewServerRunner(
			WithConfig(cfg),
			WithServer(server),
		)
		runner.Run()
		return nil
	}
}

func showVersion(cmd *cli.Command) {
	fmt.Printf(VERSION_TEMPLATE, cmd.Root().Version, time.Now().Format("2006-01-02 15:04:05"))
}
