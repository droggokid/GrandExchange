# GrandExchange

## Db stuff

`psql -h localhost -p 5432 -U postgres -d postgres

`\l                  -- list all databases
\c mydb             -- switch to a database
\dt                 -- list tables in current db
\d tablename        -- describe a table
\q                  -- quit`

As it is:
Check Redis
If hit, return it
If miss, check DB
If DB row is still fresh enough, return DB result and repopulate Redis
If DB row is stale, call API, upsert DB, update Redis, return fresh data

With pagination:
check Redis freshness marker for that item/query
if still fresh, read paginated data from DB
if stale/missing, refresh from API
upsert DB
update Redis freshness marker
return paginated DB results
