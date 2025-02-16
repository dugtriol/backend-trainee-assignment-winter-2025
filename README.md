# backend-trainee-assignment-winter-2025

## Запустить проект

```
docker compose up
```
или
```
make up
```

## curl запросы:

#### /api/auth
```
curl -X POST "http://localhost:8080/api/auth" \
     -H "Content-Type: application/json" \
     -d '{
           "username": "user",
           "password": "1234"
         }'
```

#### /api/buy/{item}
```
curl -X GET "http://localhost:8080/api/buy/{item}" \
     -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE"
```

#### /api/sendCoin
```
curl -X POST "http://localhost:8080/api/sendCoin" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE" \
     -d '{
           "toUser": "user_id",
           "amount": 2
         }'
```

#### /api/info
```
curl -X GET "http://localhost:8080/api/info" \
     -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE"
```

P.S. При запуске интеграционных тестов нужно запустить БД (ту же, в докере)
