package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"math"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()
var receipts = make(map[string]Receipt)

func uuidValidateV4(uuidStr string) bool {
	id, err := uuid.Parse(uuidStr)
	return err == nil && id.Version() == 4
}

type Item struct {
	ShortDescription string `json:"shortDescription" validate:"required"`
	Price            string `json:"price" validate:"required,regex=^[0-9]+\.[0-9]{2}$"`
}

type Receipt struct {
	Retailer     string `json:"retailer" validate:"required"`
	PurchaseDate string `json:"purchaseDate" validate:"required"`
	PurchaseTime string `json:"purchaseTime" validate:"required,regex=^([01]\d|2[0-3]):([0-5]\d)$"`
	Items        []Item `json:"items" validate:"required,min=1,dive"`
	Total        string `json:"total" validate:"required,regex=^[0-9]+\.[0-9]{2}$"`
}

func storeReceipt(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&receipt)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := validate.Struct(receipt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	receiptID := uuid.New().String()
	receipts[receiptID] = receipt

	response := map[string]interface{}{"id": receiptID}
	json.NewEncoder(w).Encode(response)
}

func calculatePoints(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	receiptID := params["id"]

	if !uuidValidateV4(receiptID) {
		http.Error(w, "id is not valid", http.StatusBadRequest)
		return
	}

	receipt, ok := receipts[receiptID]
	if !ok {
		http.Error(w, "No receipt found for that id", http.StatusNotFound)
		return
	}

	points := getPoints(receipt)

	response := map[string]int{"points": points}
	json.NewEncoder(w).Encode(response)
}

func getPoints(receipt Receipt) int {
	totalPoints := 0

	// Rule 1: One point for every alphanumeric character in the retailer name.
    totalPoints += len(strings.Join(strings.FieldsFunc(receipt.Retailer, func(r rune) bool {
        return !unicode.IsLetter(r) && !unicode.IsNumber(r)
    }), ""))

    // Rule 2: 50 points if the total is a round dollar amount with no cents.
    if val, err := strconv.ParseFloat(receipt.Total, 64); err == nil && int(val)%1 == 0 {
        totalPoints += 50
    }

	// Rule 3: 25 points if the total is a multiple of 0.25.
	if val, err := strconv.ParseFloat(receipt.Total, 64); err == nil && int(val)%25 == 0 {
		totalPoints += 25
	}

	// Rule 4: 5 points for every two items on the receipt.
	totalPoints += len(receipt.Items) / 2 * 5

	// Rule 5: If the trimmed length of the item description is a multiple of 3,
	// multiply the price by 0.2 and round up to the nearest integer.
	for _, item := range receipt.Items {
		trimmedLength := len(strings.TrimSpace(item.ShortDescription))
		if trimmedLength%3 == 0 {
			pointsForItem := int(math.Ceil(strToFloat64(item.Price) * 0.2))
			totalPoints += pointsForItem
		}
	}

	// Rule 6: 6 points if the day in the purchase date is odd.
	purchaseDay := getDayFromDateString(receipt.PurchaseDate)
	if purchaseDay%2 != 0 {
		totalPoints += 6
	}

	// Rule 7: 10 points if the time of purchase is after 2:00 pm and before 4:00 pm.
	purchaseHour := getTimeFromTimeString(receipt.PurchaseTime)
	if purchaseHour > 14 && purchaseHour < 16 {
		totalPoints += 10
	}

	return totalPoints
}

func strToFloat64(s string) float64 {
	val, _ := strconv.ParseFloat(s, 64)
	return val
}

func getDayFromDateString(dateString string) int {
	re := regexp.MustCompile(`^(\d{4})-(\d{2})-(\d{2})$`)
	match := re.FindStringSubmatch(dateString)
	day, _ := strconv.Atoi(match[3])
	return day
}

func getTimeFromTimeString(timeString string) int {
	re := regexp.MustCompile(`^(\d{2}):(\d{2})$`)
	match := re.FindStringSubmatch(timeString)
	hour, _ := strconv.Atoi(match[1])
	return hour
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/receipts/process", storeReceipt).Methods("POST")
	r.HandleFunc("/receipts/{id}/points", calculatePoints).Methods("GET")

	port := 3000
	fmt.Printf("Server is running on http://localhost:%d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}
