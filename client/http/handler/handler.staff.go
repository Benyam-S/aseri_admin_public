package handler

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"

	"github.com/Benyam-S/aseri/entity"
	"github.com/Benyam-S/aseri/tools"
)

// HandleStaffMember is a handler func that handles a request for viewing staff member's options
func (handler *StaffHandler) HandleStaffMember(w http.ResponseWriter, r *http.Request) {
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
		PageState: pageState, URI: "admin/staff"}

	err := handler.temp.ExecuteTemplate(w, "staff.html", templateOutput)
	if err != nil {
		panic(err)
	}
	return
}

// +++++++++++++++++++++++++++++++++++++ ADDING MEMBER +++++++++++++++++++++++++++++++++++++

// HandleAddStaffMember is a handler that handles the request for adding staff member to the system
func (handler *StaffHandler) HandleAddStaffMember(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	_, ok := ctx.Value(entity.Key("staff_member")).(*entity.Staff)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	newStaffMember := new(entity.Staff)
	newStaffMember.FirstName = strings.TrimSpace(r.FormValue("first_name"))
	newStaffMember.LastName = strings.TrimSpace(r.FormValue("last_name"))
	newStaffMember.Email = strings.TrimSpace(r.FormValue("email"))
	newStaffMember.PhoneNumber = strings.TrimSpace(r.FormValue("phone_number"))
	newStaffMember.Role = entity.RoleStaff

	newPassword := new(entity.Password)
	newPassword.Password = r.FormValue("password")
	vPassword := r.FormValue("vpassword")

	// validating staff member profile and cleaning up
	errMap := handler.stService.ValidateStaffMemberProfile(newStaffMember)

	if errMap != nil {
		output, _ := json.Marshal(errMap.StringMap())
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	err := handler.pwService.VerifyPassword(newPassword, vPassword)
	if err != nil {
		output, _ := json.Marshal(map[string]string{"password": err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	err = handler.stService.AddStaffMember(newStaffMember, newPassword)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	output, _ := json.Marshal(newStaffMember)
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}

// HandleInitAddAdmin is a handler that initiate the process of adding new admin
func (handler *StaffHandler) HandleInitAddAdmin(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	_, ok := ctx.Value(entity.Key("staff_member")).(*entity.Staff)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	newStaffMember := new(entity.Staff)
	newStaffMember.FirstName = strings.TrimSpace(r.FormValue("first_name"))
	newStaffMember.LastName = strings.TrimSpace(r.FormValue("last_name"))
	newStaffMember.Email = strings.TrimSpace(r.FormValue("email"))
	newStaffMember.PhoneNumber = strings.TrimSpace(r.FormValue("phone_number"))
	newStaffMember.Role = entity.RoleAdmin

	newPassword := new(entity.Password)
	newPassword.Password = r.FormValue("password")
	vPassword := r.FormValue("vpassword")

	// validating staff member profile and cleaning up
	errMap := handler.stService.ValidateStaffMemberProfile(newStaffMember)

	if errMap != nil {
		output, _ := json.Marshal(errMap.StringMap())
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	err := handler.pwService.VerifyPassword(newPassword, vPassword)
	if err != nil {
		output, _ := json.Marshal(map[string]string{"password": err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	// Generating OTP code
	otp := tools.GenerateOTP()
	// Generating unique key for identifying the OTP token
	emailNonce := uuid.Must(uuid.NewRandom())

	// Reading message body from our asset folder
	wd, _ := os.Getwd()
	dir := filepath.Join(wd, "../../assets/messages", "/message.email.otp.json")
	data, err1 := ioutil.ReadFile(dir)

	var messageEmail map[string][]string
	err2 := json.Unmarshal(data, &messageEmail)

	if err1 != nil || err2 != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	msg := messageEmail["message_body"][0] + otp + ". " + messageEmail["message_body"][1]
	err = tools.SendEmail(os.Getenv("super_admin_email"), "Admin registration process", msg)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Saving all the data to a temporary database
	tempOutput1, err1 := json.Marshal(newStaffMember)
	tempOutput2, err2 := json.Marshal(newPassword)
	if err1 != nil || err2 != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	handler.store.Add(otp+emailNonce.String(), string(tempOutput1))
	handler.store.Add(otp+newStaffMember.Email, string(tempOutput2))

	// Sending nonce to the client with the message ID, so it can be used to retrive the otp token
	output, _ := json.Marshal(map[string]string{"nonce": emailNonce.String()})
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}

// HandleFinishAddAdmin is a handler func that handles a request for constructing admin account
func (handler *StaffHandler) HandleFinishAddAdmin(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	_, ok := ctx.Value(entity.Key("staff_member")).(*entity.Staff)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	newPassword := new(entity.Password)
	newStaffMember := new(entity.Staff)

	nonce := r.FormValue("nonce")
	otp := r.FormValue("otp")

	storedStaffMember := handler.store.Get(otp + nonce)

	// un marshaling staff member data
	err := json.Unmarshal([]byte(storedStaffMember), newStaffMember)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	storedPassword := handler.store.Get(otp + newStaffMember.Email)
	err = json.Unmarshal([]byte(storedPassword), newPassword)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = handler.stService.AddStaffMember(newStaffMember, newPassword)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Cleaning store
	handler.store.Remove(otp + nonce)
	handler.store.Remove(otp + newStaffMember.Email)

	output, _ := json.Marshal(newStaffMember)
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}

// +++++++++++++++++++++++++++++++++++++ VIEW MEMBER +++++++++++++++++++++++++++++++++++++

// HandleViewStaffMember is a handler that enables admin to view a certain staff member's profile
func (handler *StaffHandler) HandleViewStaffMember(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	_, ok := ctx.Value(entity.Key("staff_member")).(*entity.Staff)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	id := mux.Vars(r)["id"]

	staffMember, err := handler.stService.FindStaffMember(id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	output, err := json.MarshalIndent(staffMember, "", "\t")
	w.Write(output)
}

// HandleAllStaffMembers is a handler that enables admin to view all staff member's with pagenation
func (handler *StaffHandler) HandleAllStaffMembers(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	_, ok := ctx.Value(entity.Key("staff_member")).(*entity.Staff)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	role := r.URL.Query().Get("role")
	pagination := r.URL.Query().Get("page")

	pageNum, _ := strconv.ParseInt(pagination, 0, 0)

	staffMembers, pageCount := handler.stService.AllStaffMembers(role, pageNum)
	output, _ := json.MarshalIndent(map[string]interface{}{
		"Result": staffMembers, "CurrentPage": pageNum, "PageCount": pageCount}, "", "\t")
	w.Write(output)
}

// HandleSearchStaffMembers is a handler that enables admin to search for staff member's with pagenation
func (handler *StaffHandler) HandleSearchStaffMembers(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	_, ok := ctx.Value(entity.Key("staff_member")).(*entity.Staff)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	key := r.URL.Query().Get("key")
	role := r.URL.Query().Get("role")
	pagination := r.URL.Query().Get("page")

	pageNum, _ := strconv.ParseInt(pagination, 0, 0)

	staffMembers, pageCount := handler.stService.SearchStaffMembers(key, role, pageNum)
	output, _ := json.MarshalIndent(map[string]interface{}{
		"Result": staffMembers, "CurrentPage": pageNum, "PageCount": pageCount}, "", "\t")
	w.Write(output)
}

// +++++++++++++++++++++++++++++++++++++ UPDATE MEMBER +++++++++++++++++++++++++++++++++++++

// HandleUploadStaffMemberPhoto is a handler func that handles a request for uploading profile pic by admin
func (handler *StaffHandler) HandleUploadStaffMemberPhoto(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	_, ok := ctx.Value(entity.Key("staff_member")).(*entity.Staff)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	userID := mux.Vars(r)["id"]
	staffMember, err := handler.stService.FindStaffMember(userID)

	if err != nil {
		output, _ := json.Marshal(map[string]string{"error": err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	// checking for multipart form data, the image has to be sent in multipart form data
	fm, fh, err := r.FormFile("profile_pic")
	if err != nil {
		output, _ := json.Marshal(map[string]string{"error": err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}
	defer fm.Close()

	// Reading the stream
	tempFile, _ := ioutil.ReadAll(fm)
	tempFileType := http.DetectContentType(tempFile)
	newBufferReader := bytes.NewBuffer(tempFile)

	// checking if the sent file is image
	if !strings.HasPrefix(tempFileType, "image") {
		output, _ := json.MarshalIndent(map[string]string{"error": "invalid format sent"}, "", "\t")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	// checking the file sent doesn't exceed the size limit
	if fh.Size > 5000000 {
		output, _ := json.MarshalIndent(map[string]string{"error": "image exceeds the file size limit, 5MB"},
			"", "\t")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	prevFileName := staffMember.ProfilePic
	newFileName, err := tools.UploadSinglePhoto(staffMember.ID, prevFileName, newBufferReader, fh)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = handler.stService.UpdateStaffMemberSingleValue(staffMember.ID, "profile_pic", newFileName)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if prevFileName != "" {
		wd, _ := os.Getwd()
		filePath := filepath.Join(wd, "../../assets/profilepics", prevFileName)
		tools.RemoveFile(filePath)
	}

	output, _ := json.MarshalIndent(map[string]string{"file_name": newFileName}, "", "\t")
	w.Write(output)
}

// HandleUpdateStaffMember is a handler func that handles a request for updating staff member's profile by admin
func (handler *StaffHandler) HandleUpdateStaffMember(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	_, ok := ctx.Value(entity.Key("staff_member")).(*entity.Staff)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	userID := mux.Vars(r)["id"]
	staffMember, err := handler.stService.FindStaffMember(userID)

	if err != nil {
		output, _ := json.Marshal(map[string]string{"error": err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	staffMember.FirstName = strings.TrimSpace(r.FormValue("first_name"))
	staffMember.LastName = strings.TrimSpace(r.FormValue("last_name"))
	staffMember.Email = strings.TrimSpace(r.FormValue("email"))
	staffMember.PhoneNumber = strings.TrimSpace(r.FormValue("phone_number"))

	errMap := handler.stService.ValidateStaffMemberProfile(staffMember)

	if errMap != nil {
		output, _ := json.MarshalIndent(errMap.StringMap(), "", "\t")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	err = handler.stService.UpdateStaffMember(staffMember)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

// HandleChangeStaffMemberPassword is a handler func that handles a request for changing staff member's passwords by admin
func (handler *StaffHandler) HandleChangeStaffMemberPassword(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	admin, ok := ctx.Value(entity.Key("staff_member")).(*entity.Staff)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	userID := mux.Vars(r)["id"]
	staffMember, err := handler.stService.FindStaffMember(userID)

	if err != nil {
		output, _ := json.Marshal(map[string]string{"error": err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	adminPassword := r.FormValue("admin_password")
	newPassword := r.FormValue("new_password")
	vPassword := r.FormValue("new_vPassword")

	// admin password is used for access since the staff member old password is unknown
	memberPassword, err := handler.pwService.FindPassword(admin.ID)
	if err != nil {
		output, _ := json.Marshal(map[string]string{"admin_password": err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	hasedPassword, _ := base64.StdEncoding.DecodeString(memberPassword.Password)
	err = bcrypt.CompareHashAndPassword(hasedPassword, []byte(adminPassword+memberPassword.Salt))
	if err != nil {
		output, _ := json.Marshal(map[string]string{"admin_password": "invalid admin password used"})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	newMemberPassword := new(entity.Password)
	newMemberPassword.ID = staffMember.ID
	newMemberPassword.Password = newPassword

	err = handler.pwService.VerifyPassword(newMemberPassword, vPassword)
	if err != nil {
		output, _ := json.MarshalIndent(map[string]string{"password": err.Error()}, "", "\t")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	err = handler.pwService.UpdatePassword(newMemberPassword)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

// +++++++++++++++++++++++++++++++++++++ REMOVE MEMBER DATA +++++++++++++++++++++++++++++++++++++

// HandleRemoveStaffMemberPhoto is a handler func that handles the request for removing staff members's profile pic by admin
func (handler *StaffHandler) HandleRemoveStaffMemberPhoto(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	_, ok := ctx.Value(entity.Key("staff_member")).(*entity.Staff)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	userID := mux.Vars(r)["id"]
	staffMember, err := handler.stService.FindStaffMember(userID)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = handler.stService.UpdateStaffMemberSingleValue(staffMember.ID, "profile_pic", "")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if staffMember.ProfilePic != "" {
		wd, _ := os.Getwd()
		filePath := filepath.Join(wd, "../../assets/profilepics", staffMember.ProfilePic)
		tools.RemoveFile(filePath)
	}
}

// HandleDeleteStaffMember is a handler that handles the request for deleting a certain staff member
func (handler *StaffHandler) HandleDeleteStaffMember(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	currentStaffMember, ok := ctx.Value(entity.Key("staff_member")).(*entity.Staff)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	id := mux.Vars(r)["id"]
	password := r.FormValue("password")

	// Can't delete current account
	if id == currentStaffMember.ID {
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}

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

	staffMember, err := handler.stService.DeleteStaffMember(id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	output, err := json.MarshalIndent(staffMember, "", "\t")
	w.Write(output)
}
