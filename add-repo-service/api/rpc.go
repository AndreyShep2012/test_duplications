package api

import (
	"context"
	"net/http"

	"gitlab.qarea.org/jiraquality/cq-web-backend/add-repo-service/addservice"
	"gitlab.qarea.org/jiraquality/cq-web-backend/add-repo-service/models"
	"gitlab.qarea.org/jiraquality/cq-web-backend/user-service/ctx"

	cqerrors "gitlab.qarea.org/jiraquality/cq-web-backend/common-lib/errors"
)

// Cfg  api config
type Cfg struct {
	HTTPClient *http.Client
	Conf       *models.Config
	Parser     *ctx.RSATokenParser
	AddService *addservice.Service
}

// API -
type API struct {
	client     *http.Client
	conf       *models.Config
	parser     *ctx.RSATokenParser
	addService *addservice.Service
}

//NewService instance entity
func NewService(c *Cfg) *API {
	return &API{
		client:     c.HTTPClient,
		conf:       c.Conf,
		parser:     c.Parser,
		addService: c.AddService,
	}
}

//Version adds profile
func (h API) Version(r *http.Request, args *struct{}, res *string) error {
	*res = h.conf.BindedVersion
	return nil
}

//AddRepo --
func (h API) AddRepo(r *http.Request, args *AddProjectSettingsRequest, reply *AddProjectSettingsResponse) error {
	if !args.CheckRequired() {
		return cqerrors.ErrRequiredEmpty
	}

	return h.parser.ParseCtxWithClaims(args.Context, func(ctx context.Context, claims ctx.Claims) error {
		err := h.addService.CheckAddProjectSettingsPermissions(args.UserID)
		if err != nil {
			return err
		}

		id, err := h.addService.AddRepo(args.UserID, args.ProjectID, args.ProjectLink, args.TrackedBranch)
		reply.ID = id
		return err
	})
}

//GetStorageKey return pre-defined RSA public key, should be added for repositore access in Read Only mode
func (h API) GetStorageKey(r *http.Request, args *GetStorageKeyRequest, reply *GetStorageKeyResponse) error {
	if !args.CheckRequired() {
		return cqerrors.ErrRequiredEmpty
	}

	return h.parser.ParseCtxWithClaims(args.Context, func(ctx context.Context, claims ctx.Claims) error {
		key, err := h.addService.GetSSHAccessKeyForUser(args.UserID, "")
		reply.Key = string(key)
		return err
	})
}
