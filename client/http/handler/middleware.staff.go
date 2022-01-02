package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/gorilla/mux"

	"github.com/Benyam-S/aseri/client/http/session"
	"github.com/Benyam-S/aseri/entity"
)

// SessionAuthentication is a middleware that validates a request cookie contain a valid aseri session value
func (handler *StaffHandler) SessionAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(os.Getenv(entity.AppCookieName))
		if err != nil {
			handler.HandleLogin(w, r)
			return
		}

		signedString := cookie.Value
		clientSession, err := session.Extract(signedString)
		if err != nil {
			handler.HandleLogin(w, r)
			return
		}

		// Creating a new response write so it can be passed to the next handler
		newW := httptest.NewRecorder()

		// Adding the client session to the context
		ctx := context.WithValue(r.Context(), entity.Key("client_session"), clientSession)
		r = r.WithContext(ctx)

		next(newW, r)

		// Obtaining all the header key value pairs from the previously created response writer
		for k, v := range newW.HeaderMap {
			w.Header()[k] = v
		}

		content := newW.Body.Bytes()

		// Meaning we didn't request Set-Cookie from the inside handlers
		if len(w.Header()["Set-Cookie"]) == 0 {
			clientSession.ExpiresAt = time.Now().Add(time.Hour * 240).Unix()
			clientSession.UpdatedAt = time.Now().Unix()
			clientSession.Save(w)
		}

		w.WriteHeader(newW.Code)
		w.Write(content)

	}
}

// Authorization is a middleware that authorize a given session has a valid	staff member
func (handler *StaffHandler) Authorization(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		clientSession, ok := ctx.Value(entity.Key("client_session")).(*session.ClientSession)

		if !ok {
			handler.HandleLogin(w, r)
			return

		}

		serverSession, err := handler.stService.FindSession(clientSession.SessionID)
		if err != nil {
			handler.HandleLogin(w, r)
			return
		}

		// Deactivated means the user has logged out from this session
		if serverSession.Deactivated {
			handler.HandleLogin(w, r)
			return
		}

		// If you want you can notify the user
		if serverSession.DeviceInfo != r.UserAgent() {
			handler.HandleLogin(w, r)
			return
		}

		staffMember, err := handler.stService.FindStaffMember(serverSession.UserID)
		if err != nil {
			handler.HandleLogin(w, r)
			return
		}

		if !IsRequestAllowed(r, staffMember) {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		ctx = context.WithValue(ctx, entity.Key("staff_member"), staffMember)
		r = r.WithContext(ctx)

		// Updating serverSessions for better user experience
		handler.stService.UpdateSession(serverSession)

		next(w, r)
	}
}

// CSRFAuthentication is a middleware that validates a request has a valid csrf token
func (handler *StaffHandler) CSRFAuthentication(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		staffMember, ok := ctx.Value(entity.Key("staff_member")).(*entity.Staff)
		if !ok {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		csrf := mux.Vars(r)["csrf"]

		// Invalid csrf token or session doesn't match
		if handler.store.Get(csrf) != staffMember.ID {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		next(w, r)
	}
}

// CSRFAuthenticationExternal is a middleware that validates a request from external system has a valid csrf token
func (handler *StaffHandler) CSRFAuthenticationExternal(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Access-Control-Allow-Origin", "*")

		csrf := mux.Vars(r)["csrf"]
		userID := mux.Vars(r)["user_id"]

		// Invalid csrf token or session doesn't match
		if handler.store.Get(csrf) != userID {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		next(w, r)
	}
}
