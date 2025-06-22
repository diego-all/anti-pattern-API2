Actualmente el proyecto solo tiene un dockerfile para la API.
Y la base de datos solo tiene un archivo /sql/init.sql

donde crea y popula la base de datos.

CREATE TABLE IF NOT EXISTS instruments (
    id SERIAL PRIMARY KEY,
    name VARCHAR(64),
    description VARCHAR(200),
    price INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL
);

INSERT INTO instruments (name, description, price, created_at, updated_at)
VALUES 
('Guitarra eléctrica', 'Guitarra Fender Stratocaster de seis cuerdas', 1200, NOW(), NOW()),
('Batería acústica', 'Set completo de batería Pearl con platillos', 2300, NOW(), NOW()),
('Teclado digital', 'Yamaha con 88 teclas contrapesadas', 850, NOW(), NOW()),
('Violín', 'Violín acústico hecho a mano con arco y estuche', 600, NOW(), NOW()),
('Saxofón alto', 'Saxofón profesional con boquilla y correa', 1500, NOW(), NOW());

Aca esta el docker compose donde se monta el archivo init.sql para que postgre lo ejecute al iniciar:

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


EN otros proyectos que tengo no tengo el problema del connection refused:

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
[+] Building 6.2s (17/17) FINISHED                                                                                                               docker:default
 => [internal] load build definition from Dockerfile                                                                                                       0.0s
 => => transferring dockerfile: 1.15kB                                                                                                                     0.0s
 => [internal] load metadata for docker.io/library/alpine:latest                                                                                           1.0s
 => [internal] load metadata for docker.io/library/golang:1.23.2-alpine                                                                                    0.9s
 => [auth] library/alpine:pull token for registry-1.docker.io                                                                                              0.0s
 => [auth] library/golang:pull token for registry-1.docker.io                                                                                              0.0s
 => [internal] load .dockerignore                                                                                                                          0.0s
 => => transferring context: 2B                                                                                                                            0.0s
 => [internal] load build context                                                                                                                          0.0s
 => => transferring context: 4.54kB                                                                                                                        0.0s
 => [builder 1/6] FROM docker.io/library/golang:1.23.2-alpine@sha256:9dd2625a1ff2859b8d8b01d8f7822c0f528942fe56cfe7a1e7c38d3b8d72d679                      0.0s
 => [stage-1 1/3] FROM docker.io/library/alpine:latest@sha256:8a1f59ffb675680d47db6337b49d22281a139e9d709335b492be023728e11715                             0.0s
 => CACHED [builder 2/6] WORKDIR /app                                                                                                                      0.0s
 => CACHED [builder 3/6] COPY go.mod go.sum ./                                                                                                             0.0s
 => CACHED [builder 4/6] RUN go mod download                                                                                                               0.0s
 => [builder 5/6] COPY . .                                                                                                                                 0.1s
 => [builder 6/6] RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go-app main.go                                                           4.9s
 => CACHED [stage-1 2/3] WORKDIR /root/                                                                                                                    0.0s
 => CACHED [stage-1 3/3] COPY --from=builder /go-app .                                                                                                     0.0s
 => exporting to image                                                                                                                                     0.0s
 => => exporting layers                                                                                                                                    0.0s
 => => writing image sha256:507276655e4247523758f1c1ee7b52f8fc7a1415498ff2736e899b5a7396ed91                                                               0.0s
 => => naming to docker.io/library/anti-pattern-api2_app                                                                                                   0.0s
Creating anti-pattern-api2_db_1 ... done
Creating anti-pattern-api2_app_1 ... done
root@pho3nix:/home/diegoall/MAESTRIA_ING/anti-pattern-API/anti-pattern-API2# docker ps -a | grep anti-pattern
0aacbd87ed1c   anti-pattern-api2_app                           "./go-app"               16 seconds ago   Exited (1) 15 seconds ago                                                                                                                   anti-pattern-api2_app_1
99310c723401   postgres:15-alpine                              "docker-entrypoint.s…"   16 seconds ago   Up 16 seconds                    0.0.0.0:5432->5432/tcp, :::5432->5432/tcp                                                                  anti-pattern-api2_db_1
root@pho3nix:/home/diegoall/MAESTRIA_ING/anti-pattern-API/anti-pattern-API2# docker logs anti-pattern-api2_app_1
2025/06/18 03:17:48 No se pudo cargar el archivo .env, usando variables del entorno
2025/06/18 03:17:48 Ping a la base de datos falló: failed to connect to `host=db user=user database=mydatabase`: dial error (dial tcp 172.28.0.2:5432: connect: connection refused)

EN otros proyectos utilizo otra estrategia de generar un dockerfile para la base de datos de la siguiente forma:

- ACa esta el Dockerfile:
FROM --platform=linux/amd64 postgres:12.5-alpine

COPY up.sql /docker-entrypoint-initdb.d/1.sql

CMD ["postgres"]

- Aca esta up.sql

\c sqli;

DROP TABLE IF EXISTS users;
    
CREATE TABLE
  public.users (
    id serial NOT NULL,
    email character varying (255) NOT NULL,
    first_name character varying (255) NOT NULL,
    last_name character varying (255) NOT NULL,
    address character varying (255) NOT NULL,
    password character varying (60) NOT NULL,
    user_active integer NOT NULL DEFAULT 0,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NOT NULL DEFAULT now()
  );

  -- ALTER TABLE
  --   public.users
  -- ADD
  --   CONSTRAINT users_pkey PRIMARY KEY (id);

-- CREATE UNIQUE INDEX users_index_2 ON "public"."users" (id);


INSERT INTO public.users (email, first_name, last_name, address, password, user_active, created_at, updated_at) 
VALUES 
  ('usuario1@example.com', 'Juan', 'Pérez', 'Calle 123', '$2a$12$IjOgt/06hlNF13IOsrb8veJemUeSDB.7X27UtSubDbjBgXuL.j5ci', 1, now(), now()),
  ('usuario2@example.com', 'María', 'Gómez', 'Avenida 456', '$2a$12$xzfjUjBa06RwrNRu.wb.M.8bWJMc2cI9GZObV9495ypXRbfjNUyPS', 1, now(), now()),
  ('usuario3@example.com', 'Luis', 'Martínez', 'Calle 789', '$2a$12$1HgyDgcSZuZQDkKbEN6elug3P5Z62Rjrrf/YQdDEBiJ3sSuxcqpWW', 1, now(), now()),
  ('usuario4@example.com', 'Ana', 'Rodríguez', 'Avenida 101112', '$2a$12$IjOgt/06hlNF13IOsrb8veJemUeSDB.7X27UtSubDbjBgXuL.j5ci', 1, now(), now()),
  ('usuario5@example.com', 'Pedro', 'López', 'Calle 131415', '$2a$12$xzfjUjBa06RwrNRu.wb.M.8bWJMc2cI9GZObV9495ypXRbfjNUyPS', 1, now(), now());


- Aca esta el docker-compose.yaml
version: '3.9'

services:
  # Start Postgres, and ensure that data is stored to a mounted volume

  # serviceName
  postgres_sqli:
    # Pending: Connect with Dockerfile
    # build:
    #  context: ./database
    #  dockerfile: ./database/Dockerfile

    # Assign the image or build the image on development process
    #image: 'diegoall1990/sqli-pg-db:0.0.1'
    image: 'diegoall1990/linux-sqli-pg-db:0.0.1'
    container_name: linux_postgres_sqli_dummy
    ports:
      - "5432:5432"
    restart: always

    environment:
      - POSTGRES_USER=${POSTGRES_USER}
      - POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
      - POSTGRES_DB=${POSTGRES_DB}

# Bad security practice
    # environment:
    #   NOMBRE=${NOMBRE}
    #   POSTGRES_USER: postgres
    #   POSTGRES_PASSWORD: password
    #   POSTGRES_DB: sqli
    volumes:
      #- ./db-data/postgres/:/var/lib/postgresql/data/
      - ./db-data/postgres/:/var/lib/postgresql/data/:rw

# Deployment artifacts
    deploy:
     mode: replicated
     replicas: 1

    # --- AGREGAR HEALTHCHECK PARA POSTGRES ---
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${POSTGRES_USER} -d ${POSTGRES_DB}"]
      interval: 5s
      timeout: 5s
      retries: 10 # Intentar 10 veces (50 segundos en total)
      start_period: 30s # Darle a Postgres 30 segundos para iniciar antes de empezar a chequear
    # --- FIN HEALTHCHECK ---

  api_sqli:
    build:
      context: .
      dockerfile: ./cmd/api/Dockerfile
    container_name: golang_api_sqli_dummy
    ports:
      - "9090:9090"
    # --- MODIFICAR depends_on para esperar el healthcheck ---
    depends_on:
      postgres_sqli:
        condition: service_healthy # Esperar a que postgres_sqli esté marcado como 'healthy'
    # --- FIN MODIFICACIÓN ---
    environment:
      - DSN=${DSN}
      # - DSN=host=postgres_sqli port=5432 user=postgres password=password dbname=sqli sslmode=disable timezone=UTC connect_timeout=5


Se tiene la estrategia de generar la imagen tanto de la base de datos como para la API a partir de dockerfiles, podrias ajustar el proyecto para utilizar esta estrategia. Quizas ayuda a solucionar el problema del timeout.

De ser necesario generar los archivos completops y dar la respuesta en español




