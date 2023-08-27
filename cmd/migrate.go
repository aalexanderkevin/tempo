package main

import (
	"context"
	"fmt"
	"os"

	"tempo/config"
	"tempo/storage"

	"github.com/spf13/cobra"
)

var (
	migrationPath  string
	rollback       bool
	versionToForce int
)

func Migrate(appProvider AppProvider) *cobra.Command {
	cliCommand := &cobra.Command{
		Use:   "migrate",
		Short: "Run the database migration",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			app, closeResourcesFn, err := appProvider.BuildContainer(ctx, buildOptions{
				MySql: true,
			})
			if err != nil {
				return err
			}
			if closeResourcesFn != nil {
				defer closeResourcesFn()
			}

			dB := app.Db()
			err = storage.MigrateMysqlDb(dB, &migrationPath, rollback, versionToForce)
			if err != nil {
				fmt.Printf("Error when migration: %s \n", err.Error())
				os.Exit(1)
			}

			fmt.Println("Finish migrating database")
			return nil
		},
	}

	cfg := config.Instance()
	cliCommand.Flags().StringVarP(&migrationPath, "path", "p", cfg.DB.Migrations.Path, "The migration folder")
	cliCommand.Flags().BoolVarP(&rollback, "rollback", "r", false, "Rollback to prev migration (-1 step)")
	cliCommand.Flags().IntVarP(&versionToForce, "force", "f", -1, "Force to specific version")
	return cliCommand
}
