package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Benyam-S/aseri/entity"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// HandleJobAttribute is a handler func that hanles a request for accessing the job attribute values
func (handler *StaffHandler) HandleJobAttribute(w http.ResponseWriter, r *http.Request) {
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
		PageState: pageState, URI: "staff/attribute"}

	handler.temp.ExecuteTemplate(w, "attribute.html", templateOutput)
	return
}

// +++++++++++++++++++++++++++++++++++++ ADD JOB ATTRIBUTE +++++++++++++++++++++++++++++++++++++

// HandleAddJobAttribute is a handler that handles the request for adding job attribute to the system
func (handler *StaffHandler) HandleAddJobAttribute(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	_, ok := ctx.Value(entity.Key("staff_member")).(*entity.Staff)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	category := mux.Vars(r)["category"]
	newAttribute := new(entity.JobAttribute)
	newAttribute.Name = strings.TrimSpace(r.FormValue("name"))
	tableName := ""

	switch category {
	case "job_sector":
		tableName = "job_sectors"
	case "job_type":
		tableName = "job_types"
	case "education_level":
		tableName = "education_levels"
	default:
		output, _ := json.MarshalIndent(map[string]string{"error": "unknown category"}, "", "\t")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	err := handler.cmService.ValidateJobAttribute(tableName, newAttribute)
	if err != nil {
		output, _ := json.MarshalIndent(map[string]string{"error": err.Error()}, "", "\t")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	err = handler.cmService.AddJobAttribute(newAttribute, tableName)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	output, _ := json.Marshal(newAttribute)
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}

// +++++++++++++++++++++++++++++++++++++ VIEW JOB ATTRIBUTES +++++++++++++++++++++++++++++++++++++

// HandleViewJobAttributes is a handler func that provides all the job attribute fields found according to their types
func (handler *StaffHandler) HandleViewJobAttributes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, ok := ctx.Value(entity.Key("staff_member")).(*entity.Staff)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	attributeContainer := make(map[string][]*entity.JobAttribute)

	// Getting initial data
	attributeContainer["JobTypes"] = handler.cmService.AllJobAttributes("job_types")
	attributeContainer["JobSectors"] = handler.cmService.AllJobAttributes("job_sectors")
	attributeContainer["EducationLevels"] = handler.cmService.AllJobAttributes("education_levels")

	output, _ := json.MarshalIndent(attributeContainer, "", "\t")
	w.Write(output)
}

// +++++++++++++++++++++++++++++++++++++ REMOVE JOB ATTRIBUTE DATA +++++++++++++++++++++++++++++++++++++

// HandleDeleteJobAttribute is a handler that handles the request for deleting a certain job attribute
func (handler *StaffHandler) HandleDeleteJobAttribute(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	_, ok := ctx.Value(entity.Key("staff_member")).(*entity.Staff)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	id := mux.Vars(r)["id"]
	category := mux.Vars(r)["category"]
	tableName := ""

	switch category {
	case "job_sector":
		tableName = "job_sectors"
	case "job_type":
		tableName = "job_types"
	case "education_level":
		tableName = "education_levels"
	default:
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	attribute, err := handler.cmService.DeleteJobAttribute(id, tableName)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	output, err := json.MarshalIndent(attribute, "", "\t")
	w.Write(output)
}
