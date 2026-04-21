# Cinema Ticket Booking

## What is this?

This project is a simple **movie and screening booking web application**.

- **Backend**: Go + Gorilla Mux
- **Frontend**: HTML/CSS/JavaScript
- **Database**: PostgreSQL
- **Docker**: Runs with Docker Compose

## Main features

- Registration and login
- Listing movies and screenings
- Creating ticket reservations
- Viewing and deleting reservations
- Deleting movies or screenings

## Quick start with Docker Compose

**If it's not present already,** create a `.env` file in the project root before starting the stack. Example:

```env
POSTGRES_USER=cinemauser
POSTGRES_PASSWORD=cinemapassword
POSTGRES_DB=cinemadb
DB_HOST=postgres
DB_PORT=5432
BACKEND_PORT=8080
FRONTEND_PORT=80
POSTGRES_PORT=5432
```

```bash
docker compose up -d --build
```

Published ports come from the `.env` file above:

- Frontend: `FRONTEND_PORT` -> container `80`
- Backend: `BACKEND_PORT` -> container `BACKEND_PORT`
- PostgreSQL: `POSTGRES_PORT` -> container `5432`

Then open the frontend in your browser:

```bash
http://localhost
```

## Makefile usage

- `make build` - build Docker images
- `make up` - start containers
- `make down` - stop containers
- `make clean` - remove containers and data
- `make test` - run backend tests
- `make backend-build` - compile backend
- `make backend-run` - run backend locally without Docker

## API endpoints

All requests and responses use `Content-Type: application/json`.

- Direct backend base URL: `http://localhost:8080`
- Frontend-proxied API base URL: `http://localhost/api`

Protected routes require `Authorization: Bearer <token>`, using the token returned by `POST /api/login`.

In the current implementation, authenticated users can create and delete movies, screenings, and reservations. `GET /api/movies` and `GET /api/screenings` are public.

---

### Authentication

#### `POST /api/register`

Create a new user account.

**Request body**

```json
{
  "email": "user@example.com",
  "password": "secret123",
  "full_name": "Jane Doe"
}
```

**Response `200 OK`**

```json
{
  "id": 1,
  "email": "user@example.com",
  "full_name": "Jane Doe",
  "is_admin": false,
  "created_at": "2026-04-19T20:00:00Z"
}
```

| Status | Meaning |
| --- | --- |
| `200 OK` | User created |
| `400 Bad Request` | Missing field or email already registered |

---

#### `POST /api/login`

Authenticate an existing user.

**Request body**

```json
{
  "email": "user@example.com",
  "password": "secret123"
}
```

**Response `200 OK`**

```json
{
  "token": "jwt_token_placeholder",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "full_name": "Jane Doe",
    "is_admin": false,
    "created_at": "2026-04-19T20:00:00Z"
  }
}
```

| Status | Meaning |
| --- | --- |
| `200 OK` | Login successful |
| `401 Unauthorized` | User not found or incorrect password |

---

### Movies

#### `GET /api/movies`

Return all movies ordered alphabetically by title.

**Response `200 OK`**

```json
[
  {
    "id": 1,
    "title": "The Dark Knight",
    "description": "Batman raises the stakes...",
    "duration": 152,
    "genre": "Action",
    "poster_url": "https://..."
  }
]
```

---

#### `POST /api/movies`

Create a new movie.

**Request headers**

| Header | Type | Description |
| --- | --- | --- |
| `Authorization` | string | **Required** - `Bearer <token>` |

**Request body**

```json
{
  "title": "Inception",
  "description": "A thief who steals corporate secrets...",
  "duration": 148,
  "genre": "Sci-Fi",
  "poster_url": "https://..."
}
```

`poster_url` is optional. All other fields are required.

**Response `201 Created`** — returns the created `Movie` object.

| Status | Meaning |
| --- | --- |
| `401 Unauthorized` | Missing or invalid bearer token |
| `201 Created` | Movie created |
| `400 Bad Request` | Missing required field |

---

#### `DELETE /api/movies/{id}`

Delete a movie and all its associated screenings (cascade).

**Request headers**

| Header | Type | Description |
| --- | --- | --- |
| `Authorization` | string | **Required** - `Bearer <token>` |

| Status | Meaning |
| --- | --- |
| `401 Unauthorized` | Missing or invalid bearer token |
| `204 No Content` | Deleted successfully |
| `404 Not Found` | Movie does not exist |

---

### Screenings

#### `GET /api/screenings`

Return all screenings ordered by start time. Filter by movie with the optional query parameter.

| Query param | Type | Description |
| --- | --- | --- |
| `movie_id` | integer | Filter screenings for a specific movie |

**Response `200 OK`**

```json
[
  {
    "id": 1,
    "movie_id": 1,
    "starts_at": "2026-04-18T18:00:00Z",
    "available_seats": 80,
    "price": 2500.00
  }
]
```

---

#### `POST /api/screenings`

Create a new screening for an existing movie.

**Request headers**

| Header | Type | Description |
| --- | --- | --- |
| `Authorization` | string | **Required** - `Bearer <token>` |

**Request body**

```json
{
  "movie_id": 1,
  "starts_at": "2026-05-01T19:00:00Z",
  "available_seats": 60,
  "price": 2500.00
}
```

`starts_at` accepts ISO 8601 (`2006-01-02T15:04:05Z`) or short datetime (`2006-01-02T15:04`).

**Response `201 Created`** — returns the created `Screening` object.

| Status | Meaning |
| --- | --- |
| `401 Unauthorized` | Missing or invalid bearer token |
| `201 Created` | Screening created |
| `400 Bad Request` | Missing field, invalid date format, or movie not found |

---

#### `DELETE /api/screenings/{id}`

Delete a screening and all its associated reservations (cascade).

**Request headers**

| Header | Type | Description |
| --- | --- | --- |
| `Authorization` | string | **Required** - `Bearer <token>` |

| Status | Meaning |
| --- | --- |
| `401 Unauthorized` | Missing or invalid bearer token |
| `204 No Content` | Deleted successfully |
| `404 Not Found` | Screening does not exist |

---

### Reservations

#### `GET /api/reservations`

Return reservations for the authenticated user. If the authenticated user has `is_admin=true`, all reservations are returned.

**Request headers**

| Header | Type | Description |
| --- | --- | --- |
| `Authorization` | string | **Required** - `Bearer <token>` |

**Response `200 OK`**

```json
[
  {
    "id": 1,
    "screening_id": 2,
    "seats": 2,
    "created_at": "2026-04-19T20:15:00Z",
    "movie_title": "The Dark Knight",
    "starts_at": "2026-04-18T18:00:00Z",
    "price": 5000.00
  }
]
```

Note: `price` is the total for the reservation (`seat_price × seats`).

| Status | Meaning |
| --- | --- |
| `401 Unauthorized` | Missing or invalid bearer token |
| `200 OK` | Reservations returned |

---

#### `POST /api/reservations`

Book seats for a screening. Decrements `available_seats` atomically.

**Request headers**

| Header | Type | Description |
| --- | --- | --- |
| `Authorization` | string | **Required** - `Bearer <token>` |

**Request body**

```json
{
  "screening_id": 1,
  "seats": 2
}
```

**Response `201 Created`**

```json
{
  "id": 5,
  "screening_id": 1,
  "seats": 2,
  "created_at": "2026-04-19T20:15:00Z",
  "movie_title": "",
  "starts_at": "0001-01-01T00:00:00Z",
  "price": 0
}
```

| Status | Meaning |
| --- | --- |
| `201 Created` | Reservation created |
| `401 Unauthorized` | Missing or invalid bearer token |
| `400 Bad Request` | Screening not found, invalid seat count, or not enough seats |

---

#### `DELETE /api/reservations/{id}`

Delete a reservation. Does **not** restore `available_seats` on the screening.

**Request headers**

| Header | Type | Description |
| --- | --- | --- |
| `Authorization` | string | **Required** - `Bearer <token>` |

| Status | Meaning |
| --- | --- |
| `401 Unauthorized` | Missing or invalid bearer token |
| `204 No Content` | Deleted successfully |
| `404 Not Found` | Reservation does not exist |

## Note

Frontend `/api` calls are proxied through Nginx to the backend, so the UI runs on its own host port.
