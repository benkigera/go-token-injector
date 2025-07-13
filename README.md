# MQTT Go Application

This project provides a Go application that interacts with an MQTT broker, stores meter data in a PostgreSQL database, and exposes an API for both retrieving meter data and injecting credit tokens.

## Features

### Meter Data API

This API allows you to fetch the latest meter data. For more details, refer to `api/updater/README.md`.

### Token Injection API

This API facilitates the injection of credit tokens into smart meters. Unlike traditional APIs that might return an immediate response from an external service, this API implements a more robust and reliable flow:

1.  **Dynamic MQTT Subscription:** Upon receiving an injection request, the service first dynamically subscribes to a specific MQTT topic associated with the meter number provided in the request. This topic is where the actual token injection response will be published by the meter.
2.  **External API Call:** The service then sends an HTTP POST request to an external API (e.g., a meter management system) to initiate the token injection. This external API is responsible for forwarding the token to the physical meter.
3.  **Asynchronous Processing:** The external API typically returns a `202 Accepted` status code, indicating that the request has been received and is being processed asynchronously. It does *not* provide the final result of the injection.
4.  **Waiting for MQTT Response:** Crucially, after successfully sending the HTTP request, the `InjectionService` then waits for a message on the dynamically subscribed MQTT topic. This is where the meter's actual response (e.g., confirmation of injected units, success/failure status) will arrive.
5.  **Synchronous API Response (from MQTT):** Only when the MQTT response is received (or a configurable timeout occurs) will the `InjectToken` function return to the client. This ensures that the API consumer receives the definitive result of the token injection directly from the meter, rather than just an acknowledgment of the request being sent.

#### Why this approach is better:

This asynchronous, MQTT-driven approach provides several advantages:

*   **Reliability:** The API response directly reflects the meter's status, not just the external system's acknowledgment. This reduces ambiguity and provides a more accurate outcome.
*   **Real-time Feedback:** The client receives real-time feedback on the injection status as soon as the meter responds, without needing to poll for updates.
*   **Decoupling:** The external API and the meter communication are decoupled, making the system more resilient to failures in any single component.
*   **Scalability:** The MQTT broker can efficiently handle a large number of asynchronous messages, improving the scalability of the injection process.

## Usage

### Token Injection Example (POST Request)

To inject a token, send a POST request to the `/api/inject_token` endpoint. Replace `[YOUR_API_HOST]` with the actual hostname or IP address where your application is deployed.

```bash
curl -X POST http://[YOUR_API_HOST]:8082/api/inject_token \
--header 'Content-Type: application/json' \
--data '{
    "meter_number": "09000030529",
    "credit_token": "54158492421732529281"
}'
```

For more details on the request and response formats, refer to `api/injector/README.md`.

## Setup and Running

(Instructions for setting up and running the application will go here, including Docker Compose and direct execution.)

