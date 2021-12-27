Documentation available with Swagger UIs Docker Container

"docker run -p 80:8080 -e SWAGGER_JSON=/app/openapi.json -v ~/Desktop/api:/app swaggerapi/swagger-ui"


redis : docker run --name redis -p 6379:6379 -d redis

redis : docker exec -i -t a46(container) redis-cli
