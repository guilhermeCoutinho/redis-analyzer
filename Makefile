deps:
	@docker-compose \
	--file docker-compose.yaml \
	--project-name=$(redis-mem-analyzer) \
	up -d --force-recreate --no-deps 
	@MODE=POPULATEREDIS go run *.go

run:
	@go run *.go