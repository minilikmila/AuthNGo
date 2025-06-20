package main

import (
	"embed"
	"flag"
	"os"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/minilikmila/standard-auth-go/cmd"
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

	log.SetOutput(os.Stdout)
	// log.SetLevel(config.LogLevel)

}

var (
	// DevMode = flag.Bool("dev", false, "development mood")
	DevMode = flag.String("mode", "dev", "development mood")
)

/*// go:embed internal/infrastructure/migrations/*.sql*/

var files embed.FS

func main() {
	// parse cli command flag
	flag.Parse()

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	// if *DevMode {
	// 	log.SetFormatter(&log.JSONFormatter{})
	// }
	log.Println("DevMode: ", *DevMode)

	args := map[string]string{
		"mode": *DevMode,
	}
	if *DevMode == "dev" {
		args["mode"] = "debug"
	}

	cmd.Run(files, args)
}
