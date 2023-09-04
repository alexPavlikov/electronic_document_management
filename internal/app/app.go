package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/alexPavlikov/electronic_document_management/internal/config"
	"github.com/alexPavlikov/electronic_document_management/internal/entity/client"
	client_db "github.com/alexPavlikov/electronic_document_management/internal/entity/client/db"
	"github.com/alexPavlikov/electronic_document_management/internal/entity/contract"
	contract_db "github.com/alexPavlikov/electronic_document_management/internal/entity/contract/db"
	"github.com/alexPavlikov/electronic_document_management/internal/entity/equipment"
	equipment_db "github.com/alexPavlikov/electronic_document_management/internal/entity/equipment/db"
	object "github.com/alexPavlikov/electronic_document_management/internal/entity/objects"
	object_db "github.com/alexPavlikov/electronic_document_management/internal/entity/objects/db"
	"github.com/alexPavlikov/electronic_document_management/internal/entity/requests"
	requests_db "github.com/alexPavlikov/electronic_document_management/internal/entity/requests/db"
	dbClient "github.com/alexPavlikov/electronic_document_management/pkg/client/postgresql"
	"github.com/alexPavlikov/electronic_document_management/pkg/logging"

	"github.com/julienschmidt/httprouter"
)

var ClientPostgreSQL dbClient.Client

func Run() {
	logger := logging.GetLogger()
	logger.Info(config.LOG_INFO, "Create router")
	router := httprouter.New()

	cfg := config.GetConfig()

	var err error

	ClientPostgreSQL, err = dbClient.NewClient(context.TODO(), cfg.Storage)
	if err != nil {
		logger.Fatalf("failed to get new client postgresql, due to err: %v", err)
	}

	logger.Info(config.LOG_INFO, " - Start requests handlers")
	rRep := requests_db.NewRepository(ClientPostgreSQL, logger)
	rSer := requests.NewService(rRep, logger)
	rHan := requests.NewHandler(rSer, logger)
	rHan.Register(router)

	logger.Info(config.LOG_INFO, " - Start client handlers")
	cRep := client_db.NewRepository(ClientPostgreSQL, logger)
	cSer := client.NewService(cRep, logger)
	cHan := client.NewHandler(cSer, logger)
	cHan.Register(router)

	ctRep := contract_db.NewRepository(ClientPostgreSQL, logger)
	ctSer := contract.NewService(ctRep, logger)
	ctHan := contract.NewHandler(ctSer, logger)
	ctHan.Register(router)

	oRep := object_db.NewRepository(ClientPostgreSQL, logger)
	oSer := object.NewService(oRep, logger)
	oHan := object.NewHandler(oSer, logger)
	oHan.Register(router)

	eRep := equipment_db.NewRepository(ClientPostgreSQL, logger)
	eSer := equipment.NewService(eRep, logger)
	eHan := equipment.NewHandler(eSer, logger)
	eHan.Register(router)

	start(router, *cfg)
}

func start(r *httprouter.Router, cfg config.Config) {
	logger := logging.GetLogger()
	logger.Info(config.LOG_INFO, "Start application")

	var listener net.Listener
	var listenerErr error

	logger.Info(config.LOG_INFO, "Listen TCP")
	listener, listenerErr = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
	logger.Infof("%s Server is listen on port: %s:%s", config.LOG_INFO, cfg.Listen.BindIP, cfg.Listen.Port)
	if listenerErr != nil {
		logger.Fatal(config.LOG_ERROR, listenerErr.Error())
	}

	server := &http.Server{
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	err := server.Serve(listener)
	if err != nil {
		logger.Fatal(config.LOG_ERROR, err.Error())
	}

}
