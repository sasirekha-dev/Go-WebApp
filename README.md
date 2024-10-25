##### build Todo web app #########
There are four endpoints to
- List
- Add a task
- delete a task from Todo
- Update a task

To run in two different data stores 
1. Memory store
2. Json file

To change mode to CLI
=========================
.\server -mode=cli insert -item="learn Go" -status="on going"

.\server -mode=cli update -item="learn Go" 

.\server -mode=cli delete -item="learn Go"  

.\server -mode=cli list 

To run server via API calls
============================
.\server -mode=api

Server runs in 8080 port
http://localhost:8080/ - single landing page to demonstrate CRUD operations on the database.


Note: By default the server runs in API mode


Testing
========
Include unit test and parallel test
To run unit test execute -  go test OR 
                            go test -v <package name> -run <test suite name>
                            go test -v  WebApp\store -run TestInMemoryStore_ParallelOperations
To run the parallel test - go test -race OR 
To run specific test case - go test -v <package name> -race <test suite name>
For example: go test WebApp\store -race TestInMemoryStore_ParallelOperations

Pending Additions
===================
* Add a V2 API to the REST API that supports multiple users
* Use Parallel test to validate the solution is concurrent safe across multiple users.
* Add a Web API that supports multiple users
* Add multi-user support to the CLI App
* Add an additional data store implementation using an external DB (PostgreSQL)

