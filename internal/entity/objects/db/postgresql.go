package object_db

import (
	"context"

	object "github.com/alexPavlikov/electronic_document_management/internal/entity/objects"
	dbClient "github.com/alexPavlikov/electronic_document_management/pkg/client/postgresql"
	"github.com/alexPavlikov/electronic_document_management/pkg/logging"
	"github.com/alexPavlikov/electronic_document_management/pkg/utils"
)

type repository struct {
	client dbClient.Client
	logger logging.Logger
}

func NewRepository(client dbClient.Client, logger *logging.Logger) object.Repository {
	return &repository{
		client: client,
		logger: *logger,
	}
}

func (r *repository) InsertObject(ctx context.Context, obj *object.Object) error {
	query := `
	INSERT INTO 
		public."Object" (name, address, work_schedule)
	VALUES 
		($1, $2, $3, $4)
	RETURNING 
		id
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	rows := r.client.QueryRow(ctx, query, &obj.Name, &obj.Address, &obj.WorkSchedule)

	err := rows.Scan(&obj.Id)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) SelectObject(ctx context.Context, id int) (obj object.Object, err error) {
	query := `
	SELECT 
		id, name, address, work_schedule 
	FROM 
		public."Object" 
	WHERE 
		id = $1
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	rows, err := r.client.Query(ctx, query, id)
	if err != nil {
		return object.Object{}, err
	}

	err = rows.Scan(&obj.Id, &obj.Name, &obj.Address, &obj.WorkSchedule)
	if err != nil {
		return object.Object{}, err
	}

	return obj, nil
}

func (r *repository) SelectObjects(ctx context.Context) (objs []object.Object, err error) {
	query := `
	SELECT 
		id, name, address, work_schedule 
	FROM 
		public."Object" 
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var obj object.Object

	for rows.Next() {

		err = rows.Scan(&obj.Id, &obj.Name, &obj.Address, &obj.WorkSchedule)
		if err != nil {
			return nil, err
		}

		objs = append(objs, obj)
	}

	return objs, nil
}

func (r *repository) UpdateObject(ctx context.Context, obj *object.Object) error {
	query := `
	UPDATE INTO 
		public."Object" 
	SET 
		name = $1, address = $2, work_schedule = $3
	WHERE 
		id = $4
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	_, err := r.client.Query(ctx, query, &obj.Name, &obj.Address, &obj.WorkSchedule, &obj.Id)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) DeleteObject(ctx context.Context, id int) error {
	query := `
	DELETE 
	FROM 
		public."Object" 
	WHERE 
		id = $1
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	_, err := r.client.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
