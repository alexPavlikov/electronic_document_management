package user

import (
	"net/http"
	"text/template"

	"github.com/alexPavlikov/electronic_document_management/internal/config"
	"github.com/alexPavlikov/electronic_document_management/internal/handlers"
	"github.com/alexPavlikov/electronic_document_management/pkg/logging"
	"github.com/julienschmidt/httprouter"
)

type handler struct {
	service *Service
	logger  *logging.Logger
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, "/edm/users", h.UsersHandler)
}

func NewHandler(service *Service, logger *logging.Logger) handlers.Handlers {
	return &handler{
		service: service,
		logger:  logger,
	}
}

func (h *handler) UsersHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./internal/html/*.html")
	if err != nil {
		h.logger.Tracef("%s - failed open UserHandler", config.LOG_ERROR)
		w.WriteHeader(http.StatusNotFound)
	}

	err = tmpl.ExecuteTemplate(w, "users", nil)
	if err != nil {
		h.logger.Tracef("%s - failed open UserHandler", config.LOG_ERROR)
		w.WriteHeader(http.StatusNotFound)
	}
}
