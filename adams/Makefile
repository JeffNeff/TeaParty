image: 
	@cd cmd && go mod tidy && GOOS=linux go build -o be main.go
	@gcloud builds submit --tag gcr.io/mineonlium/adms-prod . 

image-staging: 
	@gcloud builds submit --tag gcr.io/mineonlium/admsbro .

build: 
	@cd cmd && go mod tidy && GOOS=linux go build -o be main.go