package handler

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Benyam-S/aseri/entity"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// HandleUser is a handler func that handles a request for viewing user related options
func (handler *StaffHandler) HandleUser(w http.ResponseWriter, r *http.Request) {
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
		PageState: pageState, URI: "staff/user"}

	handler.temp.ExecuteTemplate(w, "user.html", templateOutput)
	return
}

// +++++++++++++++++++++++++++++++++++++ VIEW USER +++++++++++++++++++++++++++++++++++++

// HandleViewUser is a handler that enables staff members to view a certain user's profile
func (handler *StaffHandler) HandleViewUser(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	_, ok := ctx.Value(entity.Key("staff_member")).(*entity.Staff)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	id := mux.Vars(r)["id"]

	user, err := handler.urService.FindUser(id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	output, err := json.MarshalIndent(user, "", "\t")
	w.Write(output)
}

// HandleAllUsers is a handler that enables staff members to view all users with pagenation
func (handler *StaffHandler) HandleAllUsers(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	_, ok := ctx.Value(entity.Key("staff_member")).(*entity.Staff)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	category := r.URL.Query().Get("category")

	pagination := r.URL.Query().Get("page")
	pageNum, _ := strconv.ParseInt(pagination, 0, 0)

	users, pageCount := handler.urService.AllUsersWithPagination(category, pageNum)
	output, _ := json.MarshalIndent(map[string]interface{}{
		"Result": users, "CurrentPage": pageNum, "PageCount": pageCount}, "", "\t")
	w.Write(output)
}

// HandleSearchUsers is a handler that enables staff members to search for users with pagenation
func (handler *StaffHandler) HandleSearchUsers(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	_, ok := ctx.Value(entity.Key("staff_member")).(*entity.Staff)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	key := r.URL.Query().Get("key")
	category := r.URL.Query().Get("category")
	pagination := r.URL.Query().Get("page")

	pageNum, _ := strconv.ParseInt(pagination, 0, 0)

	users, pageCount := handler.urService.SearchUsers(key, category, pageNum)
	output, _ := json.MarshalIndent(map[string]interface{}{
		"Result": users, "CurrentPage": pageNum, "PageCount": pageCount}, "", "\t")
	w.Write(output)
}

// +++++++++++++++++++++++++++++++++++++ REMOVE USER DATA +++++++++++++++++++++++++++++++++++++

// HandleDeleteUser is a handler that handles the request for deleting a certain user
func (handler *StaffHandler) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	currentStaffMember, ok := ctx.Value(entity.Key("staff_member")).(*entity.Staff)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	id := mux.Vars(r)["id"]
	password := r.FormValue("password")

	// Checking if the password of the given user exists, it may seem redundant but it will prevent from null point exception
	userPassword, err := handler.pwService.FindPassword(currentStaffMember.ID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Comparing the hashed password with the given password
	hasedPassword, _ := base64.StdEncoding.DecodeString(userPassword.Password)
	err = bcrypt.CompareHashAndPassword(hasedPassword, []byte(password+userPassword.Salt))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	user, err := handler.urService.DeleteUser(id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	output, err := json.MarshalIndent(user, "", "\t")
	w.Write(output)
}
