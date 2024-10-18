.PHONY: *

docker_up:
	@echo "starting docker service"
	docker ps > /dev/null 2>&1 || (limactl start docker && sleep 5)
	@echo "docker service was started"
docker_down: 
	@echo "stoping docker service"
	docker ps && limactl stop docker
	@echo "docker service was stoped"
server: docker_up
	@echo "build and run go-micro services"
	docker-compose up -d --build

run: docker_up
	@echo "build and run go-micro services"
	docker-compose up --build

simple_run: docker_up build_auth build_logger build_frontend
	@echo "build and run go-micro services"
	docker-compose -f docker-compose-simple.yaml up --build

clean:
	@echo "stop docker-compose and remove docker containers"
	docker-compose down 
	docker system prune -a -f

build_auth:
	@echo build auth app...
	cd auth && make build

build_logger:
	@echo build logger app...
	cd logger && make build

build_frontend:
	@echo build front-end...
	cd front-end && make build
