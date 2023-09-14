package objects_db

import (
	"context"
	"fmt"

	"github.com/alexPavlikov/electronic_document_management/internal/entity/client"
	"github.com/alexPavlikov/electronic_document_management/internal/entity/objects"
	dbClient "github.com/alexPavlikov/electronic_document_management/pkg/client/postgresql"
	"github.com/alexPavlikov/electronic_document_management/pkg/logging"
	"github.com/alexPavlikov/electronic_document_management/pkg/utils"
)

type repository struct {
	client dbClient.Client
	logger logging.Logger
}

func NewRepository(client dbClient.Client, logger *logging.Logger) objects.Repository {
	return &repository{
		client: client,
		logger: *logger,
	}
}

func (r *repository) InsertObject(ctx context.Context, obj *objects.Object) error {
	query := `
	INSERT INTO 
		public."Objects" (name, address, work_schedule)
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

	r.logger.LogEvents("Добавлен", fmt.Sprintf("%s c id=:%d", "объект", &obj.Id))

	return nil
}

func (r *repository) SelectObject(ctx context.Context, id int) (obj objects.Object, err error) {
	query := `
	SELECT 
		id, name, address, work_schedule 
	FROM 
		public."Objects" 
	WHERE 
		id = $1
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	rows, err := r.client.Query(ctx, query, id)
	if err != nil {
		return objects.Object{}, err
	}

	err = rows.Scan(&obj.Id, &obj.Name, &obj.Address, &obj.WorkSchedule)
	if err != nil {
		return objects.Object{}, err
	}

	return obj, nil
}

func (r *repository) SelectObjects(ctx context.Context) (objs []objects.Object, err error) {
	query := `
	SELECT 
		id, name, address, work_schedule 
	FROM 
		public."Objects" 
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var obj objects.Object

	for rows.Next() {

		err = rows.Scan(&obj.Id, &obj.Name, &obj.Address, &obj.WorkSchedule)
		if err != nil {
			return nil, err
		}

		obj.Client, err = r.getClientObject(ctx, obj.Id)
		if err != nil {
			return nil, err
		}

		objs = append(objs, obj)
	}

	return objs, nil
}

func (r *repository) UpdateObject(ctx context.Context, obj *objects.Object) error {
	query := `
	UPDATE INTO 
		public."Objects" 
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

	r.logger.LogEvents("Обновлен", fmt.Sprintf("%s c id=:%d", "объект", &obj.Id))

	return nil
}

func (r *repository) DeleteObject(ctx context.Context, id int) error {
	query := `
	DELETE 
	FROM 
		public."Objects" 
	WHERE 
		id = $1
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	_, err := r.client.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	r.logger.LogEvents("Удален", fmt.Sprintf("%s c id=:%d", "объект", id))

	return nil
}

func (r *repository) getClientObject(ctx context.Context, id int) (cl client.Client, err error) {
	query := `
	SELECT 
		id, name, inn, kpp, ogrn, owner, phone, email, address, create_date, status
	FROM 
		public."Client"
	WHERE 
		id = $1
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	rows := r.client.QueryRow(ctx, query, id)
	err = rows.Scan(&cl.Id, &cl.Name, &cl.INN, &cl.KPP, &cl.OGRN, &cl.Owner, &cl.Phone, &cl.Email, &cl.Address, &cl.CreateDate, &cl.Status)
	if err != nil {
		return client.Client{}, err
	}

	return cl, nil
}
