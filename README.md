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

- CREATE TABLE users (
  id SERIAL PRIMARY KEY, -- Unique ID (auto-incremented)
  name VARCHAR(255) NOT NULL, -- User name, not null
  email VARCHAR(255) NOT NULL, -- Email, not null
  sub_id VARCHAR(50), -- Subscription ID or external identifier
  verification_status BOOLEAN, -- Boolean for verification status
  setup_status VARCHAR(50) -- Setup status (like 'completed', 'pending', etc.)
  );

- INSERT INTO users (id, name, email, sub_id, verification_status, setup_status)
  VALUES (
  2,
  'John Smith',
  'john.smith@example.com',
  'SUB123456abc',
  TRUE,
  'completed'
  );
