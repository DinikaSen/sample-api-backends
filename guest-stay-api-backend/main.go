package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// --- Data models ---

type Charge struct {
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
}

type Folio struct {
	Charges     []Charge `json:"charges"`
	TotalAmount float64  `json:"totalAmount"`
}

type Room struct {
	RoomType string  `json:"roomType"`
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
}

// --- Mock data ---

var mockFolios = map[string]Folio{
	"101-201": {
		Charges: []Charge{
			{Description: "Room charge - Night 1", Amount: 250.00, Currency: "USD"},
			{Description: "Room charge - Night 2", Amount: 250.00, Currency: "USD"},
			{Description: "Mini bar", Amount: 45.50, Currency: "USD"},
			{Description: "Spa service", Amount: 120.00, Currency: "USD"},
		},
		TotalAmount: 665.50,
	},
	"102-202": {
		Charges: []Charge{
			{Description: "Room charge - Night 1", Amount: 189.00, Currency: "USD"},
			{Description: "Restaurant - Dinner", Amount: 78.00, Currency: "USD"},
		},
		TotalAmount: 267.00,
	},
}

var mockRooms = []Room{
	{RoomType: "Standard King", Price: 189.00, Currency: "USD"},
	{RoomType: "Deluxe King", Price: 249.00, Currency: "USD"},
	{RoomType: "Junior Suite", Price: 349.00, Currency: "USD"},
	{RoomType: "Executive Suite", Price: 499.00, Currency: "USD"},
}

var mockRoomsByOfferType = map[string][]Room{
	"upsell": {
		{RoomType: "Deluxe King", Price: 249.00, Currency: "USD"},
		{RoomType: "Junior Suite", Price: 349.00, Currency: "USD"},
	},
	"upselldiscount": {
		{RoomType: "Deluxe King", Price: 199.00, Currency: "USD"},
		{RoomType: "Junior Suite", Price: 299.00, Currency: "USD"},
	},
	"smu": {
		{RoomType: "Standard King", Price: 159.00, Currency: "USD"},
		{RoomType: "Standard Double", Price: 149.00, Currency: "USD"},
	},
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

// GET /hospitality-customer/v2/guests/{guestId}/stays/{stayId}/chargesview
func getGuestFolioHandler(w http.ResponseWriter, r *http.Request) {
	guestID, err := parsePathInt(r, "guestId")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid guestId")
		return
	}
	stayID, err := parsePathInt(r, "stayId")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid stayId")
		return
	}

	key := fmt.Sprintf("%d-%d", guestID, stayID)
	folio, ok := mockFolios[key]
	if !ok {
		// Return a default folio for unknown guest/stay combos
		folio = Folio{
			Charges: []Charge{
				{Description: "Room charge", Amount: 200.00, Currency: "USD"},
			},
			TotalAmount: 200.00,
		}
	}

	log.Printf("GET folio: guestId=%d stayId=%d", guestID, stayID)
	writeJSON(w, http.StatusOK, folio)
}

// POST /hospitality-customer/v2/guests/{guestId}/stays/{stayId}/checkout
func checkoutGuestHandler(w http.ResponseWriter, r *http.Request) {
	guestID, err := parsePathInt(r, "guestId")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid guestId")
		return
	}
	stayID, err := parsePathInt(r, "stayId")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid stayId")
		return
	}

	log.Printf("POST checkout: guestId=%d stayId=%d", guestID, stayID)
	w.WriteHeader(http.StatusAccepted)
}

// GET /hospitality-customer/v2/guests/{guestId}/stays/{stayId}/rooms
func getAvailableRoomsHandler(w http.ResponseWriter, r *http.Request) {
	guestID, err := parsePathInt(r, "guestId")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid guestId")
		return
	}
	stayID, err := parsePathInt(r, "stayId")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid stayId")
		return
	}

	offerType := r.URL.Query().Get("includeOfferTypes")
	arrivalTime := r.URL.Query().Get("arrivalTime")

	log.Printf("GET rooms: guestId=%d stayId=%d offerType=%q arrivalTime=%q", guestID, stayID, offerType, arrivalTime)

	var rooms []Room
	if offerType != "" {
		filtered, ok := mockRoomsByOfferType[offerType]
		if !ok {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("unknown includeOfferTypes value: %s", offerType))
			return
		}
		rooms = filtered
	} else {
		rooms = mockRooms
	}

	writeJSON(w, http.StatusOK, rooms)
}

func main() {
	mux := http.NewServeMux()

	base := "/hospitality-customer/v2"

	mux.HandleFunc("GET "+base+"/guests/{guestId}/stays/{stayId}/chargesview", getGuestFolioHandler)
	mux.HandleFunc("POST "+base+"/guests/{guestId}/stays/{stayId}/checkout", checkoutGuestHandler)
	mux.HandleFunc("GET "+base+"/guests/{guestId}/stays/{stayId}/rooms", getAvailableRoomsHandler)

	addr := ":8080"
	log.Printf("Guest Stay API mock server running on %s", addr)
	log.Printf("Base path: %s", base)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
