package main

import (
	"net/http"
)

func (cfg *apiConfig) resetHits(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cfg.fileserverHits.Store(0)
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Platform is not set to dev, can't remove all users."))
		return
	}
	err := cfg.db.DeleteUsers(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error deleting users: " + err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0 and users table cleared."))
}
