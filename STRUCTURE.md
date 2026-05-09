# Project Structure

```
coffe-order/
│
├── backend/                          # Golang Backend
│   ├── cmd/
│   │   └── server/
│   │       └── main.go              # Entry point aplikasi
│   │
│   ├── internal/
│   │   ├── handlers/
│   │   │   └── handlers.go          # HTTP request handlers
│   │   │
│   │   ├── middleware/
│   │   │   └── logger.go            # Logging middleware
│   │   │
│   │   └── models/
│   │       └── models.go            # Data structures
│   │
│   ├── pkg/
│   │   └── dummy/
│   │       └── data.go              # Dummy data untuk development
│   │
│   ├── go.mod                        # Go dependencies
│   └── go.sum                        # Go dependencies checksum
│
├── frontend/                         # Svelte Frontend
│   ├── src/
│   │   ├── lib/
│   │   │   ├── api/
│   │   │   │   └── client.js       # API client functions
│   │   │   │
│   │   │   ├── components/
│   │   │   │   ├── Navbar.svelte   # Navigation component
│   │   │   │   └── ProductCard.svelte # Product card component
│   │   │   │
│   │   │   └── stores/
│   │   │       └── cart.js         # Shopping cart state
│   │   │
│   │   ├── routes/
│   │   │   ├── +layout.svelte      # Layout wrapper
│   │   │   ├── +page.svelte        # Dashboard page
│   │   │   │
│   │   │   ├── menu/
│   │   │   │   └── +page.svelte    # Menu page
│   │   │   │
│   │   │   ├── history/
│   │   │   │   └── +page.svelte    # Order history page
│   │   │   │
│   │   │   ├── profile/
│   │   │   │   └── +page.svelte    # Profile page
│   │   │   │
│   │   │   └── outlets/
│   │   │       └── +page.svelte    # Outlets page
│   │   │
│   │   ├── app.html                 # HTML template
│   │   └── app.css                  # Global styles
│   │
│   ├── static/                       # Static assets
│   ├── package.json                  # Node dependencies
│   ├── svelte.config.js             # Svelte configuration
│   └── vite.config.js               # Vite configuration
│
├── .gitignore                        # Git ignore rules
├── README.md                         # Project documentation
├── QUICKSTART.md                     # Quick start guide
├── STRUCTURE.md                      # This file
├── AGENTS-COFFE.md                   # Project requirements
├── SKILL.md                          # Design guidelines
└── start.sh                          # Startup script
```

## File Descriptions

### Backend Files

#### `backend/cmd/server/main.go`
Entry point aplikasi backend. Setup router, middleware, dan start HTTP server.

#### `backend/internal/handlers/handlers.go`
Berisi semua HTTP request handlers untuk:
- Products (GET all, GET by ID)
- Orders (GET all, GET by ID, POST create, PATCH update status)
- Dashboard stats
- User profile
- Outlets

#### `backend/internal/middleware/logger.go`
Middleware untuk logging HTTP requests.

#### `backend/internal/models/models.go`
Definisi data structures:
- Product
- Order & OrderItem
- User
- DashboardStats
- Outlet

#### `backend/pkg/dummy/data.go`
Dummy data untuk development:
- 8 products (coffee & pastry)
- 3 sample orders
- 2 users
- Dashboard statistics
- 3 outlets

### Frontend Files

#### `frontend/src/lib/api/client.js`
API client functions untuk berkomunikasi dengan backend:
- getProducts(), getProduct(id)
- getOrders(), createOrder(), updateOrderStatus()
- getDashboardStats()
- getUser(id)
- getOutlets()

#### `frontend/src/lib/components/Navbar.svelte`
Navigation bar component dengan links ke semua pages.

#### `frontend/src/lib/components/ProductCard.svelte`
Reusable component untuk menampilkan product card dengan:
- Product image placeholder
- Name, description, price
- Add to cart button
- Availability status

#### `frontend/src/lib/stores/cart.js`
Svelte store untuk shopping cart state management:
- addToCart()
- removeFromCart()
- updateCartQuantity()
- clearCart()
- getCartTotal()

#### `frontend/src/routes/+layout.svelte`
Layout wrapper untuk semua pages, berisi Navbar dan global styles.

#### `frontend/src/routes/+page.svelte`
Dashboard page dengan:
- Statistics cards (revenue, orders, pending, completed)
- Recent orders list

#### `frontend/src/routes/menu/+page.svelte`
Menu page dengan:
- Category filters
- Product grid
- Shopping cart panel
- Checkout functionality

#### `frontend/src/routes/history/+page.svelte`
Order history page dengan:
- Status filters
- Order timeline
- Order details

#### `frontend/src/routes/profile/+page.svelte`
Profile page dengan:
- User information
- Edit profile form
- Settings menu

#### `frontend/src/routes/outlets/+page.svelte`
Outlets page dengan:
- Search functionality
- Outlet cards with location info
- Open/closed status
- Contact buttons

## Design Patterns

### Backend
- **Handler Pattern**: Separate handlers untuk setiap endpoint
- **Middleware Pattern**: Reusable middleware untuk logging
- **Repository Pattern**: Dummy data sebagai data source (bisa diganti dengan database)

### Frontend
- **Component-Based**: Reusable Svelte components
- **Store Pattern**: Centralized state management dengan Svelte stores
- **API Client Pattern**: Centralized API calls
- **Route-Based**: File-based routing dengan SvelteKit

## Data Flow

1. User interacts dengan UI (Svelte components)
2. Component calls API client function
3. API client sends HTTP request ke backend
4. Backend handler processes request
5. Handler returns data dari dummy store
6. Frontend receives response dan updates UI
7. Svelte reactivity updates DOM

## Styling Approach

- **CSS-in-Svelte**: Scoped styles dalam `.svelte` files
- **Coffee Theme**: Brown, cream, dan gold color palette
- **Animations**: CSS transitions dan keyframe animations
- **Responsive**: Mobile-first design dengan media queries
- **Modern UI**: Gradients, shadows, rounded corners

## Future Enhancements

Struktur ini sudah siap untuk:
- Database integration (tambah `backend/internal/database/`)
- Authentication (tambah `backend/internal/auth/`)
- Testing (tambah `backend/internal/handlers/*_test.go`)
- Deployment configs (tambah `docker-compose.yml`, `Dockerfile`)
