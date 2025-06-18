# Anti-pattern-API2

API written in Golang with some golang antipatterns.

    ├── db
    │   ├── db.go
    │   ├── Dockerfile
    │   └── init.sql
    ├── go.mod
    ├── go.sum
    ├── handlers
    │   └── instrument_handler.go
    ├── main.go
    ├── models
    │   └── instrument.go
    ├── docker-compose.yml
    ├── Dockerfile
    ├── README.md
    └── request.md



## Run Database

    docker-compose up -d
    docker-compose down
    docker-compose down -v --rmi all
    docker-compose up --build -d



    docker exec -it anti-pattern-api2_db_1 psql -U user -d mydatabase


**: dial error (dial tcp 172.28.0.2:5432: connect: connection refused)** reintentos del msa, healthceck from docker-compose.yml

2025/06/18 03:16:12 Ping a la base de datos falló: failed to connect to `host=db2 user=user database=mydatabase`: hostname resolving error (lookup db2 on 127.0.0.11:53: server misbehaving)

2025/06/18 03:17:48 Ping a la base de datos falló: failed to connect to `host=db user=user database=mydatabase`: dial error (dial tcp 172.28.0.2:5432: connect: connection refused)


sera el multistager build?  por que no pasa en el de microservices? no agregar healtcheck en el docker-compose siempre ha funcionado.


      # - DSN=host=postgres_sqli port=5432 user=postgres password=password dbname=sqli sslmode=disable timezone=UTC connect_timeout=5

**Con multistage build**

REPOSITORY                                 TAG               IMAGE ID       CREATED          SIZE
anti-pattern-api2_db                       latest            574ef4010805   9 minutes ago    274MB
anti-pattern-api2_app                      latest            507276655e42   40 minutes ago   19.9MB


**Sin multistage build**

REPOSITORY                                 TAG               IMAGE ID       CREATED              SIZE
anti-pattern-api2_app                      latest            68f660ee482a   About a minute ago   468MB
anti-pattern-api2_db                       latest            574ef4010805   55 minutes ago       274MB


**despues de cambiar al bind mount la primera vez sucedio el connection refussed y la segunda vez no sucedio y conecto inmediatamenete**

RELACION VOLUMEN (Named|Bind) + Dockerfiles (Normal|Multistage) + Driver (Pool|Native)

Docker compose toma el up.sql o init de la ruta directamente en el volumen (Named Volumen):

    volumes:
      # Monta el archivo init.sql para que PostgreSQL lo ejecute al iniciar
      # - ./sql/init.sql:/docker-entrypoint-initdb.d/init.sql


2 veces , despues de creado el volumen sube mas rapido y no hay connection refused.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
