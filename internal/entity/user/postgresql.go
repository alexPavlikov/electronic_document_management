package user

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

// GetUser implements Repository.
func (r *repository) SelectUser(ctx context.Context, id string) (User, error) {
	query := `
		SELECT 
			id, email, full_name, phone, image, role
		FROM 
			User 
		WHERE 
			id = $1`

	r.logger.Tracef("SQL Query: %s", utils.FormatQuery(query))

	rows, err := r.client.Query(ctx, query, id)
	if err != nil {
		return User{}, err
	}

	var us User

	err = rows.Scan(&us.Id, &us.Email, &us.FullName, &us.Phone, &us.Image, &us.Role)
	if err != nil {
		return User{}, err
	}
	return us, nil
}

// GetUsers implements Repository.
func (r *repository) SelectUsers(ctx context.Context) ([]User, error) {
	query := `
	SELECT 
		id, email, full_name, phone, image, role 
	FROM 
		User`

	r.logger.Tracef("SQL Query: %s", utils.FormatQuery(query))

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var us User
	var users []User
	for rows.Next() {
		err = rows.Scan(&us.Id, &us.Email, &us.FullName, &us.Phone, &us.Image, &us.Role)
		if err != nil {
			return nil, err
		}
		users = append(users, us)
	}
	return users, nil
}

// GetUsersByParameter implements Repository.
func (r *repository) SelectUsersByParameter(ctx context.Context) ([]User, error) {
	panic("unimplemented")
}

// InsertUser implements Repository.
func (r *repository) InsertUser(ctx context.Context, user User) error {
	panic("unimplemented")
}

// UpdateUser implements Repository.
func (r *repository) UpdateUser(ctx context.Context, user *User) error {
	panic("unimplemented")
}

// DeleteUser implements Repository.
func (r *repository) DeleteUser(ctx context.Context, id string) error {
	panic("unimplemented")
}

func NewRepository(client dbClient.Client, logger *logging.Logger) Repository {
	return &repository{
		client: client,
		logger: *logger,
	}
}
