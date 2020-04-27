package scanner

import (
	"log"
	"path/filepath"

	"gopkg.in/mgo.v2/bson"

	gitlib "gitlab.qarea.org/jiraquality/cq-web-backend/git-lib"

	cqmodels "gitlab.qarea.org/jiraquality/cq-web-backend/common-lib/models"
)

//StartCrashQueue -
func (s Scanner) startCrashQueue() {
	go s.doStartCrashQueue()
}

func (s Scanner) doStartCrashQueue() {
	repos, err := s.storage.DBLib.GetReposForStatus(cqmodels.RepoStatusReadyForGitScan)
	if err != nil {
		log.Println("doStartCrashQueue. getting added repos error", err)
		return
	}

	s.crashQueue = repos
	if len(s.crashQueue) > 0 {
		s.scanNextRepo(s.crashQueue[0])
	} else {
		s.startCron()
	}
}

func (s Scanner) scanNextRepo(repo cqmodels.RepoObject) {
	if repo.IsFirstProcessing {
		err := s.repoFirstProcessing(repo)
		if err != nil {
			err = s.storage.DBLib.SetRepoStatus(repo.RepositoryID, cqmodels.RepoStatusGitScanError, err)
			if err != nil {
				log.Println("scanNextRepo. set repo status error", err, repo.RepositoryID)
			}

			return
		}
	}
}

func (s Scanner) repoFirstProcessing(repo cqmodels.RepoObject) error {
	commits, commiters, err := s.getAllCommitsAndCommitters(repo)
	if err != nil {
		return err
	}

	if err := s.storage.SaveCommits(commits); err != nil {
		return err
	}

	if err := s.storage.SaveCommiters(commiters); err != nil {
		return err
	}

	return nil
}

func (s Scanner) getAllCommitsAndCommitters(repo cqmodels.RepoObject) (commits, commiters []interface{}, err error) {
	gcommits, err := gitlib.GetCommitsList(s.getRepoSourcePath(repo))
	if err != nil {
		return
	}

	for _, c := range gcommits {
		commits = append(commits, createCQCommitObject(repo, c))
		if !commiterInArray(commiters, c.AuthorName, c.AuthorEmail) {
			commiters = append(commiters, createCQCommiterObject(c, repo.UserID))
		}
	}

	return
}

func (s Scanner) getRepoSourcePath(repo cqmodels.RepoObject) string {
	return filepath.Join(s.config.CommonConfig.SourcePath, repo.ProjectFCode)
}

func createCQCommitObject(repo cqmodels.RepoObject, commit gitlib.Commit) cqmodels.Commit {
	return cqmodels.Commit{
		ProjectID:    repo.ProjectID,
		RepositoryID: repo.RepositoryID,
		UserID:       repo.UserID,
		CommitHash:   commit.Hash,
		CommitAuth:   commit.AuthorEmail,
		CommitSubj:   commit.Subject,
		IsMerge:      commit.IsMerge,
		CommitTime:   commit.Time,
	}
}

func createCQCommiterObject(commit gitlib.Commit, userID bson.ObjectId) cqmodels.Commiter {
	return cqmodels.Commiter{
		Name:     commit.AuthorName,
		Email:    commit.AuthorEmail,
		IsActive: true,
		AtName:   commit.AuthorName,
		AtEmail:  commit.AuthorEmail,
		UserID:   userID,
	}
}

func commiterInArray(arr []interface{}, name, email string) bool {
	for _, c := range arr {
		if cmtr, ok := c.(cqmodels.Commiter); ok {
			if cmtr.Email == email && cmtr.Name == name {
				return true
			}
		} else {
			return false
		}
	}

	return false
}
