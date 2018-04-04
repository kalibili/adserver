# adserver

Adserver is a basic webserver built using Beego web framework of Golang.

The client can select in which all regions he/she wants to or does not want to show his/her ad's.

The whole data is loaded in memory before the server is loaded.


Prerequisites:
PostgreSQL(v.9.1)
go(v.1.8)



Steps to run the server:

Step 1: go get github.com/gorilla/mux
        go get github.com/lib/pq

Step 2: Import the schema(adserer_go.sql) into the PostgreSQL Database
        adserer_go.sql contains the SQL Schema for the Project

Step 3: There are 2 go servers:
        1. adserer_geo.go
           This server will use IP geolocation to track the location and show the respective ad to the user.
        2. adserer_api.go
           This server will take location as query parameter and show the respective ad to the user.
        To run any one of the server use the command 
        go run adserer_geo.go
                or
        go run adserer_api.go

