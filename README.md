# go-api-demo

This is a simple project that shows how I'd write an HTTP API in Go with some application layering like API layer, 
Service layer and Persistence layer.  

The API allows a User to retrieve his balance, list his transactions and create a transaction to send money to some
other User. 

## How to start the project 

Inside the api-demo folder, execute `docker-compose up` to have everything magically started!

**Testing**

Inside the api-demo folder, execute `go test ./...`

**Example**

The API is authenticated, every request should include the credentials of your user in the query params, for instance:
  - `curl localhost:8080/me?username=breno&password=1234`
  - `curl localhost:8080/me/transactions?username=breno&password=1234`

## API

The API server runs by default on port 8080, the healthcheck server, on port 8585.

## API Server

**Authentication**

As previously mentioned, the API has 2 mandatory query parameters for every request, **username** and **password**, not
providing these will lead to requests being rejected.
 
#### /me
  - **GET**: returns the balance of the current user.

#### /me/transactions
  - **GET**: returns the transactions that the current user made.
  
  - **POST**: creates a transaction, requiring the following payload: 
  ```json
    {
      "target_user_id": "STRING|UUID",
      "amount": 10.5
    }
  ```

## Healthcheck Server
 
### /healthcheck
**GET**: simply returns an OK header if the service is alive

## Test data

There are pre-created users that can be used to authenticate, as signing in and signing up was not implement.
Feel free to create your own users on the schema.sql file, just remember to run `docker-compose down` because the PG
database might still stay up with data. 

**username:password**
 - breno:1234
 - bruno:4321
 - brono:abcd
 - brano:abcdef
