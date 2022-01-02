package handler

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Benyam-S/aseri/entity"
	"github.com/Benyam-S/aseri/tools"
	"golang.org/x/crypto/bcrypt"
)

// HandleViewProfile is a handler that enables user to staff member view it's profile
func (handler *StaffHandler) HandleViewProfile(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	staffMember, ok := ctx.Value(entity.Key("staff_member")).(*entity.Staff)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	output, _ := json.MarshalIndent(staffMember, "", "\t")
	w.Write(output)
}

// HandleUploadPhoto is a handler func that handles a request for uploading profile pic for staff member
func (handler *StaffHandler) HandleUploadPhoto(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	staffMember, ok := ctx.Value(entity.Key("staff_member")).(*entity.Staff)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
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

// HandleUpdateProfile is a handler func that handles a request for updating staff member's profile
func (handler *StaffHandler) HandleUpdateProfile(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	staffMember, ok := ctx.Value(entity.Key("staff_member")).(*entity.Staff)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
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

	err := handler.stService.UpdateStaffMember(staffMember)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

// HandleChangePassword is a handler func that handles a request for changing staff member's passwords
func (handler *StaffHandler) HandleChangePassword(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	staffMember, ok := ctx.Value(entity.Key("staff_member")).(*entity.Staff)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	currentPassword := r.FormValue("current_password")
	newPassword := r.FormValue("new_password")
	vPassword := r.FormValue("new_vPassword")

	// admin password is used for access since the staff member old password is unknown
	memberPassword, err := handler.pwService.FindPassword(staffMember.ID)
	if err != nil {
		output, _ := json.Marshal(map[string]string{"current_password": err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	hasedPassword, _ := base64.StdEncoding.DecodeString(memberPassword.Password)
	err = bcrypt.CompareHashAndPassword(hasedPassword, []byte(currentPassword+memberPassword.Salt))
	if err != nil {
		output, _ := json.Marshal(map[string]string{"current_password": "invalid admin password used"})
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

// HandleRemovePhoto is a handler func that handles the request for removing staff members's profile
func (handler *StaffHandler) HandleRemovePhoto(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	staffMember, ok := ctx.Value(entity.Key("staff_member")).(*entity.Staff)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	err := handler.stService.UpdateStaffMemberSingleValue(staffMember.ID, "profile_pic", "")
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

// HandleDeleteAccount is a method that handles the request for deleting a staff member
func (handler *StaffHandler) HandleDeleteAccount(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	staffMember, ok := ctx.Value(entity.Key("staff_member")).(*entity.Staff)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	password := r.FormValue("password")

	// Checking if the password of the given staff member exists, it may seem redundant but it will prevent from null point exception
	memberPassword, err := handler.pwService.FindPassword(staffMember.ID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Comparing the hashed password with the given password
	hasedPassword, _ := base64.StdEncoding.DecodeString(memberPassword.Password)
	err = bcrypt.CompareHashAndPassword(hasedPassword, []byte(password+memberPassword.Salt))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	_, err = handler.stService.DeleteStaffMember(staffMember.ID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
