package requests_db

import (
	"context"
	"fmt"

	"github.com/alexPavlikov/electronic_document_management/internal/entity/client"
	"github.com/alexPavlikov/electronic_document_management/internal/entity/contract"
	"github.com/alexPavlikov/electronic_document_management/internal/entity/equipment"
	"github.com/alexPavlikov/electronic_document_management/internal/entity/objects"
	"github.com/alexPavlikov/electronic_document_management/internal/entity/requests"
	"github.com/alexPavlikov/electronic_document_management/internal/entity/user"
	dbClient "github.com/alexPavlikov/electronic_document_management/pkg/client/postgresql"
	"github.com/alexPavlikov/electronic_document_management/pkg/logging"
	"github.com/alexPavlikov/electronic_document_management/pkg/utils"
	"github.com/lib/pq"
)

type repository struct {
	client dbClient.Client
	logger logging.Logger
}

// InsertRequest implements Repository.
func (r *repository) InsertRequest(ctx context.Context, req *requests.Request) error {
	query := `
	INSERT INTO public."Request"
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

	r.logger.LogEvents("Добавлена", fmt.Sprintf("%s c id=:%d на сотрудника - %s", "заявка", &req.Id, fmt.Sprint(&req.Worker)))

	return nil
}

// SelectRequest implements Repository.
func (r *repository) SelectRequest(ctx context.Context, id int) (req requests.Request, err error) {
	query := `
	SELECT 
		id, title, description, priority, start_date, end_date, files
	FROM 
		public."Request"
	WHERE 
		id = $1`

	r.logger.Tracef("SQL Query: %s", utils.FormatQuery(query))

	rows, err := r.client.Query(ctx, query, id)
	if err != nil {
		return requests.Request{}, err
	}
	for rows.Next() {
		err = rows.Scan(&req.Id, &req.Title, &req.Description, &req.Priority, &req.StartDate, &req.EndDate, pq.Array(&req.Files))
		if err != nil {
			return requests.Request{}, err
		}
	}
	return req, nil
}

// SelectRequests implements Repository.
func (r *repository) SelectRequests(ctx context.Context) (reqs []requests.Request, err error) {
	query := `
	SELECT 
		id, title, client, worker, client_object, equipment, contract, description, priority, start_date, end_date, files, status
	FROM 
		public."Request"`

	r.logger.Tracef("SQL Query: %s", utils.FormatQuery(query))

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var client, worker, client_object, equipment, contract int
	var status string

	for rows.Next() {
		var rt requests.Request
		err = rows.Scan(&rt.Id, &rt.Title, &client, &worker, &client_object, &equipment, &contract, &rt.Description, &rt.Priority, &rt.StartDate, &rt.EndDate, pq.Array(&rt.Files), &status)
		if err != nil {
			return nil, err
		}

		rt.Client, err = r.getRequestClient(ctx, client)
		if err != nil {
			return nil, err
		}
		rt.Worker, err = r.getRequestWorker(ctx, worker)
		if err != nil {
			return nil, err
		}
		rt.Contract, err = r.getRequestContract(ctx, contract)
		if err != nil {
			return nil, err
		}
		rt.Equipment, err = r.getRequestEquipment(ctx, equipment)
		if err != nil {
			return nil, err
		}
		rt.ClientObject, err = r.getRequestClientObject(ctx, client_object)
		if err != nil {
			return nil, err
		}
		rt.Status, err = r.getRequestStatus(ctx, status)
		if err != nil {
			return nil, err
		}

		reqs = append(reqs, rt)
	}
	return reqs, nil
}

// UpdateRequest implements Repository.
func (r *repository) UpdateRequest(ctx context.Context, req *requests.Request) error {
	query := `
	UPDATE 
		public."Request"
	SET 
		title = $1, client = $2, worker = $3, client_object = $4, equipment = $5, contract = $6, description = $7, priority = $8, start_date = $9, end_date = $10, files = $11, status = $12
	WHERE
		id = $13`

	r.logger.Tracef("SQL Query: %s", utils.FormatQuery(query))

	_, err := r.client.Query(ctx, query, req.Title, req.Client, req.Worker, req.ClientObject, req.Equipment, req.Contract, req.Description, req.Priority, req.StartDate, req.EndDate, req.Files, req.Status, req.Id)
	if err != nil {
		return err
	}

	r.logger.LogEvents("Изменена", fmt.Sprintf("%s c id=:%d сотрудником - %s", "заявка", &req.Id, fmt.Sprint(&req.Worker)))

	return nil
}

// DeleteRequest implements Repository.
func (r *repository) CloseRequest(ctx context.Context, status string, id int) error {
	query := `
	UPDATE 
		public."Request"
	SET 
		status = $1
	WHERE
		id = $2`

	r.logger.Tracef("SQL Query: %s", utils.FormatQuery(query))

	_, err := r.client.Query(ctx, query, status, id)
	if err != nil {
		return err
	}

	r.logger.LogEvents("Закрыта", fmt.Sprintf("%s c id=:%d", "заявка", id))

	return nil
}

func NewRepository(client dbClient.Client, logger *logging.Logger) requests.Repository {
	return &repository{
		client: client,
		logger: *logger,
	}
}

// client, worker, client_object, equipment, contract
// status

func (r *repository) getRequestClient(ctx context.Context, id int) (cl client.Client, err error) {
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
func (r *repository) getRequestWorker(ctx context.Context, id int) (us user.User, err error) {
	query := `
	SELECT 
		id, email, full_name, phone, image, role 
	FROM 
		public."User"
	WHERE 
		id = $1
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	rows := r.client.QueryRow(ctx, query, id)
	err = rows.Scan(&us.Id, &us.Email, &us.FullName, &us.Phone, &us.Image, &us.Role)
	if err != nil {
		return user.User{}, err
	}

	return us, nil
}
func (r *repository) getRequestContract(ctx context.Context, id int) (ct contract.Contract, err error) {
	query := `
	SELECT 
		id, name, client, start_date, end_date, amount, file, status 
	FROM 
		public."Contract"
	WHERE 
		id = $1
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	rows := r.client.QueryRow(ctx, query, id)
	err = rows.Scan(&ct.Id, &ct.Name, &ct.Client, &ct.DataStart, &ct.DataEnd, &ct.Amount, &ct.File, &ct.Status)
	if err != nil {
		return contract.Contract{}, err
	}

	return ct, nil
}
func (r *repository) getRequestEquipment(ctx context.Context, id int) (eq equipment.Equipment, err error) {
	query := `
	SELECT 
		id, name, type, manufacturer, model, unique_number, contract, create_date
	FROM 
		public."Equipment"
	WHERE 
		id = $1
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	rows := r.client.QueryRow(ctx, query, id)
	err = rows.Scan(&eq.Id, &eq.Name, &eq.Type, &eq.Manufacture, &eq.Model, &eq.UniqueNumber, &eq.Contract, &eq.CreateDate)
	if err != nil {
		return equipment.Equipment{}, err
	}

	return eq, nil
}
func (r *repository) getRequestClientObject(ctx context.Context, id int) (cl client.ClientObject, err error) {
	query := `
	SELECT 
		id, object
	FROM 
		public."Client_objects"
	WHERE 
		id = $1
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	rows := r.client.QueryRow(ctx, query, id)
	err = rows.Scan(&cl.Id, &cl.Object.Id)
	if err != nil {
		return client.ClientObject{}, err
	}

	cl.Object, err = r.getRequestObject(ctx, cl.Object.Id)
	if err != nil {
		return client.ClientObject{}, err
	}

	return cl, nil
}
func (r *repository) getRequestStatus(ctx context.Context, id string) (rs requests.ReqStatus, err error) {
	query := `
	SELECT 
		name, color
	FROM 
		public."Request_status"
	WHERE 
		name = $1
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	rows := r.client.QueryRow(ctx, query, id)
	err = rows.Scan(&rs.Name, &rs.Color)
	if err != nil {
		return requests.ReqStatus{}, err
	}

	return rs, nil
}

func (r *repository) getRequestObject(ctx context.Context, id int) (obj objects.Object, err error) {
	query := `
	SELECT 
		id, name, address, work_schedule
	FROM 
		public."Objects"
	WHERE 
		id = $1
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	rows := r.client.QueryRow(ctx, query, id)
	err = rows.Scan(&obj.Id, &obj.Name, &obj.Address, &obj.WorkSchedule)
	if err != nil {
		return objects.Object{}, err
	}

	return obj, nil
}
