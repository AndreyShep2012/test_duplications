package storage

import (
	cqerrors "gitlab.qarea.org/jiraquality/cq-web-backend/common-lib/errors"
	cqmodels "gitlab.qarea.org/jiraquality/cq-web-backend/common-lib/models"
	cqmongodb "gitlab.qarea.org/jiraquality/cq-web-backend/common-lib/mongodb"
	dblib "gitlab.qarea.org/jiraquality/cq-web-backend/db-lib"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

//InsertSettingRecordDB -
func (db MongoDB) InsertSettingRecordDB(so cqmodels.SettingsObject) error {
	_, err := db.DBLib.GetSettingsExists(so)
	if err == nil {
		return cqerrors.ErrSettingExists
	}

	return db.cqMongo.InsertRecordsStrong([]interface{}{so}, cqmodels.CQSettingsCollection)
}

//SetRepoStatus -
func (db MongoDB) SetRepoStatus(status int, repoID, folder, err string) error {
	filter := bson.M{
		"_id": bson.ObjectIdHex(repoID),
	}

	set := bson.M{
		"status": status,
		"error":  err,
	}

	if err == "" {
		set["project_fcode"] = folder
	}

	return db.cqMongo.UpdateRecords(filter, set, cqmodels.CQSettingsCollection)
}

//AddSSHAccessKey -
func (db MongoDB) AddSSHAccessKey(key cqmodels.AccessKeyParams) error {
	return db.cqMongo.InsertRecords([]interface{}{key}, cqmodels.CQGitAccessKeysCollection)
}

//GetAddedRepos -
func (db MongoDB) GetAddedRepos() (cqmodels.SettingsList, error) {
	filter := bson.M{
		"status": cqmodels.RepoStatusAdded,
	}

	res := cqmodels.SettingsList{}
	err := db.cqMongo.GetRecords(filter, cqmodels.CQSettingsCollection, &res)
	return res, err
}
