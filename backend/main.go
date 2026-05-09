package main

import (
	"coffee-order/internal/gateway"
	"coffee-order/internal/handlers"
	"coffee-order/internal/middleware"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Initialize Tomoro API client
	tomoroClient := gateway.NewTomoroClient()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(tomoroClient)
	storeHandler := handlers.NewStoreHandler(tomoroClient)
	menuHandler := handlers.NewMenuHandler(tomoroClient)
	cartHandler := handlers.NewCartHandler(tomoroClient)
	voucherHandler := handlers.NewVoucherHandler(tomoroClient)
	orderHandler := handlers.NewOrderHandler(tomoroClient)

	// Setup router
	r := mux.NewRouter()

	// Apply middleware
	r.Use(middleware.Logger)

	// API routes
	api := r.PathPrefix("/api").Subrouter()

	// Auth routes
	api.HandleFunc("/auth/login", authHandler.Login).Methods("POST")

	// Store routes
	api.HandleFunc("/stores", storeHandler.GetStoreList).Methods("GET")
	api.HandleFunc("/stores/{code}", storeHandler.GetStoreDetail).Methods("GET")

	// Menu routes
	api.HandleFunc("/menu", menuHandler.GetMenuList).Methods("GET")

	// Cart routes
	api.HandleFunc("/cart/add", cartHandler.AddToCart).Methods("POST")
	api.HandleFunc("/cart/edit", cartHandler.EditCart).Methods("POST")
	api.HandleFunc("/cart", cartHandler.GetCart).Methods("GET")
	api.HandleFunc("/cart/all", cartHandler.GetCartAll).Methods("GET")

	// Voucher routes
	api.HandleFunc("/vouchers", voucherHandler.GetVouchers).Methods("GET")

	// Order routes
	api.HandleFunc("/orders/calculate", orderHandler.CalculateOrder).Methods("POST")
	api.HandleFunc("/orders/voucher/apply", orderHandler.ApplyVoucher).Methods("POST")
	api.HandleFunc("/orders/voucher/remove", orderHandler.RemoveVoucher).Methods("POST")
	api.HandleFunc("/orders/create", orderHandler.CreateOrder).Methods("POST")
	api.HandleFunc("/orders/pay-status", orderHandler.GetPayStatus).Methods("GET")
	api.HandleFunc("/orders/history", orderHandler.GetOrderHistory).Methods("GET")

	// CORS configuration
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-Device-Code", "X-WToken", "X-UCDE"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	log.Println("🚀 Server starting on :8080")
	log.Println("📡 Proxying to Tomoro Coffee API")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
