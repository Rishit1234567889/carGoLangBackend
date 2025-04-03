package car

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rishit1234567889/carZone/models"
	"github.com/rishit1234567889/carZone/service"
	"go.opentelemetry.io/otel"
)

type CarHandler struct {
	service service.CarServiceInterface
}

func NewCarHandler(service service.CarServiceInterface) *CarHandler {
	return &CarHandler{
		service: service,
	}
}

func (h *CarHandler) GetCarById(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	tracer := otel.Tracer("CarHandler")

	ctx, span := tracer.Start(r.Context(), "GetCarByID-Handler")

	defer span.End()

	vars := mux.Vars(r)
	id := vars["id"]

	resp, err := h.service.GetCarById(ctx, id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		span.RecordError(err)

		log.Println("Error : ", err)
		return
	}

	body, err := json.Marshal(resp)

	if err != nil {
		span.RecordError(err)
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error : ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	//write the response body
	_, err = w.Write(body)

	if err != nil {
		span.RecordError(err)
		log.Println("Error Writing Response: ", err)
	}
}

func (h *CarHandler) GetCarByBrand(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("CarHandler")

	ctx, span := tracer.Start(r.Context(), "GetCarByBrand-Handler")

	defer span.End()

	vars := mux.Vars(r)
	brand := vars["brand"]
	isEngine := r.URL.Query().Get("isEngine") == "true"

	cars, err := h.service.GetCarByBrand(ctx, brand, isEngine)
	if err != nil {
		span.RecordError(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cars)
}

func (h *CarHandler) CreateCar(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("CarHandler")

	ctx, span := tracer.Start(r.Context(), "CreateCar-Handler")

	defer span.End()

	var carReq models.CarRequest
	if err := json.NewDecoder(r.Body).Decode(&carReq); err != nil {
		span.RecordError(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdCar, err := h.service.CreateCar(ctx, &carReq)
	if err != nil {
		span.RecordError(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdCar)
}

func (h *CarHandler) UpdateCar(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("CarHandler")

	ctx, span := tracer.Start(r.Context(), "UpdateCar-Handler")

	defer span.End()
	vars := mux.Vars(r)
	id := vars["id"]

	var carReq models.CarRequest
	if err := json.NewDecoder(r.Body).Decode(&carReq); err != nil {
		span.RecordError(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	carID, err := uuid.Parse(id) // Capture both return values
	if err != nil {
		span.RecordError(err)
		http.Error(w, "Invalid ID format", http.StatusBadRequest) // Handle parsing error
		return
	}

	updatedCar, err := h.service.UpdateCar(ctx, carID, &carReq)
	if err != nil {
		span.RecordError(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedCar)
}

func (h *CarHandler) DeleteCar(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("CarHandler")

	ctx, span := tracer.Start(r.Context(), "DeleteCar-Handler")

	defer span.End()
	vars := mux.Vars(r)
	id := vars["id"]

	deletedCar, err := h.service.DeleteCar(ctx, id)
	if err != nil {
		span.RecordError(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deletedCar)
}
