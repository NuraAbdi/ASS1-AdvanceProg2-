# Order & Payment Microservices

## Overview

This project implements two microservices:

* Order Service
* Payment Service

They communicate via REST API.

---

## Architecture

Each service follows Clean Architecture:

* domain
* usecase
* repository
* transport

---

## Communication

Order Service sends HTTP request to Payment Service.

If payment is successful → order becomes "Paid".
If failed → order becomes "Failed".

---

## Database

Order Service uses PostgreSQL.

Table: orders

* id
* customer_id
* item_name
* amount
* status
* created_at

---

## Endpoints

### Order Service

* POST /orders
* GET /orders/{id}
* PATCH /orders/{id}/cancel

### Payment Service

* POST /payments

---

## Business Rules

* Amount must be > 0
* If amount > 100000 → payment declined
* Only Pending orders can be cancelled

---

## Failure Handling

If Payment Service is unavailable → Order Service returns 503 error.

---

## Testing

Tested using Postman:

* Create order
* Get order
* Cancel order
* Database verification (PostgreSQL)

---

## Conclusion

The system demonstrates microservices architecture with REST communication and database persistence.
