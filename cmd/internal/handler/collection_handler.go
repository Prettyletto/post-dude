package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Prettyletto/post-dude/cmd/internal/model"
	"github.com/Prettyletto/post-dude/cmd/internal/service"
)

type CollectionHandler struct {
	service service.CollectionService
}

func NewCollectionHandler(service service.CollectionService) *CollectionHandler {
	return &CollectionHandler{service: service}
}

func (h *CollectionHandler) CreateCollectionHandler(w http.ResponseWriter, r *http.Request) {
	var collection model.Collection
	if err := json.NewDecoder(r.Body).Decode(&collection); err != nil {
		http.Error(w, "Error in the payload", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateCollection(&collection); err != nil {
		http.Error(w, fmt.Sprintf("Failed to create collection: %v", err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *CollectionHandler) GetAllCollectionsHandler(w http.ResponseWriter, r *http.Request) {
	collections, err := h.service.RetrieveAllCollections()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve collections: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(collections); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *CollectionHandler) GetCollectionHandler(w http.ResponseWriter, r *http.Request) {
	id, atoierr := strconv.Atoi(r.PathValue("id"))
	if atoierr != nil {
		http.Error(w, fmt.Sprintf("Failed in the id request: %v", atoierr), http.StatusBadRequest)
		return
	}

	collection, err := h.service.RetrieveCollectionById(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve collection: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(collection); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *CollectionHandler) UpdateCollectionHandler(w http.ResponseWriter, r *http.Request) {
	id, atoierr := strconv.Atoi(r.PathValue("id"))
	if atoierr != nil {
		http.Error(w, fmt.Sprintf("Failed in the id request: %v", atoierr), http.StatusBadRequest)
		return
	}

	var newCollection model.Collection
	if err := json.NewDecoder(r.Body).Decode(&newCollection); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode payload: %v", err), http.StatusBadRequest)
	}

	updated, err := h.service.UpdateCollection(id, &newCollection)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update collection: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	updated.ID = id
	if err := json.NewEncoder(w).Encode(updated); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode collection: %v", err), http.StatusInternalServerError)
		return
	}

}

func (h *CollectionHandler) DeleteCollectionHandler(w http.ResponseWriter, r *http.Request) {
	id, atoierr := strconv.Atoi(r.PathValue("id"))
	if atoierr != nil {
		http.Error(w, fmt.Sprintf("failed in convert id: %v", atoierr), http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteCollection(id); err != nil {
		http.Error(w, fmt.Sprintf("failed to delete collection: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)

}
