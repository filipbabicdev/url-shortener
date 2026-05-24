# URL Shortener

A URL shortener written in Go. Takes a long URL, returns a short code, and
redirects on lookup. Built as a portfolio project to demonstrate idiomatic
Go backend structure and a well-known system-design pattern (Base62 ID encoding).

**Live demo:** https://url-shortener-70oc.onrender.com

> Deployed on Render (Docker) with a managed PostgreSQL database on Neon.
> Free tier, so the first request after idle may take 30-50s to wake (cold start).

## Tech Stack

- **Language**: Go 1.25
- **Router**: [chi/v5](https://github.com/go-chi/chi)
- **Database**: PostgreSQL 16 (via [pgx/v5](https://github.com/jackc/pgx))
- **Infrastructure**: Docker + docker-compose; deployed on Render with managed Postgres (Neon)

## How it works

When you shorten a URL, the row is inserted into Postgres and receives an
auto-incrementing `BIGSERIAL` id. That id is then **Base62-encoded** to produce
the short code:

```
id 1        -> "1"
id 125      -> "21"
id 100000   -> "q0U"
```

Lookup reverses the path: the short code maps directly to a row via the stored
`short_code`, and the handler issues a `302` redirect to the original URL.

### Why Base62 over a random/hash code?

- **No collisions.** The code is derived from a unique primary key, so two URLs
  can never produce the same code. A random generator has to check-and-retry on
  collision; hashing has the same risk plus truncation tradeoffs.
- **Short and dense.** Base62 (`0-9A-Za-z`) packs more values per character than
  hex or Base36. A 7-character Base62 code covers ~3.5 trillion ids.
- **Deterministic and cheap.** Pure arithmetic, no extra storage, no lookup table.

### Known tradeoff: enumeration

Because codes are sequential ids in disguise, they are guessable — `/2` follows
`/1`. For a public production service you would typically add an offset, a
non-sequential id source (e.g. Snowflake), or a random salt per code to prevent
scraping. This MVP intentionally keeps the simple encoding and documents the
tradeoff rather than hiding it.

## Getting Started

### Prerequisites

- Docker + docker-compose

### Run locally

```bash
git clone https://github.com/filipbabicdev/url-shortener.git
cd url-shortener
cp .env.example .env
docker compose up --build
```

The schema is applied automatically on first run (the migration is mounted into
Postgres's `docker-entrypoint-initdb.d`). The server starts on `localhost:8090`.

### Endpoints

| Method | Endpoint   | Request body            | Response |
|--------|------------|-------------------------|----------|
| `POST` | `/shorten` | `{"url": "https://..."}`| `{"short_url": "http://localhost:8090/21"}` |
| `GET`  | `/{code}`  | —                       | `302` redirect to the original URL |
| `GET`  | `/health`  | —                       | `200 OK` |

Example:

```bash
# shorten
curl -X POST http://localhost:8090/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://github.com/filipbabicdev"}'

# follow the redirect
curl -iL http://localhost:8090/1
```

## Design Decisions

**Insert-then-encode in one transaction** — the short code depends on the id,
but the id only exists after insert. `Create` inserts the row, reads back the
generated id via `RETURNING`, encodes it, and updates `short_code` — all inside a
single transaction so a failure can't leave a row without its code.

**Repository pattern** — handlers depend on a `URLRepository`, not on the pool
directly. Storage can change without touching HTTP code.

**Pure-Go driver (pgx) + static binary** — `CGO_ENABLED=0` yields a static
binary on a minimal alpine image, which keeps the deployed container small and
the build simple (no C toolchain).

**Healthcheck-gated startup** — in docker-compose the app waits until Postgres
reports healthy before connecting, avoiding a startup race.

## Roadmap

- [x] Deploy to Render with managed Postgres (Neon)
- [ ] Graceful shutdown (handle `SIGTERM`, drain in-flight requests)
- [ ] Click analytics (hit counter per code)
- [ ] Rate limiting on `/shorten`
- [ ] Custom/vanity codes

## Author

Filip Babić — [GitHub](https://github.com/filipbabicdev) · [LinkedIn](https://linkedin.com/in/rooky)