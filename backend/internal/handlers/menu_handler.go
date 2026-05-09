package handlers

import (
	"coffee-order/internal/gateway"
	"coffee-order/internal/models"
	"coffee-order/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

// MenuHandler handles menu-related requests
type MenuHandler struct {
	client *gateway.TomoroClient
}

// NewMenuHandler creates a new menu handler
func NewMenuHandler(client *gateway.TomoroClient) *MenuHandler {
	return &MenuHandler{client: client}
}

// GetMenuList handles menu list request
func (h *MenuHandler) GetMenuList(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	storeCode := r.URL.Query().Get("storeCode")

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

	// Build query string (mainMenuType=1 for regular menu)
	path := fmt.Sprintf("/portal/app/basic/menu/getMenuList?mainMenuType=1&storeCode=%s", storeCode)

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
		utils.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get menu: %v", err))
		return
	}

	// Parse response
	var menuResp models.MenuListResponse
	if err := json.Unmarshal(respBody, &menuResp); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to parse menu response")
		return
	}

	utils.WriteSuccess(w, menuResp.Data)
}
