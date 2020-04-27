package addservice

import (
	"strings"

	"gitlab.qarea.org/jiraquality/cq-web-backend/common-lib/helpers"
	cqmodels "gitlab.qarea.org/jiraquality/cq-web-backend/common-lib/models"

	"gopkg.in/mgo.v2/bson"
)

//GetSSHAccessKeyForUser retrieve pre-generated RSA key for requested user and
//generate new one if it doesn't exist
func (s Service) GetSSHAccessKeyForUser(userID, repoURL string) ([]byte, error) {
	if s.isDebugGitServer(repoURL) {
		return []byte(s.config.CommonConfig.DevSSHPublicKey), nil
	}

	keyObject, err := s.storage.DBLib.GetSSHAccessKeyForUser(userID)
	if err != nil {
		return nil, err
	}

	//check if exists
	if !keyObject.ID.Valid() {
		return s.addNewKey(userID)
	}

	return keyObject.PublicKey, nil
}

func (s Service) addNewKey(userID string) ([]byte, error) {
	pubKey, privKey, err := helpers.PrepareNewKeys()
	if err != nil {
		return nil, err
	}

	ko := cqmodels.AccessKeyParams{
		UserID:     bson.ObjectIdHex(userID),
		PublicKey:  pubKey,
		PrivateKey: privKey,
	}

	err = s.storage.AddSSHAccessKey(ko)
	if err != nil {
		return nil, err
	}

	return pubKey, nil
}

func (s Service) isDebugGitServer(repoURL string) bool {
	for _, server := range s.config.DevGitServers {
		if strings.Contains(repoURL, server) {
			return true
		}
	}

	return false
}
