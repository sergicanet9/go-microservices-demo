# go-microservices-demo
A lightweight **Go microservices** demo showcasing gRPC and HTTP communication, built on top of the [scv-go-tools](https://github.com/sergicanet9/scv-go-tools) library and using [go-hexagonal-api](https://github.com/sergicanet9/go-hexagonal-api) as a backend service.

<!-- TODO draft -->
<!-- ## 🚀 Services
### 1. go-hexagonal-api
- gRPC + HTTP API for user management.
- Exposes endpoints via gRPC and automatically via HTTP (gRPC-Gateway).
- Runs in its own Docker container.

### 2. Task Tracker
- HTTP API for managing tasks per user.
- Communicates internally with `go-hexagonal-api` via gRPC for user authentication.
- Endpoints:
  - `GET /tasks` – List tasks for the authenticated user.
  - `POST /tasks` – Create a new task for the authenticated user.
  - `DELETE /tasks/{id}` – Delete a task.

## 🏁 Getting Started
### Run with Docker Compose
```
make up
```
This launches two containers:
go-hexagonal-api
task-tracker

## Access the Services
- **Swagger UI for go-hexagonal-api**: [http://localhost:8082/swagger/index.html](http://localhost:8082/swagger/index.html)  
- **gRPC UI for go-hexagonal-api**: [http://localhost:8082/grpcui/](http://localhost:8082/grpcui/)  
- **Task Tracker HTTP API**: [http://localhost:8085](http://localhost:8085)

## ⚙️ Authentication
Both services require a valid JWT token for user-specific operations.  
The Task Tracker service validates the token by calling `go-hexagonal-api` internally via gRPC.

## 🛠️ Development
- **Task Tracker** service is located in the `task-tracker/` folder.  
- **go-hexagonal-api** is pulled as a Docker image from your local build or a registry.

### Build & Run Task Tracker
```
cd task-tracker
go build -o task-tracker ./cmd
./task-tracker --port=8085 --hexagonal-api-grpc=go-hexagonal-api:50051
```

## Access the Services
- **Swagger UI for go-hexagonal-api**: [http://localhost:8082/swagger/index.html](http://localhost:8082/swagger/index.html)  
- **gRPC UI for go-hexagonal-api**: [http://localhost:8082/grpcui/](http://localhost:8082/grpcui/)  
- **Task Tracker HTTP API**: [http://localhost:8085](http://localhost:8085)


## ⚙️ Authentication
Both services require a valid JWT token for user-specific operations.  
The Task Tracker service validates the token by calling `go-hexagonal-api` internally via gRPC.

## 🛠️ Development
- **Task Tracker** service is located in the `task-tracker/` folder.  
- **go-hexagonal-api** is pulled as a Docker image from your local build or a registry.

## Build & Run Task Tracker
```bash
cd task-tracker
go build -o task-tracker ./cmd
./task-tracker --port=8085 --hexagonal-api-grpc=go-hexagonal-api:50051
```

## 📚 References
* [scv-go-tools](https://github.com/sergicanet9/scv-go-tools)
* [go-hexagonal-api](https://github.com/sergicanet9/go-hexagonal-api)

## ✍️ Author
Sergi Canet Vela

## ⚖️ License
This project is licensed under the terms of the MIT license. -->