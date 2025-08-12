package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

var theParking *parking

func main() {
	theParking = newParking() //new empety parking
	fmt.Println(theParking)

	chiRouter := chi.NewRouter()
	chiRouter.Post("/parkingEntery", handleParkingEntery)
	chiRouter.Post("/parkingExit", handleParkingExit)

	server := &http.Server{
		Addr:    ":3000",
		Handler: chiRouter,
	}

	if err := server.ListenAndServe(); err != nil {
		fmt.Println("there was an error while serving and listening!")
	}
}

func handleParkingEntery(responseWriter http.ResponseWriter, request *http.Request) {
	type parkingEntery struct {
		Car_type           string `json:"car_type"`
		VehiclePlateNumber string `json:"vehiclePlateNumber"`
	}

	var parkingEnteryData parkingEntery
	err := json.NewDecoder(request.Body).Decode(&parkingEnteryData)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
	} else {
		theParking.processCustomer(parkingEnteryData.Car_type, parkingEnteryData.VehiclePlateNumber)
		fmt.Print(theParking)
	}
}

func handleParkingExit(responseWriter http.ResponseWriter, request *http.Request) {
	type parkingExitRequestData struct {
		VehiclePlateNumber string `json:"vehiclePlateNumber"`
	}

	var parkingExitData parkingExitRequestData
	err := json.NewDecoder(request.Body).Decode(&parkingExitData)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
	} else {
		parkingBillData := &parkingBill{VehiclePlateNumber: parkingExitData.VehiclePlateNumber}
		theParking.exitParking(parkingBillData)

		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)
		json.NewEncoder(responseWriter).Encode(parkingBillData)

		fmt.Print(theParking)
	}
}
