build:
	docker-compose build

# Remove cahche before build
buildcache:
	docker-compose build --no-cache

#detached mode
up:
	docker-compose up -d

down:
	docker-compose down

# Remvoe caching
cache:
	docker builder prune

# Remove Docker image (paste image name after rm below)
remove img:
	docker image rm 

# Start interactive shell inside container
docker-shell:
	docker-compose run --rm app sh




# Alternative comamnds to start docker - change port dynamically
dbldapp:
	docker build -t webapi-go-app .  

dblrun:
	docker run -p 8080:8080 -v $PWD:/app webapi-go-app

migrate:
#	 migrate -database postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME} -path repository/database/migrations up
	migrate -database postgres://webservice_dev_user:yourpassword@db:5432/webservice_dev?sslmode=disable -path repository/database/migrations up