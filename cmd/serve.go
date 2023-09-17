package cmd

import (
	"geo-service/config"
	"geo-service/controller"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/zeynab-sb/geoolocation"
	"github.com/zeynab-sb/geoolocation/database"
)

var serveCMD = &cobra.Command{
	Use:   "serve",
	Short: "serve API",
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

func serve() {
	cfg := &database.DBConfig{
		Driver:      config.C.Database.Driver,
		Host:        config.C.Database.Host,
		Port:        config.C.Database.Port,
		DB:          config.C.Database.DB,
		User:        config.C.Database.User,
		Password:    config.C.Database.Password,
		Location:    config.C.Database.Location,
		MaxConn:     config.C.Database.MaxConn,
		IdleConn:    config.C.Database.IdleConn,
		Timeout:     config.C.Database.Timeout,
		DialRetry:   config.C.Database.DialRetry,
		DialTimeout: config.C.Database.DialTimeout,
	}

	geo, err := geoolocation.New(cfg)
	if err != nil {
		log.Errorf("Cannot make geo location instance: %s", err)
	}

	e := echo.New()

	locationController := controller.Location{Geo: geo}

	e.GET("/locations/:ip", locationController.GetByIP)

	// Start server
	e.Logger.Fatal(e.Start(config.C.Address))
}
