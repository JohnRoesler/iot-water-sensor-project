# iot-water-sensor-project

## local development

```
docker-compose build
docker-compose up

curl localhost:8888 -v -d '{"timeStamp": 51300052528, "symbol": "dff", "volume": 277, "temperature": 235}'
```

## future enhancements

- metrics:
  - overall: errors, tps
  - request level: errors by code, per endpoint
- logging w/levels
- configuration from files, allow setting any relevant parameters like log level
- authentication: ideally oauth
- replace API with message queue allowing for message processing to take place async with processor retries, DLQ, etc.
