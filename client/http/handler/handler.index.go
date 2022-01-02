package handler

import (
	"html/template"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Benyam-S/aseri/common"
	"github.com/Benyam-S/aseri/entity"
	"github.com/Benyam-S/aseri/feedback"
	"github.com/Benyam-S/aseri/job"
	"github.com/Benyam-S/aseri/jobapplication"
	"github.com/Benyam-S/aseri/log"
	"github.com/Benyam-S/aseri/password"
	"github.com/Benyam-S/aseri/staff"
	"github.com/Benyam-S/aseri/subscription"
	"github.com/Benyam-S/aseri/tools"
	"github.com/Benyam-S/aseri/user"
)

// StaffHandler is a type that defines a staff handler
type StaffHandler struct {
	sync.Mutex
	stService staff.IService
	pwService password.IService
	urService user.IService
	jbService job.IService
	jaService jobapplication.IService
	sbService subscription.IService
	fdService feedback.IService
	cmService common.IService
	temp      *template.Template
	store     tools.IStore
	pushChan  chan string
	pq        common.IPushQueue
	logger    *log.Logger
}

// TemplateContainer is a struct that holds every thing that gets inside a template
type TemplateContainer struct {
	StaffMember    *entity.Staff
	CSRFToken      string
	DisplaySection string
	PageState      string
	URI            string // URI is used for setting user's current uri address
}

// Initiator is a struct that holds a job initiator
type Initiator struct {
	ID          string
	UserName    string
	PhoneNumber string
	Category    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ReferralWithVisitCount struct {
	Referral   *entity.Referral
	VisitCount int64
}

// NewStaffHandler is a func that returns a new staff handler type
func NewStaffHandler(staffService staff.IService,
	passwordService password.IService, userService user.IService,
	jobService job.IService, jobApplicationService jobapplication.IService,
	subscriptionService subscription.IService, feedbackService feedback.IService,
	commonService common.IService, template *template.Template,
	store tools.IStore, pushChannel chan string, pushQueue common.IPushQueue, log *log.Logger) *StaffHandler {
	return &StaffHandler{stService: staffService, pwService: passwordService,
		urService: userService, jbService: jobService, jaService: jobApplicationService,
		sbService: subscriptionService, fdService: feedbackService,
		cmService: commonService, temp: template, store: store, pq: pushQueue,
		pushChan: pushChannel, logger: log}
}

// IsRequestAllowed is a function that checks if a given request is allowed to be performed by the staff member
func IsRequestAllowed(r *http.Request, staffMember *entity.Staff) bool {

	uri := r.RequestURI

	switch {

	case strings.Contains(uri, "/admin/"):
		return staffMember.Role == entity.RoleAdmin

	case strings.Contains(uri, "/staff/"):
		return true

	case uri == "/":
		return true
	}

	return false
}
