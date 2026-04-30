# Go E-Commerce API (Production-Ready Backend Service)

A production-style e-commerce backend built with **Go**, **PostgreSQL**, and **Docker**, featuring **CI/CD with GitHub Actions**, **container registry with AWS ECR**, and **cloud deployment on AWS EC2**.

This project demonstrates real-world backend engineering practices including **layered architecture**, **database migrations**, **transactional workflows**, and **cloud-native deployment**.

---

## 🚀 Features

* RESTful API for products and orders
* Multi-item order support
* Transactional order processing
* Concurrency-safe stock management
* Automatic database migrations on startup
* Dockerised services for local and cloud environments
* CI/CD pipeline building and pushing images to AWS ECR
* Deployed on AWS EC2

---

## Tech Stack

* **Backend:** Go (net/http, chi router)
* **Database:** PostgreSQL
* **ORM/Query Layer:** sqlc (type-safe SQL)
* **Migrations:** Goose
* **Containerisation:** Docker & Docker Compose
* **CI/CD:** GitHub Actions
* **Cloud:** AWS EC2 + ECR

---

## Architecture

```text
Client → HTTP API (Go)
        → Service Layer (business logic)
        → Repository Layer (sqlc)
        → PostgreSQL
```

* **Handler layer** → HTTP + request validation
* **Service layer** → business logic + transactions
* **Repository layer** → SQL queries generated via sqlc

---

## 🔄 Order Flow

1. Client sends `POST /orders`
2. Service starts a database transaction
3. For each item:

   * Validate stock using atomic update
   * Insert order item
4. Commit transaction
5. Return structured order response

---

## Running Locally

### Prerequisites

* Docker
* Docker Compose

### Run

```bash
docker compose up --build
```

API will be available at:

```text
http://localhost:8080
```

---

## Deployment (AWS)

### CI/CD Flow

1. Push to `main`
2. GitHub Actions:

   * Builds Docker image
   * Pushes to AWS ECR
3. EC2 instance:

   * Pulls latest image
   * Runs via Docker Compose

---

## Environment Variables

```env
DSN=host=db user=postgres password=postgres dbname=ecom sslmode=disable
```

---

## API Endpoints

### Products

```http
GET /products
GET /products/{id}
```

### Orders

```http
POST /orders
GET /orders/{id}
```

### Health

```http
GET /health
```

---

## Example Request

```json
POST /orders
{
  "customer_id": 1,
  "items": [
    {
      "product_id": 1,
      "quantity": 2,
      "price_cents": 1500
    }
  ]
}
```

---

## 📈 Future Improvements

* Authentication (JWT)
* Pagination for product listing
* Request validation layer
* Observability (logging, metrics)
* Infrastructure as Code (Terraform)
* Kubernetes deployment

---
