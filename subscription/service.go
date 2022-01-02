package subscription

import "github.com/Benyam-S/aseri/entity"

// IService is an interface that defines all the service methods of a job subscription struct
type IService interface {
	AddSubscription(newSubscription *entity.Subscription) error
	ValidateSubscription(subscription *entity.Subscription) entity.ErrMap
	FindSubscription(id string) (*entity.Subscription, error)
	FindMultipleSubscriptions(userID string) []*entity.Subscription
	FindSubscriptionMatch(jobType, jobSector, educationLevel, experience string) []*entity.Subscription
	TotalSubscribers() int64
	UpdateSubscription(subscription *entity.Subscription) error
	DeleteSubscription(id string) (*entity.Subscription, error)
	DeleteMultipleSubscriptions(userID string) []*entity.Subscription
}
