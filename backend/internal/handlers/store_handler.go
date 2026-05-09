package handlers

import (
	"coffee-order/internal/gateway"
	"coffee-order/internal/models"
	"coffee-order/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// StoreHandler handles store-related requests
type StoreHandler struct {
	client *gateway.TomoroClient
}

// NewStoreHandler creates a new store handler
func NewStoreHandler(client *gateway.TomoroClient) *StoreHandler {
	return &StoreHandler{client: client}
}

// GetStoreList handles store list request
func (h *StoreHandler) GetStoreList(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	lat := r.URL.Query().Get("lat")
	lng := r.URL.Query().Get("lng")
	search := r.URL.Query().Get("search")
	pageStr := r.URL.Query().Get("page")
	sizeStr := r.URL.Query().Get("size")

	// Parse coordinates
	latitude, _ := strconv.ParseFloat(lat, 64)
	longitude, _ := strconv.ParseFloat(lng, 64)

	// Parse pagination
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	size, _ := strconv.Atoi(sizeStr)
	if size < 1 {
		size = 20
	}

	// Get auth from headers
	token := r.Header.Get("Authorization")
	if token != "" {
		// Remove "Bearer " prefix if present
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}
	}

	deviceCode := r.Header.Get("X-Device-Code")
	wToken := r.Header.Get("X-WToken")
	ucde := r.Header.Get("X-UCDE")

	// Build query string
	path := fmt.Sprintf("/portal/app/basic/storeInfo/getStoreList/v3?centerPointLatitude=%f&centerPointLongitude=%f&pageNo=%d&pageSize=%d&storeName=%s",
		latitude, longitude, page, size, search)

	// Create request context
	ctx := gateway.RequestContext{
		Token:      token,
		DeviceCode: deviceCode,
		WToken:     wToken,
		UCDE:       ucde,
		Latitude:   latitude,
		Longitude:  longitude,
	}

	// Call Tomoro API
	respBody, err := h.client.Request(ctx, "GET", path, nil)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get stores: %v", err))
		return
	}

	// Parse response
	var storeResp models.StoreListResponse
	if err := json.Unmarshal(respBody, &storeResp); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to parse store response")
		return
	}

	utils.WriteSuccess(w, storeResp.Data)
}

// GetStoreDetail handles store detail request
func (h *StoreHandler) GetStoreDetail(w http.ResponseWriter, r *http.Request) {
	// Get store code from URL
	vars := mux.Vars(r)
	storeCode := vars["code"]

	if storeCode == "" {
		utils.WriteError(w, http.StatusBadRequest, "Store code is required")
		return
	}

	// Get auth from headers
	token := r.Header.Get("Authorization")
	if token != "" && len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	deviceCode := r.Header.Get("X-Device-Code")
	wToken := r.Header.Get("X-WToken")
	ucde := r.Header.Get("X-UCDE")

	// Build query string
	path := fmt.Sprintf("/portal/app/basic/storeInfo/getStoreDetail/v2?storeCode=%s", storeCode)

	// Create request context
	ctx := gateway.RequestContext{
		Token:      token,
		DeviceCode: deviceCode,
		WToken:     wToken,
		UCDE:       ucde,
	}

	// Call Tomoro API
	respBody, err := h.client.Request(ctx, "GET", path, nil)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get store detail: %v", err))
		return
	}

	// Parse and return response
	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to parse response")
		return
	}

	utils.WriteSuccess(w, result["data"])
}
