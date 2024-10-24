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

Note: By default the server runs in API mode