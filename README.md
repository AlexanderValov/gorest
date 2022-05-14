# gorest 

### REST API by Golang

### Run:
- you need to open postresql connection and save URL to it into settings or you can use .env. 
- you need to add a table to the database. 
- run main go file.

### Endpoints:
- .../rag/v1/users/create (POST) - create user 
- .../rag/v1/users/update (PUT) - update user
- .../rag/v1/users/delete/{id} (DELETE) - delete user
- .../rag/v1/users/all (GET) - get all users
- .../rag/v1/users/{id} (GET) - get user
- .../rag/v1/check (GET) - check server
