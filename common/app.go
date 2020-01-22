package common

import (
	"context"
	"github.com/urfave/cli"
	"log"
	"os"
	"os/signal"
	"time"
)

func App(serviceName string, serviceUsage string, httpServer HttpServer, serviceConfig ServiceConfig) (app *cli.App,
	err error) {
	var config string
	var mode = "release"
	app = &cli.App{
		Name:  serviceName,
		Usage: serviceUsage,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config, c",
				Usage:       "Load config from `FILE`",
				Destination: &config,
			},
			&cli.StringFlag{
				Name:        "mode, m",
				Usage:       "run mode, can be debug or release",
				Destination: &mode,
			},
		},
		Commands: []*cli.Command{
			{
				Name:         "httpserver",
				Aliases:      nil,
				Usage:        "server over http",
				UsageText:    "",
				Description:  "",
				ArgsUsage:    "",
				Category:     "",
				BashComplete: nil,
				Before:       nil,
				After:        nil,
				Action: func(c *cli.Context) error {
					log.Println("Start Server [default : release]:", mode)
					httpReal := httpServer.Init(serviceConfig, config)
					go httpReal.Serve(mode)
					// Wait for interrupt signal to gracefully shutdown the server with
					// a timeout of 5 seconds.
					quit := make(chan os.Signal)
					signal.Notify(quit, os.Interrupt, os.Kill)
					<-quit
					httpReal.Shutdown()
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					defer cancel()
					ctx.Done()
					log.Println("Server Shutdown OK")
					return nil
				},
				OnUsageError:       nil,
				Subcommands:        nil,
				Flags:              nil,
				SkipFlagParsing:    false,
				HideHelp:           false,
				Hidden:             false,
				HelpName:           "",
				CustomHelpTemplate: "",
			},
			{
				Name:         "worker-server",
				Aliases:      nil,
				Usage:        "server of worker",
				UsageText:    "",
				Description:  "",
				ArgsUsage:    "",
				Category:     "",
				BashComplete: nil,
				Before:       nil,
				After:        nil,
				Action: func(c *cli.Context) error {
					log.Println("Start Server [default : release]:", mode)
					httpReal := httpServer.Init(serviceConfig, config)
					go httpReal.ServeWorker(mode)

					// Wait for interrupt signal to gracefully shutdown the server with
					// a timeout of 5 seconds.
					quit := make(chan os.Signal)
					signal.Notify(quit, os.Interrupt)
					<-quit
					log.Println("Shutdown Server")
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					defer cancel()
					ctx.Done()
					log.Println("Server Worker Exit OK")
					return nil
				},
				OnUsageError:       nil,
				Subcommands:        nil,
				Flags:              nil,
				SkipFlagParsing:    false,
				HideHelp:           false,
				Hidden:             false,
				HelpName:           "",
				CustomHelpTemplate: "",
			},
			{
				Name:         "swagger-api",
				Aliases:      nil,
				Usage:        "init swagger apis",
				UsageText:    "",
				Description:  "",
				ArgsUsage:    "",
				Category:     "",
				BashComplete: nil,
				Before:       nil,
				After:        nil,
				Action: func(c *cli.Context) error {
					print("Start swagger api create or init")
					return nil
				},
				OnUsageError:       nil,
				Subcommands:        nil,
				Flags:              nil,
				SkipFlagParsing:    false,
				HideHelp:           false,
				Hidden:             false,
				HelpName:           "",
				CustomHelpTemplate: "",
			},
		},
	}
	return

}
