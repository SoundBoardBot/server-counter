package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/SoundBoardBot/server-counter/config"
	"github.com/SoundBoardBot/server-counter/db"
	"github.com/SoundBoardBot/server-counter/http"
	"github.com/SoundBoardBot/server-counter/tasks"
	"github.com/SoundBoardBot/server-counter/utils"

	"github.com/go-co-op/gocron/v2"
)

func main() {
	config.Parse()

	lerr := utils.Configure(nil, config.Conf.JsonLogs, config.Conf.LogLevel)
	if lerr != nil {
		panic(fmt.Errorf("failed to create zap logger: %w", lerr))
	}
	db.Init()

	s, s_err := gocron.NewScheduler()
	if s_err != nil {
		panic(fmt.Errorf("failed to create cron scheduler: %w", s_err))
	}

	utils.Logger.Sugar().Infof("top.gg Enabled: %s", IsEnabled(config.Conf.Auth.TopGG))
	utils.Logger.Sugar().Infof("discordbotlist.com Enabled: %s", IsEnabled(config.Conf.Auth.DiscordBotList))
	utils.Logger.Sugar().Infof("botlist.me Enabled: %s", IsEnabled(config.Conf.Auth.BotListMe))
	utils.Logger.Sugar().Infof("voidbots.net Enabled: %s", IsEnabled(config.Conf.Auth.VoidBots))

	if config.Conf.OneShot {
		utils.Logger.Info("Running Job - Bot Stats")
		tasks.UpdateBotStats()
		return
	}

	// add hourly fresh bot data
	s.NewJob(gocron.CronJob("*/5 * * * *", false), gocron.NewTask(func() {
		utils.Logger.Debug("Running Job - Bot Stats")
		tasks.UpdateBotStats()
	}))

	// start the scheduler
	utils.Logger.Info("Starting cron jobs")
	s.Start()

	// Start HTTP server...
	if config.Conf.Http.Enabled {
		utils.Logger.Info("Starting HTTP server...")
		server := http.NewServer(utils.Logger, config.Conf)
		server.RegisterRoutes()
		server.Start()
	}

	// keep alive until shutdown signal
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)
	<-shutdownCh

	utils.Logger.Info("Received shutdown signal")
}

func IsEnabled(value string) string {
	if value == "" {
		return "No"
	}
	return "Yes"
}
