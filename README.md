# URL Shortener

- A `CLI` written in `Go` that shortens `URL`s in a scalable, performant, and **collision-free** way.

## Flow

**1. Create**

```mermaid
sequenceDiagram
    participant C as Operator
    participant S as CLI
    participant R as Redis
    participant P as Postgres

    C->>S: shorten (original_url)
    S->>R: INCR counter
    R-->>S: sequential id
    S->>S: Hash ID encode (Base62 alphabet + salt)
    S->>P: INSERT short_code, original_url
    P-->>S: OK
    S-->>C: short URL
```

**2. Read**

```mermaid
sequenceDiagram
    participant C as Browser
    participant S as HTTP Service
    participant P as Postgres

    Note over C,P: First request
    C->>S: GET /{code}
    S->>S: memory cache miss
    S->>P: SELECT original_url
    P-->>S: original_url
    S->>S: store in memory cache
    S-->>C: 302 Redirect + Cache-Control 1 year

    alt Browser cached the 302
        Note over C: Redirects locally, server is not hit
    else Client hits the server again
        C->>S: GET /{code}
        S->>S: memory cache hit
        S-->>C: 302 Redirect
        Note over P: Postgres is not queried
    end
```

## Why this design

`Redis INCR` atomically increments a counter and returns a unique integer.
Concurrent requests cannot race and receive the same `ID`.

That `ID` is then encoded with `Hash ID` using a [Base62](https://base62.org/) alphabet and a secret salt.
The salt shuffles the encoding so sequential `IDs` do not produce sequential, guessable short codes.

Only after the code is generated is it inserted into Postgres.
Because IDs are unique by construction, **the service never queries the database to check whether a short code already exists**.
That kind of lookup would be a bottleneck and a source of contention in a distributed system.

## Getting started

[Docker](https://docs.docker.com/get-docker/) is the only requirement. `Go` never runs on the host — it
compiles inside the image build stage.

**1. Configure**

```bash
cp env.local .env
```

Two fields in `.env` need attention:

- `POSTGRES_PASSWORD` ships empty. Fill it, or the `Postgres` container refuses to initialize.
- `HASHIDS_SALT` is what makes sequential `IDs` produce unguessable codes. Use a long random secret.

`POSTGRES_*`, `REDIS_URL`, and `APP_DOMAIN` work as they come for local use.

**2. Create the schema**

```bash
make migrate
```

Starts `Postgres`, waits until it reports healthy, then applies every file in `migrations/`.
This comes first: on an empty database the `CLI` fails with `relation "urls" does not exist`.

**3. Start the stack**

```bash
make up
```

Starts `Redis` and the `HTTP` service, building the app image on the first run.

**4. Confirm it works**

```bash
make list     # No shortened URLs
make shorten  # prompts for a URL, prints shortener.localhost/{code}
```

Opening the printed `URL` answers `302` to the original address, and an unknown code answers `404`.

## Commands

```bash
make help       # prints CLI usage
make shorten    # prompts for a URL, prints shortener.localhost/{code}
make load-test  # prompts for a URL and how many times to shorten it
make list       # lists the 20 most recent codes
make down       # stops containers; volumes (data) stay
```

## After changing the code

`make up` reuses the image that already exists, so a new binary needs a rebuild.
Skip it and the old binary keeps running against your new code.

```bash
make image
make up
```

## Starting over

`make down` keeps the volumes, so the data survives. To drop both the rows and the `Redis`
counter and come back to a clean slate:

```bash
docker compose down -v
make migrate
make up
```