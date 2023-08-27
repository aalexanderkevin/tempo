package main

import (
	"context"

	"tempo/controller"
	"tempo/helper"

	"github.com/segmentio/ksuid"
	"github.com/spf13/cobra"
)

func Server(appProvider AppProvider) *cobra.Command {
	cliCommand := &cobra.Command{
		Use:   "server",
		Short: "Run the REST API server",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := helper.ContextWithRequestId(context.Background(), ksuid.New().String())
			logger := helper.GetLogger(ctx).WithField("method", "server")

			app, closeResourcesFn, err := appProvider.BuildContainer(ctx, buildOptions{
				MySql: true,
			})
			if err != nil {
				panic(err)
			}
			if closeResourcesFn != nil {
				defer closeResourcesFn()
			}

			// Start Http Server
			err = controller.NewHttpServer(app).Start()
			if err != nil {
				logger.WithError(err).Error("Error starting web server")
				return err
			}

			return nil
		},
	}
	return cliCommand
}
