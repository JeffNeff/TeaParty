
# TeaParty

`TeaParty` is a unique digital asset exchange that operates without the need for liquidity pools.

## Prerequisites

Before you begin, ensure you have the following development tools installed:

- **Go**: Version 1.19. You can download it from [here](https://golang.org/dl/).
- **Node**: Version 16. Download from [Node.js official site](https://nodejs.org/).
- **Docker**: Used for containerization. Install from [Docker's official site](https://www.docker.com/get-started).
- **Make**: A build automation tool.
- **gcc**: The GNU Compiler Collection.
- **ko**: A tool for building and deploying Golang applications. Install using:
  ```bash
  go install github.com/google/ko@latest
  ```

## Build

### Building the Application

To build the `TeaParty` application:

1. Navigate to the `adams` directory:
   ```bash
   cd adams
   ```
2. Use the `make` command to build:
   ```bash
   make build
   ```

## Run

### Running the Application

1. Ensure the `docker-compose.yaml` file is populated with the necessary environment variables.
2. Use Docker Compose to run the application:
   ```bash
   docker-compose up
   ```

## Deploy

### Building the Container Image

[Provide instructions on how to build the container image for deployment.]

### Deploying the Application

Deployment manifests are provided to help you deploy the `TeaParty` application:

- **`/config/1-infra`**: Contains deployment manifests for various crypto RPC nodes.
- **`/config/2-party`**: Contains the necessary manifests for deploying the `TeaParty` backend.
- **`/config/3-tools`**: Features an example of a containerized `tea` frontend application. Note: This might not work out-of-the-box and may require additional configuration.


