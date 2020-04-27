package models

import (
	cqmodels "gitlab.qarea.org/jiraquality/cq-web-backend/common-lib/models"
)

// Config config
type Config struct {
	CommonConfig cqmodels.Config

	ExecPath      string
	BindedVersion string
}
