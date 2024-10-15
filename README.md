# Library API

API for a simple library, written in Go.

![Swagger UI](./imgs/swagger.png)

## Technologies Used

- **Go**: Programming language used for the API. (v1.23.2)
- **PostgreSQL**: Database used for storing data.
- **Swagger**: For API documentation.
- **GitHub Actions**: For CI/CD pipeline for generating swagger documentation.

## Requirements
- Go 1.23 or higher
- PostgreSQL database
- `migrate` tool for database migrations
## Local Installation

To run the project locally, follow these steps:

1. **Clone the repository**:
    ```sh
    git clone https://github.com/spin311/library-api.git
    cd library-api
    ```
2. **Create a `.env` file**:
    - Create a `.env` file in the root directory.
    - Add the following environment variables to the file:
    ```
    DBHOST=localhost
    DBPORT=5432
    DBUSER=your_db_username
    DBPASSWORD=your_db_password
    DBNAME=your_db_name
    SERVER_PORT=:8080
    ```
    - Replace values with your database credentials.

3. **Set up the database**:
    - Ensure PostgreSQL is installed and running.
    - Create a database for the project.
    - Install the `migrate` tool:
    ```sh
   go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
   ````
    - Run command with db values from .env file:
   ```sh 
    migrate -path ./migration -database postgres://dbusername:dbpassword@localhost:dbport/dbname?sslmode=disable up
   ```

4. **Install dependencies**:
    ```sh
    go mod tidy
    ```
5. **Run the application**:
    ```sh
    go run cmd/api/main.go
    ```
6. **Access the API**:
    - The API will be available at `http://localhost:8080`.
    - Swagger documentation can be accessed at `http://localhost:8080/swagger/index.html`.

## Endpoints

### User Endpoints

- **Create User**: `POST /users`
    - Request Body: `{ "first_name": "John", "last_name": "Doe" }`

- **Get All Users**: `GET /users`

- **Get User by ID**: `GET /users/{userId}`

### Book Endpoints

- **Get All Books**: `GET /books`

- **Get Book by ID**: `GET /books/{bookId}`

- **Borrow Book**: `POST /users/{userId}/books/{bookId}/borrow`

- **Return Book**: `PUT /users/{userId}/books/{bookId}/return`

## Project Structure

The project is organized into several directories to maintain a clean and modular structure:

- `cmd/api/`: Contains the entry point of the application (`main.go`).
- `config/`: Holds configuration settings (`config.go`).
- `docs/swagger/`: Contains Swagger documentation.
- `internal/app/`: Includes the core application logic, divided into `handlers` for HTTP handlers, `helpers` for utility functions, and `services` for business logic.
- `internal/repository/models/`: Defines the database models.
- `migration/`: Contains database migration files.
- `pkg/config/`: Provides configuration-related packages.
- `.env`: Stores environment variables.
- `.github/workflows/`: Contains GitHub Actions workflows (`openapi.yml`).
- `swagger.json`: Defines the Swagger API specification.

## CI/CD Pipeline
The CI/CD pipeline for this project is implemented using GitHub Actions. 
It is configured to trigger on pushes and pull requests to the main branch. 
The pipeline includes steps to check out the repository, set up the Go environment, install Swagger, generate Swagger documentation, and commit and push the updated documentation back to the repository.
The pipeline configuration is defined in the `openapi.yml` file located in the `.github/workflows/` directory.

## API Documentation

The API documentation is auto-generated using Swagger. To view it, navigate to:

- `http://localhost:8080/swagger/index.html`
