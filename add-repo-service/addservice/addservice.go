package addservice

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"gitlab.qarea.org/jiraquality/cq-web-backend/add-repo-service/cloner"
	"gitlab.qarea.org/jiraquality/cq-web-backend/add-repo-service/models"
	"gitlab.qarea.org/jiraquality/cq-web-backend/add-repo-service/storage"
	cqerrors "gitlab.qarea.org/jiraquality/cq-web-backend/common-lib/errors"
	"gitlab.qarea.org/jiraquality/cq-web-backend/common-lib/helpers"
	cqmodels "gitlab.qarea.org/jiraquality/cq-web-backend/common-lib/models"
	"gopkg.in/mgo.v2/bson"
)

//Cfg -
type Cfg struct {
	Storage *storage.MongoDB
	Config  models.Config
}

//Service -
type Service struct {
	storage    *storage.MongoDB
	config     models.Config
	crashQueue cqmodels.SettingsList
}

//New -
func New(cfg Cfg) *Service {
	return &Service{
		storage: cfg.Storage,
		config:  cfg.Config,
	}
}

//CheckAddProjectSettingsPermissions -
func (s Service) CheckAddProjectSettingsPermissions(userID string) error {
	usr, err := s.storage.DBLib.UserByID(userID)
	if err != nil {
		return err
	}

	if !usr.Permissions.AllowToManageRepos {
		return cqerrors.ErrPermissionDeny
	}

	list, err := s.storage.DBLib.GetActiveSettingsByUserID(bson.ObjectIdHex(userID))
	if err != nil {
		return err
	}

	if len(list) >= usr.Limitations.ReposLimit {
		return cqerrors.ErrPermissionDeny
	}

	return nil
}

//AddRepo -
func (s Service) AddRepo(userID, projectID, repoLink, branch string) (string, error) {
	if ok := helpers.IsValidGitLink(repoLink); !ok {
		return "", cqerrors.ErrBrokenGitLink
	}

	personalSC := hasPersonalScanner(userID, s.config.PScannersList)
	rID := bson.NewObjectId()
	settings := cqmodels.SettingsObject{
		ProjectID:        bson.ObjectIdHex(projectID),
		RepoLink:         repoLink,
		TrackedBranch:    branch,
		RepositoryID:     rID,
		UserID:           bson.ObjectIdHex(userID),
		RepositoryRating: 0,
		Status:           cqmodels.RepoStatusAdded,
		HasPersonalSC:    personalSC,
		BannerActive:     true,
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	}

	err := s.storage.InsertSettingRecordDB(settings)
	if err == nil {
		s.startClone(userID, rID.Hex(), repoLink, branch, s.cloneCallback)
	}

	return rID.Hex(), err
}

func (s Service) startClone(userID, repoID, link, branch string, callback cloner.Callback) {
	op := cloner.Options{
		Link:     link,
		Branch:   branch,
		RepoID:   repoID,
		Callback: callback,
	}

	folder := fmt.Sprintf("%s-%s-%s", helpers.RandomString(6), helpers.RandomString(6), helpers.RandomString(6))
	exPath, err := helpers.GetExecutablePath()
	if err != nil {
		callback(err, op)
		return
	}

	path := filepath.Join(exPath, s.config.CommonConfig.SourcePath, folder)
	op.Folder = folder
	op.Path = path

	p, err := s.storage.DBLib.GetSSHAccessKeyForUser(userID)
	if err != nil {
		callback(err, op)
		return
	}

	op.SSHKey = p.PublicKey
	cloner.StartClone(op)
}

func (s Service) cloneCallback(cloneErr error, op cloner.Options) {
	status := cqmodels.RepoStatusReadyForGitScan
	errText := ""
	if cloneErr != nil {
		status = cqmodels.RepoStatusFailedToClone
		errText = cloneErr.Error()
	}

	err := s.storage.SetRepoStatus(status, op.RepoID, op.Folder, errText)
	if err != nil {
		log.Println("set repo ", op.Folder, "status in clone callback error", err)
	}
}

func hasPersonalScanner(instanceToken string, list []string) bool {
	hasScanner, _ := helpers.InArray(instanceToken, list)
	return hasScanner
}
