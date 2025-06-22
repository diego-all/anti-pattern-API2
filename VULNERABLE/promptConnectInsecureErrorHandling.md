
Validando las credenciales me permite acceder a la DB y corroborar que si se estan insertando los datos:

root@pho3nix:/home/diegoall/MAESTRIA_ING/anti-pattern-API/anti-pattern-API2# docker exec -it anti-pattern-api2_db_1 psql -U user -d mydatabase
psql (15.13)
Type "help" for help.

mydatabase=# \dt
          List of relations
 Schema |    Name     | Type  | Owner 
--------+-------------+-------+-------
 public | instruments | table | user
(1 row)

mydatabase=# select * from instruments
mydatabase-# select * from instruments;
ERROR:  syntax error at or near "select"
LINE 2: select * from instruments;
        ^
mydatabase=# select * from instruments;
 id |        name        |                   description                   | price |         created_at         |         updated_at         
----+--------------------+-------------------------------------------------+-------+----------------------------+----------------------------
  1 | Guitarra eléctrica | Guitarra Fender Stratocaster de seis cuerdas    |  1200 | 2025-06-18 02:37:36.771265 | 2025-06-18 02:37:36.771265
  2 | Batería acústica   | Set completo de batería Pearl con platillos     |  2300 | 2025-06-18 02:37:36.771265 | 2025-06-18 02:37:36.771265
  3 | Teclado digital    | Yamaha con 88 teclas contrapesadas              |   850 | 2025-06-18 02:37:36.771265 | 2025-06-18 02:37:36.771265
  4 | Violín             | Violín acústico hecho a mano con arco y estuche |   600 | 2025-06-18 02:37:36.771265 | 2025-06-18 02:37:36.771265
  5 | Saxofón alto       | Saxofón profesional con boquilla y correa       |  1500 | 2025-06-18 02:37:36.771265 | 2025-06-18 02:37:36.771265
(5 rows)

Aca esta el archivo docker-compose.yaml

version: '3.8'

services:

  db:
    image: postgres:15-alpine # Usa una imagen oficial de PostgreSQL
    ports:
      - "5432:5432" # Opcional: Mapea el puerto de la DB al host (útil para herramientas externas)
    environment:
      POSTGRES_USER: user # Nombre de usuario de la base de datos
      POSTGRES_PASSWORD: password # Contraseña del usuario
      POSTGRES_DB: mydatabase # Nombre de la base de datos
    volumes:
      # Monta el archivo init.sql para que PostgreSQL lo ejecute al iniciar
      - ./sql/init.sql:/docker-entrypoint-initdb.d/init.sql
      # Persiste los datos de la base de datos en un volumen nombrado
      - pgdata:/var/lib/postgresql/data
  app:
    build:
      context: . # El contexto de construcción es el directorio actual
      dockerfile: Dockerfile # Usa el Dockerfile que acabamos de crear
    ports:
      - "8080:8080" # Mapea el puerto 8080 del host al puerto 8080 del contenedor
    environment:
      # Asegúrate de que estos valores coincidan con los de tu .env o sean los que esperas
      PORT: 8080
      DATABASE_URL: postgres://user:password@db:5432/mydatabase?sslmode=disable
    depends_on:
      - db # Asegura que la base de datos se inicie antes que la aplicación
    # Para desarrollo, puedes montar el código fuente para recargas en caliente si usas un observador
    # volumes:
    #   - .:/app # Descomenta para montar el volumen para desarrollo


volumes:
  pgdata: # Define el volumen nombrado para la persistencia de datos de PostgreSQL

Aca estan los logs:

diegoall@pho3nix:~/MAESTRIA_ING/anti-pattern-API/anti-pattern-API2$ sudo -s
[sudo] password for diegoall: 
root@pho3nix:/home/diegoall/MAESTRIA_ING/anti-pattern-API/anti-pattern-API2# docker-compose down -v --rmi all
/usr/lib/python3/dist-packages/paramiko/transport.py:237: CryptographyDeprecationWarning: Blowfish has been deprecated and will be removed in a future release
  "class": algorithms.Blowfish,
Stopping anti-pattern-api2_db_1 ... done
Removing anti-pattern-api2_app_1 ... done
Removing anti-pattern-api2_db_1  ... done
Removing network anti-pattern-api2_default
Removing volume anti-pattern-api2_pgdata
Removing image postgres:15-alpine
Removing image anti-pattern-api2_app
root@pho3nix:/home/diegoall/MAESTRIA_ING/anti-pattern-API/anti-pattern-API2# docker-compose up --build -d
/usr/lib/python3/dist-packages/paramiko/transport.py:237: CryptographyDeprecationWarning: Blowfish has been deprecated and will be removed in a future release
  "class": algorithms.Blowfish,
Creating network "anti-pattern-api2_default" with the default driver
Creating volume "anti-pattern-api2_pgdata" with default driver
Pulling db (postgres:15-alpine)...
15-alpine: Pulling from library/postgres
fe07684b16b8: Already exists
2777460b63f4: Pull complete
642e176e7683: Pull complete
b4dcca6808e5: Pull complete
77b69ff8bb36: Pull complete
45886f8a09ca: Pull complete
331cba96f288: Pull complete
6380a3c9c68c: Pull complete
f2ee91c57ab1: Pull complete
8e7dfe758b13: Pull complete
639ffb3d4c66: Pull complete
Digest: sha256:2985f77749c75e90d340b8538dbf55d4e5b2c5396b2f05b7add61a7d8cd50a99
Status: Downloaded newer image for postgres:15-alpine
Building app
[+] Building 6.9s (17/17) FINISHED                                                                                                               docker:default
 => [internal] load build definition from Dockerfile                                                                                                       0.0s
 => => transferring dockerfile: 1.43kB                                                                                                                     0.0s
 => [internal] load metadata for docker.io/library/alpine:latest                                                                                           0.9s
 => [internal] load metadata for docker.io/library/golang:1.23.2-alpine                                                                                    0.9s
 => [auth] library/golang:pull token for registry-1.docker.io                                                                                              0.0s
 => [auth] library/alpine:pull token for registry-1.docker.io                                                                                              0.0s
 => [internal] load .dockerignore                                                                                                                          0.0s
 => => transferring context: 2B                                                                                                                            0.0s
 => [stage-1 1/3] FROM docker.io/library/alpine:latest@sha256:8a1f59ffb675680d47db6337b49d22281a139e9d709335b492be023728e11715                             0.0s
 => [builder 1/6] FROM docker.io/library/golang:1.23.2-alpine@sha256:9dd2625a1ff2859b8d8b01d8f7822c0f528942fe56cfe7a1e7c38d3b8d72d679                      0.0s
 => [internal] load build context                                                                                                                          0.0s
 => => transferring context: 3.88kB                                                                                                                        0.0s
 => CACHED [builder 2/6] WORKDIR /app                                                                                                                      0.0s
 => CACHED [builder 3/6] COPY go.mod ./                                                                                                                    0.0s
 => CACHED [builder 4/6] RUN go mod download                                                                                                               0.0s
 => [builder 5/6] COPY . .                                                                                                                                 0.1s
 => [builder 6/6] RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go-app main.go                                                           5.5s
 => CACHED [stage-1 2/3] WORKDIR /root/                                                                                                                    0.0s
 => CACHED [stage-1 3/3] COPY --from=builder /go-app .                                                                                                     0.0s
 => exporting to image                                                                                                                                     0.0s
 => => exporting layers                                                                                                                                    0.0s
 => => writing image sha256:acec2a5248cd645815631e8a4c98ec7607ce347cc8d0dfba4c147cff555b4941                                                               0.0s
 => => naming to docker.io/library/anti-pattern-api2_app                                                                                                   0.0s
Creating anti-pattern-api2_db_1 ... done
Creating anti-pattern-api2_app_1 ... done
root@pho3nix:/home/diegoall/MAESTRIA_ING/anti-pattern-API/anti-pattern-API2# docker ps | grep anti-pattern
93206a0449e6   postgres:15-alpine                       "docker-entrypoint.s…"   13 seconds ago   Up 12 seconds             0.0.0.0:5432->5432/tcp, :::5432->5432/tcp                                                                  anti-pattern-api2_db_1
root@pho3nix:/home/diegoall/MAESTRIA_ING/anti-pattern-API/anti-pattern-API2# docker ps -a | grep anti-pattern
989a7090cf6a   anti-pattern-api2_app                           "./go-app"               32 seconds ago   Exited (1) 31 seconds ago                                                                                                                anti-pattern-api2_app_1
93206a0449e6   postgres:15-alpine                              "docker-entrypoint.s…"   32 seconds ago   Up 31 seconds                 0.0.0.0:5432->5432/tcp, :::5432->5432/tcp                                                                  anti-pattern-api2_db_1
root@pho3nix:/home/diegoall/MAESTRIA_ING/anti-pattern-API/anti-pattern-API2# docker logs anti-pattern-api2_db_1
The files belonging to this database system will be owned by user "postgres".
This user must also own the server process.

The database cluster will be initialized with locale "en_US.utf8".
The default database encoding has accordingly been set to "UTF8".
The default text search configuration will be set to "english".

Data page checksums are disabled.

fixing permissions on existing directory /var/lib/postgresql/data ... ok
creating subdirectories ... ok
selecting dynamic shared memory implementation ... posix
selecting default max_connections ... 100
selecting default shared_buffers ... 128MB
selecting default time zone ... UTC
creating configuration files ... ok
running bootstrap script ... ok
sh: locale: not found
2025-06-18 02:43:35.862 UTC [35] WARNING:  no usable system locales were found
performing post-bootstrap initialization ... ok
initdb: warning: enabling "trust" authentication for local connections
initdb: hint: You can change this by editing pg_hba.conf or using the option -A, or --auth-local and --auth-host, the next time you run initdb.
syncing data to disk ... ok


Success. You can now start the database server using:

    pg_ctl -D /var/lib/postgresql/data -l logfile start

waiting for server to start....2025-06-18 02:43:36.176 UTC [41] LOG:  starting PostgreSQL 15.13 on x86_64-pc-linux-musl, compiled by gcc (Alpine 14.2.0) 14.2.0, 64-bit
2025-06-18 02:43:36.177 UTC [41] LOG:  listening on Unix socket "/var/run/postgresql/.s.PGSQL.5432"
2025-06-18 02:43:36.181 UTC [44] LOG:  database system was shut down at 2025-06-18 02:43:36 UTC
2025-06-18 02:43:36.185 UTC [41] LOG:  database system is ready to accept connections
 done
server started
CREATE DATABASE


/usr/local/bin/docker-entrypoint.sh: running /docker-entrypoint-initdb.d/init.sql
CREATE TABLE
INSERT 0 5


waiting for server to shut down....2025-06-18 02:43:36.316 UTC [41] LOG:  received fast shutdown request
2025-06-18 02:43:36.317 UTC [41] LOG:  aborting any active transactions
2025-06-18 02:43:36.318 UTC [41] LOG:  background worker "logical replication launcher" (PID 47) exited with exit code 1
2025-06-18 02:43:36.319 UTC [42] LOG:  shutting down
2025-06-18 02:43:36.319 UTC [42] LOG:  checkpoint starting: shutdown immediate
2025-06-18 02:43:36.340 UTC [42] LOG:  checkpoint complete: wrote 931 buffers (5.7%); 0 WAL file(s) added, 0 removed, 0 recycled; write=0.011 s, sync=0.009 s, total=0.022 s; sync files=304, longest=0.001 s, average=0.001 s; distance=4256 kB, estimate=4256 kB
2025-06-18 02:43:36.344 UTC [41] LOG:  database system is shut down
 done
server stopped

PostgreSQL init process complete; ready for start up.

2025-06-18 02:43:36.435 UTC [1] LOG:  starting PostgreSQL 15.13 on x86_64-pc-linux-musl, compiled by gcc (Alpine 14.2.0) 14.2.0, 64-bit
2025-06-18 02:43:36.436 UTC [1] LOG:  listening on IPv4 address "0.0.0.0", port 5432
2025-06-18 02:43:36.436 UTC [1] LOG:  listening on IPv6 address "::", port 5432
2025-06-18 02:43:36.437 UTC [1] LOG:  listening on Unix socket "/var/run/postgresql/.s.PGSQL.5432"
2025-06-18 02:43:36.439 UTC [59] LOG:  database system was shut down at 2025-06-18 02:43:36 UTC
2025-06-18 02:43:36.443 UTC [1] LOG:  database system is ready to accept connections
root@pho3nix:/home/diegoall/MAESTRIA_ING/anti-pattern-API/anti-pattern-API2# docker ps -a | grep anti-pattern
989a7090cf6a   anti-pattern-api2_app                           "./go-app"               57 seconds ago   Exited (1) 56 seconds ago                                                                                                                anti-pattern-api2_app_1
93206a0449e6   postgres:15-alpine                              "docker-entrypoint.s…"   57 seconds ago   Up 57 seconds                 0.0.0.0:5432->5432/tcp, :::5432->5432/tcp                                                                  anti-pattern-api2_db_1
root@pho3nix:/home/diegoall/MAESTRIA_ING/anti-pattern-API/anti-pattern-API2# docker logs anti-pattern-api2_app_1
2025/06/18 02:43:35 No se pudo cargar el archivo .env, usando variables del entorno
2025/06/18 02:43:35 Ping a la base de datos falló: failed to connect to `user=user database=mydatabase`: 172.28.0.2:5432 (db): dial error: dial tcp 172.28.0.2:5432: connect: connection refused


Podrias ayudarme a identificar el error y a corregirlo para poder conectar la APi a la base de datos desde los contenedores.
De ser necesario generar archivos y correciones en su totalidad.
Respuesta en español
