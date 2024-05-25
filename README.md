# Payment App

The Payment App is a simple application that implements wallet management and discount services. It is a testing app and is not meant for production use.

## How to Build and Run

### Prerequisites
- Go version 1.21 or higher
- Docker & docker compose for running the program inside the docker
- PostgreSQL for running the compiled binary


### Configuration
The configuration is specified in the configs/config.yml file. Ensure the database connection configurations are correctly set in this file.

#### Compiling the binary
`Before compiling & running the application, ensure that a PostgreSQL database is up and running.`
1. Run the migrate tool :
```shell
make migrate
```
2. Build the application :
```shell
make build
```
3. Running the application :
```shell
./payment
```
`You can modify the configuration in the configs/config.yml or pass config.yml from another directory using the -config flag.`


#### Running with docker compose
To run the application with Docker Compose, simply execute the following command:
```shell
docker-compose up -d
```
This will start the payment app, which will listen and serve on port 8080.

