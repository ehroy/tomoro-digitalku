# 📊 Analisis API Tomoro Coffee

## 🔍 Overview

File `api-service.tomoro-coffee.id_05-08-2026-23-57-21.proxymanlogv2` adalah Proxyman log yang berisi request/response dari API Tomoro Coffee.

## 📡 API Endpoint yang Terdeteksi

### **GET Store List**
```
GET /portal/app/basic/storeInfo/getStoreList/v3
Host: api-service.tomoro-coffee.id
```

**Query Parameters:**
- `centerPointLatitude`: -6.57398290625041
- `centerPointLongitude`: 110.6845192956732
- `pageNo`: 1
- `pageSize`: 20
- `storeName`: t (search query)

## 📦 Response Structure

### **Success Response (200 OK)**

```json
{
  "data": {
    "records": [
      {
        "storeCode": "JPR15A-M",
        "storeName": "Saudara Swalayan Tahunan",
        "storePicture": "",
        "storePhone": "08871341878",
        "storeAddress": "Indonesia/Central Java/Semarang/...",
        "longitude": 110.7038604,
        "latitude": -6.6319047,
        "isDelivery": 1,
        "businessStatus": 0,
        "distance": 6785
      }
    ],
    "total": 360,
    "size": 20,
    "current": 1,
    "pages": 18
  },
  "code": 0,
  "msg": "success",
  "success": true
}
```

## 🏗️ Data Model Analysis

### **Store Object Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `storeCode` | string | Unique store identifier |
| `storeName` | string | Store display name |
| `storePhone` | string | Contact phone number |
| `storeAddress` | string | Full address |
| `longitude` | float | GPS longitude |
| `latitude` | float | GPS latitude |
| `isDelivery` | int | Delivery available (0/1) |
| `businessStatus` | int | Store status (0=Open) |
| `distance` | int | Distance from center (meters) |

### **Status Codes:**

**businessStatus:**
- `0` = Open/Active
- `1` = Closed
- `2` = Temporarily Closed

## 🎯 Key Insights

1. **Pagination**: 20 items per page, 360 total stores
2. **Search**: By store name, case-insensitive
3. **Location**: GPS-based, sorted by distance
4. **Delivery**: Store delivery + third-party options

## 🔄 Implementation Recommendations

Backend perlu diupdate untuk match struktur API real ini.
