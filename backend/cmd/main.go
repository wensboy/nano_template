package main

import (
	"context"
	"os"

	"example.com/nano_template/pkg/util"

	_ "example.com/nano_template/cmd/docs"
	"github.com/urfave/cli/v3"
)

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

var (
	APP_NAME        = "nanot"
	APP_VERSION     = "v0.1.0"
	APP_LAST_MODIFY = "unknown"
)

func main() {

	cmd := &cli.Command{
		Name:    APP_NAME,
		Version: APP_VERSION,
	}

	MountCommands(cmd)

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		util.Error(err.Error())
	}

}
