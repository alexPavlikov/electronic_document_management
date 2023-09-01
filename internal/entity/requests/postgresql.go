package requests

import (
	"context"

	dbClient "github.com/alexPavlikov/electronic_document_management/pkg/client/postgresql"
	"github.com/alexPavlikov/electronic_document_management/pkg/logging"
	"github.com/alexPavlikov/electronic_document_management/pkg/utils"
)

type repository struct {
	client dbClient.Client
	logger logging.Logger
}

// InsertRequest implements Repository.
func (r *repository) InsertRequest(ctx context.Context, req Request) error {
	query := `
	INSERT INTO public.Request
	(title, client, worker, client_object, equipment, contract, description, priority, start_date, end_date, files, status)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
	RETURNING id
	`

	r.logger.Tracef("SQL Query: %s", utils.FormatQuery(query))

	rows := r.client.QueryRow(ctx, query, &req.Title, &req.Client, &req.Worker, &req.ClientObject, &req.Equipment, &req.Contract, &req.Description, &req.Priority, &req.StartDate, &req.EndDate, &req.Files, &req.Status)
	err := rows.Scan(req.Id)
	if err != nil {
		return err
	}
	return nil
}

// SelectRequest implements Repository.
func (r *repository) SelectRequest(ctx context.Context, id string) (req Request, err error) {
	query := `
	SELECT 
		id, title, client, worker, client_object, equipment, contract, description, priority, start_date, end_date, files, status
	FROM 
		public.Request
	WHERE 
		id = $1`

	r.logger.Tracef("SQL Query: %s", utils.FormatQuery(query))

	rows, err := r.client.Query(ctx, query, id)
	if err != nil {
		return Request{}, err
	}
	for rows.Next() {
		err = rows.Scan(&req.Id, &req.Title, &req.Client, &req.Worker, &req.ClientObject, &req.Equipment, &req.Contract, &req.Description, &req.Priority, &req.StartDate, &req.EndDate, &req.Files, &req.Status)
		if err != nil {
			return Request{}, err
		}
	}
	return req, nil
}

// SelectRequests implements Repository.
func (r *repository) SelectRequests(ctx context.Context) (reqs []Request, err error) {
	query := `
	SELECT 
		id, title, client, worker, client_object, equipment, contract, description, priority, start_date, end_date, files, status
	FROM 
		public.Request`

	r.logger.Tracef("SQL Query: %s", utils.FormatQuery(query))

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var r Request
		err = rows.Scan(&r.Id, &r.Title, &r.Client, &r.Worker, &r.ClientObject, &r.Equipment, &r.Contract, &r.Description, &r.Priority, &r.StartDate, &r.EndDate, &r.Files, &r.Status)
		if err != nil {
			return nil, err
		}
		reqs = append(reqs, r)
	}
	return reqs, nil
}

// UpdateRequest implements Repository.
func (r *repository) UpdateRequest(ctx context.Context, req *Request) error {
	query := `
	UPDATE 
		public.Request
	SET 
		title = $1, client = $2, worker = $3, client_object = $4, equipment = $5, contract = $6, description = $7, priority = $8, start_date = $9, end_date = $10, files = $11, status = $12
	WHERE
		id = $13`

	r.logger.Tracef("SQL Query: %s", utils.FormatQuery(query))

	_, err := r.client.Query(ctx, query, req.Title, req.Client, req.Worker, req.ClientObject, req.Equipment, req.Contract, req.Description, req.Priority, req.StartDate, req.EndDate, req.Files, req.Status, req.Id)
	if err != nil {
		return err
	}
	return nil
}

// DeleteRequest implements Repository.
func (r *repository) CloseRequest(ctx context.Context, status string, id string) error {
	query := `
	UPDATE 
		public.Request
	SET 
		status = $1
	WHERE
		id = $2`

	r.logger.Tracef("SQL Query: %s", utils.FormatQuery(query))

	_, err := r.client.Query(ctx, query, status, id)
	if err != nil {
		return err
	}
	return nil
}

func NewRepository(client dbClient.Client, logger *logging.Logger) Repository {
	return &repository{
		client: client,
		logger: *logger,
	}
}
