package storage

import (
	cqmodels "gitlab.qarea.org/jiraquality/cq-web-backend/common-lib/models"
	cqmongodb "gitlab.qarea.org/jiraquality/cq-web-backend/common-lib/mongodb"
	dblib "gitlab.qarea.org/jiraquality/cq-web-backend/db-lib"
	"gopkg.in/mgo.v2"
)

// Cfg config
type Cfg struct {
	Session      *mgo.Session
	DatabaseName string
}

//MongoDB structure implementation
type MongoDB struct {
	cqMongo *cqmongodb.MongoDB
	DBLib   *dblib.DBLib
}

//New instance entity
func New(conf *Cfg) *MongoDB {
	m := &MongoDB{
		cqMongo: cqmongodb.New(&cqmongodb.Cfg{
			Session:      conf.Session,
			DatabaseName: conf.DatabaseName}),
		DBLib: dblib.New(&dblib.Cfg{
			Session:      conf.Session.Copy(),
			DatabaseName: conf.DatabaseName,
		}),
	}

	return m
}

//SaveCommits -
func (db MongoDB) SaveCommits(commits []interface{}) error {
	return db.cqMongo.InsertRecords(commits, cqmodels.CQCommitsCollection)
}

//SaveCommiters -
func (db MongoDB) SaveCommiters(commiters []interface{}) error {
	return db.cqMongo.InsertRecords(commiters, cqmodels.CQCommittersCollection)
}
