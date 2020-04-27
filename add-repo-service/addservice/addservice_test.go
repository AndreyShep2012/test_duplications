package addservice_test

import (
	"log"
	"testing"
	"time"

	"gopkg.in/mgo.v2/bson"

	"gitlab.qarea.org/jiraquality/cq-web-backend/add-repo-service/addservice"
	"gitlab.qarea.org/jiraquality/cq-web-backend/add-repo-service/models"
	"gitlab.qarea.org/jiraquality/cq-web-backend/add-repo-service/storage"
	"gitlab.qarea.org/jiraquality/cq-web-backend/common-lib/clients"
)

func TestAddRepo(t *testing.T) {
	s := initMongoDB()
	c := models.Config{}

	serv := addservice.New(addservice.Cfg{
		Storage: s,
		Config:  c,
	})

	u := bson.NewObjectId()
	p := bson.NewObjectId()
	_, err := serv.AddRepo(u.Hex(), p.Hex(), "https://github.com/roman-vrm/qatestlab.git", "master")
	if err != nil {
		t.Fatal(err)
	}

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
