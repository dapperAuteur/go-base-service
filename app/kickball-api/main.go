package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ardanlabs/conf"
	"github.com/pkg/errors"
)

// build is the git version of this program. It is set using build flags in the makefile.
var build = "develop"

func main() {
	log := log.New(os.Stdout, "KICKBALL : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	if err := run(log); err != nil {
		log.Println("main: error: ", err)
		os.Exit(1)
	}
}

func run(log *log.Logger) error {

	// ==
	// Configuration

	var cfg struct {
		conf.Version
		Web struct {
			Address         string        `conf:"default:0.0.0.0:3000"`
			Debug           string        `conf:"default:0.0.0.0:4000"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:5s"`
			ShutdownTimeout time.Duration `conf:"default:5s"`
		}
	}

	cfg.Version.SVN = build
	cfg.Version.Desc = "copyright information here"

	if err := conf.Parse(os.Args[1:], "KICKBALL", &cfg); err != nil {
		switch err {
		case conf.ErrHelpWanted:
			usage, err := conf.Usage("KICKBALL", &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return nil
		case conf.ErrVersionWanted:
			version, err := conf.VersionString("KICKBALL", &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config version")
			}
			fmt.Println(version)
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}

	return nil
}
