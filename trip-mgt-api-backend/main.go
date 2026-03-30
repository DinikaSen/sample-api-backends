package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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
		{TripID: 101, Name: "New York City Break", StartDate: "2025-08-15", EndDate: "2025-08-18"},
		{TripID: 102, Name: "Safari Adventure Nairobi", StartDate: "2025-11-10", EndDate: "2025-11-20"},
	},
	2: {
		{TripID: 201, Name: "Miami Beach Getaway", StartDate: "2025-09-05", EndDate: "2025-09-08"},
		{TripID: 202, Name: "Holiday in Maldives", StartDate: "2025-12-20", EndDate: "2025-12-30"},
	},
	3: {
		{TripID: 301, Name: "Washington DC Cultural Tour", StartDate: "2025-10-10", EndDate: "2025-10-14"},
		{TripID: 302, Name: "Conference Tokyo", StartDate: "2026-03-05", EndDate: "2026-03-09"},
	},
	4: {
		{TripID: 401, Name: "Miami Art Basel", StartDate: "2025-12-04", EndDate: "2025-12-08"},
		{TripID: 402, Name: "Ski Trip Aspen", StartDate: "2026-01-05", EndDate: "2026-01-10"},
	},
	5: {
		{TripID: 501, Name: "Washington DC Cherry Blossom", StartDate: "2026-03-28", EndDate: "2026-04-02"},
		{TripID: 502, Name: "Honeymoon Paris", StartDate: "2025-06-10", EndDate: "2025-06-20"},
	},
}

var nextTripSeq = map[int64]int64{
	1: 3,
	2: 3,
	3: 3,
	4: 3,
	5: 3,
}

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

	log.Printf("GET trips: guestId=%d", guestID)

	trips, ok := mockTrips[guestID]
	if !ok {
		trips = []Trip{}
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

	seq := nextTripSeq[guestID]
	trip.TripID = guestID*100 + seq
	nextTripSeq[guestID] = seq + 1
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
			writeJSON(w, http.StatusOK, updated)
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
