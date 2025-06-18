# Anti-pattern-API2

API written in Golang with some golang antipatterns.


## Run Database

    docker-compose up -d
    docker-compose down
    docker-compose down -v --rmi all
    docker-compose up --build -d



    docker exec -it anti-pattern-api2_db_1 psql -U user -d mydatabase


2025/06/18 03:16:12 Ping a la base de datos falló: failed to connect to `host=db2 user=user database=mydatabase`: hostname resolving error (lookup db2 on 127.0.0.11:53: server misbehaving)

2025/06/18 03:17:48 Ping a la base de datos falló: failed to connect to `host=db user=user database=mydatabase`: dial error (dial tcp 172.28.0.2:5432: connect: connection refused)


sera el multistager build?  por que no pasa en el de microservices? no agregar healtcheck en el docker-compose siempre ha funcionado.


      # - DSN=host=postgres_sqli port=5432 user=postgres password=password dbname=sqli sslmode=disable timezone=UTC connect_timeout=5
