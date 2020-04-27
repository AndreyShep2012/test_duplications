package scanner

import (
	cqmodels "gitlab.qarea.org/jiraquality/cq-web-backend/common-lib/models"
	"gitlab.qarea.org/jiraquality/cq-web-backend/git-scanner-service/models"
	"gitlab.qarea.org/jiraquality/cq-web-backend/git-scanner-service/storage"
)

//Cfg -
type Cfg struct {
	Storage *storage.MongoDB
	Config  models.Config
}

//Scanner -
type Scanner struct {
	storage    *storage.MongoDB
	config     models.Config
	crashQueue cqmodels.ReposList
}

//New -
func New(cfg Cfg) *Scanner {
	return &Scanner{
		storage: cfg.Storage,
		config:  cfg.Config,
	}
}

//Start -
func (s Scanner) Start() {
	s.startCrashQueue()
}

func (s Scanner) startCron() {

}
