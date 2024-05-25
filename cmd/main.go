package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"payment/internal/discounts"
	"payment/internal/wallets"
	"payment/pkg/config"
	"payment/pkg/db"
	"payment/pkg/utils"
	"runtime"
	"strings"
	"time"
)

// Program Info
var (
	version = "1.0.0"
	build   = "Custom"
	name    = "Payment Service"
)

// Version returns version
func Version() string {
	return version
}

// VersionStatement returns a list of strings representing the full version info.
func VersionStatement() string {
	return strings.Join([]string{
		name, " ", Version(), " ", build, " (", runtime.Version(), " ", runtime.GOOS, "/", runtime.GOARCH, ")",
	}, "")
}

// Global variables for injecting dependencies into routers
var (
	logger        *log.Logger
	configuration *config.Config
	database      *db.DB
)

func main() {
	fmt.Println(VersionStatement())
	var err error

	logger = log.StandardLogger()
	logger.SetFormatter(&log.JSONFormatter{})

	configFilePath := flag.String("config", "configs/config.yaml", "Path to the YAML configuration file")
	flag.Parse()

	configuration, err = config.LoadConfig(*configFilePath)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("Configuration loaded")

	dsn := db.LoadDSN(configuration)
	database, err = db.New(dsn)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("Connected to database")

	var walletHandler = wallets.NewHandler(database, logger)
	var discountHandler = discounts.NewHandler(&discounts.Config{
		CreditExpiration: time.Duration(configuration.DiscountConfig.ExpireTime) * time.Minute,
		CodeLength:       configuration.DiscountConfig.CodeLength},
		logger, database)

	r := mux.NewRouter()
	http.Handle("/", utils.RecoverHandler(r))

	walletHandler.RegisterRoutes(r)
	discountHandler.RegisterRoutes(r)

	logger.Infof("%s is listening on port %d", name, configuration.ServerPort)
	err = http.ListenAndServe(fmt.Sprintf(":%d", configuration.ServerPort), nil)
	if err != nil {
		panic(err)
	}
}
