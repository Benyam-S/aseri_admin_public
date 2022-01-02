package service

import (
	"errors"
	"regexp"
	"strings"

	"github.com/Benyam-S/aseri/common"
	"github.com/Benyam-S/aseri/entity"
	"github.com/Benyam-S/aseri/subscription"
)

// Service is a type that defines a job subscription service
type Service struct {
	subscriptionRepo subscription.ISubscriptionRepository
	cmService        common.IService
}

// NewSubscriptionService is a function that returns a new job subscription service
func NewSubscriptionService(subscriptionRepository subscription.ISubscriptionRepository,
	commonService common.IService) subscription.IService {
	return &Service{subscriptionRepo: subscriptionRepository, cmService: commonService}
}

// AddSubscription is a method that adds a new job subscription to the system
func (service *Service) AddSubscription(newSubscription *entity.Subscription) error {

	err := service.subscriptionRepo.Create(newSubscription)
	if err != nil {
		return errors.New("unable to add new subscription")
	}

	return nil
}

// ValidateSubscription is a method that validates a job subscription entries.
// It checks if the subscription has a valid entries or not and return map of errors if any.
func (service *Service) ValidateSubscription(subscription *entity.Subscription) entity.ErrMap {

	errMap := make(map[string]error)

	var isValidJobType bool
	var isValidJobSector bool
	var isValidEducationLevel bool
	var isValidWorkExperience bool

	validJobTypes := service.cmService.GetValidJobTypesForSubscription()
	validJobSectors := service.cmService.GetValidJobSectorsForSubscription()
	validEducationLevels := service.cmService.GetValidEducationLevelsForSubscription()
	validWorkExperiences := service.cmService.GetValidWorkExperiencesForSubscription()

	for _, validJobType := range validJobTypes {
		if subscription.Type == validJobType.ID || subscription.Type == validJobType.Name {
			isValidJobType = true
			// Lastly saving the job type name
			subscription.Type = validJobType.Name
			break
		}
	}

	for _, validJobSector := range validJobSectors {
		if subscription.Sector == validJobSector.ID || subscription.Sector == validJobSector.Name {
			isValidJobSector = true
			// Lastly saving the job sector name
			subscription.Sector = validJobSector.Name
			break
		}
	}

	for _, validEducationLevel := range validEducationLevels {
		if subscription.EducationLevel == validEducationLevel.ID ||
			subscription.EducationLevel == validEducationLevel.Name {
			isValidEducationLevel = true
			// Lastly saving the education level name
			subscription.EducationLevel = validEducationLevel.Name
			break
		}
	}

	for _, validWorkExperience := range validWorkExperiences {
		if subscription.Experience == validWorkExperience {
			isValidWorkExperience = true
			break
		}
	}

	if !isValidJobType {
		errMap["type"] = errors.New("invalid job type used")
	}

	if !isValidJobSector {
		errMap["sector"] = errors.New("invalid job sector used")
	}

	if !isValidEducationLevel {
		errMap["education_level"] = errors.New("invalid education level used")
	}

	if !isValidWorkExperience {
		errMap["experience"] = errors.New("invalid work experience used")
	}

	if isValidJobSector && isValidJobType && isValidEducationLevel && isValidWorkExperience {
		prevSubscriptions := service.subscriptionRepo.FindMultiple(subscription.UserID)

		for _, prevSubscription := range prevSubscriptions {
			if prevSubscription.ID != subscription.ID &&
				strings.ToLower(prevSubscription.Sector) == strings.ToLower(strings.TrimSpace(subscription.Sector)) &&
				strings.ToLower(prevSubscription.Type) == strings.ToLower(strings.TrimSpace(subscription.Type)) &&
				strings.ToLower(prevSubscription.EducationLevel) == strings.ToLower(strings.TrimSpace(subscription.EducationLevel)) &&
				strings.ToLower(prevSubscription.Experience) == strings.ToLower(strings.TrimSpace(subscription.Experience)) {

				errMap["error"] = errors.New("subscription already exists")
				break
			}
		}
	}

	if len(errMap) > 0 {
		return errMap
	}

	return nil
}

// FindSubscription is a method that find and return a job subscription that matches the id value
func (service *Service) FindSubscription(id string) (*entity.Subscription, error) {

	empty, _ := regexp.MatchString(`^\s*$`, id)
	if empty {
		return nil, errors.New("no subscription found")
	}

	subscription, err := service.subscriptionRepo.Find(id)
	if err != nil {
		return nil, errors.New("no subscription found")
	}
	return subscription, nil
}

// FindMultipleSubscriptions is a method that find and return multiple job subscriptions that matchs the userID value
func (service *Service) FindMultipleSubscriptions(userID string) []*entity.Subscription {

	empty, _ := regexp.MatchString(`^\s*$`, userID)
	if empty {
		return []*entity.Subscription{}
	}

	return service.subscriptionRepo.FindMultiple(userID)
}

// FindSubscriptionMatch is a method that find and return multiple subscription that matches the given job sector and typ
func (service *Service) FindSubscriptionMatch(jobSector, jobType, educationLevel, experience string) []*entity.Subscription {

	emptySector, _ := regexp.MatchString(`^\s*$`, jobSector)
	if emptySector {
		return []*entity.Subscription{}
	}

	jobSectors := strings.Split(strings.TrimSpace(jobSector), ",")
	jobTypes := append(strings.Split(strings.TrimSpace(jobType), ","), "Any")
	educationLevels := []string{educationLevel, "Any"}
	experiences := []string{experience, "Any"}

	// If empty job type then match it for subscribes with 'Any' subscription
	emptyType, _ := regexp.MatchString(`^\s*$`, jobType)
	if emptyType {
		jobTypes = []string{"Any"}
	}

	// If empty education level then match it for subscribes with 'Any' subscription
	emptyEducationLevel, _ := regexp.MatchString(`^\s*$`, jobType)
	if emptyEducationLevel {
		educationLevels = []string{"Any"}
	}

	// If empty work experience then match it for subscribes with 'Any' subscription
	emptyExperience, _ := regexp.MatchString(`^\s*$`, jobType)
	if emptyExperience {
		experiences = []string{"Any"}
	}

	subscriptions := service.subscriptionRepo.Match(jobSectors, jobTypes, educationLevels, experiences)
	uniqueSubscribers := make([]*entity.Subscription, 0)

	for _, subscription := range subscriptions {
		isAdded := false
		for _, uniqueSubscriber := range uniqueSubscribers {
			if uniqueSubscriber.UserID == subscription.UserID {
				isAdded = true
			}
		}

		if !isAdded {
			uniqueSubscribers = append(uniqueSubscribers, subscription)
		}
	}

	return uniqueSubscribers
}

// TotalSubscribers is a method that returns the total number of subscribers for bot push notifications
func (service *Service) TotalSubscribers() int64 {
	return service.subscriptionRepo.Total()
}

// UpdateSubscription is a method that updates a job subscription in the system
func (service *Service) UpdateSubscription(subscription *entity.Subscription) error {
	err := service.subscriptionRepo.Update(subscription)
	if err != nil {
		return errors.New("unable to update subscription")
	}

	return nil
}

// DeleteSubscription is a method that deletes a job subscription from the system using an id
func (service *Service) DeleteSubscription(id string) (*entity.Subscription, error) {

	subscription, err := service.subscriptionRepo.Delete(id)
	if err != nil {
		return nil, errors.New("unable to delete subscription")
	}

	return subscription, nil
}

// DeleteMultipleSubscriptions is a method that deletes multiple job subscriptions from the system that match the given userID
func (service *Service) DeleteMultipleSubscriptions(userID string) []*entity.Subscription {
	return service.subscriptionRepo.DeleteMultiple(userID)
}
