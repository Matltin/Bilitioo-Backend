# Bilitioo Backend

This is the backend service for **Bilitioo**, a ticket reservation platform, built with **Go**, **Gin**, **PostgreSQL**, **Redis**, and **Asynq** for background tasks.

## 📦 Features

- 🔐 User Authentication (Sign-In, Log-In, Email Verification)
- 👤 Profile Management (Update/Get User Profile)
- 🏙️ City & Ticket Search
- 🎫 Ticket Reservation
- 💳 Payment Processing
- 🧾 Penalty Management
- 📊 Report Management (Create, Answer, Review)
- 🧵 Background Job Distribution with Asynq
- 🚀 Redis Integration

## 🛠️ Technologies Used

- **Golang**
- **Gin** – HTTP Web Framework
- **PostgreSQL** – Relational Database
- **sqlc** – Type-safe DB queries
- **Redis** – Session/Cache storage
- **Paseto** – Secure token handling
- **Asynq** – Background job queue
- **Docker** – Deployment and environment

## 📁 Project Structure (partial)

```
├── api/              # Gin router setup & handlers
├── db/               # SQL queries & database logic
├── redis/            # Redis client configuration
├── token/            # Token maker (Paseto)
├── util/             # Configuration & helper functions
├── worker/           # Background task processing
```

## 🚀 Getting Started

### Prerequisites

- Go 1.20+
- PostgreSQL
- Redis
- Docker (optional but recommended)

### Installation

```bash
git clone https://github.com/Matltin/Bilitioo-Backend.git
cd Bilitioo-Backend
go mod tidy
```

### Environment Setup

Create a `.env` file with the following configuration:

```makefile
DB_SOURCE=postgresql://user:password@localhost:5432/bilitioo
REDIS_ADDRESS=localhost:6379
TOKEN_SYMMETRIC_KEY=your-very-secret-key
```

Load config using `util.LoadConfig()` (already implemented).

### Run the Server

```bash
go run main.go
```

Server will start on `localhost:8080` (or the address you configure).

## 🧪 API Endpoints (sample)

### Public Routes

| Method | Endpoint       | Description         |
|--------|----------------|---------------------|
| POST   | /sign-in       | Register user       |
| POST   | /log-in        | Log in user         |
| GET    | /verify-email  | Verify email token  |

### Authenticated Routes

| Method | Endpoint                        | Description                       |
|--------|----------------------------------|-----------------------------------|
| PUT    | /profile                        | Update user profile               |
| GET    | /profile                        | Get user profile                  |
| GET    | /city                           | Get city list                     |
| POST   | /city                           | Search tickets by city            |
| POST   | /reservation                    | Create reservation                |
| GET    | /completedReservation           | Get completed reservations        |
| GET    | /allReservation                 | Get all reservations              |
| POST   | /payment                        | Handle payment                    |
| GET    | /ticket-penalties/:ticket_id    | Get ticket penalties              |
| PUT    | /penalty/:ticket_id             | Cancel reservation                |
| POST   | /report                         | Create report                     |
| PUT    | /report                         | Answer report                     |
| GET    | /report                         | Get reports                       |
| PUT    | /manage-report                  | Manage reports                    |
| GET    | /completed-tickets              | Get all completed tickets         |
| GET    | /notcompleted-tickets           | Get not completed tickets         |
| GET    | /tickets                        | Get all tickets                   |

## 🧩 Contributing

Contributions are welcome! Please fork the repo and submit a pull request.

## 📄 License

MIT License © 2025 [Matin](https://github.com/Matltin)


## 🔌 How to Connect Components

This backend is composed of multiple connected components that must be running for full functionality:

### 1. PostgreSQL Database

Ensure PostgreSQL is installed and running. You can configure the DB connection string in `.env` as:

```env
DB_SOURCE=postgresql://your_user:your_password@localhost:5432/bilitioo
```

### 2. Redis Server

Install Redis and start it on the default port `6379`, or change it in the `.env`:

```env
REDIS_ADDRESS=localhost:6379
```

### 3. HTTP Server

Configured via:

```env
HTTP_SERVER_ADDRESS=localhost:8080
```

The HTTP server starts and listens on this address using:

```go
server := api.NewServer(config, taskDistributor, Queries, redis)
server.Start(config.HTTPServerAddress)
```

### 4. Asynq Worker (Task Processor)

The background job processor for sending emails and other tasks is run with:

```go
go runTaskProcessor(config, redisOpt, Queries)
```

This uses:

- Redis for queue
- Mail config from `.env`

### 5. Asynq Scheduler (Cron Jobs)

A cron job runs every minute to clean up expired reservations:

```go
go runScheduler(redisOpt, taskDistributor)
```

This uses:

```go
scheduler.Register("* * * * *", worker.NewCleanExpiredReservationsTask())
```

---

Make sure Redis, PostgreSQL, and your `.env` config are correct before running `main.go`.


# 📘 Auth API with Redis Support

This document explains how the `registerUserRedis` and `loginUserRedis` endpoints work using Redis for caching and rate-limiting.

---

## 🔐 POST `/auth/register-redis`

Registers a new user using email or phone number. Uses Redis to **prevent rapid repeated attempts** and **cache user data**.

### ✅ Request JSON
```json
{
  "email": "user@example.com",
  "phone_number": "09123456789",
  "password": "yourStrongPassword"
}
```

- At least one of `email` or `phone_number` **must** be provided.
- `password` must be at least **8 characters**.

### 🛠 How It Works

1. Checks Redis if a recent signup attempt was made (rate-limiting).
2. Validates the email and phone number format.
3. Hashes the password securely.
4. Creates the user in PostgreSQL.
5. Initializes the user's profile.
6. Sends a verification email using Asynq.
7. Caches the user in Redis for 5 minutes.

### 🔁 Response (Success)
```http
HTTP/1.1 200 OK
```

### 🔁 Response (Error Example)
```json
{
  "error": "invalid email format"
}
```

---

## 🔐 POST `/auth/login-redis`

Logs in the user using either email or phone number. First checks Redis for faster login; falls back to database if cache is missed.

### ✅ Request JSON
```json
{
  "email": "user@example.com",
  "phone_number": "",
  "password": "yourStrongPassword"
}
```

- Provide either `email` or `phone_number`.
- `password` is required.

### 🛠 How It Works

1. Attempts to find user in Redis cache.
2. If found:
   - Verifies that email or phone is confirmed.
   - Checks password.
   - Returns JWT access token.
3. If not in cache:
   - Queries the database.
   - Performs same verification.
   - Caches the user in Redis for 5 minutes.

### 🔁 Response (Success)
```json
{
  "user": {
    "id": 1,
    "email": "user@example.com",
    "phone_number": "09123456789",
    "hashed_password": "hashed..."
  },
  "access_token": "your.jwt.token"
}
```

### 🔁 Response (Error Example)
```json
{
  "error": "verify your email first"
}
```

---

## 🧠 Redis Keys Used

- `signup:attempt:<email>:<phone>` – Temporary key for rate-limiting (TTL: 20s)
- `user:<email>:<phone>` – Cached user data (TTL: 5 minutes)

> ⚠️ Note: Users must verify their email or phone number before logging in.

---

# 📘 Redis-Enhanced Ticket Search API

This README explains how the `searchTickets` and `getTicketDetails` API endpoints work using Redis for caching and improving performance.

---

## 🚍 POST /tickets/search

Search for available tickets based on origin, destination, date, and vehicle type. Redis is used to cache search results and speed up repeated requests.

### ✅ Request JSON Format
```json
{
  "origin_city_id": 1,
  "destination_city_id": 2,
  "departure_date": "2025-06-25",
  "vehicle_type": "BUS"
}
```
- `origin_city_id`, `destination_city_id`: Must be valid city IDs (required).
- `departure_date`: Must be in the format `YYYY-MM-DD` (required).
- `vehicle_type`: One of `BUS`, `TRAIN`, or `AIRPLANE` (required).

### 🛠 How It Works
1. Builds a Redis key based on search parameters.
2. Tries to fetch results from Redis.
3. If cache hit:
   - Returns the cached ticket list.
4. If cache miss:
   - Queries the PostgreSQL database.
   - Saves the result to Redis with a 5-minute expiration.
   - Returns the fresh result.

### 🔁 Response JSON Format (Success)
```json
[
  {
    "ticket_id": 1,
    "origin": "Tehran",
    "destination": "Shiraz",
    "departure_time": "2025-06-25T08:00:00Z",
    "vehicle_type": "BUS",
    ...
  },
  ...
]
```

### 🔁 Response JSON Format (Error)
```json
{
  "error": "invalid date format, expected YYYY-MM-DD"
}
```

---

## 🎫 GET /tickets/:ticket_id/details

Get detailed information for a single ticket by ticket ID. Uses Redis to reduce load on the database.

### ✅ Request
- Path parameter: `ticket_id` (integer, required)
- Example: `/tickets/42/details`

### 🛠 How It Works
1. Generates a Redis cache key from the ticket ID.
2. Checks if ticket detail exists in Redis.
3. If found in Redis:
   - Returns the cached ticket detail.
4. If not found:
   - Queries the PostgreSQL database.
   - Builds a detailed response (based on vehicle type).
   - Caches it in Redis for 10 minutes.

### 🔁 Response JSON Format (Success)
```json
{
  "origin": "Tehran",
  "destination": "Mashhad",
  "departureTime": "2025-06-25T09:00:00Z",
  "arrivalTime": "2025-06-25T17:00:00Z",
  "amount": 350000,
  "capacity": 40,
  "vehicle_type": "AIRPLANE",
  "feature": "Economy",
  "status": "available",
  "flight_class": "ECONOMY",
  "airplane_name": "Boeing 737"
}
```

### 🔁 Response JSON Format (Error)
```json
{
  "error": "ticket not found"
}
```

---

## 🧠 Redis Cache Keys
- `search:<origin>:<destination>:<date>:<vehicle>` → ticket list (5 min)
- `ticket_details:<ticket_id>` → ticket detail (10 min)

> ⚠️ This caching improves performance but does not replace proper database consistency.

---


## 🎟️ Reservation API

These endpoints handle ticket reservations, including creation and retrieval of user-specific reservation data.

---

### 🔐 POST `/reservation`

**Description:**  
Creates a reservation for one or more unreserved tickets and processes the associated payment.

#### ✅ Request Body

```json
{
  "tickets": [101, 102, 103]
}
```

- `tickets`: Array of ticket IDs (must not be already reserved)

#### ⚙️ How It Works

1. Verifies ticket IDs and checks their reservation status.
2. Calculates the total price.
3. Creates a payment (deducts amount from the user and sends to the system).
4. Creates a reservation for each ticket with the associated payment ID.

#### 🔁 Success Response

```json
{
  "reservations": [
    {
      "user_id": 5,
      "ticket_id": 101,
      "payment_id": 12
    },
    {
      "user_id": 5,
      "ticket_id": 102,
      "payment_id": 12
    }
  ],
  "payment": {
    "id": 12,
    "from_account": 5,
    "to_account": "myself",
    "amount": 150000
  }
}
```

#### ❌ Error Responses

- `400 Bad Request` – Empty ticket list or already reserved ticket
- `403 Forbidden` – Duplicate reservation or payment issue
- `404 Not Found` – Ticket not found
- `500 Internal Server Error` – Server or database error

---

### 🔐 GET `/allReservation`

**Description:**  
Retrieves **all reservations** (both completed and pending) for the authenticated user.

#### 🔁 Success Response

```json
[
  {
    "id": 1,
    "user_id": 5,
    "ticket_id": 101,
    "payment_id": 12,
    "status": "PENDING"
  },
  {
    "id": 2,
    "user_id": 5,
    "ticket_id": 102,
    "payment_id": 12,
    "status": "COMPLETED"
  }
]
```

#### ❌ Error Responses

- `404 Not Found` – No reservations found
- `500 Internal Server Error` – Database query failed

---

### 🔐 GET `/completedReservation`

**Description:**  
Returns only the **completed reservations** for the authenticated user.

#### 🔁 Success Response

```json
[
  {
    "id": 2,
    "user_id": 5,
    "ticket_id": 102,
    "payment_id": 12,
    "status": "COMPLETED"
  }
]
```

#### ❌ Error Responses

- `404 Not Found` – No completed reservations
- `500 Internal Server Error` – Server/database error

---

> ⚠️ **Note**: All endpoints require authentication. Include the JWT access token in the `Authorization` header:
>
> ```
> Authorization: Bearer <your_jwt_token>
> ```

---


## 🧾 Report Management API

These endpoints allow users to create reports and admins to answer and manage them. They also handle updates to ticket reservations based on report outcomes.

---

### 📌 GET /report

**Description**:  
Fetch all submitted reports (admin access required).

**Response**:  
Returns a list of report objects.

```json
[
  {
    "id": 1,
    "reservation_id": 123,
    "request_type": "CANCELATION",
    "request_text": "Please cancel this ticket.",
    "response_text": null,
    "user_id": 10,
    "admin_id": null
  }
]
```

---

### 📌 POST /report

**Description**:  
Create a new report for a reservation.

**Request JSON**:

```json
{
  "request_text": "Please cancel my ticket.",
  "request_type": "CANCELATION",
  "reservation_id": 123
}
```

**Response**:  
Returns the created report object.

```json
{
  "id": 1,
  "reservation_id": 123,
  "request_type": "CANCELATION",
  "request_text": "Please cancel my ticket.",
  "user_id": 10
}
```

---

### 📌 PUT /report

**Description**:  
Admin answers a user report.

**Request JSON**:

```json
{
  "id": 1,
  "response_text": "Your request is accepted and your ticket will be canceled."
}
```

**Response**:  
Returns the updated report with response.

```json
{
  "id": 1,
  "response_text": "Your request is accepted and your ticket will be canceled.",
  "admin_id": 1001
}
```

---

### 📌 PUT /manage-report

**Description**:  
Admin updates a reservation’s status in response to a report.

**Request JSON**:

```json
{
  "reserevation_id": 123,
  "to_status_reservation": "CANCELED"
}
```

**How it Works**:

- Retrieves current reservation status.
- Retrieves reservation details including amount and user info.
- Based on the new status:
  - If refund is applicable, it updates the user's wallet and payment info.
  - Updates the reservation's status.
  - Logs the status change in `change_reservation`.

**Response**:  
Returns HTTP 200 OK with no body on success.

---

### 🔄 Internal Helper Logic

#### changeadd()

- Adds refund amount to user wallet.
- Updates payment record.
- Calls `chageWithOutAdd()` to finish status update and logging.

#### chageWithOutAdd()

- Updates reservation status in DB.
- Creates a new `change_reservation` entry to log the transition.

# Bilitioo Backend – API Endpoints

This document explains the functionality of selected API endpoints from the `api` package in the Bilitioo backend. Each method is described with its purpose, request parameters, and response structure.

---

## 1. `GET /tickets/:ticket_id/penalties`

### **Method:** `getTicketPenalties`

### **Purpose:**
Fetches the penalty details (e.g., percentage to refund) for canceling a ticket.

### **Request URI Parameter:**

| Name     | Type  | Required | Description     |
|----------|-------|----------|-----------------|
| ticket_id| int64 | Yes      | The ID of the ticket |

### **Response (200 OK):**
Returns penalty information for the specified ticket (usually includes fields like `BeforeDay` and `AfterDay`).

```json
{
  "ticket_id": 123,
  "before_day": 10,
  "after_day": 30
}
```

### **Errors:**
- `400 Bad Request` – Invalid or missing ticket ID.
- `404 Not Found` – No penalty data found for the ticket.
- `500 Internal Server Error` – Server or database error.

---

## 2. `DELETE /reservations/:ticket_id`

### **Method:** `cancelReservation`

### **Purpose:**
Cancels a reserved ticket and issues a refund after applying the applicable cancellation penalty.

### **Request URI Parameter:**

| Name     | Type  | Required | Description     |
|----------|-------|----------|-----------------|
| ticket_id| int64 | Yes      | ID of the reserved ticket |

### **Authentication Required:** ✅  
Must be authorized with a valid token (`UserID` from token payload is used to verify ownership).

### **Business Rules:**
- Only reservations with status `"RESERVED"` can be canceled.
- Cannot cancel if the departure time has already passed.
- Penalty amount depends on how close the cancellation is to the departure time.

### **Response (200 OK):**

```json
{
  "message": "CANCELED",
  "amount_refunded": 135000,
  "ticket_id": 123,
  "change_reservation": {
    "id": 456,
    "reservation_id": 123,
    "admin_id": 1,
    "user_id": 12,
    "from_status": "RESERVED",
    "to_status": "CANCELED",
    "created_at": "2025-06-12T09:30:00Z"
  }
}
```

### **Errors:**
- `400 Bad Request` – Invalid ticket ID, already canceled/resolved ticket, or missed cancellation deadline.
- `401 Unauthorized` – Token is valid but user does not own the ticket.
- `404 Not Found` – Ticket or reservation not found.
- `403 Forbidden` – Unique violation during change reservation creation.
- `500 Internal Server Error` – Error processing refund or database issue.

---

## 3. `GET /verify_email?id={id}&secret_code={code}`

### **Method:** `verifyEmail`

### **Purpose:**
Verifies a user's email using a secret code received via email.

### **Query Parameters:**

| Name        | Type   | Required | Description              |
|-------------|--------|----------|--------------------------|
| id          | int64  | Yes      | Verification record ID   |
| secret_code | string | Yes      | Secret code sent via email |

### **Process:**
1. Looks up the verification record using `id` and `secret_code`.
2. If found, marks the email as verified in the user record.
3. Invalidates any cached user data.

### **Response (200 OK):**

```json
{
  "message": "email successfully verified"
}
```

### **Errors:**
- `400 Bad Request` – Missing or invalid query parameters.
- `400 Bad Request` – Verification link is invalid or expired.
- `500 Internal Server Error` – Server error or failed database update.

---

## Notes

- All endpoints are implemented using [Gin](https://github.com/gin-gonic/gin) as the HTTP framework.
- Errors are formatted using `errorResponse(err)` (not included in this snippet).
- Internal DB logic handled using SQLC-generated query interfaces.

---


# Bilitioo Backend API

## 📦 `POST /pay-payment`

### 🔧 Description

The `payPayment` endpoint is responsible for updating a payment’s status and type, updating related reservation statuses, adjusting ticket statuses, logging changes, and (if applicable) crediting a user's wallet. It's commonly used after a user attempts to finalize a payment (e.g., by paying via wallet or credit card), and all associated entities (reservations, tickets, user activity) need to be updated in the database and cache.

---

### 🧾 Request Format

#### Endpoint
```
POST /pay-payment
```

#### Request Body (JSON)

```json
{
  "payment_id": 123,
  "type": "WALLET",
  "payment_status": "COMPLETED",
  "reservation_status": "RESERVED",
  "user_activity_id": 456
}
```

| Field              | Type   | Required | Description                                                                 |
|-------------------|--------|----------|-----------------------------------------------------------------------------|
| `payment_id`       | int64  | ✅ Yes   | ID of the payment record to be updated.                                    |
| `type`             | string | ✅ Yes   | Type of payment. Must be one of: `CASH`, `CREDIT_CARD`, `WALLET`, `BANK_TRANSFER`, `CRYPTO`. |
| `payment_status`   | string | ✅ Yes   | New payment status. Must be one of: `PENDING`, `COMPLETED`, `FAILED`, `REFUNDED`. |
| `reservation_status` | string | ✅ Yes | New reservation status. Must be one of: `RESERVED`, `RESERVING`, `CANCELED`, `CANCELED-BY-TIME`. |
| `user_activity_id` | int64  | ❌ No    | (Optional) If given, the user activity associated with this payment will also be updated to `PURCHASED`. |

---

### 📤 Response Format

```json
{
  "payment": {
    "id": 123,
    "status": "COMPLETED",
    "type": "WALLET",
    ...
  },
  "reservations": [
    {
      "id": 101,
      "status": "RESERVED",
      "ticket_id": 55,
      ...
    }
  ],
  "user_activity": {
    "id": 456,
    "status": "PURCHASED",
    ...
  }
}
```

| Field           | Type     | Description                                                   |
|----------------|----------|---------------------------------------------------------------|
| `payment`       | object   | Updated payment record.                                       |
| `reservations`  | array    | List of updated reservations associated with the payment.     |
| `user_activity` | object   | (If applicable) updated user activity record.                 |

---

### ⚙️ What This Method Does

1. **Validates Input:**
   - Ensures `payment_type`, `payment_status`, and `reservation_status` are valid.
   - Parses and binds incoming JSON payload.

2. **Updates Payment:**
   - Updates the payment record with the provided `type` and `status`.

3. **Fetches Related Reservations:**
   - Retrieves all reservations tied to the `payment_id`.

4. **Processes Each Reservation:**
   - Checks if the reservation is in `RESERVING` state.
   - Updates the reservation status to the new one provided (`RESERVED`, etc.).
   - Updates the associated ticket status.
   - Invalidates Redis cache for:
     - Ticket details.
     - All ticket search results (`search:*`).
   - Records the status change in a `change_reservation` log table.

5. **Wallet Logic (if applicable):**
   - If the payment type is `WALLET`, credits the user’s wallet balance with the total value of updated tickets.

6. **Updates User Activity (if provided):**
   - If `user_activity_id` is passed, sets its status to `PURCHASED`.

---

### 🛑 Error Responses

| HTTP Status Code | Reason                                                   |
|------------------|----------------------------------------------------------|
| `400 Bad Request` | Invalid input (e.g., invalid enum values).              |
| `404 Not Found`   | Payment or reservation not found.                       |
| `500 Internal Server Error` | Database or internal logic failure.            |

---

### ✅ Example Use Case

User completes a ticket purchase using their wallet. This endpoint:
- Marks the payment as `COMPLETED`
- Sets all the related reservations as `RESERVED`
- Changes ticket statuses accordingly
- Deducts wallet balance
- Updates user activity to `PURCHASED`

---