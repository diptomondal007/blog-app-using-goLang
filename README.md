# A simple blog app developed with golang
 
---
### To run this app clone this repo and navigate to the repo's directory with terminal and run command:
##  Step 1:
    Create a psql database which is running on port 5432 with following username and password and dbname
    host     = "localhost"
    port     = 5432
    user     = "postgres"
    password = "12345"
    dbname   = "blogDB"
    or run command in docker
    ```
    docker run -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=12345 -e POSTGRES_DB=blogDB library/postgres
    ```
## Step 2:
   Create database tables with the query in file create_tables.txt
## Step 3:
    go run main.go
## Step 4:
   Open a browser and navigate to localhost:8000
