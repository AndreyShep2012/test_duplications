package scanner

import (
	"log"
	"testing"
	"time"

	"gitlab.qarea.org/jiraquality/cq-web-backend/common-lib/clients"
	"gitlab.qarea.org/jiraquality/cq-web-backend/git-scanner-service/models"
	"gitlab.qarea.org/jiraquality/cq-web-backend/git-scanner-service/storage"
)

func TestFirstProcessing(t *testing.T) {
	s := initMongoDB()
	c := models.Config{}

	c.CommonConfig.SourcePath = "/Users/andy/work/test/testcqrepos"
	serv := New(Cfg{
		Storage: s,
		Config:  c,
	})
	serv.Start()
	time.Sleep(time.Second * 25)
}

func initMongoDB() *storage.MongoDB {
	log.Println("open mongo database...")

	mgoSession := clients.GetMongoSession(&clients.MongoConf{
		Host: "127.0.0.1",
		Auth: "admin",
	})

	cfg := &storage.Cfg{
		Session:      mgoSession,
		DatabaseName: "cq_data",
	}
	ds := storage.New(cfg)
	return ds
}
