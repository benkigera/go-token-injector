# Token Injection API

This API allows you to inject a credit token for a specific meter number.

## Endpoint

`POST /api/inject_token`

## Request Body

```json
{
    "meter_number": "<YOUR_METER_NUMBER>",
    "credit_token": "<YOUR_CREDIT_TOKEN>"
}
```

## Example cURL Request

```bash
curl -X POST http://localhost:8082/api/inject_token \
--header 'Content-Type: application/json' \
--data '{
    "meter_number": "09000030529",
    "credit_token": "42369210743477933997"
}'
```

## Response

The API will return the MQTT response received after the token injection, which typically includes details like `injected-units` and `credit-token-ack`.

```json
[
  {
    "n": "1P-Energy-Meter",
    "v": "09000030529",
    "t": 1755231168
  },
  {
    "n": "injected-units",
    "v": 2
  },
  {
    "n": "credit-token-ack",
    "v": 1
  }
]
```

