# iot-water-sensor-project

## stack

- [Go 1.17](https://go.dev/doc/go1.17), [GORM](https://gorm.io/docs/), & [Chi](https://github.com/go-chi/chi)
- [Docker](https://docs.docker.com/desktop/) & [Docker Compose](https://docs.docker.com/compose/)
- Shell
- [Postgres 10](https://www.postgresql.org/about/news/postgresql-10-released-1786/)

## local development

```bash
# builds the water sensor server's Dockerfile
docker-compose build

# brings up the postgres and water server (use `up -d` to run in the background)
docker-compose up

# single water reading reported via curl
curl localhost:8888 -v -d '{"timeStamp": 51300052528, "symbol": "dff", "volume": 277, "temperature": 235}'

## example response
...
< HTTP/1.1 201 Created
...
{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"timeStamp":51300052528,"symbol":"dff","volume":277,"temperature":235}

# script with input csv (populate input.csv with additional data) to mimic many sensors
./send_data.sh
```

## future enhancements

With enough time and if there was a need to build this into a production ready system, I would consider some additional items highlighted below.

- metrics:
  - overall: errors, tps
  - request level: errors by code, per endpoint
- tracing cost/benefit could be considered and would provide more benefit the more dependencies and interconnections come to be
- logging w/levels to keep logs to reasonable amount (in general errors only)
- configuration from files, allow setting any relevant parameters like log level
- authentication: ideally oauth
- replace API with message queue allowing for message processing to take place async with processor retries, DLQ, etc.
