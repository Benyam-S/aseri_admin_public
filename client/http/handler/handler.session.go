package handler

import (
	"encoding/base64"
	"net/http"

	"github.com/Benyam-S/aseri/client/http/session"
	"github.com/Benyam-S/aseri/entity"
	"golang.org/x/crypto/bcrypt"
)

// HandleLogin is a handler func that handle login request
func (handler *StaffHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		handler.temp.ExecuteTemplate(w, "login.html", nil)
		return

	} else if r.Method == http.MethodPost {

		identifier := r.FormValue("identifier")
		password := r.FormValue("password")

		// Checking if the user exists
		staffMember, err := handler.stService.FindStaffMember(identifier)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// Checking if the password of the given user exists, it may seem redundant but it will prevent from null point exception
		userPassword, err := handler.pwService.FindPassword(staffMember.ID)
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

		// Creating client side session and saving cookie
		userSession := session.Create(staffMember.ID)
		err = userSession.Save(w)
		if err != nil {
			http.Error(w, "unable to set cookie", http.StatusInternalServerError)
			return
		}

		// Creating server side session and saving to the system
		err = handler.stService.AddSession(userSession, staffMember, r)
		if err != nil {
			http.Error(w, "unable to set session", http.StatusInternalServerError)
			return
		}
	}
}

// HandleLogout is a handler func that handles a logout request
func (handler *StaffHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	clientSession, ok := ctx.Value(entity.Key("client_session")).(*session.ClientSession)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	serverSession, err := handler.stService.FindSession(clientSession.SessionID)
	if err == nil {
		// Instead of deleting the session deactivate them for further use
		serverSession.Deactivated = true
		handler.stService.UpdateSession(serverSession)
	}

	// Removing client cookie
	clientSession.Remove(w)
}
