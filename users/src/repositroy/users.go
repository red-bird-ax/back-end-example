package repositroy

import (
	"fmt"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/google/uuid"

	"github.com/red-bird-ax/poster/users/src/repositroy/model"
	"github.com/red-bird-ax/poster/utils/data"
)

type Users struct {
	connection *dbx.DB
}

func NewUsers(connection *dbx.DB) Users {
	return Users{connection: connection}
}

func (repo Users) Create(user model.User) error {
	return repo.connection.Model(&user).Insert()
}

func (repo Users) Get(id uuid.UUID) (*model.User, error) {
	user := new(model.User)
	err := repo.connection.Select().Model(id, user)
	return user, err
}

func (repo Users) GetAll(options *data.Options) ([]model.User, error) {
	users := make([]model.User, 0)
	query := repo.connection.Select()

	if options != nil {
		query = query.Offset(options.Pagination.Offset).Limit(options.Pagination.Limit).OrderBy(options.OrderBy)
	}

	err := query.All(&users)
	return users, err
}

func (repo Users) GetByName(name string) (*model.User, error) {
	user := new(model.User)
	err := repo.connection.Select().Where(dbx.HashExp{"user_name": name}).One(user)
	return user, err
}

func (repo Users) SearchFor(queryText string, options *data.Options) ([]model.User, error) {
	users := make([]model.User, 0)
	query := fmt.Sprintf(
		"SELECT * FROM %s WHERE LOWER(user_name) LIKE LOWER('%%%s%%') OR LOWER(full_name) LIKE LOWER('%%%s%%')",
		model.UsersTableName,
		queryText,
		queryText,
	)
	if options != nil {
		query += fmt.Sprintf(
			" ORDER BY %s OFFSET %d LIMIT %d",
			options.OrderBy,
			options.Pagination.Offset,
			options.Pagination.Limit,
		)
	}
	query += ";"

	err := repo.connection.NewQuery(query).All(&users)
	return users, err
}

func (repo Users) Update(user model.User) error {
	return repo.connection.Model(&user).Update()
}

func (repo Users) Delete(id uuid.UUID) error {
	return repo.connection.Model(&model.User{ID: id}).Delete()
}