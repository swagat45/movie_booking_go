# Movie Ticket Booking Backend (Go Implementation)

> **Note:** The original assignment asks for a Django/DRF project.  
> This repository implements the **same API design and flow** using **Go (Gin + GORM)** instead of Django, with JWT auth and Swagger documentation available at `/swagger/`.

---

## 📌 Overview

This project is a **Movie Ticket Booking System** backend with:

- User registration & login using **JWT authentication**
- APIs to:
  - List movies
  - List shows for a movie
  - Book a seat for a show
  - View user’s bookings
  - Cancel a booking
- **Business rules**:
  - No overbooking beyond total seats
  - No double booking of a seat for a show
  - Only the booking owner can cancel their booking
- **Swagger documentation** exposed at:  
  `http://localhost:8080/swagger/index.html`  
  (i.e. `/swagger/` on the server)

---

## 🛠 Tech Stack

- **Language**: Go
- **Framework**: Gin
- **ORM**: GORM
- **Database**: SQLite (pure Go driver `github.com/glebarez/sqlite`)
- **Auth**: JWT (`github.com/golang-jwt/jwt/v5`)
- **Password hashing**: bcrypt
- **API Docs**: Swagger (`swaggo`)

---

## 🚀 Setup Instructions (How to Run the Project)

### 1. Clone the repository

```bash
git clone <your-github-repo-url>
cd movie-booking-go
```

### 2. Install Go dependencies

```bash
go mod tidy
```

> This will download all required Go packages listed in `go.mod`.

### 3. Run the server

```bash
go run .
```

By default, the server starts at:

```text
http://localhost:8080
```

### 4. Health Check

You can verify the server is up using:

```bash
curl http://localhost:8080/health
```

Expected response:

```json
{"status":"ok"}
```

---

## 🔐 How to Generate JWT Tokens and Call APIs

All protected APIs use **JWT-based Bearer authentication**.

### Step 1: Sign up a user

**Endpoint**

```http
POST /api/signup
```

**Request body (JSON)**

```json
{
  "name": "Test User",
  "email": "test@example.com",
  "password": "secret123"
}
```

**Example (curl)**

```bash
curl -X POST http://localhost:8080/api/signup   -H "Content-Type: application/json"   -d '{
    "name": "Test User",
    "email": "test@example.com",
    "password": "secret123"
  }'
```

---

### Step 2: Login and obtain JWT token

**Endpoint**

```http
POST /api/login
```

**Request body (JSON)**

```json
{
  "email": "test@example.com",
  "password": "secret123"
}
```

**Example (curl)**

```bash
curl -X POST http://localhost:8080/api/login   -H "Content-Type: application/json"   -d '{
    "email": "test@example.com",
    "password": "secret123"
  }'
```

**Sample response**

```json
{
  "token": "<JWT_TOKEN_HERE>"
}
```

Copy this `token` value — you’ll use it in the `Authorization` header for protected APIs.

---

### Step 3: Calling protected APIs with JWT

For **all** protected endpoints, add this header:

```http
Authorization: Bearer <JWT_TOKEN_HERE>
```

Example:

```bash
curl -X GET http://localhost:8080/api/my-bookings   -H "Authorization: Bearer <JWT_TOKEN_HERE>"
```

---

## 📚 API Flow & Endpoints

Below is the typical flow and the corresponding endpoints.

### 1️⃣ View all movies (Public)

```http
GET /api/movies
```

Example:

```bash
curl http://localhost:8080/api/movies
```

---

### 2️⃣ View shows for a movie (Public)

```http
GET /api/movies/{movie_id}/shows
```

Example:

```bash
curl http://localhost:8080/api/movies/1/shows
```

---

### 3️⃣ Book a seat for a show (Protected)

```http
POST /api/shows/{show_id}/book
```

**Headers**

```http
Authorization: Bearer <JWT_TOKEN_HERE>
Content-Type: application/json
```

**Body**

```json
{
  "seat_number": 10
}
```

**Example**

```bash
curl -X POST http://localhost:8080/api/shows/1/book   -H "Authorization: Bearer <JWT_TOKEN_HERE>"   -H "Content-Type: application/json"   -d '{
    "seat_number": 10
  }'
```

**Business rules enforced**

- `seat_number` must be within `1..TotalSeats` for that show.
- Cannot book a seat that is already booked for that show.
- Cannot book more seats than the total capacity.

---

### 4️⃣ View user’s own bookings (Protected)

```http
GET /api/my-bookings
```

**Headers**

```http
Authorization: Bearer <JWT_TOKEN_HERE>`
```

**Example**

```bash
curl http://localhost:8080/api/my-bookings   -H "Authorization: Bearer <JWT_TOKEN_HERE>"
```

Returns a list of bookings for the logged-in user, with show and movie details preloaded.

---

### 5️⃣ Cancel a booking (Protected)

```http
POST /api/bookings/{booking_id}/cancel
```

**Headers**

```http
Authorization: Bearer <JWT_TOKEN_HERE>
```

**Example**

```bash
curl -X POST http://localhost:8080/api/bookings/1/cancel   -H "Authorization: Bearer <JWT_TOKEN_HERE>"
```

**Rules**

- Only the user who created the booking can cancel it.
- Only bookings with status `"booked"` can be cancelled.
- On cancellation, status is updated to `"cancelled"`.

---

## 📄 Swagger Documentation (Available at `/swagger/`)

Swagger UI is integrated using **swaggo** and exposes the API documentation.

- **URL:**  
  `http://localhost:8080/swagger/index.html`

This corresponds to the `/swagger/` path described in the assignment.

From Swagger UI you can:

- See all endpoints and their request/response schemas
- View models (SignupRequest, LoginRequest, LoginResponse, BookSeatRequest)
- Manually try out APIs (you can paste the JWT token in the `Authorization` header field)

---

## 🧪 Seed Data

On first run, the app seeds:

- **Movie**:  
  - Title: `Inception`  
  - Duration: `148` minutes  

- **Show** (for that movie):  
  - Screen: `Screen 1`  
  - Start time: `now + 2 hours`  
  - `TotalSeats = 50`

You can immediately test:

- `GET /api/movies`
- `GET /api/movies/1/shows`
- Then login and book seats for `show_id = 1`.

---

## ✅ Summary

This project implements the required **movie booking flow**, **JWT authentication**, and **Swagger documentation** (available at `/swagger/`) as per the assignment’s expectations, using Go as the backend stack.
