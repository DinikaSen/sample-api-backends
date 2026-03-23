package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// --- Data models ---

type Trip struct {
	TripID    int64  `json:"tripId"`
	Name      string `json:"name"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

// --- Mock data ---

var mockTrips = map[int64][]Trip{
	1: {
		{TripID: 101, Name: "Summer Vacation", StartDate: "2025-07-01", EndDate: "2025-07-10"},
		{TripID: 102, Name: "Business Trip NYC", StartDate: "2025-08-15", EndDate: "2025-08-18"},
	},
	2: {
		{TripID: 201, Name: "Anniversary Getaway", StartDate: "2025-09-05", EndDate: "2025-09-08"},
	},
}

var nextTripID int64 = 301

// --- Helpers ---

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func parsePathInt(r *http.Request, name string) (int64, error) {
	val := r.PathValue(name)
	return strconv.ParseInt(val, 10, 64)
}

// --- Handlers ---

// GET /guests/{guestId}/trips
func listTripsHandler(w http.ResponseWriter, r *http.Request) {
	guestID, err := parsePathInt(r, "guestId")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid guestId")
		return
	}

	tripName := r.URL.Query().Get("tripName")
	log.Printf("GET trips: guestId=%d tripName=%q", guestID, tripName)

	trips, ok := mockTrips[guestID]
	if !ok {
		trips = []Trip{}
	}

	if tripName != "" {
		var filtered []Trip
		for _, t := range trips {
			if strings.Contains(strings.ToLower(t.Name), strings.ToLower(tripName)) {
				filtered = append(filtered, t)
			}
		}
		if filtered == nil {
			filtered = []Trip{}
		}
		trips = filtered
	}

	writeJSON(w, http.StatusOK, trips)
}

// POST /guests/{guestId}/trips
func createTripHandler(w http.ResponseWriter, r *http.Request) {
	guestID, err := parsePathInt(r, "guestId")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid guestId")
		return
	}

	var trip Trip
	if err := json.NewDecoder(r.Body).Decode(&trip); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	trip.TripID = nextTripID
	nextTripID++
	mockTrips[guestID] = append(mockTrips[guestID], trip)

	log.Printf("POST create trip: guestId=%d tripId=%d name=%q", guestID, trip.TripID, trip.Name)
	writeJSON(w, http.StatusCreated, trip)
}

// GET /guests/{guestId}/trips/{tripId}
func getTripHandler(w http.ResponseWriter, r *http.Request) {
	guestID, err := parsePathInt(r, "guestId")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid guestId")
		return
	}
	tripID, err := parsePathInt(r, "tripId")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid tripId")
		return
	}

	log.Printf("GET trip: guestId=%d tripId=%d", guestID, tripID)

	for _, t := range mockTrips[guestID] {
		if t.TripID == tripID {
			writeJSON(w, http.StatusOK, t)
			return
		}
	}

	writeError(w, http.StatusNotFound, "trip not found")
}

// PUT /guests/{guestId}/trips/{tripId}
func updateTripHandler(w http.ResponseWriter, r *http.Request) {
	guestID, err := parsePathInt(r, "guestId")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid guestId")
		return
	}
	tripID, err := parsePathInt(r, "tripId")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid tripId")
		return
	}

	var updated Trip
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	updated.TripID = tripID

	trips := mockTrips[guestID]
	for i, t := range trips {
		if t.TripID == tripID {
			trips[i] = updated
			mockTrips[guestID] = trips
			log.Printf("PUT update trip: guestId=%d tripId=%d", guestID, tripID)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	writeError(w, http.StatusNotFound, "trip not found")
}

// DELETE /guests/{guestId}/trips/{tripId}
func deleteTripHandler(w http.ResponseWriter, r *http.Request) {
	guestID, err := parsePathInt(r, "guestId")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid guestId")
		return
	}
	tripID, err := parsePathInt(r, "tripId")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid tripId")
		return
	}

	trips := mockTrips[guestID]
	for i, t := range trips {
		if t.TripID == tripID {
			mockTrips[guestID] = append(trips[:i], trips[i+1:]...)
			log.Printf("DELETE trip: guestId=%d tripId=%d", guestID, tripID)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	writeError(w, http.StatusNotFound, "trip not found")
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /guests/{guestId}/trips", listTripsHandler)
	mux.HandleFunc("POST /guests/{guestId}/trips", createTripHandler)
	mux.HandleFunc("GET /guests/{guestId}/trips/{tripId}", getTripHandler)
	mux.HandleFunc("PUT /guests/{guestId}/trips/{tripId}", updateTripHandler)
	mux.HandleFunc("DELETE /guests/{guestId}/trips/{tripId}", deleteTripHandler)

	addr := ":8080"
	log.Printf("Trip Management API mock server running on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
