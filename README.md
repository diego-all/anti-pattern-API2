# Anti-pattern-API2

API written in Golang with some golang antipatterns.

> Bad Golang Scaffolding "The Standard Go Project Structure"

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
    docker exec -it anti-pattern-api2_app_1 /bin/sh


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


status code
https://go.dev/src/net/http/status.go



“Accept interfaces, return structs.”
https://medium.com/capital-one-tech/doing-well-by-doing-bad-writing-bad-code-with-go-part-1-2dbb96ce079a
https://medium.com/capital-one-tech/doing-well-by-doing-bad-writing-bad-code-with-go-part-2-e270d305c9f7


repository
error handling

coupling cohession

Common Anti-Patterns in Go Web Applications
https://threedots.tech/post/common-anti-patterns-in-go-web-applications/


DEPLOYMENT

Docker-compose

artifact gcR

Autenticación de Docker en la VM: El comando gcloud auth configure-docker se ejecuta en el runner de GitHub Actions, no en la instancia de GCE. La instancia de GCE necesita su propia forma de autenticarse con Docker. Esto se logra mediante los scopes y la cuenta de servicio.


docker exec -it runner-db-1 psql -U user -d mydatabase



# SQLi

fetch("https://35.227.95.135/instruments/vulnerable-sqli?id=3' OR ''='", {
  method: "DELETE",
})
.then(res => res.text())
.then(data => console.log(data))
.catch(err => console.error(err));

🔹 RESTer (Firefox)
🔹 Postman Web o Extensión Chrome


Instancia 'go-hardenized-app-instance' no existe. Creando nueva instancia con IP estática.
WARNING: You have selected a disk size of under [200GB]. This may result in poor I/O performance. For more information, see: https://developers.google.com/compute/docs/disks#performance.
Created [https://www.googleapis.com/compute/v1/projects/rare-lambda-415802/zones/us-east1-b/instances/go-hardenized-app-instance].
WARNING: Some requests generated warnings:
 - Disk size: '20 GB' is larger than image size: '10 GB'. You might need to resize the root repartition manually if the operating system does not support automatic resizing. See https://cloud.google.com/compute/docs/disks/add-persistent-disk#resize_pd for details.
NAME                        ZONE        MACHINE_TYPE  PREEMPTIBLE  INTERNAL_IP  EXTERNAL_IP     STATUS
go-hardenized-app-instance  us-east1-b  e2-small                   10.142.0.7   35.211.130.133  RUNNING