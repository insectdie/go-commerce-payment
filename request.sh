curl --location 'localhost:9990/gateway-service/signup' \
--header 'Content-Type: application/json' \
--data '{
    "email": "",
    "username": "fatannajuda",
    "password": "fatannajuda",
    "role": "Admin",
    "address": "Jakarta"
}'

curl --location 'localhost:9990/gateway-service/signin' \
--header 'Content-Type: application/json' \
--data '{
    "username":"fatannajuda",
    "password":"test"
}'