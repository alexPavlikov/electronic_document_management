package equipment

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

	router.ServeFiles("/assets/*filepath", http.Dir("assets"))

	router.HandlerFunc(http.MethodGet, "/edm/requests", h.RequestsHandler)
}

func NewHandler(service *Service, logger *logging.Logger) handlers.Handlers {
	return &handler{
		service: service,
		logger:  logger,
	}
}

func (h *handler) RequestsHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./internal/html/*.html")
	if err != nil {
		h.logger.Tracef("%s - failed open RequestsHandler", config.LOG_ERROR)
		w.WriteHeader(http.StatusNotFound)
	}

	data := map[string]interface{}{}
	header := map[string]string{"Title": "ЭДМ - Заявки"}

	err = tmpl.ExecuteTemplate(w, "header", header)
	if err != nil {
		h.logger.Tracef("%s - failed open RequestsHandler", config.LOG_ERROR)
		w.WriteHeader(http.StatusNotFound)
	}

	err = tmpl.ExecuteTemplate(w, "request", data)
	if err != nil {
		h.logger.Tracef("%s - failed open RequestsHandler", config.LOG_ERROR)
		w.WriteHeader(http.StatusNotFound)
	}

}
