# ActionLogger

User activity monitoring system

Goal: Implement a gRPC API for a user activity logging system with filtering and streaming support.

Requirements:

1. Service using Connect RPC:
Create a proto file describing the ActionLogger service.
Implement the server using connect-go.

2. PostgreSQL + pgx:  
Create a `user_actions` table: id, user_id, action_type, timestamp, details (JSONB).  
Use pgx v5 to work with the database.

3. Filtering:  
Implement aggregate queries with a combination of filters:
- by user_id  
- by action_type  
- by time range  
- search by details (JSONB field)

4. Streaming output:  
- For GetActions use server-side streaming  
- Limit output in batches of 100 records

5. Monitoring:  
- WatchActions must track new events (use PostgreSQL LISTEN/NOTIFY)

# Installation

```sh
cp .env.example .env
make install-deps
docker compose up
make local-migration-up
```

# Testing

```sh
go run cmd/server/main.go
go run cmd/client/main.go
```

 The [client/main.go](cmd/client/main.go) script runs sequentially:
 - LogAction
 - GetActions without filter (result: 10 records – 9 from migration [seed_user_actions](migrations/20250720001325_seed_user_actions.sql) and 1 logged by LogAction just now)
 - GetActions with filter (3 records from the migration match the filter, marked with comments in the [migration file](migrations/20250720001325_seed_user_actions.sql))
 - WatchActions

To test WatchActions run [LogAction.http](test/LogAction.http). Newly created records will be printed to the console.