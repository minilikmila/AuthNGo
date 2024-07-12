package cmd

import (
	"embed"
	"fmt"
	"net/http"
	"os"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/minilikmila/goAuth/config"
	"github.com/minilikmila/goAuth/db"
	"github.com/minilikmila/goAuth/routes"
	"github.com/sirupsen/logrus"
)

// var logger = log.New()

func Run(files embed.FS) {

	// Initialize Database connection
	db, err := db.InitDatabase(files)
	if err != nil {
		fmt.Println("db conn error ", err)
		os.Exit(1)
	}

	routes := routes.InitRoute(db)

	host := fmt.Sprintf("%s:%s", config.Config.App.Host, config.Config.App.Port)

	logrus.WithField("host", "http://"+host).Info("started goAuth server")

	logrus.Fatalln(http.ListenAndServe(fmt.Sprintf("%s:%s", config.Config.App.Host, config.Config.App.Port), routes))
	// routes.Run(":"+config.Config.App.Port)

}
