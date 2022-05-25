# Ninja GO template

This solution aims to be the ported version of Ninja Rest API. It supports:
- Project structure, following the standard project layout, adapting to NRA.
- Developing in a Docker container, with live loading.
- A dummy Celery worker to simulate the AI pipeline.
- Global config using `viper`.
- A Worker Service to trigger Celery Tasks via Message Queue.
- 2 core entities just like the NRA (`job` and `user`).
- ORM implemented using `gorm`, separated into services and repositories.
- JWT authentication, including refresh mechanism.
- Common middlewares (CORS, Request ID, GZIP).
- Swagger documentation dynamically generated.
- Payload field validation using `validator/v10`.
- Some demo unit tests.

## Getting Started
```
make build && make up
```

Login first with:
```
curl --location --request POST 'http://localhost:8080/api/v1/auth/login' \
--form 'username="surge@paragon.com"' \
--form 'password="paragon"'
```

then start creating jobs:
```
curl --location --request POST 'http://localhost:8080/api/v1/jobs/' \
--header 'Authorization: Bearer <your_token>' \
--header 'Content-Type: application/json' \
--data-raw '{
    "user_id": 1
}'
```
or creates a new user.

## Testing
```
make test
```

## Known issues
- The second task sent by gocelery will be held by the Web API’s worker instance until it’s gracefully killed.
- No lazy operation for gocelery, and the framework is so outdated with open issues.
- No implementation and validation yet for uploading files, since the app implementation is just to check the feasibility of working with Celery workers.
- Mocking in Go behaves differently.
- No implementation for filter yet.

## References
- [Gin Examples](https://github.com/gin-gonic/examples)
