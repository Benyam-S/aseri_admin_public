package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Benyam-S/aseri/entity"
	"github.com/google/uuid"
)

// HandleDashboard is a handler func that hanles a request for access the dashboard
func (handler *StaffHandler) HandleDashboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	staffMember, ok := ctx.Value(entity.Key("staff_member")).(*entity.Staff)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	pageState := r.URL.Query().Get("page-state")
	displaySection := r.URL.Query().Get("display-section")

	// storing the csrf token
	csrfToken := uuid.Must(uuid.NewRandom())
	handler.store.Add(csrfToken.String(), staffMember.ID)

	templateOutput := TemplateContainer{StaffMember: staffMember,
		CSRFToken: csrfToken.String(), DisplaySection: displaySection,
		PageState: pageState, URI: "staff/dashboard"}

	handler.temp.ExecuteTemplate(w, "dashboard.html", templateOutput)
	return
}

// HandleInitDashboard is a handler func that provides all the needed data for the dashboard
func (handler *StaffHandler) HandleInitDashboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, ok := ctx.Value(entity.Key("staff_member")).(*entity.Staff)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	dashboardContainer := make(map[string]interface{})

	// Getting job data
	totalUsers := make(map[string]interface{})
	totalUsers["aseri"] = handler.urService.TotalUsers(entity.UserCategoryAseri)
	totalUsers["agent"] = handler.urService.TotalUsers(entity.UserCategoryAgent)
	totalUsers["job_seeker"] = handler.urService.TotalUsers(entity.UserCategoryJobSeeker)
	totalUsers["total"] = handler.urService.TotalUsers(entity.UserCategoryAny)

	// Getting user data
	totalJobs := make(map[string]interface{})
	totalJobs["pending"] = handler.jbService.TotalJobs(entity.JobStatusPending)
	totalJobs["opened"] = handler.jbService.TotalJobs(entity.JobStatusOpened)
	totalJobs["closed"] = handler.jbService.TotalJobs(entity.JobStatusClosed)
	totalJobs["declined"] = handler.jbService.TotalJobs(entity.JobStatusDecelined)
	totalJobs["total"] = handler.jbService.TotalJobs(entity.JobStatusAny)

	// Getting subscription data
	totalSubscribers := handler.sbService.TotalSubscribers()

	// Getting all users
	allUsers := handler.urService.AllUsers()

	// Getting all jobs
	allJobs := handler.jbService.AllJobs()

	dashboardContainer["jobs_total"] = totalJobs
	dashboardContainer["users_total"] = totalUsers
	dashboardContainer["subscribers_total"] = totalSubscribers
	dashboardContainer["all_users"] = allUsers
	dashboardContainer["all_jobs"] = allJobs

	// Getting referral data
	referralsWithVistCount := make([]*ReferralWithVisitCount, 0)
	referrals := handler.cmService.AllReferrals()
	referralCount := handler.cmService.CountAllReferralVists()

	for _, referral := range referrals {
		count := handler.cmService.CountReferralVists(referral.Code)
		referralsWithVistCount = append(referralsWithVistCount,
			&ReferralWithVisitCount{Referral: referral, VisitCount: count})
	}

	dashboardContainer["referrals"] = referralsWithVistCount
	dashboardContainer["referral_count"] = referralCount

	output, _ := json.MarshalIndent(dashboardContainer, "", "\t")
	w.Write(output)
}
