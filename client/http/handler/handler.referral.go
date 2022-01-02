package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Benyam-S/aseri/entity"
	"github.com/gorilla/mux"
)

// HandleNewVisitReferral is a handler func that handles a request for updating referral for new visit
func (handler *StaffHandler) HandleNewVisitReferral(w http.ResponseWriter, r *http.Request) {

	// Allowing CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")

	referralCode := mux.Vars(r)["code"]
	referral, err := handler.cmService.FindReferral(referralCode)
	if err != nil {
		output, _ := json.MarshalIndent(map[string]string{"error": err.Error()}, "", "\t")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	referralVist := &entity.ReferralVisit{Code: referral.Code}
	err = handler.cmService.AddReferralVisit(referralVist)
	if err != nil {
		output, _ := json.MarshalIndent(map[string]string{"error": err.Error()}, "", "\t")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	output, _ := json.Marshal(referral)
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}
