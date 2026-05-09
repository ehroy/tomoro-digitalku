# Coffee Order System

Coffee Order adalah aplikasi order kopi mobile-first dengan frontend SvelteKit dan backend Golang yang mem-proxy Tomoro Coffee API.

## Overview

- Frontend selalu mengakses backend lokal di `http://localhost:8080/api`
- Backend meneruskan request ke Tomoro API dan menormalkan response untuk frontend
- Checkout memakai `mobileDeeplinkCheckoutUrl`
- Halaman status order melakukan polling payment status dari Tomoro
- UI dibuat mobile-only, desktop menampilkan layar penjelasan

## Main Flow

1. User login via modal login.
2. User pilih menu dan checkout.
3. Backend membuat order ke Tomoro melalui `/api/orders/create`.
4. Response create menyimpan `mobileDeeplinkCheckoutUrl` dan `tradeOrderCode`.
5. Frontend pindah ke `/checkout/status`.
6. Status page melakukan refresh otomatis ke:
   - `/api/orders/pay-status?tradeOrderCode=...`
   - `/api/orders/history`
7. Jika payment sudah selesai, tombol bayar diganti pesan selesai.

## Tech Stack

Backend:
- Go
- Gorilla Mux
- CORS

Frontend:
- SvelteKit
- Vite

## Project Structure

```text
coffe-order/
├── backend/
│   ├── config/
│   ├── internal/
│   │   ├── gateway/
│   │   ├── handlers/
│   │   ├── middleware/
│   │   └── models/
│   └── main.go
└── frontend/
    └── src/
        ├── lib/
        │   ├── api/
        │   ├── components/
        │   └── stores/
        └── routes/
            ├── +page.svelte
            ├── +layout.svelte
            ├── checkout/
            ├── history/
            ├── menu/
            ├── outlets/
            └── profile/
```

## Key Pages

- `/` - home dashboard
- `/menu` - browse and add items
- `/checkout` - calculate, apply voucher, create order
- `/checkout/status` - payment status screen
- `/history` - order and voucher history
- `/outlets` - outlet list
- `/profile` - profile screen

## Backend API

Base URL: `http://localhost:8080/api`

### Auth

- `POST /auth/login`

### Store

- `GET /stores`
- `GET /stores/{code}`

### Menu

- `GET /menu`

### Cart

- `GET /cart?storeCode=...&mainMenuType=1`
- `GET /cart/all?storeCode=...`
- `POST /cart/add`
- `POST /cart/edit`

### Voucher

- `GET /vouchers`

### Orders

- `POST /orders/calculate`
- `POST /orders/voucher/apply`
- `POST /orders/voucher/remove`
- `POST /orders/create`
- `GET /orders/pay-status?tradeOrderCode=...`
- `GET /orders/history`

## Checkout Details

The checkout page reads these values from the create response:

- `data.code` as `tradeOrderCode`
- `data.mobileDeeplinkCheckoutUrl` or `data.walletsChargesResultVo.mobileDeeplinkCheckoutUrl`
- `data.guidePicture` or `data.walletsChargesResultVo.guidePicture`

The status page stores checkout state in `sessionStorage` so refresh does not lose the payment link.

## Login UI

- Mobile-only experience
- No test credentials shown in UI
- Login modal includes a purchase option for users without an account
- Button links to `https://digitalku-murah.com/`

## Mobile Only Behavior

- Main app shell is hidden on desktop
- Desktop users see a notice to open the app on a phone
- App layout is optimized for touch use and bottom navigation

## Local Setup

### Backend

```bash
cd backend
go mod download
go run main.go
```

Backend runs on `http://localhost:8080`

### Frontend

```bash
cd frontend
npm install
npm run dev
```

Frontend runs on `http://localhost:5173`

## Verification

```bash
cd backend && go test ./...
cd frontend && npm run build
```

## VPS Deployment

Target contoh: Ubuntu 22.04, `backend` di `:8080`, `frontend` di-proxy via Nginx.

### 1. Install dependencies

```bash
sudo apt update
sudo apt install -y git nginx
```

Install Go dan Node.js LTS sesuai kebutuhan server.

### 2. Clone project

```bash
cd /var/www
sudo git clone <repo-url> coffe-order
sudo chown -R $USER:$USER coffe-order
```

### 3. Build backend

```bash
cd /var/www/coffe-order/backend
go test ./...
go build -o coffe-order
```

### 4. Build frontend

```bash
cd /var/www/coffe-order/frontend
npm install
npm run build
```

### 5. Run backend with systemd

Create `/etc/systemd/system/coffe-order.service`:

```ini
[Unit]
Description=Coffee Order Backend
After=network.target

[Service]
WorkingDirectory=/var/www/coffe-order/backend
ExecStart=/var/www/coffe-order/backend/coffe-order
Restart=always
RestartSec=5
User=www-data
Environment=GOMAXPROCS=2

[Install]
WantedBy=multi-user.target
```

Enable it:

```bash
sudo systemctl daemon-reload
sudo systemctl enable coffe-order
sudo systemctl start coffe-order
sudo systemctl status coffe-order
```

### 6. Run frontend service

For the current setup, run SvelteKit in preview mode on a local port.

```bash
cd /var/www/coffe-order/frontend
npm run build
npm run preview -- --host 127.0.0.1 --port 4173
```

You can keep it alive with systemd or PM2.

Example systemd unit `/etc/systemd/system/coffe-order-frontend.service`:

```ini
[Unit]
Description=Coffee Order Frontend
After=network.target

[Service]
WorkingDirectory=/var/www/coffe-order/frontend
ExecStart=/usr/bin/npm run preview -- --host 127.0.0.1 --port 4173
Restart=always
RestartSec=5
User=www-data

[Install]
WantedBy=multi-user.target
```

Enable it:

```bash
sudo systemctl daemon-reload
sudo systemctl enable coffe-order-frontend
sudo systemctl start coffe-order-frontend
sudo systemctl status coffe-order-frontend
```

### 7. Nginx reverse proxy

Example `/etc/nginx/sites-available/coffe-order`:

```nginx
server {
  listen 80;
  server_name your-domain.com;

  location /api/ {
    proxy_pass http://127.0.0.1:8080/api/;
    proxy_http_version 1.1;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
  }

  location / {
    proxy_pass http://127.0.0.1:4173;
    proxy_http_version 1.1;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
  }
}
```

Enable site:

```bash
sudo ln -s /etc/nginx/sites-available/coffe-order /etc/nginx/sites-enabled/coffe-order
sudo nginx -t
sudo systemctl reload nginx
```

### 8. SSL (recommended)

Use Certbot for HTTPS:

```bash
sudo apt install -y certbot python3-certbot-nginx
sudo certbot --nginx -d your-domain.com
```

### 9. Deployment checklist

- Backend service is running
- `GET /api/orders/history` returns data
- `POST /api/orders/create` returns `mobileDeeplinkCheckoutUrl`
- Frontend build exists and is served by Nginx
- `/checkout/status` can refresh after browser reload

## Notes

- Frontend should keep using backend proxy, not Tomoro directly.
- Payment status should come from Tomoro pay status API.
- Order history is normalized by backend for display.
