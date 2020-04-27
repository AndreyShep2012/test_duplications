package api

import (
	"gitlab.qarea.org/jiraquality/cq-web-backend/user-service/ctx"
)

//AddProjectSettingsRequest -
type AddProjectSettingsRequest struct {
	Context       ctx.Context
	ProjectID     string `json:"project_id"`
	ProjectLink   string `json:"project_link"`
	TrackedBranch string `json:"tracked_branch"`
	UserID        string `json:"user_id"`
}

//CheckRequired -
func (a AddProjectSettingsRequest) CheckRequired() bool {
	return a.Context.Token != "" && a.ProjectID != "" &&
		a.ProjectLink != "" && a.TrackedBranch != "" &&
		a.UserID != ""
}

//AddProjectSettingsResponse -
type AddProjectSettingsResponse struct {
	ID string `json:"repository_id"`
}

//GetStorageKeyRequest -
type GetStorageKeyRequest struct {
	Context ctx.Context
	UserID  string `json:"user_id"`
}

//GetStorageKeyResponse -
type GetStorageKeyResponse struct {
	Key string
}

//CheckRequired -
func (g GetStorageKeyRequest) CheckRequired() bool {
	return g.Context.Token != "" && g.UserID != ""
}
