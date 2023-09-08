# TeaParty

`Teaparty` an automated digital escrow system.

## Development Tools
 * Go 1.19
 * Node 16
 * Docker
 * Make
 * gcc
 * ko (go install github.com/google/ko@latest)

## Setting up a local Docker Desktop enviorment

From the root of this git:
```
docker compose up -d
```
This will build all of the projects and bring them up in your local enviorment providing the following local services to be accessed:

* `http://localhost:8080` - A dockerized `tea` application
* `http://localhost:8089` - `reddis commander` a debugging tool that allows one to easily view the state of the redis containers.



