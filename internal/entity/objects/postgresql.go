package object

import (
	dbClient "github.com/alexPavlikov/electronic_document_management/pkg/client/postgresql"
	"github.com/alexPavlikov/electronic_document_management/pkg/logging"
)

type repository struct {
	client dbClient.Client
	logger logging.Logger
}

func NewRepository(client dbClient.Client, logger *logging.Logger) Repository {
	return &repository{
		client: client,
		logger: *logger,
	}
}
