package lib

import (
	"encoding/json"
	"net/http"
)

func HandleLatestLaunch(client *SpaceXClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		launch, err := client.GetLatestLaunch()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(launch)
	}
}

func HandleRocket(client *SpaceXClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rocketID := r.URL.Query().Get("id")
		if rocketID == "" {
			http.Error(w, "rocket ID is required", http.StatusBadRequest)
			return
		}

		rocket, err := client.GetRocket(rocketID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rocket)
	}
}

func HandleListRockets(client *SpaceXClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rockets, err := client.GetAllRockets()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rockets)
	}
}

func HandleRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		endpoints := map[string]string{
			"/":                  "Shows this list of available endpoints",
			"/api/latest-launch": "Get the latest SpaceX launch",
			"/api/rocket":        "Get a specific rocket by ID (use ?id=[rocket_id])",
			"/api/rockets":       "Get a list of all SpaceX rockets",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(endpoints)
	}
}
