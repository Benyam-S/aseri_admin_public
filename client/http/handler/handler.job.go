package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Benyam-S/aseri/entity"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// HandleJob is a handler func that handles a request for viewing job related options
func (handler *StaffHandler) HandleJob(w http.ResponseWriter, r *http.Request) {
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
		PageState: pageState, URI: "staff/job"}

	handler.temp.ExecuteTemplate(w, "job.html", templateOutput)
	return
}

// +++++++++++++++++++++++++++++++++++++ ADD JOB +++++++++++++++++++++++++++++++++++++

// HandleInitAddJob is a handler that gets all the needed data for adding new job
func (handler *StaffHandler) HandleInitAddJob(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	_, ok := ctx.Value(entity.Key("staff_member")).(*entity.Staff)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	validJobSectors := handler.cmService.GetValidJobSectors()
	validJobTypes := handler.cmService.GetValidJobTypes()
	validEducationLevels := handler.cmService.GetValidEducationLevels()
	validWorkExperiences := handler.cmService.GetValidWorkExperiences()
	contactTypes := handler.cmService.GetValidContactTypes()

	type InitJobContainer struct {
		ValidJobSectors      []*entity.JobAttribute
		ValidJobTypes        []*entity.JobAttribute
		ValidEducationLevels []*entity.JobAttribute
		ValidWorkExperiences []string
		ContactTypes         []string
	}

	container := InitJobContainer{ValidJobSectors: validJobSectors, ValidJobTypes: validJobTypes,
		ValidEducationLevels: validEducationLevels, ValidWorkExperiences: validWorkExperiences,
		ContactTypes: contactTypes}

	output, _ := json.MarshalIndent(container, "", "\t")
	w.Write(output)
}

// HandleAddJob is a handler that handles the request for adding job to the system
func (handler *StaffHandler) HandleAddJob(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	staffMember, ok := ctx.Value(entity.Key("staff_member")).(*entity.Staff)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	newJob := new(entity.Job)
	newJob.Title = strings.TrimSpace(r.FormValue("title"))
	newJob.Description = strings.TrimSpace(r.FormValue("description"))
	newJob.Employer = strings.TrimSpace(r.FormValue("employer"))
	newJob.Type = strings.TrimSpace(r.FormValue("job_type"))
	newJob.Sector = strings.TrimSpace(r.FormValue("job_sector"))
	newJob.EducationLevel = strings.TrimSpace(r.FormValue("education_level"))
	newJob.Experience = strings.TrimSpace(r.FormValue("work_experience"))
	newJob.Gender = strings.TrimSpace(r.FormValue("gender"))
	newJob.ContactType = strings.TrimSpace(r.FormValue("contact_type"))
	newJob.PostType = strings.TrimSpace(r.FormValue("post_type"))
	newJob.Link = strings.TrimSpace(r.FormValue("link"))
	newJob.ContactInfo = strings.TrimSpace(r.FormValue("contact_info"))
	dueDateString := strings.TrimSpace(r.FormValue("due_date"))

	if dueDateString != "" {
		layout := "2006-01-02T15:04:05.000Z"
		dueDate, err := time.Parse(layout, dueDateString)
		if err == nil {
			newJob.DueDate = &dueDate
		}
	}

	newJob.Status = entity.JobStatusOpened
	newJob.InitiatorID = staffMember.ID

	errMap := handler.jbService.ValidateJob(newJob)
	if errMap != nil {
		output, _ := json.Marshal(errMap.StringMap())
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	err := handler.jbService.AddJob(newJob)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	output, _ := json.Marshal(newJob)
	w.WriteHeader(http.StatusOK)
	w.Write(output)

	// Adding data to push queue
	request := new(entity.ChannelRequest)
	request.Type = entity.PushForApproval
	request.Value = newJob.ID

	handler.pq.AddToQueue(request)
	handler.pushChan <- entity.StartPush
}

// HandleAddMultipleJobs is a handler that handles the request for adding multiple jobs at once to the system
func (handler *StaffHandler) HandleAddMultipleJobs(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	staffMember, ok := ctx.Value(entity.Key("staff_member")).(*entity.Staff)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	jobsErrMap := make(map[int64]map[string]string)
	newJobs := make(map[int64]*entity.Job)

	jsonJobs := r.FormValue("jobs")

	err := json.Unmarshal([]byte(jsonJobs), &newJobs)
	if err != nil {
		jobsErrMap[0] = map[string]string{"error": "unable to parse jobs"}
		output, _ := json.Marshal(jobsErrMap)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	for key, newJob := range newJobs {
		// Triming white space
		newJob.Title = strings.TrimSpace(newJob.Title)
		newJob.Description = strings.TrimSpace(newJob.Description)
		newJob.Employer = strings.TrimSpace(newJob.Employer)
		newJob.EducationLevel = strings.TrimSpace(newJob.EducationLevel)
		newJob.Experience = strings.TrimSpace(newJob.Experience)
		newJob.Type = strings.TrimSpace(newJob.Type)
		newJob.Sector = strings.TrimSpace(newJob.Sector)
		newJob.Gender = strings.TrimSpace(newJob.Gender)
		newJob.Link = strings.TrimSpace(newJob.Link)
		newJob.ContactInfo = strings.TrimSpace(newJob.ContactInfo)
		newJob.ContactType = strings.TrimSpace(newJob.ContactType)
		newJob.PostType = strings.TrimSpace(newJob.PostType)
		// Due date will be automatically parsed from json
		newJob.InitiatorID = staffMember.ID

		errMap := handler.jbService.ValidateJob(newJob)
		if errMap != nil {
			jobsErrMap[key] = errMap.StringMap()
		}
	}

	if len(jobsErrMap) != 0 {
		output, _ := json.Marshal(jobsErrMap)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	for key, newJob := range newJobs {

		newJob.Status = entity.JobStatusOpened
		err := handler.jbService.AddJob(newJob)
		if err != nil {
			jobsErrMap[key] = map[string]string{"error": err.Error()}
		}
	}

	if len(jobsErrMap) != 0 {

		output, _ := json.Marshal(jobsErrMap)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(output)

		for key, newJob := range newJobs {

			isAdded := true
			for errorKey := range jobsErrMap {
				if errorKey == key {
					isAdded = false
					break
				}
			}

			if isAdded {
				// Adding data to push queue
				request := new(entity.ChannelRequest)
				request.Type = entity.PushForApproval
				request.Value = newJob.ID

				handler.pq.AddToQueue(request)
			}
		}

		handler.pushChan <- entity.StartPush

	} else {

		output, _ := json.Marshal(newJobs)
		w.WriteHeader(http.StatusOK)
		w.Write(output)

		for _, newJob := range newJobs {

			// Adding data to push queue
			request := new(entity.ChannelRequest)
			request.Type = entity.PushForApproval
			request.Value = newJob.ID

			handler.pq.AddToQueue(request)
		}

		handler.pushChan <- entity.StartPush
	}
}

// +++++++++++++++++++++++++++++++++++++ VIEW JOB +++++++++++++++++++++++++++++++++++++

// HandleViewJob is a handler that enables staff members to view a certain job's detail
func (handler *StaffHandler) HandleViewJob(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	_, ok := ctx.Value(entity.Key("staff_member")).(*entity.Staff)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	id := mux.Vars(r)["id"]

	job, err := handler.jbService.FindJob(id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	output, _ := json.MarshalIndent(job, "", "\t")
	w.Write(output)
}

// HandleViewJobInitiator is a handler that enables staff members to view a a certain job initiator profile
func (handler *StaffHandler) HandleViewJobInitiator(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	_, ok := ctx.Value(entity.Key("staff_member")).(*entity.Staff)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	id := mux.Vars(r)["id"]
	initiator := new(Initiator)

	user, err := handler.urService.FindUser(id)
	if err != nil {
		staff, err := handler.stService.FindStaffMember(id)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		initiator.ID = staff.ID
		initiator.UserName = staff.FirstName + " " + staff.LastName
		initiator.PhoneNumber = staff.PhoneNumber
		initiator.Category = staff.Role
		initiator.CreatedAt = staff.CreatedAt
		initiator.UpdatedAt = staff.UpdatedAt
	} else {
		initiator.ID = user.ID
		initiator.UserName = user.UserName
		initiator.PhoneNumber = user.PhoneNumber
		initiator.Category = user.Category
		initiator.CreatedAt = user.CreatedAt
		initiator.UpdatedAt = user.UpdatedAt
	}

	output, err := json.MarshalIndent(initiator, "", "\t")
	w.Write(output)
}

// HandleAllJobs is a handler that enables staff members to view all jobs with pagenation
func (handler *StaffHandler) HandleAllJobs(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	_, ok := ctx.Value(entity.Key("staff_member")).(*entity.Staff)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	status := r.URL.Query().Get("status")
	pagination := r.URL.Query().Get("page")
	pageNum, _ := strconv.ParseInt(pagination, 0, 0)

	jobs, pageCount := handler.jbService.AllJobsWithPagination(status, pageNum)
	output, _ := json.MarshalIndent(map[string]interface{}{
		"Result": jobs, "CurrentPage": pageNum, "PageCount": pageCount}, "", "\t")
	w.Write(output)
}

// HandleSearchJobs is a handler that enables staff members to search for jobs with pagenation
func (handler *StaffHandler) HandleSearchJobs(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	_, ok := ctx.Value(entity.Key("staff_member")).(*entity.Staff)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	key := r.URL.Query().Get("key")
	status := r.URL.Query().Get("status")
	pagination := r.URL.Query().Get("page")

	pageNum, _ := strconv.ParseInt(pagination, 0, 0)

	jobs, pageCount := handler.jbService.SearchJobs(key, status, pageNum)
	output, _ := json.MarshalIndent(map[string]interface{}{
		"Result": jobs, "CurrentPage": pageNum, "PageCount": pageCount}, "", "\t")
	w.Write(output)
}

// +++++++++++++++++++++++++++++++++++++ UPDATE JOB +++++++++++++++++++++++++++++++++++++

// HandleChangeJobStatus is a handler that handles the request for modifying job status
func (handler *StaffHandler) HandleChangeJobStatus(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	_, ok := ctx.Value(entity.Key("staff_member")).(*entity.Staff)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	id := mux.Vars(r)["id"]
	status := r.FormValue("status")

	switch status {
	case "approve":
		status = entity.JobStatusOpened
	case "decline":
		status = entity.JobStatusDecelined
	case "close":
		status = entity.JobStatusClosed
	}

	job, err := handler.jbService.ChangeJobStatus(id, status)
	if err != nil {
		if err.Error() == "unable to update job" {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		output, _ := json.MarshalIndent(map[string]string{"error": err.Error()}, "", "\t")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return

	}

	output, _ := json.MarshalIndent(job, "", "\t")
	w.Write(output)

	// Adding data to push queue
	request := new(entity.ChannelRequest)
	request.Type = entity.PushForApproval
	request.Value = job.ID

	handler.pq.AddToQueue(request)
	handler.pushChan <- entity.StartPush
}
