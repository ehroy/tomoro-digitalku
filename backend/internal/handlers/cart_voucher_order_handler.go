package handlers

import (
	"coffee-order/config"
	"coffee-order/internal/gateway"
	"coffee-order/internal/models"
	"coffee-order/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func requestContextFromHeaders(r *http.Request) gateway.RequestContext {
	token := r.Header.Get("Authorization")
	if token != "" && len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	return gateway.RequestContext{
		Token:      token,
		DeviceCode: r.Header.Get("X-Device-Code"),
		WToken:     r.Header.Get("X-WToken"),
		UCDE:       r.Header.Get("X-UCDE"),
		Latitude:   config.DefaultLatitude,
		Longitude:  config.DefaultLongitude,
	}
}

func normalizeOrderPayload(body map[string]interface{}) map[string]interface{} {
	normalized := make(map[string]interface{}, len(body)+4)
	for key, value := range body {
		normalized[key] = value
	}

	if businessCode, ok := normalized["businessCode"].(string); ok && businessCode != "" {
	} else if voucherCode, ok := normalized["voucherCode"].(string); ok && voucherCode != "" {
		normalized["businessCode"] = voucherCode
	} else if memberCouponCode, ok := normalized["memberCouponCode"].(string); ok && memberCouponCode != "" {
		normalized["businessCode"] = memberCouponCode
	}

	if _, exists := normalized["bizChannel"]; !exists || normalized["bizChannel"] == "" {
		normalized["bizChannel"] = "APP"
	}
	if _, exists := normalized["businessType"]; !exists || normalized["businessType"] == nil {
		normalized["businessType"] = 1
	}
	if _, exists := normalized["deliverType"]; !exists || normalized["deliverType"] == nil {
		normalized["deliverType"] = 1
	}
	if _, exists := normalized["payChannelCode"]; !exists || normalized["payChannelCode"] == "" {
		normalized["payChannelCode"] = "1007"
	}

	return normalized
}

func normalizeCalculationPayload(body map[string]interface{}) map[string]interface{} {
	normalized := make(map[string]interface{}, len(body)+2)
	for key, value := range body {
		normalized[key] = value
	}

	if businessCode, ok := normalized["businessCode"].(string); ok && businessCode != "" {
		// keep as-is
	} else if voucherCode, ok := normalized["voucherCode"].(string); ok && voucherCode != "" {
		normalized["businessCode"] = voucherCode
	} else if memberCouponCode, ok := normalized["memberCouponCode"].(string); ok && memberCouponCode != "" {
		normalized["businessCode"] = memberCouponCode
	}

	if _, exists := normalized["bizChannel"]; !exists || normalized["bizChannel"] == "" {
		normalized["bizChannel"] = "APP"
	}
	if _, exists := normalized["deliverType"]; !exists || normalized["deliverType"] == nil {
		normalized["deliverType"] = 1
	}

	delete(normalized, "businessType")
	delete(normalized, "payChannelCode")

	return normalized
}

func normalizeRemoveVoucherPayload(body map[string]interface{}) map[string]interface{} {
	normalized := make(map[string]interface{}, len(body)+2)
	for key, value := range body {
		normalized[key] = value
	}

	delete(normalized, "businessCode")
	delete(normalized, "voucherCode")
	delete(normalized, "memberCouponCode")
	delete(normalized, "businessType")

	if _, exists := normalized["bizChannel"]; !exists || normalized["bizChannel"] == "" {
		normalized["bizChannel"] = "APP"
	}
	if _, exists := normalized["deliverType"]; !exists || normalized["deliverType"] == nil {
		normalized["deliverType"] = 1
	}

	return normalized
}

func normalizeApplyVoucherPayload(body map[string]interface{}) map[string]interface{} {
	normalized := normalizeCalculationPayload(body)
	normalized["businessType"] = 1
	return normalized
}

func toFloat64(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case json.Number:
		if n, err := v.Float64(); err == nil {
			return n, true
		}
	case string:
		if n, err := strconv.ParseFloat(strings.TrimSpace(v), 64); err == nil {
			return n, true
		}
	}
	return 0, false
}

func normalizeCalculatedSummary(raw map[string]interface{}) map[string]interface{} {
	if raw == nil {
		return map[string]interface{}{}
	}

	out := make(map[string]interface{}, len(raw))
	for key, value := range raw {
		if num, ok := toFloat64(value); ok && (strings.Contains(strings.ToLower(key), "money") || strings.Contains(strings.ToLower(key), "amount") || strings.Contains(strings.ToLower(key), "price")) {
			out[key] = num
			continue
		}
		out[key] = value
	}

	if total, ok := toFloat64(out["paymentMoney"]); ok {
		out["paymentMoney"] = total
	}
	if discount, ok := toFloat64(out["discountMoney"]); ok {
		out["discountMoney"] = discount
	}
	if subtotal, ok := toFloat64(out["itemTotalMoney"]); ok {
		out["itemTotalMoney"] = subtotal
	}

	return out
}

func normalizeVoucherData(raw interface{}) map[string]interface{} {
	if raw == nil {
		return map[string]interface{}{"coupons": []interface{}{}}
	}

	if data, ok := raw.(map[string]interface{}); ok {
		if coupons, exists := data["coupons"]; exists {
			return map[string]interface{}{"coupons": coupons}
		}
		if records, exists := data["records"]; exists {
			return map[string]interface{}{"coupons": records}
		}
		if couponList, exists := data["couponList"]; exists {
			return map[string]interface{}{"coupons": couponList}
		}
	}

	return map[string]interface{}{"coupons": raw}
}

func hoistCheckoutFields(data map[string]interface{}) map[string]interface{} {
	if data == nil {
		return nil
	}

	wallets, ok := data["walletsChargesResultVo"].(map[string]interface{})
	if !ok {
		return data
	}

	if url, ok := wallets["mobileDeeplinkCheckoutUrl"].(string); ok && url != "" {
		data["mobileDeeplinkCheckoutUrl"] = url
	}
	if guidePicture, ok := wallets["guidePicture"].(string); ok && guidePicture != "" {
		data["guidePicture"] = guidePicture
	}
	if url, ok := wallets["mobileWebCheckoutUrl"].(string); ok && url != "" {
		data["mobileWebCheckoutUrl"] = url
	}

	return data
}

// CartHandler handles cart-related requests
type CartHandler struct {
	client *gateway.TomoroClient
}

// NewCartHandler creates a new cart handler
func NewCartHandler(client *gateway.TomoroClient) *CartHandler {
	return &CartHandler{client: client}
}

func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	storeCode := r.URL.Query().Get("storeCode")
	mainMenuType := r.URL.Query().Get("mainMenuType")
	if mainMenuType == "" {
		mainMenuType = "1"
	}

	ctx := requestContextFromHeaders(r)

	path := fmt.Sprintf("/portal/app/cart/getCart/v3?mainMenuType=%s&storeCode=%s", url.QueryEscape(mainMenuType), url.QueryEscape(storeCode))
	respBody, err := h.client.Request(ctx, "GET", path, nil)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get cart: %v", err))
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to parse response")
		return
	}

	if data, ok := result["data"].(map[string]interface{}); ok {
		data = hoistCheckoutFields(data)
		utils.WriteSuccess(w, data)
		return
	}

	utils.WriteSuccess(w, result["data"])
}

func (h *CartHandler) getCartByMenuType(ctx gateway.RequestContext, storeCode, mainMenuType string) (map[string]interface{}, error) {
	path := fmt.Sprintf("/portal/app/cart/getCart/v3?mainMenuType=%s&storeCode=%s", url.QueryEscape(mainMenuType), url.QueryEscape(storeCode))
	respBody, err := h.client.Request(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}

	if data, ok := result["data"].(map[string]interface{}); ok {
		return data, nil
	}

	if data, ok := result["data"].([]interface{}); ok {
		return map[string]interface{}{"records": data}, nil
	}

	if result["data"] == nil {
		return map[string]interface{}{"records": []interface{}{}}, nil
	}

	return map[string]interface{}{"records": []interface{}{result["data"]}}, nil
}

func mergeCartResults(results ...map[string]interface{}) map[string]interface{} {
	merged := make([]interface{}, 0)
	seen := map[string]bool{}

	for _, data := range results {
		if data == nil {
			continue
		}

		for _, key := range []string{"cartItemList", "records", "cartItems", "items", "list"} {
			items, ok := data[key].([]interface{})
			if !ok {
				continue
			}

			for _, rawItem := range items {
				item, ok := rawItem.(map[string]interface{})
				if !ok {
					merged = append(merged, rawItem)
					continue
				}

				productID := fmt.Sprintf("%v", item["itemCode"])
				if productID == "<nil>" || productID == "" {
					productID = fmt.Sprintf("%v", item["productCode"])
				}
				sizeCode := fmt.Sprintf("%v", item["pluCode"])
				if sizeCode == "<nil>" {
					sizeCode = ""
				}
				cartKey := fmt.Sprintf("%v", item["cartKey"])
				if cartKey == "<nil>" || cartKey == "" {
					cartKey = fmt.Sprintf("%s:%s", productID, sizeCode)
				}

				if seen[cartKey] {
					continue
				}
				seen[cartKey] = true
				merged = append(merged, rawItem)
			}
		}
	}

	return map[string]interface{}{"records": merged}
}

// GetCartAll handles cart requests across all menu buckets
func (h *CartHandler) GetCartAll(w http.ResponseWriter, r *http.Request) {
	storeCode := r.URL.Query().Get("storeCode")

	ctx := requestContextFromHeaders(r)

	menuTypes := []string{"1", "2", "3", "4", "5"}
	results := make([]map[string]interface{}, 0, len(menuTypes))
	for _, menuType := range menuTypes {
		data, err := h.getCartByMenuType(ctx, storeCode, menuType)
		if err != nil {
			continue
		}
		results = append(results, data)
	}

	if len(results) == 0 {
		utils.WriteSuccess(w, map[string]interface{}{"records": []interface{}{}})
		return
	}

	utils.WriteSuccess(w, mergeCartResults(results...))
}

// AddToCart handles add to cart request
func (h *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	ctx := requestContextFromHeaders(r)

	// Call Tomoro API
	respBody, err := h.client.Request(ctx, "POST", "/portal/app/cart/addCart/v3", body)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to add to cart: %v", err))
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

// EditCart handles cart item edit/remove request
func (h *CartHandler) EditCart(w http.ResponseWriter, r *http.Request) {
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	ctx := requestContextFromHeaders(r)

	respBody, err := h.client.Request(ctx, "POST", "/portal/app/cart/editCart/v3", body)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to edit cart: %v", err))
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to parse response")
		return
	}

	if data, ok := result["data"].(map[string]interface{}); ok {
		utils.WriteSuccess(w, data)
		return
	}

	utils.WriteSuccess(w, result["data"])
}

// VoucherHandler handles voucher-related requests
type VoucherHandler struct {
	client *gateway.TomoroClient
}

// NewVoucherHandler creates a new voucher handler
func NewVoucherHandler(client *gateway.TomoroClient) *VoucherHandler {
	return &VoucherHandler{client: client}
}

func normalizeImageURL(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") || strings.HasPrefix(raw, "data:") {
		return raw
	}
	if strings.HasPrefix(raw, "/") {
		return "https://api-service.tomoro-coffee.id" + raw
	}
	return "https://api-service.tomoro-coffee.id/" + raw
}

func formatTradeOrderStatus(status int) string {
	switch status {
	case 6:
		return "Dibatalkan"
	case 4, 5:
		return "Selesai"
	case 2, 3:
		return "Diproses"
	default:
		return "Diproses"
	}
}

func isPaidOrder(paymentOrderCode, transCode string) bool {
	if paymentOrderCode == "" || paymentOrderCode == "0" {
		return false
	}
	if transCode == "" || transCode == "0" {
		return false
	}
	return true
}

func formatPaymentStatus(orderStatus int, paymentOrderCode, transCode string) string {
	switch orderStatus {
	case 1:
		return "Menunggu pembayaran"
	case 6:
		return "Dibatalkan"
	default:
		if isPaidOrder(paymentOrderCode, transCode) {
			return "Sudah dibayar"
		}
		return "Sudah dibayar"
	}
}

// GetVouchers handles get vouchers request
func (h *VoucherHandler) GetVouchers(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	storeCode := r.URL.Query().Get("storeCode")

	ctx := requestContextFromHeaders(r)

	if storeCode != "" {
		// Store-specific vouchers for checkout
		path := fmt.Sprintf("/portal/app/coupon/getStoreAvailableCoupon?storeCode=%s", url.QueryEscape(storeCode))
		respBody, err := h.client.Request(ctx, "GET", path, nil)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get vouchers: %v", err))
			return
		}

		var result map[string]interface{}
		if err := json.Unmarshal(respBody, &result); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to parse response")
			return
		}

		utils.WriteSuccess(w, normalizeVoucherData(result["data"]))
		return
	}

	// Member vouchers for history tab
	body := map[string]interface{}{
		"pageNo":   1,
		"pageSize": 20,
	}
	respBody, err := h.client.Request(ctx, "POST", "/portal/app/coupon/getCouponMemberList", body)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get vouchers: %v", err))
		return
	}

	var voucherResp models.MemberVoucherListResponse
	if err := json.Unmarshal(respBody, &voucherResp); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to parse response")
		return
	}

	utils.WriteSuccess(w, map[string]interface{}{
		"records": voucherResp.Data.Records,
		"total":   len(voucherResp.Data.Records),
		"size":    20,
		"current": 1,
		"pages":   1,
	})
}

// OrderHandler handles order-related requests
type OrderHandler struct {
	client *gateway.TomoroClient
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(client *gateway.TomoroClient) *OrderHandler {
	return &OrderHandler{client: client}
}

// CalculateOrder handles order calculation request
func (h *OrderHandler) CalculateOrder(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	ctx := requestContextFromHeaders(r)

	h.writeCalculatedOrderWithPayload(w, ctx, body)
}

// ApplyVoucher handles voucher apply request
func (h *OrderHandler) ApplyVoucher(w http.ResponseWriter, r *http.Request) {
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	ctx := requestContextFromHeaders(r)

	h.writeCalculatedOrderWithPayload(w, ctx, body)
}

// RemoveVoucher handles voucher remove request
func (h *OrderHandler) RemoveVoucher(w http.ResponseWriter, r *http.Request) {
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	ctx := requestContextFromHeaders(r)

	h.writeCalculatedOrderWithPayload(w, ctx, body)
}

func (h *OrderHandler) writeCalculatedOrder(w http.ResponseWriter, ctx gateway.RequestContext, body map[string]interface{}) {
	h.writeCalculatedOrderWithPayload(w, ctx, normalizeOrderPayload(body))
}

func (h *OrderHandler) writeCalculatedOrderWithPayload(w http.ResponseWriter, ctx gateway.RequestContext, payload map[string]interface{}) {
	respBody, err := h.client.Request(ctx, "POST", "/portal/app/order/calcTradeOrderAgain", payload)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to calculate order: %v", err))
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to parse response")
		return
	}

	if data, ok := result["data"].(map[string]interface{}); ok {
		utils.WriteSuccess(w, normalizeCalculatedSummary(data))
		return
	}

	utils.WriteSuccess(w, result["data"])
}

// CreateOrder handles order creation request
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	delete(body, "items")
	delete(body, "payment")
	delete(body, "voucherCode")
	body = normalizeOrderPayload(body)

	// Get auth from headers
	ctx := requestContextFromHeaders(r)

	// Call Tomoro API
	respBody, err := h.client.Request(ctx, "POST", "/portal/app/order/createTradeAndPayment", body)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create order: %v", err))
		return
	}

	// Parse and return response
	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to parse response")
		return
	}

	if success, ok := result["success"].(bool); ok && !success {
		msg, _ := result["msg"].(string)
		if msg == "" {
			msg = "Failed to create order"
		}
		utils.WriteError(w, http.StatusBadRequest, msg)
		return
	}

	if data, ok := result["data"].(map[string]interface{}); ok {
		utils.WriteSuccess(w, data)
		return
	}

	utils.WriteSuccess(w, result["data"])
}

// GetPayStatus handles payment status request
func (h *OrderHandler) GetPayStatus(w http.ResponseWriter, r *http.Request) {
	tradeOrderCode := r.URL.Query().Get("tradeOrderCode")
	if tradeOrderCode == "" {
		utils.WriteError(w, http.StatusBadRequest, "tradeOrderCode is required")
		return
	}

	ctx := requestContextFromHeaders(r)
	path := fmt.Sprintf("/portal/app/pay/getPayStatus?tradeOrderCode=%s", url.QueryEscape(tradeOrderCode))
	respBody, err := h.client.Request(ctx, "GET", path, nil)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get pay status: %v", err))
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to parse response")
		return
	}

	paymentData := map[string]interface{}{}
	if data, ok := result["data"].(map[string]interface{}); ok {
		paymentData = data
	}

	extractString := func(keys ...string) string {
		for _, key := range keys {
			if value, ok := paymentData[key].(string); ok && value != "" {
				return value
			}
		}
		return ""
	}

	redirectUrl := extractString(
		"redirectUrl",
		"checkoutUrl",
		"mobileDeeplinkCheckoutUrl",
		"payUrl",
		"url",
		"deeplinkUrl",
	)

	status := "Menunggu pembayaran"
	if paid, _ := paymentData["success"].(bool); paid {
		status = "Sudah dibayar"
	}

	utils.WriteSuccess(w, map[string]interface{}{
		"tradeOrderCode": tradeOrderCode,
		"paymentStatus":  status,
		"isPaid":         paymentData["success"],
		"redirectUrl":    redirectUrl,
		"payMoney":       paymentData["payMoney"],
		"currencyUnit":    paymentData["currencyUnit"],
		"currencyCode":    paymentData["currencyCode"],
		"payChannelName":  paymentData["payChannelName"],
		"pointNum":        paymentData["pointNum"],
		"raw":             paymentData,
	})
}

// GetOrderHistory handles order history request
func (h *OrderHandler) GetOrderHistory(w http.ResponseWriter, r *http.Request) {
	pageNo := 1
	pageSize := 20

	// Get auth from headers
	ctx := requestContextFromHeaders(r)

	body := map[string]interface{}{
		"pageNo":   pageNo,
		"pageSize": pageSize,
	}

	respBody, err := h.client.Request(ctx, "POST", "/portal/app/order/getTradeOrderPage", body)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get order history: %v", err))
		return
	}

	var result models.TradeOrderPageResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to parse response")
		return
	}

	normalized := models.NormalizedOrderHistoryResponse{
		Records: make([]models.OrderHistory, 0, len(result.Data.Records)),
		Total:   result.Data.RecordsTotal,
		Size:    result.Data.PageSize,
		Current: result.Data.PageNo,
		Pages:   result.Data.PagesTotal,
	}

	for _, order := range result.Data.Records {
		items := make([]models.OrderHistoryItem, 0, len(order.TradeOrderItemVOList))
		for _, item := range order.TradeOrderItemVOList {
			price := item.PaymentPrice
			if price == 0 {
				price = item.PluPrice
			}
			subtotal := item.PluTotalMoney
			if subtotal == 0 {
				subtotal = price * float64(item.Amount)
			}

			items = append(items, models.OrderHistoryItem{
				ProductCode: item.ItemCode,
				ProductName: item.ItemName,
				Image:       normalizeImageURL(item.ItemPictureUrls),
				Quantity:    item.Amount,
				Price:       price,
				Subtotal:    subtotal,
			})
		}

		paymentStatus := formatPaymentStatus(order.Status, order.PaymentOrderCode, order.TransCode)

			normalized.Records = append(normalized.Records, models.OrderHistory{
				OrderCode:      order.Code,
				StoreCode:      order.StoreCode,
				StoreName:      order.StoreName,
				Status:         formatTradeOrderStatus(order.Status),
			StatusCode:     order.Status,
			PaymentStatus:  paymentStatus,
			CreatedAt:      time.UnixMilli(order.CreateTime).Format(time.RFC3339),
			TotalAmount:    order.PaymentMoney,
			DiscountAmount: order.DiscountMoney,
			PaymentAmount:  order.PaymentMoney,
			VoucherCode:    order.MemberCouponCode,
			Payment:        "gopay",
			PaymentOrderCode: order.PaymentOrderCode,
			TransCode:      order.TransCode,
			PayType:        order.PayType,
			PayChannelCode: order.PayChannelCode,
			Items:          items,
		})
	}

	utils.WriteSuccess(w, normalized)
}
