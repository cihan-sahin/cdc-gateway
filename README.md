# CDC Optimization Gateway

A high-performance Go service that sits between **Debezium** and **Kafka**, reducing CDC noise through:

- 🔄 **Last-state coalescing**
- 📦 **Micro-batching**
- 🧠 **Merge strategies**
- 🔧 **Runtime policy management via gRPC**
- 📊 **Prometheus metrics**

This project solves a real-world scaling problem in distributed systems:  
**“Debezium sends too many CDC events, overwhelming downstream consumers.”**

The gateway intelligently reduces CDC events while keeping semantic correctness.

---

## 📐 Project Directory Structure

├── cmd/
│ └── gateway/
├── deploy/
│ └── docker-compose.yml
├── gen/
│ └── cdcgateway/v1/
├── internal/
│ ├── api/
│ │ └── grpc/
│ ├── config/
│ ├── kafka/
│ ├── policy/
│ ├── pipeline/
│ ├── metrics/
├── pkg/
│ └── cdcmodel/
├── proto/
│ └── cdcgateway/v1/
├── tests/
│
├── config.yaml
├── go.mod
└── README.md

---

## 🔁 High-Level Data Flow

               ┌────────────────────────────────────┐
               │          PostgreSQL (Source)        │
               └───────────────┬────────────────────┘
                               │ Changes (INSERT/UPDATE/DELETE)
                 Debezium CDC │
                               ▼
               ┌────────────────────────────────────┐
               │     Debezium Connect → Kafka       │
               │   Topic: dbserver1.inventory.*     │
               └───────────────┬────────────────────┘
                               │ Raw CDC JSON event
                               ▼
               ┌────────────────────────────────────┐
               │         CDC Gateway (this)         │
               │  1) Kafka Consumer (input topic)   │
               │  2) Router → find table policy     │
               │  3) Coalescer / Batch Engine       │
               │  4) Merge Strategy (replace/merge) │
               │  5) Kafka Producer (target topic)  │
               └───────────────┬────────────────────┘
                               │ Reduced / Optimized events
                               ▼
               ┌────────────────────────────────────┐
               │     Optimized Kafka Topics         │
               └────────────────────────────────────┘



**Effect:**  
100 noisy events → **1 optimized event**  
or  
100 noisy events → **1 micro-batch**.

---

## 🚀 Getting Started

You can run the entire environment using Docker Compose:

- Kafka
- Zookeeper
- Postgres (Debezium example DB)
- Debezium Connect
- CDC Gateway (this service)

---

# 🐳 Run with Docker Compose (Recommended)

### 1️⃣ Build the Gateway Image

From project root:

```bash
docker build -t cdc-gateway:dev .

