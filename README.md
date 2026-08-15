What we're building — URL Shortener
POST /shorten        → send a long URL, get back a short code (e.g. abc123)
GET  /:code          → visits the short URL → redirects to original long URL
GET  /stats/:code    → see how many times that short URL was clicked (real-time counter)

The real-time part you asked about: every time someone visits a short URL, Redis instantly increments a hit counter. You can hit /stats/:code and see live click counts updating in real time — that's Redis doing sub-millisecond reads/writes, way faster than Postgres for this use case.

Why each tool is used here
Tool	Role	Why
Go + Gin	API server	same as before
Postgres	Store URLs permanently	code → original_url mapping lives here forever
Redis	Hit counters + cache	super fast in-memory store — perfect for counters and caching frequently visited URLs
Nginx	Reverse proxy	sits in front of Go, handles all incoming traffic
Docker Compose	Run everything together	one command spins up Go + Postgres + Redis + Nginx simultaneously
How Redis fits in — the real-time flow
User visits abc123
       ↓
Nginx forwards to Go
       ↓
Go checks Redis cache first (microseconds)
  → found in cache → redirect immediately (never hits Postgres)
  → not in cache   → fetch from Postgres → store in Redis → redirect
       ↓
Redis increments hit counter for abc123
(INCR abc123:hits → atomic, instant, thread-safe)
       ↓
GET /stats/abc123 reads counter from Redis
→ returns live click count
Project structure
url-shortener/
├── main.go
├── config/
│   └── db.go          ← Postgres connection
│   └── redis.go       ← Redis connection (NEW)
├── models/
│   └── url.go
├── repository/
│   └── url_repo.go    ← Postgres queries
├── cache/
│   └── url_cache.go   ← Redis operations (NEW layer)
├── services/
│   └── url_service.go
├── handlers/
│   └── url_handler.go
├── utils/
│   └── shortcode.go   ← generates random short codes
├── routes/
│   └── routes.go
├── docker-compose.yml ← NEW — runs everything
├── nginx/
│   └── nginx.conf     ← NEW — Nginx config
├── Dockerfile         ← NEW — containerizes Go app
├── go.mod
└── .env
Build phases — exactly like before, one layer at a time
Phase 1 → Understand Docker Compose + write docker-compose.yml
Phase 2 → Nginx config (nginx.conf)
Phase 3 → Dockerfile for Go app
Phase 4 → Config layer (Postgres + Redis connections)
Phase 5 → Models + utils (shortcode generator)
Phase 6 → Repository (Postgres queries)
Phase 7 → Cache layer (Redis operations)
Phase 8 → Services (business logic — glues repo + cache)
Phase 9 → Handlers + Routes
Phase 10 → Run everything + test real-time hit counter
New concepts you'll learn per phase
Phase 1 → docker-compose.yml syntax
           services, networks, volumes, depends_on, environment vars

Phase 2 → nginx.conf
           proxy_pass, upstream, location blocks

Phase 3 → Dockerfile
           multi-stage builds, scratch images, COPY, EXPOSE

Phase 4 → Redis client in Go (go-redis)
           connecting, ping check

Phase 6 → Postgres raw SQL (like library project)
           INSERT, SELECT for URL storage

Phase 7 → Redis commands in Go
           SET with expiry (caching)
           INCR (atomic counter — the real-time part)
           GET (read counter)

Phase 8 → Cache-aside pattern
           check cache → miss → hit DB → populate cache
Tech stack locked
go
// go.mod packages
github.com/gin-gonic/gin          // HTTP framework
github.com/lib/pq                 // Postgres driver (raw SQL like library project)
github.com/redis/go-redis/v9      // Redis client
github.com/joho/godotenv          // .env loader
Docker Compose mental model (before we write code)

Think of docker-compose.yml like this — instead of running 4 separate docker run commands manually:

bash
# WITHOUT docker compose (what you've been doing)
docker run --name store-postgres ...
docker run --name redis ...
docker run --name nginx ...
docker run --name go-app ...

# WITH docker compose (one command)
docker compose up

All 4 containers start together, talk to each other on the same internal network, and shut down together with docker compose down. That's the entire point.

.env
env
PORT=8080
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=urlpass
DB_NAME=urldb
REDIS_HOST=redis
REDIS_PORT=6379
BASE_URL=http://localhost

Notice DB_HOST=postgres and REDIS_HOST=redis — not localhost anymore. Inside Docker Compose, containers talk to each other by service name, not localhost.
