Cabtrips data api [![Coverage Status](https://codecov.io/gh/springwiz/cabtrips-data-api/branch/master/graph/badge.svg)](https://codecov.io/gh/springwiz/cabtrips-data-api)

Package implements a Cabtrips api which exposes Rest Apis for exposing the publicly available cab trips dataset. It relies on Package gorilla/mux to manage routes and handle routing. The name mux stands for "HTTP request multiplexer". Like the standard http.ServeMux, mux.Router matches incoming requests against a list of registered routes and calls a handler for the route that matches the URL or other conditions.

The main features of the Cabtrips api are:
1. Exposes the Rest Apis for the publicly available cab trips dataset.
2. It is scalable and runs each HTTP Request in its own Go_Routine.
3. Is thread safe and employes locking to ensure the data integrity.
4. The Data Repository is currently on mysql but is pluggable and could handle multiple database types.
5. Employes allegro/bigcache to improve the data access time for the last week data.
6. It uses go mod to build the solution.
7. The .yaml files are used to reconfigure the solution. 

Installation
1. Install the Go lang runtime on the local machine.
2. Create the path on the local machine at github.com/springwiz.
3. Change into springwiz and Run git clone https://github.com/springwiz/cabtrips-data-api.git.
4. Change into cabtrips-data-api and run go build.
5. Change into cabtrips-data-api/client and run go build.
6. Run ./cabtrips-data-api and ./client from respective folders.
7. Run the test-data/test-commands. 

Assumptions
1. The package relies on MySql and allegro/bigcache.
2. Throws back the errors encountered to the front end as JSON Messages.
3. The data set given for this solution is a historical dataset, thus the last week data from
   24/12/2013 -> 31/12/2013 is loaded into the cache. Normally only the latest data from the
   last week will be loaded into the cache at definite intervals of time during the data.

Implementation details

1. The errors encountered on the backend are thrown from the REST Endpoints as JSON Messages with appropriate HTTP Error codes.
2. The package relies on a mix of OOB Go Unit Testing and Postman for testing the solution.
3. Bigcache was chosen among bigcache, freecache and map. This was on the past experience with bigcache and freecache.
4. spf13/viper is used to read the config yamls.

The GoLang stack was chosen from the following:
1. NodeJS
2. Java/SpringBoot
3. Golang/Gorilla Mux

The following are the considered factors: The NodeJS stack is more suitable for business processes which are I/O bound. The NodeJS Event Loop is not meant to run longer computations. The Java/SpringBoot based services are heavy weight and have much larger memory footprint. They have much higher boot time and are not preferable for web scale applications. The Java's model of 1 thread per request fails for web scale applications as thread context switches take most of the time/memory. The Golang services are light weight and offer rich support for concurrency via go routines and channels. The go routines are light weight and offer rich support for parellelism.

narays12-JSS1531:client narays12$ ./client get --medallion=9A80FE5419FEA4F44DB8E67F29D84A0F
INFO[0000] http://localhost:8080/cabtrip/9A80FE5419FEA4F44DB8E67F29D84A0F
INFO[0000] "{\"medallion\":\"9A80FE5419FEA4F44DB8E67F29D84A0F\",\"tripCount\":2}"

narays12-JSS1531:client narays12$ ./client get --medallion=9A80FE5419FEA4F44DB8E67F29D84A0F --cache=true
INFO[0000] http://localhost:8080/cabtrip/9A80FE5419FEA4F44DB8E67F29D84A0F?cache=true
INFO[0000] "{\"medallion\":\"9A80FE5419FEA4F44DB8E67F29D84A0F\",\"tripCount\":1}"

narays12-JSS1531:client narays12$ ./client get --medallion=9A80FE5419FEA4F44DB8E67F29D84A0F --date=20131231 --cache=true
INFO[0000] http://localhost:8080/cabtrip/9A80FE5419FEA4F44DB8E67F29D84A0F/date/20131231?cache=true
INFO[0000] "{\"medallion\":\"9A80FE5419FEA4F44DB8E67F29D84A0F\",\"tripCount\":1}"

narays12-JSS1531:client narays12$ ./client get --medallion=9A80FE5419FEA4F44DB8E67F29D84A0F --date=20131231
INFO[0000] http://localhost:8080/cabtrip/9A80FE5419FEA4F44DB8E67F29D84A0F/date/20131231
INFO[0000] "{\"medallion\":\"9A80FE5419FEA4F44DB8E67F29D84A0F\",\"tripCount\":1}"

narays12-JSS1531:client narays12$ ./client refresh
INFO[0000] http://localhost:8080/cache/refresh_cache
INFO[0135] HTTP/1.1 200 OK

Improvements
1. Implement more Apis to offer all the CRUD Operations and Data Query Operations.
2. Minimize explicit locking.
3. Offer connectors to various databases such as mongodb.
4. Employ scheduling to periodically refresh the cache so that cache always contains latest data.
5. Improve unit test code coverage.
6. Performance testing.
