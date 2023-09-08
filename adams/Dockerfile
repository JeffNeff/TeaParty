FROM ubuntu:jammy
RUN apt-get update && apt-get install -y ca-certificates
EXPOSE 8080
COPY /cmd/be /be

ENTRYPOINT ["/be"]