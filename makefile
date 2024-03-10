startDocker:
	docker compose up

stopAll:
	docker stop postgres
	docker stop postgres-adminer-1
