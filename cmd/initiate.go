package cmd

import (
	"embed"
	"fmt"
	"net/http"
	"os"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/minilikmila/standard-auth-go/config"
	database_ "github.com/minilikmila/standard-auth-go/internal/database"
	"github.com/minilikmila/standard-auth-go/internal/routes"

	"github.com/sirupsen/logrus"
)

// var logger = log.New()

func Run(files embed.FS, args map[string]string) {

	mode, ok := args["mode"]
	if !ok {
		mode = "debug"
	}

	config, err := config.New("config.json")
	if err != nil {
		fmt.Println("db conn error ", err)
		os.Exit(1)
	}
	// Initialize Database connection
	db, err := database_.InitDatabase(config)
	if err != nil {
		fmt.Println("db conn error ", err)
		os.Exit(1)
	}

	routes := routes.InitRoute(db, config, mode)

	host := fmt.Sprintf("%s:%v", config.Host, config.Port)

	logrus.WithField("host", "http://"+host).Info("started Go authentication server")

	logrus.Fatalln(http.ListenAndServe(fmt.Sprintf("%s:%v", config.Host, config.Port), routes))
	// routes.Run(":"+config.Config.App.Port)

}
