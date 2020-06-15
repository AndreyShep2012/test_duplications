package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/StanislavKH/rpc/v2"
	"github.com/StanislavKH/rpc/v2/json2"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"

	"gitlab.qarea.org/jiraquality/cq-web-backend/add-repo-service/addservice"
	"gitlab.qarea.org/jiraquality/cq-web-backend/add-repo-service/api"
	"gitlab.qarea.org/jiraquality/cq-web-backend/add-repo-service/models"
	"gitlab.qarea.org/jiraquality/cq-web-backend/add-repo-service/storage"
	"gitlab.qarea.org/jiraquality/cq-web-backend/user-service/ctx"

	"gitlab.qarea.org/jiraquality/cq-web-backend/common-lib/clients"
	"gitlab.qarea.org/jiraquality/cq-web-backend/common-lib/helpers"
	"gitlab.qarea.org/jiraquality/cq-web-backend/common-lib/logger"
)

var conf models.Config
var test int

func init() {
	exPath, err := helpers.GetExecutablePath()
	if err != nil {
		log.Fatal("can't detect executable path:", err)
	}

	confPath := filepath.Join(exPath, "config", "config.toml")
	if _, err := toml.DecodeFile(confPath, &conf); err != nil {
		log.Fatal("can't process configuration for application: ", err)
	}

	commonConfPath := filepath.Join(exPath, "../", "common-config", "config.toml")
	if _, err := toml.DecodeFile(commonConfPath, &conf.CommonConfig); err != nil {
		log.Fatal("can't process common configuration for application: ", err)
	}

	conf.ExecPath = exPath
}

func main() {
	initLogger()

	as := initAddService(initMongoDB())
	as.StartCrashQueue()
	serviceCfg := &api.Cfg{
		Conf:       &conf,
		HTTPClient: clients.CreateHTTPClient(),
		Parser:     initTokenParser(),
		AddService: as,
	}

	if testFunction(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11) == 9 {
		log.Fatal("asdasdasd")
	}

	if testFunction12(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11) == 9 {
		log.Fatal("asdasdasd")
	}

	if testFunction13(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11) == 9 {
		log.Fatal("asdasdasd")
	}

	if testFunction14(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11) == 9 {
		log.Fatal("asdasdasd")
	}

	service := api.NewService(serviceCfg)
	startService(service)
}

func initLogger() {
	logger.RedirectLog(&logger.Conf{
		Filename:   "log/log.log",
		MaxSize:    conf.CommonConfig.LogMaxSize,
		MaxBackups: conf.CommonConfig.LogMaxBackups,
		MaxAge:     conf.CommonConfig.LogMaxAge,
	}, nil)
}

func testFunction(p1, p2, p3, p4, p5, p6, p7, p8, p9, p10, p11, p12 int) int {
	return 10
}

func testFunction12(p1, p2, p3, p4, p5, p6, p7, p8, p9, p10, p11, p12 int) int {
	return 10
}

func testFunction13(p1, p2, p3, p4, p5, p6, p7, p8, p9, p10, p11, p12 int) int {
	return 10
}

func testFunction13(p1, p2, p3, p4, p5, p6, p7, p8, p9, p10, p11, p12 int) int {
	return 10
}

func testFunction14(p1, p2, p3, p4, p5, p6, p7, p8, p9, p10, p11, p12 int) int {
	return 10
}

func initAddService(storage *storage.MongoDB) *addservice.Service {
	return addservice.New(
		addservice.Cfg{
			Storage: storage,
			Config:  conf,
		},
	)
}

func initTokenParser() *ctx.RSATokenParser {
	tokenParser, err := ctx.NewRSATokenParser([]byte(conf.CommonConfig.AuthRSAPublicKey))
	if err != nil {
		log.Fatal("can't init token parser: ", err)
	}

	return tokenParser
}

func initMongoDB() *storage.MongoDB {
	log.Println("open mongo database...")

	mgoSession := clients.GetMongoSession(&clients.MongoConf{
		Host: conf.CommonConfig.DB.Host,
		Auth: conf.CommonConfig.DB.Auth,
		User: conf.CommonConfig.DB.User,
		Pass: conf.CommonConfig.DB.Pass,
	})

	cfg := &storage.Cfg{
		Session:      mgoSession,
		DatabaseName: conf.CommonConfig.DB.DatabaseName,
	}
	ds := storage.New(cfg)
	return ds
}

func startService(service *api.API) {
	ServiceEndpoint := "/rpc"
	configuredHost := fmt.Sprintf("%s:%d", conf.CommonConfig.AddRepoService.IP, conf.CommonConfig.AddRepoService.Port)
	log.Println("Starting RPC Collector Server on", configuredHost)

	encoder := json2.NewCodec()

	r := mux.NewRouter()
	s := rpc.NewServer()
	s.RegisterCodec(encoder, "application/json")
	s.RegisterCodec(encoder, "application/json; charset=UTF-8")
	s.RegisterBeforeFunc(helpers.EndpointRequestedBefore)
	s.RegisterService(service, "")

	if conf.CommonConfig.AddRepoService.Endpoint != "" {
		ServiceEndpoint = conf.CommonConfig.AddRepoService.Endpoint
	}

	r.Handle(ServiceEndpoint, s)
	r.Handle(ServiceEndpoint+"/metrics", promhttp.Handler()).Methods("GET")

	if conf.ForcedCORS {
		handler := cors.Default().Handler(r)
		log.Fatal(http.ListenAndServe(configuredHost, handler))
	} else {
		log.Fatal(http.ListenAndServe(configuredHost, r))
	}
}
