package handlers

import (
	"coffee-order/config"
	"coffee-order/internal/gateway"
	"coffee-order/internal/models"
	"coffee-order/internal/utils"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	client *gateway.TomoroClient
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(client *gateway.TomoroClient) *AuthHandler {
	return &AuthHandler{client: client}
}

// hashMD5 creates MD5 hash of input string
func hashMD5(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

// Login handles user login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Parse request
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if req.Phone == "" || req.PIN == "" {
		utils.WriteError(w, http.StatusBadRequest, "Phone and PIN are required")
		return
	}

	// Generate device code and tokens
	deviceCode := utils.GenerateDeviceCode()
	wToken, err := utils.GenerateWToken()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to generate tokens")
		return
	}

	ucde, err := utils.GenerateUCDE()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to generate UCDE")
		return
	}

	// Hash PIN with MD5
	hashedPIN := hashMD5(req.PIN)

	// Prepare login request to Tomoro API
	loginBody := map[string]interface{}{
		"phone":     req.Phone,
		"phoneArea": config.PhoneArea,
		"password":  hashedPIN, // Send MD5 hashed PIN
	}

	// Create request context
	ctx := gateway.RequestContext{
		DeviceCode: deviceCode,
		WToken:     wToken,
		UCDE:       ucde,
		Latitude:   config.DefaultLatitude,
		Longitude:  config.DefaultLongitude,
	}

	// Call Tomoro API
	respBody, err := h.client.Request(ctx, "POST", "/portal/app/member/v2/loginPhone", loginBody)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Sprintf("Login failed: %v", err))
		return
	}

	// Parse response
	var loginResp models.LoginResponse
	if err := json.Unmarshal(respBody, &loginResp); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to parse login response")
		return
	}

	// Check if login was successful
	if !loginResp.Success || loginResp.Code != 0 {
		utils.WriteError(w, http.StatusUnauthorized, loginResp.Msg)
		return
	}

	// Prepare response
	authData := models.AuthData{
		Token:      loginResp.Data.Token,
		DeviceCode: deviceCode,
		WToken:     wToken,
		UCDE:       ucde,
		User: models.User{
			MemberCode: loginResp.Data.AccountCode,
			Name:       loginResp.Data.Nickname,
			Phone:      loginResp.Data.Mobile,
		},
	}

	utils.WriteSuccess(w, authData)
}
