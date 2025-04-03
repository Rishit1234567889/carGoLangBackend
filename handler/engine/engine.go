package engine

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

type EngineHandler struct {
	service service.EngineServiceInterface
}

func NewEngineHandler(service service.EngineServiceInterface) *EngineHandler {
	return &EngineHandler{
		service: service,
	}
}

func (e *EngineHandler) GetEngineById(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("EngineHandler")

	ctx, span := tracer.Start(r.Context(), "GetEngineById-Handler")

	defer span.End()
	vars := mux.Vars(r)

	id := vars["id"]

	resp, err := e.service.GetEngineById(ctx, id)

	if err != nil {
		span.RecordError(err)
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	body, err := json.Marshal(resp)
	if err != nil {
		span.RecordError(err)
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(body)

	if err != nil {
		span.RecordError(err)
		log.Println("error writing response :", err)
	}

}

func (e *EngineHandler) CreateEngine(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("EngineHandler")

	ctx, span := tracer.Start(r.Context(), "CreateEngine-Handler")

	defer span.End()
	var engineReq models.EngineRequest
	if err := json.NewDecoder(r.Body).Decode(&engineReq); err != nil {
		span.RecordError(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdEngine, err := e.service.CreateEngine(ctx, &engineReq)
	if err != nil {
		span.RecordError(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(createdEngine)
	if err != nil {
		span.RecordError(err)
		http.Error(w, "Error marshalling  response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(body)
	if err != nil {
		span.RecordError(err)
		log.Println("Error writing response :", err)
	}
}

// UpdateEngine handles PUT requests to update an existing engine.
func (e *EngineHandler) UpdateEngine(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("EngineHandler")

	ctx, span := tracer.Start(r.Context(), "UpdateEngine-Handler")

	defer span.End()

	vars := mux.Vars(r)
	id := vars["id"]

	var engineReq models.EngineRequest
	if err := json.NewDecoder(r.Body).Decode(&engineReq); err != nil {
		span.RecordError(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	engineID, err := uuid.Parse(id) // Capture both return values
	if err != nil {
		span.RecordError(err)
		http.Error(w, "Invalid ID format", http.StatusBadRequest) // Handle parsing error
		return
	}

	updatedEngine, err := e.service.UpdateEngine(ctx, engineID, &engineReq)
	if err != nil {
		span.RecordError(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(updatedEngine)
	if err != nil {
		span.RecordError(err)
		http.Error(w, "Error marshalling response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(body)
	if err != nil {
		span.RecordError(err)
		log.Println("Error writing response:", err)
	}
}

// DeleteEngine handles DELETE requests to remove an engine by its ID.
func (e *EngineHandler) DeleteEngine(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("EngineHandler")

	ctx, span := tracer.Start(r.Context(), "DeleteEngine-Handler")

	defer span.End()

	vars := mux.Vars(r)
	id := vars["id"]

	deletedEngine, err := e.service.DeleteEngine(ctx, id)
	if err != nil {
		span.RecordError(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(deletedEngine)
	if err != nil {
		span.RecordError(err)
		http.Error(w, "Error marshalling response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(body)
	if err != nil {
		span.RecordError(err)
		log.Println("Error writing response:", err)
	}
}
