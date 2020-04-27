package models

import (
	cqmodels "gitlab.qarea.org/jiraquality/cq-web-backend/common-lib/models"
)

// Config config
type Config struct {
	CommonConfig cqmodels.Config

	ForcedCORS    bool     `toml:"forced_cors_headers"`
	PScannersList []string `toml:"personal_scanners_list"`
	DevGitServers []string `toml:"dev_git_servers"`

	ExecPath      string
	BindedVersion string
}
