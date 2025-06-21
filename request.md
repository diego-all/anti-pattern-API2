# Request




curl -X GET http://localhost:8080/instruments


curl -X GET http://localhost:8080/instruments/6


curl -X POST \
  http://localhost:8080/instruments \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "Flauta travesera",
    "description": "Flauta de plata con estuche rígido",
    "price": 450
  }'


curl -X PUT \
  http://localhost:8080/instruments/6 \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "Guitarra eléctrica actualizada",
    "description": "Fender Stratocaster deluxe, color azul",
    "price": 1350
  }'



curl -X DELETE http://localhost:8080/instruments/6 