package main

import (
	"fmt"

	"github.com/alpha-supsys/go-common/app/web"
	"github.com/alpha-supsys/go-common/config/env"
	"github.com/alpha-supsys/go-common/log"
	"github.com/alpha-supsys/laya-demo-game-backend/src/controller"
)

func main() {
	if cfg, err := env.LoadAllWithoutPrefix("LDG_"); err == nil {
		logger := log.NewCommon(log.Debug)

		app := web.New(cfg, logger)

		app.HandleController(controller.NewWebsocketController())

		err = app.Run(log.Debug)
		if err != nil {
			fmt.Println(err)
		}
	}
}
