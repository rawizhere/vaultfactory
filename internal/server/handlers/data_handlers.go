package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/tempizhere/vaultfactory/internal/shared/interfaces"
	"github.com/tempizhere/vaultfactory/internal/shared/models"
)

type DataHandler struct {
	dataService interfaces.DataService
}

func NewDataHandler(dataService interfaces.DataService) *DataHandler {
	return &DataHandler{
		dataService: dataService,
	}
}

type CreateDataRequest struct {
	Type     models.DataType `json:"type"`
	Name     string          `json:"name"`
	Metadata string          `json:"metadata"`
	Data     json.RawMessage `json:"data"`
}

type UpdateDataRequest struct {
	Name     string          `json:"name"`
	Metadata string          `json:"metadata"`
	Data     json.RawMessage `json:"data"`
}

type DataResponse struct {
	ID        string          `json:"id"`
	Type      models.DataType `json:"type"`
	Name      string          `json:"name"`
	Metadata  string          `json:"metadata"`
	Data      json.RawMessage `json:"data"`
	CreatedAt string          `json:"created_at"`
	UpdatedAt string          `json:"updated_at"`
	Version   int64           `json:"version"`
}

func (h *DataHandler) CreateData(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)

	var req CreateDataRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Type == "" || req.Name == "" {
		http.Error(w, "Type and name are required", http.StatusBadRequest)
		return
	}

	dataItem, err := h.dataService.CreateData(r.Context(), user.ID, req.Type, req.Name, req.Metadata, []byte(req.Data))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := DataResponse{
		ID:        dataItem.ID.String(),
		Type:      dataItem.Type,
		Name:      dataItem.Name,
		Metadata:  dataItem.Metadata,
		Data:      req.Data,
		CreatedAt: dataItem.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: dataItem.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		Version:   dataItem.Version,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *DataHandler) GetData(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)
	vars := mux.Vars(r)
	dataID := vars["id"]

	dataIDUUID, err := uuid.Parse(dataID)
	if err != nil {
		http.Error(w, "Invalid data ID", http.StatusBadRequest)
		return
	}

	dataItem, err := h.dataService.GetData(r.Context(), user.ID, dataIDUUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := DataResponse{
		ID:        dataItem.ID.String(),
		Type:      dataItem.Type,
		Name:      dataItem.Name,
		Metadata:  dataItem.Metadata,
		Data:      json.RawMessage("{}"),
		CreatedAt: dataItem.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: dataItem.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		Version:   dataItem.Version,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *DataHandler) GetUserData(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)

	dataType := r.URL.Query().Get("type")

	var items []*models.DataItem
	var err error

	if dataType != "" {
		items, err = h.dataService.GetUserDataByType(r.Context(), user.ID, models.DataType(dataType))
	} else {
		items, err = h.dataService.GetUserData(r.Context(), user.ID)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var responses []DataResponse
	for _, item := range items {
		response := DataResponse{
			ID:        item.ID.String(),
			Type:      item.Type,
			Name:      item.Name,
			Metadata:  item.Metadata,
			CreatedAt: item.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: item.UpdatedAt.Format("2006-01-02T15:04:05Z"),
			Version:   item.Version,
		}
		responses = append(responses, response)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

func (h *DataHandler) UpdateData(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)
	vars := mux.Vars(r)
	dataID := vars["id"]

	var req UpdateDataRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	dataIDUUID, err := uuid.Parse(dataID)
	if err != nil {
		http.Error(w, "Invalid data ID", http.StatusBadRequest)
		return
	}

	dataItem, err := h.dataService.UpdateData(r.Context(), user.ID, dataIDUUID, req.Name, req.Metadata, []byte(req.Data))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := DataResponse{
		ID:        dataItem.ID.String(),
		Type:      dataItem.Type,
		Name:      dataItem.Name,
		Metadata:  dataItem.Metadata,
		Data:      req.Data,
		CreatedAt: dataItem.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: dataItem.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		Version:   dataItem.Version,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *DataHandler) DeleteData(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)
	vars := mux.Vars(r)
	dataID := vars["id"]

	dataIDUUID, err := uuid.Parse(dataID)
	if err != nil {
		http.Error(w, "Invalid data ID", http.StatusBadRequest)
		return
	}

	if err := h.dataService.DeleteData(r.Context(), user.ID, dataIDUUID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *DataHandler) SyncData(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)

	lastSyncStr := r.URL.Query().Get("last_sync")
	if lastSyncStr == "" {
		http.Error(w, "last_sync parameter is required", http.StatusBadRequest)
		return
	}

	lastSync, err := time.Parse(time.RFC3339, lastSyncStr)
	if err != nil {
		http.Error(w, "Invalid last_sync format", http.StatusBadRequest)
		return
	}

	items, err := h.dataService.SyncData(r.Context(), user.ID, lastSync)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var responses []DataResponse
	for _, item := range items {
		response := DataResponse{
			ID:        item.ID.String(),
			Type:      item.Type,
			Name:      item.Name,
			Metadata:  item.Metadata,
			CreatedAt: item.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: item.UpdatedAt.Format("2006-01-02T15:04:05Z"),
			Version:   item.Version,
		}
		responses = append(responses, response)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}
