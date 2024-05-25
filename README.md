# Payment App

The Payment App is a simple application that implements wallet management and discount code generation services. It is a testing app and is not meant for production use.

## How to Build and Run

### Prerequisites
- Go version 1.21 or higher
- Docker & Docker Compose for running the program inside the docker
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
`You can modify the configuration in the configs/config.yaml or pass config.yaml from another directory using the -config flag.`


#### Running with docker compose
To run the application with Docker Compose, simply execute the following command:
```shell
docker-compose up -d
```
This will start the payment app, which will listen and serve on port 8080.


### Available Routes

#### Wallet Service Routes
- POST /wallet/register: Register a new wallet.
- PUT /wallet/{phoneNumber}: Perform a transaction.
- DELETE /wallet/{phoneNumber}: Delete a wallet.
- GET /wallet/{phoneNumber}: Get wallet details by phone number.
#### Discount Service Routes
- POST /discount: Create a new discount.
- GET /discount/usages: Get discount usages.
- GET /discount/apply: Apply a discount.


#### Examples

Register a new wallet
```shell
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"phone": "PhoneNumber", "amount": 100000}' \
  http://localhost:8080/wallet/register
```
Perform a transaction
```shell
curl -X PUT \
  -H "Content-Type: application/json" \
  -d '{"amount": 1000, "description": "Withdrawal for groceries", "type": "withdrawal"}' \
  http://localhost:8080/wallet/PhoneNumber
````
Delete a wallet
```shell
curl -X DELETE http://localhost:8080/wallet/PhoneNumber
```
Get wallet details by phone number
```shell
curl http://localhost:8080/wallet/PhoneNumber
```

Create a new discount code
```shell
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"usage_limit": 1000, "description": "Voucher for cup league", "amount": 1000000, "type": "voucher"}' \
  http://localhost:8080/discount
```

Get discount code transactions
````shell
curl http://localhost:8080/discount/usages?code=WS9DE6CH
````

Apply a discount
```shell
curl http://localhost:8080/discount/apply?code=WS9DE6CH&phone=PhoneNumber
```

#### Swagger Documentation
For detailed information on the available routes and request/response schemas, refer to the Swagger documentation provided in the [docs/swagger.yaml](docs/swagger.yaml) file.


