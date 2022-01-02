package entity

import (
	"net/http"
	"time"
)

// Staff is a type that defines a staff member
type Staff struct {
	ID          string `gorm:"primary_key; unique; not null"`
	FirstName   string
	LastName    string
	PhoneNumber string `gorm:"unique; not null"`
	Email       string `gorm:"unique; not null"`
	ProfilePic  string
	Role        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Password is a type that defines a user password
type Password struct {
	ID        string `gorm:"primary_key; unique; not null"`
	Password  string
	Salt      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// User is a type that defines the user group
type User struct {
	ID          string `gorm:"primary_key; unique; not null"`
	UserName    string
	PhoneNumber string `gorm:"unique; not null"`
	Category    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Job is a type that defines job to post
type Job struct {
	ID             string `gorm:"primary_key; unique; not null"`
	Employer       string
	Title          string
	Description    string `gorm:"type:text;"`
	Type           string
	Sector         string
	EducationLevel string
	Experience     string
	Gender         string
	ContactType    string
	ContactInfo    string // Only used if the PostType is 'Internal'
	Status         string
	PostType       string
	Link           string
	InitiatorID    string // To logging who created the job
	DueDate        *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// JobApplication is type that defines the relationship between job and jobseeker
// JobSeeker cannot apply for the same job twice so we use both JobID and JobSeekerID as primary key
type JobApplication struct {
	JobID       string `gorm:"primary_key"`
	JobSeekerID string `gorm:"primary_key"`
}

// Subscription is a type that defines job subscription
type Subscription struct {
	ID             string `gorm:"primary_key; unique; not null"`
	UserID         string
	Sector         string
	Type           string
	EducationLevel string
	Experience     string
	CreatedAt      time.Time
}

// Feedback is a type that defines user feedback
type Feedback struct {
	ID        string `gorm:"primary_key; unique; not null"`
	UserID    string
	Comment   string `gorm:"type:text;"`
	Seen      bool
	CreatedAt time.Time
}

// JobAttribute is a type that defines a job attribute like job type or job sector
type JobAttribute struct {
	ID   string `gorm:"primary_key; unique; not null"`
	Name string
}

// ChannelRequest is a type that defines a request that is set through a bot channel
type ChannelRequest struct {
	Type   string
	Value  string
	ChatID int64 // used for passing request in bot handler
	Extra  string
}

// Referral is a type that defines custom referral links
type Referral struct {
	Code      string `gorm:"primary_key; unique; not null"`
	Owner     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ReferralVisit is a type that defines custom referral visit stamp
type ReferralVisit struct {
	ID        int64  `gorm:"primary_key; unique; auto_increment; not null"`
	Code      string `gorm:"not null"`
	CreatedAt time.Time
}

// Key is a type that defines a key type that can be used a key value in context
type Key string

// ErrMap is a type that defines a map with string identifier and it's error
type ErrMap map[string]error

// StringMap is a method that returns string map corresponding to the ErrMap where the error type is converted to a string
func (errMap ErrMap) StringMap() map[string]string {
	stringMap := make(map[string]string)
	for key, value := range errMap {
		stringMap[key] = value.Error()
	}

	return stringMap
}

// Middleware is a type that defines a function that takes a handler func and return a new handler func type
type Middleware func(http.HandlerFunc) http.HandlerFunc
