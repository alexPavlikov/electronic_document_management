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
	"github.com/alexPavlikov/electronic_document_management/internal/entity/objects"
	objects_db "github.com/alexPavlikov/electronic_document_management/internal/entity/objects/db"
	"github.com/alexPavlikov/electronic_document_management/internal/entity/requests"
	requests_db "github.com/alexPavlikov/electronic_document_management/internal/entity/requests/db"
	"github.com/alexPavlikov/electronic_document_management/internal/entity/services"
	services_db "github.com/alexPavlikov/electronic_document_management/internal/entity/services/db"
	"github.com/alexPavlikov/electronic_document_management/internal/entity/user"
	user_db "github.com/alexPavlikov/electronic_document_management/internal/entity/user/db"
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

	logger.Info(config.LOG_INFO, " - Start contract handlers")
	ctRep := contract_db.NewRepository(ClientPostgreSQL, logger)
	ctSer := contract.NewService(ctRep, logger)
	ctHan := contract.NewHandler(ctSer, logger)
	ctHan.Register(router)

	logger.Info(config.LOG_INFO, " - Start object handlers")
	oRep := objects_db.NewRepository(ClientPostgreSQL, logger)
	oSer := objects.NewService(oRep, logger)
	oHan := objects.NewHandler(oSer, logger)
	oHan.Register(router)

	logger.Info(config.LOG_INFO, " - Start equipment handlers")
	eRep := equipment_db.NewRepository(ClientPostgreSQL, logger)
	eSer := equipment.NewService(eRep, logger)
	eHan := equipment.NewHandler(eSer, logger)
	eHan.Register(router)

	logger.Info(config.LOG_INFO, " - Start user handlers")
	uRep := user_db.NewRepository(ClientPostgreSQL, logger)
	uSer := user.NewService(uRep, logger)
	uHan := user.NewHandler(uSer, logger)
	uHan.Register(router)

	sRep := services_db.NewRepository(ClientPostgreSQL, logger)
	sSer := services.NewService(sRep, logger)
	sHan := services.NewHandler(sSer, logger)

	sHan.Register(router)

	requests.ClientsDTO, err = cSer.GetClients(context.TODO())
	if err != nil {
		logger.Fatalf("%s - %s", config.LOG_ERROR, err)
	}

	//requests.ClientObjectDTO

	requests.ContractsDTO, err = ctSer.GetContracts(context.TODO())
	if err != nil {
		logger.Fatalf("%s - %s", config.LOG_ERROR, err)
	}

	requests.EquipmentDTO, err = eSer.GetEquipments(context.TODO())
	if err != nil {
		logger.Fatalf("%s - %s", config.LOG_ERROR, err)
	}

	requests.WorkerDTO, err = uSer.GetUsers(context.TODO())
	if err != nil {
		logger.Fatalf("%s - %s", config.LOG_ERROR, err)
	}

	//requests.StatusDTO

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
		fmt.Println("ERRRRRRRRRRRRRRROR - ", err)
		logger.Fatal(config.LOG_ERROR, err.Error())
	}

}
