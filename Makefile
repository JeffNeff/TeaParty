up:
	@cd adams && make build
	@docker compose up --build

staging:
	@docker compose -f docker-compose.staging.yaml up -d

down:
	@docker compose down

nodes:
	@docker compose -f nodes.yaml up -d

infra: 
	@docker compose -f infra.yaml up -d

sn: 
	@docker compose -f nodes.yaml down

# run `sipper
sip:
	@cd sipper && go run .

debug: 
	@make up
	@docker logs teaparty-adams-1 --tail 50 -f

# Build fresh versions of the docker containers, start a local stack, and run `sipper` against it. 
test: 
	@make build
	@make up
	@make sip

attach:
	@docker logs teaparty-adams-1 --tail 50 -f

ngrok: 
	@ngrok http --subdomain=teaparty-adams 192.168.50.5:8080

tunnel:
	@ngrok tunnel --label edge=edghts_2NNjzClflrerVsg3M7L7yHWtl7d http://192.168.50.5:8080

watch:
	@kubectl logs -l  app=adams -f