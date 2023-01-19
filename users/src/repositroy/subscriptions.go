package repositroy

import (
	"database/sql"
	
	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/google/uuid"

	"github.com/red-bird-ax/poster/users/src/repositroy/model"
	"github.com/red-bird-ax/poster/utils/data"
)

type Subscriptions struct {
	connection *dbx.DB
}

func NewSubscriptions(connection *dbx.DB) Subscriptions {
	return Subscriptions{connection: connection}
}

func (repo Subscriptions) Create(subscription model.Subscription) error {
	return repo.connection.Model(&subscription).Insert()
}

func (repo Subscriptions) GetByUser(userID uuid.UUID, options *data.Options) ([]model.Subscription, error) {
	subscriptions := make([]model.Subscription, 0)
	query := repo.connection.Select().Where(dbx.HashExp{"user_id": userID})

	if options != nil {
		query = query.Offset(options.Pagination.Offset).Limit(options.Pagination.Limit)
	}

	err := query.All(&subscriptions)
	return subscriptions, err
}

func (repo Subscriptions) GetBySubscriber(subscriberID uuid.UUID, options *data.Options) ([]model.Subscription, error) {
	subscriptions := make([]model.Subscription, 0)
	query := repo.connection.Select().Where(dbx.HashExp{"subscriber_id": subscriberID})

	if options != nil {
		query = query.Offset(options.Pagination.Offset).Limit(options.Pagination.Limit)
	}

	err := query.All(&subscriptions)
	return subscriptions, err
}

func (repo Subscriptions) GetAll(options *data.Options) ([]model.Subscription, error) {
	subscriptions := make([]model.Subscription, 0)
	query := repo.connection.Select()

	if options != nil {
		query = query.Offset(options.Pagination.Offset).Limit(options.Pagination.Limit)
	}

	err := query.All(&subscriptions)
	return subscriptions, err
}

func (repo Subscriptions) Has(subscription model.Subscription) (bool, error) {
	var sub model.Subscription
	err := repo.connection.
		Select().
		Where(dbx.HashExp{"subscriber_id": subscription.SubscriberID, "user_id": subscription.UserID}).
		One(&sub)

	switch err {
	case sql.ErrNoRows:
		return false, nil
	case nil:
		return true, nil
	default:
		return false, err
	}
}

func (repo Subscriptions) Delete(subscription model.Subscription) error {
	_, err := repo.connection.Delete(
		model.SubscriptionsTableName,
		dbx.HashExp{
			"subscriber_id": subscription.SubscriberID,
			"user_id": subscription.UserID,
		},
	).Execute()
	return err
}