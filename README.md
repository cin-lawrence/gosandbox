# Ninja GO template

This solution aims to be the ported version of Ninja Rest API. It supports:
- Project structure for easy maintenance and extending, adopting Golang project layout standard.
- Debugging in a Docker container with live reloading.
- Interactions with worker services written in both Go and Python.
- A global-wide configuration for services.
- ORM and Auto Migration from `gorm`.

More cool features to add on:
- JWT authentication with refresh token.
- Generating Swagger documentation.
- Payload validation using `validator/v10`.
- Testing, natively supported.

## Getting Started
```
docker-compose build && docker-compose up -d
```

Create a user first with:
```
curl --location --request POST 'http://localhost:8080/api/v1/users/' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "surge@paragon.com",
    "name": "Surge"
}'
```

then start creating jobs:
```
curl --location --request POST 'http://localhost:8080/api/v1/jobs/' \
--header 'Content-Type: application/javascript' \
--data-raw '{
    "user_id": 1
}'
```

## Known issues
- The 2nd job is held by `gocelery` due to some misconfiguration. It might take a long time to investigate.

## References
- [Gin Examples](https://github.com/gin-gonic/examples)
