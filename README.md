# 🚀 CDC Optimization Gateway

<div align="center">

**Lightweight CDC optimization layer between Debezium and Kafka, built with Go.**  
Reduces noisy CDC streams using **last-state coalescing**, **micro-batching**, and **policy-based routing**.

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://go.dev)
[![Kafka](https://img.shields.io/badge/Kafka-Required-black?style=flat&logo=apache-kafka)](https://kafka.apache.org)
[![Prometheus](https://img.shields.io/badge/Metrics-Prometheus-orange?style=flat&logo=prometheus)](https://prometheus.io)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

</div>

---

## 💡 What Is This?

Debezium can produce extremely noisy CDC streams, especially for high-update tables.  
This gateway reduces noise using:

- 🔄 **Last-state coalescing**
- 📦 **Micro-batching**
- 🧠 **Merge strategies**
- 🎛 **Table policies**
- 📊 **Prometheus metrics**

Aimed at reducing Kafka load and stabilizing downstream services.

---

## ✨ Features

- 🔄 Multiple UPDATEs → **one final event**
- 📦 Time/size-based micro-batching
- 🧠 Merge strategies: replace, merge_fields, append_history
- 🎛 Policy-based routing with conditions
- 📊 Prometheus metrics
- 🧩 Modular pluggable architecture

---

## 🚀 Quick Start

### Prerequisites

- Go 1.22+
- Kafka
- Docker (optional but easiest way)

### Clone & Build

```bash
git clone git@github.com:cihan-sahin/cdc-gateway.git
cd cdc-gateway
go mod download
go build -o bin/cdc-gateway ./cmd/gateway
```
 
### 🐳 Run with Docker Compose (Recommended)
```bash
docker-compose up -d
```

### ▶️ Run Manually
```bash
./bin/cdc-gateway
```

---

##  🌐 Access Points
 - Health: http://localhost:8080/health
 - Stats: http://localhost:8080/stats
 - Policies: http://localhost:8080/policies
 - Metrics: http://localhost:8080/metrics

 Example stats output:

```bash
{"events_in": 20, "events_out": 3}
```

---

## Send multiple UPDATE events on the same key:
```bash
go run tests/scenarios/test-producer.go \
  --broker kafka:9092 \
  --topic dbserver1.inventory.customers \
  --count 20 \
  --keys 1 \
  --sleep-ms 100
```


Check stats:
```bash
curl http://localhost:8080/stats
```

Expected:
  - events_in ≈ 20
  - events_out ≈ 1–5 (coalesced)

---

## ✅ Test 2: Micro-Batching
```bash
go run tests/scenarios/test-producer.go \
  --topic dbserver1.inventory.device_metrics \
  --count 100 \
  --keys 5
```
----

✅ Test 3: Condition Filtering
```bash
--flip-status
```

Only events matching policy conditions should pass.

---

## ⚙ Example config.yaml
```bash
app:
  name: "cdc-gateway"
  env: "dev"
  http_addr: ":8080"

kafka:
  brokers:
    - "kafka:9092"
  input_topic: "dbserver1.inventory.customers"
  consumer_group: "cdc-gateway-dev"

policies:
  - table: "inventory.customers"
    mode: "last_state"
    window_ms: 200
    max_batch_size: 50
    merge_strategy: "replace"
    target_topic: "inventory.customers.optimized"
    enabled: true
    conditions:
      - field: "status"
        op: "eq"
        value: "ACTIVE"

  - table: "inventory.orders"
    mode: "passthrough"
    window_ms: 0
    max_batch_size: 0
    merge_strategy: "replace"
    target_topic: "inventory.orders.raw"
    enabled: true

  - table: "inventory.device_metrics"
    mode: "micro_batch"
    window_ms: 200
    max_batch_size: 20
    merge_strategy: "replace"
    target_topic: "inventory.device_metrics.batched"
    enabled: true
    conditions:
      - field: "type"
        op: "eq"
        value: "SNMP_POLL"
```
---

## 🏗 Architecture
```bash
┌─────────────────────┐
│   PostgreSQL / DB   │
└─────────┬───────────┘
          │ CDC Stream
          ▼
┌─────────────────────┐
│      Debezium       │
└─────────┬───────────┘
          │ Kafka Messages
          ▼
┌──────────────────────────────────────────────┐
│           CDC Optimization Gateway            │
│                                              │
│  Consumer → Router → Coalescer → Merge → Producer
│                                              │
│       Prometheus · HTTP API · gRPC API       │
└──────────────────────────────────────────────┘
          ▼
┌─────────────────────┐
│  Optimized Topics   │
└─────────────────────┘
```

---

## 📁 Project Structure
```bash
cdc-gateway/
├── cmd/gateway/         
├── config/
│   └── config.yaml
├── internal/
│   ├── api/
│   ├── config/
│   ├── kafka/
│   ├── policy/
│   ├── pipeline/
│   ├── metrics/
├── pkg/cdcmodel/
├── tests/scenarios/
└── proto/cdcgateway/v1/
```
--- 

## 🎯 Use Cases

  - Reduce noisy CDC streams
  - Protect microservices from update storms
  - Reduce Kafka storage & compute cost
  - Batch metrics/telemetry tables
  - Apply per-table intelligent routing

## 🤝 Contributing
```bash
fork → branch → commit → PR
```
---

## 📄 License

MIT License

<div align="center">

✨ If this project helps you optimize CDC pipelines, please ⭐ the repo!
Made with ❤️ using Go, Kafka, and Debezium.

</div>