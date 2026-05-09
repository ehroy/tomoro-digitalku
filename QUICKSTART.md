# Quick Start Guide

## Prerequisites

- **Golang** 1.21 atau lebih tinggi
- **Node.js** 18.x atau lebih tinggi
- **npm** atau **yarn**

## Installation Steps

### 1. Clone atau Download Project

```bash
cd /Users/roy/Desktop/golang/coffe-order
```

### 2. Setup Backend

```bash
# Navigate ke folder backend
cd backend

# Download dependencies
go mod download

# Run server
go run cmd/server/main.go
```

Backend akan berjalan di: **http://localhost:8080**

### 3. Setup Frontend (Terminal Baru)

```bash
# Navigate ke folder frontend
cd frontend

# Install dependencies
npm install

# Run development server
npm run dev
```

Frontend akan berjalan di: **http://localhost:5173**

### 4. Atau Gunakan Script Otomatis

```bash
# Dari root folder project
./start.sh
```

Script ini akan menjalankan backend dan frontend secara bersamaan.

## Testing API

Anda bisa test API menggunakan curl atau Postman:

```bash
# Get all products
curl http://localhost:8080/api/products

# Get dashboard stats
curl http://localhost:8080/api/dashboard/stats

# Get all orders
curl http://localhost:8080/api/orders

# Create new order
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": "USR001",
    "items": [
      {
        "product_id": "1",
        "name": "Espresso",
        "quantity": 2,
        "price": 25000,
        "subtotal": 50000
      }
    ]
  }'
```

## Accessing the Application

1. Buka browser
2. Navigate ke **http://localhost:5173**
3. Explore fitur-fitur:
   - Dashboard - Lihat statistik
   - Menu - Browse dan order produk
   - History - Lihat riwayat pesanan
   - Profile - Manage profile
   - Outlets - Cari lokasi outlet

## Troubleshooting

### Port Already in Use

Jika port 8080 atau 5173 sudah digunakan:

**Backend:**
Edit `backend/cmd/server/main.go` dan ubah port:
```go
http.ListenAndServe(":8081", handler)
```

**Frontend:**
Edit `frontend/vite.config.js` dan ubah port:
```js
server: {
  port: 5174
}
```

### CORS Issues

Pastikan backend CORS configuration sudah benar di `backend/cmd/server/main.go`:
```go
AllowedOrigins: []string{"http://localhost:5173"}
```

### Module Not Found

**Backend:**
```bash
cd backend
go mod tidy
```

**Frontend:**
```bash
cd frontend
rm -rf node_modules
npm install
```

## Development Tips

1. **Hot Reload**: Kedua server support hot reload
2. **API Testing**: Gunakan browser DevTools Network tab
3. **Dummy Data**: Edit `backend/pkg/dummy/data.go` untuk mengubah data
4. **Styling**: Edit file `.svelte` untuk mengubah tampilan

## Next Steps

- Implementasi database (SQLite/PostgreSQL)
- Tambahkan authentication
- Integrasikan payment gateway
- Deploy ke production

Selamat coding! ☕
