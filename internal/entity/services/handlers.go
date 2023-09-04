package services

import (
	"context"
	"net/http"
	"text/template"

	"github.com/alexPavlikov/electronic_document_management/internal/handlers"
	"github.com/alexPavlikov/electronic_document_management/pkg/logging"
	"github.com/julienschmidt/httprouter"
)

type handler struct {
	service *Service
	logger  *logging.Logger
}

func (h *handler) Register(router *httprouter.Router) {

	//router.ServeFiles("/assets/*filepath", http.Dir("assets"))

	router.HandlerFunc(http.MethodGet, "/edm/services", h.ServicesHandler)
}

func NewHandler(service *Service, logger *logging.Logger) handlers.Handlers {
	return &handler{
		service: service,
		logger:  logger,
	}
}

func (h *handler) ServicesHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("./internal/html/*.html")
	if err != nil {
		http.NotFound(w, r)
	}

	services, err := h.service.GetServices(context.TODO())
	if err != nil {
		http.NotFound(w, r)
	}

	title := map[string]string{"Title": "ЭДО - Услуги"}
	data := map[string]interface{}{"Services": services}

	err = tmpl.ExecuteTemplate(w, "header", title)
	if err != nil {
		http.NotFound(w, r)
	}

	err = tmpl.ExecuteTemplate(w, "services", data)
	if err != nil {
		http.NotFound(w, r)
	}
}
