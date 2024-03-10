startDocker:
	docker compose up

stopAll:
	docker stop postgres
	docker stop postgres-adminer-1

startPostgres:
	docker start postgres-products

openLocalPostgres:
	docker exec -it postgres-products psql -U postgres -d Price-Tracker

openDockerPostgres:
	docker exec -it postgres-products psql -U db -d Price-Tracker
