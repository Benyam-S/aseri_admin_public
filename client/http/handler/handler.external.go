package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Benyam-S/aseri/entity"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// HandleInitAddJobExternal is a handler that gets all the needed data for adding new job from external system
func (handler *StaffHandler) HandleInitAddJobExternal(w http.ResponseWriter, r *http.Request) {

	// Allowing CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")

	userID := mux.Vars(r)["user_id"]
	accessToken := mux.Vars(r)["access_token"]

	// Invalid csrf token or session doesn't match
	if handler.store.Get(accessToken) != userID {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
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
		CSRFToken            string
	}

	// storing the csrf token
	csrfToken := uuid.Must(uuid.NewRandom())
	handler.store.Add(csrfToken.String(), userID)

	container := InitJobContainer{ValidJobSectors: validJobSectors, ValidJobTypes: validJobTypes,
		ValidEducationLevels: validEducationLevels, ValidWorkExperiences: validWorkExperiences,
		ContactTypes: contactTypes, CSRFToken: csrfToken.String()}

	output, _ := json.MarshalIndent(container, "", "\t")
	w.Write(output)
}

// HandleAddJobExternal is a handler that handles the request for adding job to the system from external system
func (handler *StaffHandler) HandleAddJobExternal(w http.ResponseWriter, r *http.Request) {

	// Allowing CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")

	newJob := new(entity.Job)
	newJob.Title = strings.TrimSpace(r.FormValue("title"))
	newJob.Description = strings.TrimSpace(r.FormValue("description"))
	newJob.Employer = strings.TrimSpace(r.FormValue("employer_id"))
	newJob.InitiatorID = strings.TrimSpace(r.FormValue("employer_id"))
	newJob.Type = strings.TrimSpace(r.FormValue("job_type"))
	newJob.Sector = strings.TrimSpace(r.FormValue("job_sector"))
	newJob.EducationLevel = strings.TrimSpace(r.FormValue("education_level"))
	newJob.Experience = strings.TrimSpace(r.FormValue("work_experience"))
	newJob.Gender = strings.TrimSpace(r.FormValue("gender"))
	newJob.ContactType = strings.TrimSpace(r.FormValue("contact_type"))
	dueDateString := strings.TrimSpace(r.FormValue("due_date"))

	if dueDateString != "" {
		layout := "2006-01-02T15:04:05.000Z"
		dueDate, err := time.Parse(layout, dueDateString)
		if err == nil {
			newJob.DueDate = &dueDate
		}
	}

	newJob.Status = entity.JobStatusPending

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
}
