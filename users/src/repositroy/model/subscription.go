package model

import "github.com/google/uuid"

const SubscriptionsTableName = "subscriptions"

type Subscription struct {
	SubscriberID uuid.UUID `db:"subscriber_id"`
	UserID       uuid.UUID `db:"user_id"`
}

func (Subscription) TableName() string {
	return SubscriptionsTableName
}