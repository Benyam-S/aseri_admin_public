package repository

import (
	"fmt"
	"strings"

	"github.com/Benyam-S/aseri/entity"
	"github.com/Benyam-S/aseri/subscription"
	"github.com/Benyam-S/aseri/tools"
	"github.com/jinzhu/gorm"
)

// SubscriptionRepository is a type that defines a job subscription repository type
type SubscriptionRepository struct {
	conn *gorm.DB
}

// NewSubscriptionRepository is a function that creates a new job subscription repository type
func NewSubscriptionRepository(connection *gorm.DB) subscription.ISubscriptionRepository {
	return &SubscriptionRepository{conn: connection}
}

// Create is a method that adds a new job subscription to the database
func (repo *SubscriptionRepository) Create(newSubscription *entity.Subscription) error {
	totalNumOfSubscriptions := tools.CountMembers("subscriptions", repo.conn)
	newSubscription.ID = fmt.Sprintf("SB-%s%d", tools.RandomStringGN(7), totalNumOfSubscriptions+1)

	for !tools.IsUnique("id", newSubscription.ID, "subscriptions", repo.conn) {
		totalNumOfSubscriptions++
		newSubscription.ID = fmt.Sprintf("SB-%s%d", tools.RandomStringGN(7), totalNumOfSubscriptions+1)
	}

	err := repo.conn.Create(newSubscription).Error
	if err != nil {
		return err
	}
	return nil
}

// Find is a method that finds a certain job subscription from the database using an identifier,
// also Find() uses only id as a key for selection
func (repo *SubscriptionRepository) Find(identifier string) (*entity.Subscription, error) {

	subscription := new(entity.Subscription)
	err := repo.conn.Model(subscription).Where("id = ?", identifier).First(subscription).Error

	if err != nil {
		return nil, err
	}
	return subscription, nil
}

// FindMultiple is a method that finds multiple job subscriptions from the database the matches the given identifier
// In FindMultiple() only user_id is used as a key
func (repo *SubscriptionRepository) FindMultiple(identifier string) []*entity.Subscription {

	var subscriptions []*entity.Subscription
	err := repo.conn.Model(entity.Subscription{}).Where("user_id = ?", identifier).Find(&subscriptions).Error

	if err != nil {
		return []*entity.Subscription{}
	}
	return subscriptions
}

// Total is a method that retruns the total number of subscribers for job push notifications
func (repo *SubscriptionRepository) Total() int64 {

	subscribers := make([]string, 0)
	repo.conn.Raw("SELECT DISTINCT user_id FROM subscriptions").Scan(&subscribers)
	return int64(len(subscribers))
}

// Match is a method that finds a multiple subscriptions that match the given job type and sector
func (repo *SubscriptionRepository) Match(jobSectors, jobTypes, educationLevels, experiences []string) []*entity.Subscription {

	var whereStmt1 []string
	var whereStmt2 []string
	var whereStmt3 []string
	var whereStmt4 []string
	var sqlValues []interface{}

	subscriptions := make([]*entity.Subscription, 0)

	for _, jobSector := range jobSectors {
		whereStmt1 = append(whereStmt1, " sector = ? ")
		sqlValues = append(sqlValues, jobSector)
	}

	for _, jobType := range jobTypes {
		whereStmt2 = append(whereStmt2, " type = ? ")
		sqlValues = append(sqlValues, jobType)
	}

	for _, educationLevel := range educationLevels {
		whereStmt3 = append(whereStmt3, " education_level = ? ")
		sqlValues = append(sqlValues, educationLevel)
	}

	for _, experience := range experiences {
		whereStmt4 = append(whereStmt4, " experience = ? ")
		sqlValues = append(sqlValues, experience)
	}

	repo.conn.Raw("SELECT * FROM subscriptions WHERE ("+strings.Join(whereStmt1, "||")+") && ("+strings.Join(whereStmt2, "||")+") && ("+strings.Join(whereStmt3, "||")+") && ("+strings.Join(whereStmt4, "||")+") ", sqlValues...).Scan(&subscriptions)
	return subscriptions
}

// Update is a method that updates a certain job subscription entries in the database
func (repo *SubscriptionRepository) Update(subscription *entity.Subscription) error {

	prevSubscription := new(entity.Subscription)
	err := repo.conn.Model(prevSubscription).Where("id = ?", subscription.ID).First(prevSubscription).Error

	if err != nil {
		return err
	}

	/* --------------------------- can change layer if needed --------------------------- */
	subscription.CreatedAt = prevSubscription.CreatedAt
	/* -------------------------------------- end --------------------------------------- */

	err = repo.conn.Save(subscription).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete is a method that deletes a certain job subscription from the database using an identifier.
// In Delete() id is only used as an key
func (repo *SubscriptionRepository) Delete(identifier string) (*entity.Subscription, error) {
	subscription := new(entity.Subscription)
	err := repo.conn.Model(subscription).Where("id = ?", identifier).First(subscription).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(subscription)
	return subscription, nil
}

// DeleteMultiple is a method that deletes a set of job subscriptions from the database using an identifier.
// In Delete() user_id is only used as an key
func (repo *SubscriptionRepository) DeleteMultiple(identifier string) []*entity.Subscription {
	var subscriptions []*entity.Subscription
	repo.conn.Model(subscriptions).Where("user_id = ?", identifier).
		Find(&subscriptions)

	for _, subscription := range subscriptions {
		repo.conn.Delete(subscription)
	}

	return subscriptions
}
