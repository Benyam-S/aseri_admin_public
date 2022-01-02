package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Benyam-S/aseri/entity"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// HandleFeedback is a handler func that handles a request for viewing feedback related options
func (handler *StaffHandler) HandleFeedback(w http.ResponseWriter, r *http.Request) {
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
		PageState: pageState, URI: "staff/feedback"}

	handler.temp.ExecuteTemplate(w, "feedback.html", templateOutput)
	return
}

// +++++++++++++++++++++++++++++++++++++ VIEW FEEDBACK +++++++++++++++++++++++++++++++++++++

// HandleViewFeedback is a handler that enables staff members to view a certain feedback's detail
func (handler *StaffHandler) HandleViewFeedback(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	_, ok := ctx.Value(entity.Key("staff_member")).(*entity.Staff)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	id := mux.Vars(r)["id"]

	feedback, err := handler.fdService.FindFeedback(id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	output, _ := json.MarshalIndent(feedback, "", "\t")
	w.Write(output)
}

// HandleAllFeedbacks is a handler that enables staff members to view all feedbacks with pagenation
func (handler *StaffHandler) HandleAllFeedbacks(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	_, ok := ctx.Value(entity.Key("staff_member")).(*entity.Staff)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	status := r.URL.Query().Get("status")
	pagination := r.URL.Query().Get("page")
	pageNum, _ := strconv.ParseInt(pagination, 0, 0)

	feedbacks, pageCount := handler.fdService.AllFeedbacks(status, pageNum)
	output, _ := json.MarshalIndent(map[string]interface{}{
		"Result": feedbacks, "CurrentPage": pageNum, "PageCount": pageCount}, "", "\t")
	w.Write(output)
}

// HandleSearchFeedbacks is a handler that enables staff members to search for feedbacks with pagenation
func (handler *StaffHandler) HandleSearchFeedbacks(w http.ResponseWriter, r *http.Request) {

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

	feedbacks, pageCount := handler.fdService.SearchFeedbacks(key, status, pageNum)
	output, _ := json.MarshalIndent(map[string]interface{}{
		"Result": feedbacks, "CurrentPage": pageNum, "PageCount": pageCount}, "", "\t")
	w.Write(output)
}

// +++++++++++++++++++++++++++++++++++++ UPDATE FEEDBACK +++++++++++++++++++++++++++++++++++++

// HandleMarkFeedbackAsSeen is a handler that handles the request for marking feedback as seen
func (handler *StaffHandler) HandleMarkFeedbackAsSeen(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	_, ok := ctx.Value(entity.Key("staff_member")).(*entity.Staff)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	id := mux.Vars(r)["id"]

	err := handler.fdService.MarkAsSeen(id)
	if err != nil {
		if err.Error() == "unable to update feedback" {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		output, _ := json.MarshalIndent(map[string]string{"error": err.Error()}, "", "\t")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return

	}
}
