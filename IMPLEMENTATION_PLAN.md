# 🚀 Implementation Plan - Tomoro Coffee API Integration

## 📋 Overview

Mengintegrasikan API Tomoro Coffee ke aplikasi dengan:
- **Backend**: Golang sebagai API Gateway/Proxy
- **Frontend**: Svelte dengan real API consumption
- **Auth**: Real authentication dengan token
- **Images**: Load dari API URLs

---

## 📡 API Endpoints Discovered

Dari analisis Proxyman logs, ditemukan endpoint-endpoint berikut:

### 1. **Configuration**
```
GET /portal/app/module/config/getTabBarImg
```
- Get tab bar images/config

### 2. **Store Management**
```
GET /portal/app/basic/storeInfo/getStoreList/v3
  Query: centerPointLatitude, centerPointLongitude, pageNo, pageSize, storeName
  
GET /portal/app/basic/storeInfo/getStoreDetail/v2
  Query: storeCode
```

### 3. **Menu**
```
GET /portal/app/basic/menu/getMenuList
  Query: mainMenuType, storeCode
```

### 4. **Authentication**
```
POST /portal/app/member/v2/loginPhone
  Body: phone, password/otp
```

### 5. **Cart**
```
POST /portal/app/cart/addCart/v3
  Body: product details, quantity, storeCode
```

### 6. **Voucher/Coupon**
```
GET /portal/app/coupon/getStoreAvailableCoupon
  Query: storeCode
  
POST /portal/app/coupon/getCouponMemberQuantity
  Body: coupon details
```

### 7. **Order**
```
POST /portal/app/order/calcTradeOrderAgain
  Body: order calculation details
  
POST /portal/app/order/createTradeAndPayment
  Body: complete order details
```

---

## 🏗️ Architecture Design

### Backend (Golang) - API Gateway Pattern

```
┌─────────────┐
│   Frontend  │
│   (Svelte)  │
└──────┬──────┘
       │ HTTP
       ▼
┌─────────────┐
│   Golang    │
│   Gateway   │
│  (Port 8080)│
└──────┬──────┘
       │ HTTPS
       ▼
┌─────────────┐
│  Tomoro API │
│ (Real API)  │
└─────────────┘
```

**Benefits:**
- Hide API keys/tokens
- Add caching layer
- Transform responses
- Add logging/monitoring
- Handle CORS properly

### Frontend (Svelte) - Clean Architecture

```
┌──────────────────────────────────┐
│         Components               │
│  (Pages, UI Components)          │
└────────────┬─────────────────────┘
             │
┌────────────▼─────────────────────┐
│         Stores                   │
│  (State Management)              │
└────────────┬─────────────────────┘
             │
┌────────────▼─────────────────────┐
│       API Client                 │
│  (HTTP Requests)                 │
└────────────┬─────────────────────┘
             │
┌────────────▼─────────────────────┐
│      Backend Gateway             │
└──────────────────────────────────┘
```

---

## 📦 Data Models

### Store
```go
type Store struct {
    StoreCode       string  `json:"storeCode"`
    StoreName       string  `json:"storeName"`
    StorePicture    string  `json:"storePicture"`
    StorePhone      string  `json:"storePhone"`
    StoreAddress    string  `json:"storeAddress"`
    Longitude       float64 `json:"longitude"`
    Latitude        float64 `json:"latitude"`
    IsDelivery      int     `json:"isDelivery"`
    BusinessStatus  int     `json:"businessStatus"`
    Distance        int     `json:"distance"`
}
```

### Menu/Product
```go
type Product struct {
    ProductCode     string  `json:"productCode"`
    ProductName     string  `json:"productName"`
    ProductPicture  string  `json:"productPicture"`
    Price           float64 `json:"price"`
    OriginalPrice   float64 `json:"originalPrice"`
    Category        string  `json:"category"`
    Description     string  `json:"description"`
    IsAvailable     bool    `json:"isAvailable"`
}
```

### Cart Item
```go
type CartItem struct {
    ProductCode string  `json:"productCode"`
    ProductName string  `json:"productName"`
    Quantity    int     `json:"quantity"`
    Price       float64 `json:"price"`
    StoreCode   string  `json:"storeCode"`
}
```

### Voucher
```go
type Voucher struct {
    CouponCode      string  `json:"couponCode"`
    CouponName      string  `json:"couponName"`
    DiscountType    int     `json:"discountType"`
    DiscountValue   float64 `json:"discountValue"`
    MinOrderAmount  float64 `json:"minOrderAmount"`
    ExpiryDate      string  `json:"expiryDate"`
}
```

### Order
```go
type Order struct {
    OrderCode       string      `json:"orderCode"`
    StoreCode       string      `json:"storeCode"`
    Items           []CartItem  `json:"items"`
    TotalAmount     float64     `json:"totalAmount"`
    DiscountAmount  float64     `json:"discountAmount"`
    FinalAmount     float64     `json:"finalAmount"`
    Status          string      `json:"status"`
    CreatedAt       string      `json:"createdAt"`
}
```

---

## 🔐 Authentication Flow

### 1. Login Process
```
User Input (Phone + Password/OTP)
    ↓
Frontend → Backend Gateway
    ↓
Backend → Tomoro API (POST /portal/app/member/v2/loginPhone)
    ↓
Response: { token, userInfo }
    ↓
Store token in:
  - Backend: Session/Redis
  - Frontend: LocalStorage + Svelte Store
```

### 2. Authenticated Requests
```
Every API call includes:
  Headers:
    - token: <user_token>
    - deviceCode: <device_id>
    - ucde: <user_code>
    - wToken: <security_token>
```

---

## 🛠️ Implementation Steps

### Phase 1: Backend Gateway Setup

#### 1.1 Project Structure
```
backend/
├── cmd/server/
│   └── main.go
├── internal/
│   ├── gateway/
│   │   ├── client.go      # HTTP client to Tomoro API
│   │   ├── auth.go        # Auth middleware
│   │   └── transformer.go # Response transformation
│   ├── handlers/
│   │   ├── store.go
│   │   ├── menu.go
│   │   ├── cart.go
│   │   ├── voucher.go
│   │   └── order.go
│   ├── models/
│   │   └── models.go
│   └── middleware/
│       ├── cors.go
│       ├── logger.go
│       └── auth.go
└── config/
    └── config.go
```

#### 1.2 Gateway Client
```go
type TomoroClient struct {
    BaseURL    string
    HTTPClient *http.Client
    Token      string
    DeviceCode string
}

func (c *TomoroClient) Request(method, path string, body interface{}) (*Response, error)
```

#### 1.3 Handlers Implementation
- Store handlers (list, detail, search)
- Menu handlers (get menu by store)
- Cart handlers (add, update, remove)
- Voucher handlers (list, apply)
- Order handlers (calculate, create, history)

### Phase 2: Frontend Integration

#### 2.1 API Client
```javascript
// src/lib/api/tomoro.js
class TomoroAPI {
  constructor(baseURL) {
    this.baseURL = baseURL;
    this.token = null;
  }
  
  async login(phone, password) { }
  async getStores(lat, lng, search) { }
  async getStoreDetail(storeCode) { }
  async getMenu(storeCode) { }
  async addToCart(item) { }
  async getVouchers(storeCode) { }
  async createOrder(orderData) { }
}
```

#### 2.2 Stores (State Management)
```javascript
// src/lib/stores/auth.js
export const authStore = writable({
  token: null,
  user: null,
  isAuthenticated: false
});

// src/lib/stores/cart.js
export const cartStore = writable([]);

// src/lib/stores/voucher.js
export const voucherStore = writable(null);
```

#### 2.3 Components Update
- Login page/modal
- Store list with real data
- Menu page with store selection
- Cart with voucher application
- Checkout flow

### Phase 3: Feature Implementation

#### 3.1 Store Search & List
- Geolocation permission
- Search by name
- Filter by distance
- Show on map (optional)

#### 3.2 Menu per Store
- Select store first
- Load menu by storeCode
- Category filtering
- Product images from API

#### 3.3 Cart Management
- Add items with storeCode
- Update quantity
- Remove items
- Calculate total

#### 3.4 Voucher System
- List available vouchers
- Apply voucher to cart
- Show discount calculation
- Validate min order amount

#### 3.5 Order Creation
- Review order
- Apply voucher
- Calculate final amount
- Create order
- Show order confirmation

---

## 🔄 API Request/Response Examples

### Login
**Request:**
```json
POST /api/auth/login
{
  "phone": "08123456789",
  "password": "password123"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGc...",
    "user": {
      "id": "USR001",
      "name": "Budi Santoso",
      "phone": "08123456789"
    }
  }
}
```

### Get Stores
**Request:**
```
GET /api/stores?lat=-6.573982&lng=110.684519&search=saudara&page=1&size=20
```

**Response:**
```json
{
  "success": true,
  "data": {
    "records": [...],
    "total": 360,
    "pages": 18
  }
}
```

### Get Menu
**Request:**
```
GET /api/menu?storeCode=JPR15A-M
```

**Response:**
```json
{
  "success": true,
  "data": {
    "products": [
      {
        "productCode": "PROD001",
        "productName": "Espresso",
        "price": 25000,
        "productPicture": "https://..."
      }
    ]
  }
}
```

### Add to Cart
**Request:**
```json
POST /api/cart/add
{
  "productCode": "PROD001",
  "quantity": 2,
  "storeCode": "JPR15A-M"
}
```

### Create Order
**Request:**
```json
POST /api/orders/create
{
  "storeCode": "JPR15A-M",
  "items": [...],
  "voucherCode": "DISC10",
  "totalAmount": 50000
}
```

---

## 🎯 Success Criteria

- [ ] User can login with phone number
- [ ] User can search stores by location
- [ ] User can view store details
- [ ] User can browse menu per store
- [ ] User can add items to cart
- [ ] User can view available vouchers
- [ ] User can apply voucher to order
- [ ] User can create order
- [ ] User can view order history
- [ ] Images load from API
- [ ] All responses properly formatted
- [ ] Error handling implemented
- [ ] Loading states shown

---

## 📝 Next Steps

1. Extract all API logs untuk detail request/response
2. Implement backend gateway handlers
3. Update frontend API client
4. Implement authentication flow
5. Update UI components dengan real data
6. Test end-to-end flow
7. Add error handling & loading states
8. Optimize performance (caching, lazy loading)

---

## ⚠️ Important Notes

- **Security**: Never expose real API keys in frontend
- **Token Management**: Store securely, refresh when expired
- **Error Handling**: Handle network errors, API errors
- **Loading States**: Show spinners during API calls
- **Caching**: Cache store list, menu data
- **Geolocation**: Request permission properly
- **Images**: Lazy load, show placeholders

