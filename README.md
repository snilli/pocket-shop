# Pocket Shop

Instant gift card ordering service — สั่งซื้อ gift card ผ่าน EZ API แบบ instant

## Prerequisites

-   Go 1.26+
-   PostgreSQL 16+
-   Docker & Docker Compose (optional)

## Quick Start

### 1. ด้วย Docker Compose (แนะนำ)

```bash
# จาก root project (ez/)
cp pocket-shop/.env.example pocket-shop/.env
# แก้ค่า EZ_API_KEY, EZ_ACCESS_TOKEN, EZ_SKU ใน .env

docker compose up pocket-shop
```

จะ start ทั้ง PostgreSQL + API server ให้อัตโนมัติที่ `http://localhost:8080`

### 2. รัน Local

```bash
cd pocket-shop/

# ติดตั้ง dependencies
go mod download

# copy และแก้ไข env
cp .env.example .env
# แก้ POSTGRES_HOST=localhost และค่า EZ credentials

# รัน API server
make run-api
```

## Environment Variables

ดูตัวอย่างทั้งหมดใน `.env.example`

| Variable                            | Description                                     | Default                  |
| ----------------------------------- | ----------------------------------------------- | ------------------------ |
| `POSTGRES_HOST`                     | PostgreSQL host                                 | `localhost`              |
| `POSTGRES_PORT`                     | PostgreSQL port                                 | `5432`                   |
| `POSTGRES_USER`                     | PostgreSQL user                                 | `ez`                     |
| `POSTGRES_PASSWORD`                 | PostgreSQL password                             | `ez`                     |
| `POSTGRES_DB`                       | PostgreSQL database                             | `ez`                     |
| `SERVER_PORT`                       | API server port                                 | `8080`                   |
| `SERVER_MODE`                       | `debug` (console log) / `production` (JSON log) | `debug`                  |
| `EZ_BASE_URL`                       | EZ API base URL                                 | `https://api.ezcards.io` |
| `EZ_API_KEY`                        | EZ API key                                      | **required**             |
| `EZ_ACCESS_TOKEN`                   | EZ access token                                 | **required**             |
| `EZ_SKU`                            | EZ product SKU                                  | **required**             |
| `ORDER_FULFILLMENT_TIMEOUT_SECONDS` | Timeout สำหรับรอ order สำเร็จ                   | `60`                     |
| `DISCOVER_INTERVAL_SECONDS`         | Interval สำหรับ discover job (วินาที)           | `5`                      |

## API Endpoints

| Method | Path                 | Description                 |
| ------ | -------------------- | --------------------------- |
| POST   | `/api/v1/orders`     | สร้าง order ใหม่            |
| GET    | `/api/v1/orders/:id` | ดูสถานะ order + redeem code |
| GET    | `/swagger`           | Swagger UI                  |

## Development

### Commands

```bash
make run-api           # รัน API server
make build-api         # Build binary
make swagger           # Generate Swagger docs
make mock              # Generate mocks (mockery)
make test              # รัน unit tests
make test-integration  # รัน integration tests
make test-coverage     # รัน tests + coverage report
```

### Project Structure

```
pocket-shop/
├── cmd/api/main.go              # Entry point
├── config/                      # App configuration
├── docs/                        # Swagger docs (auto-generated)
├── internal/
│   ├── core/order/              # Business logic
│   │   ├── domain/              # Entity (Order, OrderStatus)
│   │   ├── repository.go        # Repository interfaces
│   │   ├── service.go           # Service interface
│   │   ├── repository/db/       # DB implementations
│   │   └── service/ordersvc/    # Service implementation
│   ├── delivery/http/handler/   # HTTP handlers
│   └── infrastructure/          # Database, Fiber server, EZ client
├── mock/                        # Generated mocks
├── Dockerfile
└── Makefile
```

### Tech Stack

-   **Web:** Fiber v3
-   **ORM:** Ent
-   **DI:** Uber fx
-   **Logging:** zerolog
-   **Testing:** Ginkgo/Gomega + Testcontainers
-   **Mocking:** Mockery v3

### AI Assistance

-   Boilerplate/scaffolding — สร้าง Dockerfile, docker-compose, .env.example
-   Debugging — ช่วยวิเคราะห์ log และหา root cause (เช่น discover job poll ซ้ำ)
-   Code review — ตรวจสอบ implementation ว่าครบตาม requirement

### Key Decisions (ตัดสินใจเอง)

1. **PostgreSQL + Ent ORM แทน in-memory** — เลือกใช้ persistent storage เพราะต้องการให้ available pool และ order history อยู่รอดข้าม restart ได้ แม้โจทย์อนุญาต in-memory โดยเลือก Ent เพราะเป็น ORM ระดับ code generation ที่ให้ type-safe query ทั้งหมด — ไม่ต้องเขียน raw SQL หรือ string-based query ลด runtime error ได้มาก
2. **Swagger Codegen สำหรับ EZ API client** — ใช้ Swagger Codegen (docker compose service) generate Go client จาก EZ OpenAPI spec แทนการเขียน HTTP client เอง ทำให้ได้ typed models/methods ที่ตรงกับ API spec เสมอ และสามารถ regenerate ได้ทันทีเมื่อ spec เปลี่ยน รองรับ generate client ได้หลายภาษา (Go, TypeScript, Java ฯลฯ) จาก spec เดียวกัน
3. **`FOR UPDATE ... SKIP LOCKED` สำหรับ reservation pool** — ใช้ row-level locking ของ PostgreSQL แทน application-level mutex เพื่อรองรับ concurrent requests โดยไม่ต้อง distributed lock
4. **Background discover job แยกจาก Create Order** — แทนที่จะ scan หา available EZ orders ตอน create ทุกครั้ง เลือกทำเป็น background routine ที่รันต่อเนื่อง เพื่อให้ Create Order path เร็วที่สุด (แค่ Pull จาก pool)
5. **Lazy status check ใน GET /orders/:id** — ตอน GET ถ้า order ยัง PROCESSING จะเช็ค EZ status + timeout ทันที แทนที่จะรอแค่ background job อัพเดท ทำให้ client ได้ข้อมูลล่าสุดเสมอ
6. **Clean Architecture แยก layer ชัดเจน** — domain/service/repository/delivery/infrastructure แยกออกจากกัน ทำให้ mock ง่าย test ง่าย และเปลี่ยน storage หรือ external API ได้โดยไม่กระทบ business logic

### Assumptions & Trade-offs

**Assumptions:**

-   **Single-process runtime** — ตามที่โจทย์กำหนด ไม่ต้องทำ distributed locking แต่ใช้ DB-level lock เผื่อไว้แล้วถ้าต้อง scale ในอนาคต
-   **EZ order = 1 code** — ระบบดึงแค่ redeem code แรกจาก EZ Get Codes API (first match) เพราะโจทย์เป็น instant order ที่มี 1 code ต่อ 1 order
-   **CANCELLED เป็น local concept** — การ cancel ไม่ได้ไปยกเลิก EZ order จริง เป็นแค่ "ระบบเลิกรอ" ตามที่โจทย์ระบุ

**Trade-offs:**

-   **PostgreSQL vs SQLite/in-memory** — เพิ่ม dependency (ต้องมี Postgres) แต่ได้ row locking, persistence, และ concurrent safety มาแลก
-   **Synchronous polling ใน Create Order** — client ต้องรอจน order complete หรือ timeout (สูงสุด 60s) ทำให้ response time นาน แต่ simple กว่าการทำ async + webhook
-   **Discover job poll ทุก order ที่ยังไม่ used** — ถ้ามี order ค้างเยอะจะ poll EZ API บ่อย แต่สำหรับ scope ของ assignment นี้ถือว่ารับได้
