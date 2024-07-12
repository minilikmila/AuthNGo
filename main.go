package main

import (
	"embed"
	"flag"
	"os"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/minilikmila/goAuth/cmd"
	log "github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
	// // Only log the warning severity or above.
	// log.SetLevel(log.WarnLevel)
	log.SetFormatter(&log.JSONFormatter{})

	log.SetOutput(os.Stdout)
	// log.SetLevel(config.LogLevel)

}

var (
	DevMode = flag.Bool("dev", false, "development mood")
)

//go:embed migrations/*.sql
var files embed.FS

func main() {
	// parse cli command flag
	flag.Parse()

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	if *DevMode {
		log.SetFormatter(&log.JSONFormatter{})
	}

	cmd.Run(files)
}
