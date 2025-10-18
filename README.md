# Go Weather Realtime Check + Open Telemetry

This is a version of [go-weather](https://github.com/ivandersr/go-weather), decoupling the CEP validation from the CEP and Weather request module.

This app checks the realtime weather for the queried city, using its CEP.

### Local environment execution steps

Using docker compose, execute `docker compose up -d` in the root level of this project.

It is included a zipkin container service in order to collect the telemetry related to the services.

To query for the desired CEP, use a POST request to `http://localhost:8080/weather`, with the request body following the example:

```json
{
  "cep": "87225000"
}
```

The execution statistics will be available within zipkin. To access it, use the URL `http://localhost:9411`.
