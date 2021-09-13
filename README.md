# Example Exchange Rate API

## Usage:
`docker-compose build && docker-compose up` -- Build / pull the requisite images and start 3 containers:
* DB: An instance of the Postgres database
* API: The API responsible for creating new projects and serving exchange rate information
* Ratefetch: Standalone service responsible for retrieving exchange rates and updating the database on an hourly basis
