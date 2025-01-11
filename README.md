# webservice-v2

## DB access and setup

1. sudo -i -u postgres
2. psql
3. CREATE DATABASE webservice_dev;
4. CREATE USER webservice_dev_user WITH PASSWORD 'yourpassword';
5. GRANT ALL PRIVILEGES ON DATABASE webservice_dev TO webservice_dev_user;

SELECT \* FROM pg_catalog.pg_tables;

## docker stuff

https://www.reddit.com/r/golang/comments/18qwci9/how_to_create_a_containerised_application_using/
https://thebugshots.dev/packaging-a-golang-application-using-multi-stage-docker-builds
https://www.mitrais.com/news-updates/how-to-dockerize-a-restful-api-with-golang-and-postgres/
https://ramadhansalmanalfarisi8.medium.com/how-to-dockerize-your-api-with-go-postgresql-gin-docker-9a2b16548520
https://dev.to/muhammedarifp/creating-a-simple-hello-world-web-application-with-docker-and-golang-1e14

# access db from the db container

- docker ps
- docker exec -it <container_id> bash
- psql -U <postgres_user> -d <database_name>
- SELECT \* FROM users;
- list schemas \dnS

# Create migration files

migrate create -ext sql -dir repository/database/migrations -seq name_of_the_migration

- CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL,
  sub_id VARCHAR(50),
  verification_status BOOLEAN,
  setup_status VARCHAR(50),
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now()
  );

- INSERT INTO users (
  name,
  email,
  sub_id,
  verification_status,
  setup_status,
  created_at,
  updated_at
  ) VALUES (
  'John Doe',  
   'john.doe@example.com',  
   'SUB987654321',  
   TRUE,  
   'in_progress',  
   NOW(),  
   NOW()  
  );

## Things to look into

- context. Do I need to create several different contexts, e.g. shutdownCtx (in main) and pool, err := pgxpool.NewWithConfig(context.Background(), config) and ctx, cancel := context.WithTimeout(context.Background(), 5\*time.Second) in dbConfig?

- Do I really need var wg sync.WaitGroup?

- Where and why should channel response operator; <-quit

- should db be in a goroutine? Why and why not?
