# Microservices-Based E-Commerce System Development Assessment

## Objective
Create a fully functional microservices-based e-commerce system using Go, gRPC, protobuf, and PostgreSQL. This system should include multiple services and should be containerized using Docker and Docker Compose.

## Project Requirements
You are required to build an e-commerce system that includes the following mandatory microservices:
- **User Service**
- **Product Service**
- **Order Service**
- **Payment Service**

Additionally, you should include other services based on your own decision and design.

Each service should have its own gRPC server with the following CRUD operations

All services should communicate with a PostgreSQL database.

## Tasks

### Setup Project Structure
- Organize the project directory and structure by your own decision.

### Define Protobuf Files
- Create protobuf definitions for each service.

### Generate Go Code from Protobuf
- Use `protoc` to generate Go code from the protobuf definitions.

### Implement gRPC Servers
- Implement the gRPC servers for each service with the corresponding CRUD operations.

### Setup PostgreSQL Database
- Define the database schema for each service and ensure the services interact with PostgreSQL.

### Create Dockerfiles
- Create Dockerfiles for each service to containerize the applications.

### Docker Compose Configuration
- Create a `docker-compose.yml` file to define and manage multi-container Docker applications, including all services and PostgreSQL.

### Unit Testing
- Write unit tests for each service to ensure the correctness of CRUD operations.

## Detailed Instructions

### 1. Project Structure
Organize the project directory and structure by your own decision. Ensure that the structure is logical and supports easy management and scaling of the application.

### 2. Define Protobuf Files
Define protobuf files for each service. Ensure that the files define the necessary message types and gRPC service methods for CRUD operations.

### 3. Generate Go Code from Protobuf
Generate the Go code for gRPC using `protoc`. Ensure the generated code is placed in the appropriate directories for each service.

### 4. Implement gRPC Servers
Implement gRPC servers for each service. Ensure each server includes the CRUD operations and interacts with the PostgreSQL database.

### 5. Setup PostgreSQL Database
Define the PostgreSQL database schema for each service. Ensure that each service can connect to the database and perform CRUD operations.

### 6. Create Dockerfiles
Create Dockerfiles for each service to containerize the applications. Ensure that each Dockerfile sets up the service correctly, including copying the necessary files and building the Go application.

### 7. Docker Compose Configuration
Create a `docker-compose.yml` file to define and manage multi-container Docker applications. Ensure that the file includes definitions for all services and the PostgreSQL database, and that services are correctly configured to depend on the database.

### 8. Unit Testing
Write unit tests for each service to ensure the correctness of CRUD operations. Ensure that tests cover creating, reading, updating, and deleting records, and that they check for the correct behavior and responses.

## Submission
- Provide the complete source code.
- Include the Dockerfiles for each service.
- Provide the `docker-compose.yml` file.
- Include unit tests for each service.
- Provide instructions on how to build and run the application using Docker and Docker Compose.
  
