package subscription

import "github.com/Benyam-S/aseri/entity"

// ISubscriptionRepository is an interface that defines all the repository methods of a job subscription struct
type ISubscriptionRepository interface {
	Create(newSubscription *entity.Subscription) error
	Find(identifier string) (*entity.Subscription, error)
	FindMultiple(identifier string) []*entity.Subscription
	Total() int64
	Match(jobTypes, jobSectors, educationLevels, experiences []string) []*entity.Subscription
	Update(subscription *entity.Subscription) error
	Delete(identifier string) (*entity.Subscription, error)
	DeleteMultiple(identifier string) []*entity.Subscription
}
