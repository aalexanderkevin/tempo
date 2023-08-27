package main

import (
	"context"
	"os"
	_ "time/tzdata"

	"tempo/config"
	"tempo/container"
	"tempo/repository/mysqlrepo"
	"tempo/storage"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var rootCmd = &cobra.Command{
	Use: "tempo",
}

func init() {
	loadConfig()
}

func main() {
	Execute()
}

func Execute() {
	rootCmd := registerCommands(&defaultAppProvider{})
	if err := rootCmd.Execute(); err != nil {
		logrus.Error(err.Error())
		os.Exit(1)
	}
}

func loadConfig() {
	err := config.Load()
	if err != nil {
		logrus.Errorf("Config error: %s", err.Error())
		os.Exit(1)
	}
}

func registerCommands(appProvider AppProvider) *cobra.Command {
	rootCmd.AddCommand(Server(appProvider))
	rootCmd.AddCommand(Migrate(appProvider))

	return rootCmd
}

type AppProvider interface {
	BuildContainer(ctx context.Context, options buildOptions) (*container.Container, func(), error)
}

type buildOptions struct {
	MySql bool
}

type defaultAppProvider struct {
}

func (defaultAppProvider) BuildContainer(ctx context.Context, options buildOptions) (*container.Container, func(), error) {
	var db *gorm.DB
	cfg := config.Instance()

	appContainer := container.NewContainer()
	appContainer.SetConfig(cfg)

	if options.MySql {
		db = storage.GetMySqlDb()
		appContainer.SetDb(db)

		userRepo := mysqlrepo.NewUserRepository(db)
		appContainer.SetUserRepo(userRepo)
	}

	deferFn := func() {
		if db != nil {
			storage.CloseDB(db)
		}
	}

	return appContainer, deferFn, nil
}
