package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"parking_app/store"

	"github.com/go-chi/chi/v5"
)

func (app *application) mount() *chi.Mux {
	chiRouter := chi.NewRouter()

	chiRouter.Get("/stat", app.handleStatusRequest)

	chiRouter.Post("/parkingEntery", app.handleParkingEntery)
	chiRouter.Post("/parkingExit", app.handleParkingExit)

	chiRouter.Delete("/deleteVehicle", app.handleDeleteTicket)

	fmt.Println("")
	return chiRouter
}

func (app *application) handleStatusRequest(responseWriter http.ResponseWriter, request *http.Request) {
	parkingStat := app.getParkingStat(request.Context())

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(responseWriter).Encode(parkingStat)
}

func (app *application) handleParkingEntery(responseWriter http.ResponseWriter, request *http.Request) {
	var vehicleData store.Vehicle
	err := json.NewDecoder(request.Body).Decode(&vehicleData)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
	} else {
		ticket := app.processCustomer(&vehicleData, request.Context())
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)
		json.NewEncoder(responseWriter).Encode(ticket)
	}
}

func (app *application) handleParkingExit(responseWriter http.ResponseWriter, request *http.Request) {
	var parkingBillData store.ParkingBill
	err := json.NewDecoder(request.Body).Decode(&parkingBillData)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
	} else {
		app.exitParking(&parkingBillData, request.Context())
		
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)
		json.NewEncoder(responseWriter).Encode(parkingBillData)
	}
}

func (app *application) handleDeleteTicket(responseWriter http.ResponseWriter, request *http.Request) {
	vehiclePlateNumber := request.URL.Query().Get("vehiclePlateNumber")
	fmt.Println(vehiclePlateNumber)
	result := app.deleteTicket(vehiclePlateNumber, request.Context())
	responseWriter.Header().Set("Content-Type", "text/plain")
	responseWriter.Write([]byte(result))
}