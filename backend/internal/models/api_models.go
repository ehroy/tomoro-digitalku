package models

// LoginRequest represents login request
type LoginRequest struct {
	Phone string `json:"phone"`
	PIN   string `json:"pin"`
}

// LoginResponse represents login response from Tomoro API
type LoginResponse struct {
	Data struct {
		Token       string `json:"token"`
		AccountCode string `json:"accountCode"`
		Nickname    string `json:"nickname"`
		Mobile      string `json:"mobile"`
	} `json:"data"`
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
}

// AuthData represents authentication data to return to frontend
type AuthData struct {
	Token      string `json:"token"`
	DeviceCode string `json:"deviceCode"`
	WToken     string `json:"wToken"`
	UCDE       string `json:"ucde"`
	User       User   `json:"user"`
}

// User represents user information
type User struct {
	MemberCode string `json:"memberCode"`
	Name       string `json:"name"`
	Phone      string `json:"phone"`
}

// StoreListResponse represents store list response from Tomoro API
type StoreListResponse struct {
	Data struct {
		Records []Store `json:"records"`
		Total   int     `json:"total"`
		Size    int     `json:"size"`
		Current int     `json:"current"`
		Pages   int     `json:"pages"`
	} `json:"data"`
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
}

// Store represents a store/outlet
type Store struct {
	StoreCode                   string  `json:"storeCode"`
	StoreName                   string  `json:"storeName"`
	StorePicture                string  `json:"storePicture"`
	StorePhone                  string  `json:"storePhone"`
	StoreEmail                  *string `json:"storeEmail"`
	StoreAddress                string  `json:"storeAddress"`
	Longitude                   float64 `json:"longitude"`
	Latitude                    float64 `json:"latitude"`
	DeliveryRegionConfig        int     `json:"deliveryRegionConfig"`
	DeliveryRegionValue         string  `json:"deliveryRegionValue"`
	IsDelivery                  int     `json:"isDelivery"`
	IsThirdDelivery             int     `json:"isThirdDelivery"`
	BusinessStatus              int     `json:"businessStatus"`
	DeliveryBusinessStatus      int     `json:"deliveryBusinessStatus"`
	ThirdDeliveryBusinessStatus int     `json:"thirdDeliveryBusinessStatus"`
	Distance                    int     `json:"distance"`
}

// MenuListResponse represents menu list response from Tomoro API
type MenuListResponse struct {
	Data struct {
		MenuVos []MenuCategory `json:"menuVos"`
	} `json:"data"`
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
}

// MenuCategory represents a menu category
type MenuCategory struct {
	MenuCode    string    `json:"menuCode"`
	MenuName    string    `json:"menuName"`
	MenuDesc    *string   `json:"menuDesc"`
	MenuPicture string    `json:"menuPicture"`
	MenuSort    int       `json:"menuSort"`
	Items       []Product `json:"items"`
}

// Product represents a menu item/product
type Product struct {
	Code              string         `json:"code"`
	Name              string         `json:"name"`
	Desc              string         `json:"desc"`
	Type              int            `json:"type"`
	Status            int            `json:"status"`
	IsSellOut         int            `json:"isSellOut"`
	PictureUrls       string         `json:"pictureUrls"`
	PictureMiniUrls   string         `json:"pictureMiniUrls"`
	PictureMaxUrls    string         `json:"pictureMaxUrls"`
	PluCode           string         `json:"pluCode"`
	PluCodes          []ProductPlu   `json:"pluCodes"`
	Price             float64        `json:"price"`
	LinePrice         float64        `json:"linePrice"`
	CurrencyUnit      string         `json:"currencyUnit"`
	IsCanAddCart      int            `json:"isCanAddCart"`
	MenuCode          string         `json:"menuCode"`
	MenuSecondaryTag  *string        `json:"menuSecondaryTag"`
	IsCustomize       *int           `json:"isCustomize"`
	DefaultSelectName *string        `json:"defaultSelectName"`
	Picture           ProductPicture `json:"picture"`
}

// ProductPlu represents a product price/size option
type ProductPlu struct {
	Code      string  `json:"code"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	LinePrice float64 `json:"linePrice"`
}

// ProductPicture represents product images
type ProductPicture struct {
	Main  string `json:"main"`
	Mini  string `json:"mini"`
	Max   string `json:"max"`
	Color string `json:"color"`
}

// CartRequest represents add to cart request
type CartRequest struct {
	ProductCode string  `json:"productCode"`
	ProductName string  `json:"productName"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	StoreCode   string  `json:"storeCode"`
}

// VoucherListResponse represents voucher list response
type VoucherListResponse struct {
	Data struct {
		Coupons []Voucher `json:"coupons"`
	} `json:"data"`
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
}

// MemberVoucherListResponse represents member voucher response from Tomoro API
type MemberVoucherListResponse struct {
	Data struct {
		Records []MemberVoucher `json:"records"`
	} `json:"data"`
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
}

// MemberVoucher represents a voucher owned by the user
type MemberVoucher struct {
	AccountCouponCode      string   `json:"accountCouponCode"`
	CouponCode             string   `json:"couponCode"`
	CouponName             string   `json:"couponName"`
	CouponDesc             string   `json:"couponDesc"`
	CouponPicture          *string  `json:"couponPicture"`
	CouponActiveStartTime  string   `json:"couponActiveStartTime"`
	CouponActiveEndTime    string   `json:"couponActiveEndTime"`
	DiscountEffectivePrice *float64 `json:"discountEffectivePrice"`
	DiscountEffectiveQty   *int     `json:"discountEffectiveQty"`
	CurrencyUnit           string   `json:"currencyUnit"`
	CouponContext          string   `json:"couponContext"`
	CouponStatus           int      `json:"couponStatus"`
	IsValidity             bool     `json:"isValidity"`
}

// Voucher represents a coupon/voucher
type Voucher struct {
	CouponCode     string  `json:"couponCode"`
	CouponName     string  `json:"couponName"`
	DiscountType   int     `json:"discountType"`
	DiscountValue  float64 `json:"discountValue"`
	MinOrderAmount float64 `json:"minOrderAmount"`
	ExpiryDate     string  `json:"expiryDate"`
}

// OrderRequest represents order creation request
type OrderRequest struct {
	StoreCode   string      `json:"storeCode"`
	Items       []OrderItem `json:"items"`
	VoucherCode string      `json:"voucherCode,omitempty"`
	Payment     string      `json:"payment"`
}

// OrderItem represents an item in order
type OrderItem struct {
	ProductCode string  `json:"productCode"`
	ProductName string  `json:"productName"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
}

// OrderResponse represents order creation response
type OrderResponse struct {
	Data struct {
		OrderCode   string  `json:"orderCode"`
		TotalAmount float64 `json:"totalAmount"`
		Status      string  `json:"status"`
	} `json:"data"`
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
}

// TradeOrderPageResponse represents the order history response from Tomoro API
type TradeOrderPageResponse struct {
	Data struct {
		PageNo       int          `json:"pageNo"`
		PageSize     int          `json:"pageSize"`
		PagesTotal   int          `json:"pagesTotal"`
		RecordsTotal int          `json:"recordsTotal"`
		Records      []TradeOrder `json:"records"`
	} `json:"data"`
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
}

// TradeOrder represents a single order history entry
type TradeOrder struct {
	ID                   int64              `json:"id"`
	Code                 string             `json:"code"`
	Status               int                `json:"status"`
	CancelReasonMsg      string             `json:"cancelReasonMsg"`
	StoreCode            string             `json:"storeCode"`
	StoreName            string             `json:"storeName"`
	UserCode             string             `json:"userCode"`
	UserName             string             `json:"userName"`
	ItemTotalMoney       float64            `json:"itemTotalMoney"`
	DiscountMoney        float64            `json:"discountMoney"`
	PaymentMoney         float64            `json:"paymentMoney"`
	CreateTime           int64              `json:"createTime"`
	UpdateTime           int64              `json:"updateTime"`
	MemberCouponCode     string             `json:"memberCouponCode"`
	PaymentOrderCode     string             `json:"paymentOrderCode"`
	TransCode            string             `json:"transCode"`
	PayType              int                `json:"payType"`
	PayChannelCode       string             `json:"payChannelCode"`
	TradeOrderItemVOList []TradeOrderItemVO `json:"tradeOrderItemVOList"`
}

// TradeOrderItemVO represents an order line item from history
type TradeOrderItemVO struct {
	ItemCode        string  `json:"itemCode"`
	ItemName        string  `json:"itemName"`
	ItemPictureUrls string  `json:"itemPictureUrls"`
	Amount          int     `json:"amount"`
	PluPrice        float64 `json:"pluPrice"`
	DiscountPrice   float64 `json:"discountPrice"`
	PluTotalMoney   float64 `json:"pluTotalMoney"`
	PaymentPrice    float64 `json:"paymentPrice"`
}

// NormalizedOrderHistoryResponse is returned to the frontend
type NormalizedOrderHistoryResponse struct {
	Records []OrderHistory `json:"records"`
	Total   int            `json:"total"`
	Size    int            `json:"size"`
	Current int            `json:"current"`
	Pages   int            `json:"pages"`
}

// OrderHistory is a frontend-friendly order object
type OrderHistory struct {
	OrderCode      string             `json:"orderCode"`
	StoreCode      string             `json:"storeCode"`
	StoreName      string             `json:"storeName"`
	Status         string             `json:"status"`
	StatusCode     int                `json:"statusCode"`
	PaymentStatus  string             `json:"paymentStatus"`
	CreatedAt      string             `json:"createdAt"`
	TotalAmount    float64            `json:"totalAmount"`
	DiscountAmount float64            `json:"discountAmount"`
	PaymentAmount  float64            `json:"paymentAmount"`
	VoucherCode    string             `json:"voucherCode,omitempty"`
	Payment        string             `json:"payment,omitempty"`
	PaymentOrderCode string            `json:"paymentOrderCode,omitempty"`
	TransCode      string             `json:"transCode,omitempty"`
	PayType        int                `json:"payType,omitempty"`
	PayChannelCode string             `json:"payChannelCode,omitempty"`
	Items          []OrderHistoryItem `json:"items"`
}

// OrderHistoryItem is a frontend-friendly order item
type OrderHistoryItem struct {
	ProductCode string  `json:"productCode"`
	ProductName string  `json:"productName"`
	Image       string  `json:"image,omitempty"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	Subtotal    float64 `json:"subtotal"`
}
