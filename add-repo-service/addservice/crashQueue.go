package addservice

import (
	"log"

	"gitlab.qarea.org/jiraquality/cq-web-backend/add-repo-service/cloner"
	cqmodels "gitlab.qarea.org/jiraquality/cq-web-backend/common-lib/models"
)

//StartCrashQueue -
func (s Service) StartCrashQueue() {
	go s.doStartCrashQueue()
}

func (s Service) doStartCrashQueue() {
	repos, err := s.storage.GetAddedRepos()
	if err != nil {
		log.Println("doStartCrashQueue. getting added repos error", err)
		return
	}

	s.crashQueue = repos
	if len(s.crashQueue) > 0 {
		s.cloneNextRepo(s.crashQueue[0])
	}
}

func (s Service) cloneNextRepo(repo cqmodels.SettingsObject) {
	s.startClone(repo.UserID.Hex(), repo.RepositoryID.Hex(), repo.RepoLink, repo.TrackedBranch, s.queueCloneCallback)
}

func (s Service) queueCloneCallback(cloneErr error, op cloner.Options) {
	s.cloneCallback(cloneErr, op)

	s.crashQueue = s.crashQueue[1:]
	if len(s.crashQueue) > 0 {
		s.cloneNextRepo(s.crashQueue[0])
	}
}
