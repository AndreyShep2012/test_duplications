package main

import (
	"log"
	"path/filepath"

	"github.com/BurntSushi/toml"

	"gitlab.qarea.org/jiraquality/cq-web-backend/git-scanner-service/models"
	"gitlab.qarea.org/jiraquality/cq-web-backend/git-scanner-service/scanner"
	"gitlab.qarea.org/jiraquality/cq-web-backend/git-scanner-service/storage"

	"gitlab.qarea.org/jiraquality/cq-web-backend/common-lib/clients"
	"gitlab.qarea.org/jiraquality/cq-web-backend/common-lib/helpers"
	"gitlab.qarea.org/jiraquality/cq-web-backend/common-lib/logger"
)

var conf models.Config

func init() {
	exPath, err := helpers.GetExecutablePath()
	if err != nil {
		log.Fatal("can't detect executable path:", err)
	}

	confPath := filepath.Join(exPath, "config", "config.toml")
	if _, err := toml.DecodeFile(confPath, &conf); err != nil {
		log.Fatal("can't process configuration for application: ", err)
	}

	commonConfPath := filepath.Join(exPath, "../", "common-config", "config.toml")
	if _, err := toml.DecodeFile(commonConfPath, &conf.CommonConfig); err != nil {
		log.Fatal("can't process common configuration for application: ", err)
	}

	conf.ExecPath = exPath
}

func main() {
	initLogger()

	sc := initScanner(initMongoDB())
	sc.Start()
	select {}
}

func initLogger() {
	logger.RedirectLog(&logger.Conf{
		Filename:   "log/log.log",
		MaxSize:    conf.CommonConfig.LogMaxSize,
		MaxBackups: conf.CommonConfig.LogMaxBackups,
		MaxAge:     conf.CommonConfig.LogMaxAge,
	}, nil)
}

func initScanner(storage *storage.MongoDB) *scanner.Scanner {
	return scanner.New(
		scanner.Cfg{
			Storage: storage,
			Config:  conf,
		},
	)
}

func initMongoDB() *storage.MongoDB {
	log.Println("open mongo database...")

	mgoSession := clients.GetMongoSession(&clients.MongoConf{
		Host: conf.CommonConfig.DB.Host,
		Auth: conf.CommonConfig.DB.Auth,
		User: conf.CommonConfig.DB.User,
		Pass: conf.CommonConfig.DB.Pass,
	})

	cfg := &storage.Cfg{
		Session:      mgoSession,
		DatabaseName: conf.CommonConfig.DB.DatabaseName,
	}
	ds := storage.New(cfg)
	return ds
}
