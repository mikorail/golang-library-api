Thank you for the clarification! Below is the corrected version of the API documentation where the routes for borrowing and returning books are specified using the `GET` method.

---

# Products API with JWT

## Overview

The **Products API with JWT** is a RESTful API built using Go (Golang) and the Gin framework. This API enables users to manage products while ensuring secure access through JSON Web Tokens (JWT). It utilizes PostgreSQL as the primary database and includes functionality for rate limiting and token invalidation upon logout.

## Features

- **User Authentication**: Secure login and logout processes using JWT.
- **CRUD Operations**: Create, Read, Update, and Delete functionalities for managing products.
- **Borrowing and Returning**: Users can borrow and return books.
- **Secure Endpoints**: All sensitive endpoints require JWT authentication.
- **Rate Limiting**: Protect endpoints from excessive requests.
- **Token Invalidation**: Ensure that JWTs cannot be reused after logout.
- **Swagger Documentation**: Easy exploration and testing of the API via autogenerated documentation.

## Technologies Used

- **Go**: The programming language used for building the API.
- **Gin**: A web framework for Go, facilitating HTTP request handling.
- **GORM**: An ORM for Go, used for PostgreSQL database interactions.
- **Swagger**: A tool for API documentation and testing.
- **Rate Limiter**: Integrated to prevent endpoint abuse.

## Getting Started

### Prerequisites

Ensure you have the following installed:

- [Go](https://golang.org/dl/) (version 1.16 or later)
- PostgreSQL database (ensure it's running)

### Installation

1. **Clone the repository**:
   ```bash
   git clone ...
   cd products-api-with-jwt
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Set up the environment**: Create a `.env` file in the root directory and set the required environment variables, including your database connection details and JWT secret key.

   Example `.env`:
   ```
   ENVSecretKey=my_secret_key
   SECRET_KEY = "M4K4N-514ng-GR4T15"
   APP_ENV="development"
   PORT="8080" 
   DB_URL="localhost"
   DB_NAME="user-message-api"
   DB_PORT="5432"
   DB_PASSWORD="th3password"
   DB_USER="user-message-api"
   DB_SSL="disable"
   ```

4. **Run the application**:
   ```bash
   go run main.go
   ```

   The server will start on `http://localhost:8080`.

### API Endpoints

#### Authentication

- **Login**
  - **Endpoint**: `/auth/login`
  - **Method**: `POST`
  - **Request Body**:
    ```json
    {
      "username": "admin",
      "password": "password123",
      "rememberMe": true
    }
    ```
  - **Response**:
    ```json
    {
      "status": "success",
      "code": 200,
      "message": "Login successful",
      "data": {
        "token": "your_jwt_token_here"
      }
    }
    ```

- **Logout**
  - **Endpoint**: `/auth/logout`
  - **Method**: `POST`
  - **Headers**: `Authorization: Bearer <jwt_token>`
  - **Response**:
    ```json
    {
      "status": "success",
      "code": 200,
      "message": "Logout successful",
      "data": null
    }
    ```

  Upon logout, the JWT is invalidated, ensuring it cannot be reused to access secure endpoints.

#### Products

All product-related endpoints require a valid JWT token in the `Authorization` header.

- **Get All Products**
  - **Endpoint**: `/books`
  - **Method**: `GET`
  - **Response**:
    ```json
    {
      "status": "success",
      "code": 200,
      "message": "Products retrieved successfully",
      "data": [ ... ]
    }
    ```

- **Get Product by ID**
  - **Endpoint**: `/books/:id`
  - **Method**: `GET`
  - **Response**:
    ```json
    {
      "status": "success",
      "code": 200,
      "message": "Product retrieved successfully",
      "data": { ... }
    }
    ```

- **Create Product**
  - **Endpoint**: `/books`
  - **Method**: `POST`
  - **Request Body**:
    ```json
    {
      "title": "Produk A",
      "description": "Deskripsi Produk A",
      "stock": 1000,
      "stok": 10
    }
    ```
  - **Response**:
    ```json
    {
      "status": "success",
      "code": 201,
      "message": "Product created successfully",
      "data": { ... }
    }
    ```

- **Update Product**
  - **Endpoint**: `/books/:id`
  - **Method**: `PUT`
  - **Request Body**: Same as Create Product
  - **Response**: Similar to Create Product response

- **Delete Product**
  - **Endpoint**: `/books/:id`
  - **Method**: `DELETE`
  - **Response**:
    ```json
    {
      "status": "success",
      "code": 204,
      "message": "Product deleted successfully",
      "data": null
    }
    ```

#### Borrowing and Returning Books

All borrowing and returning endpoints require a valid JWT token in the `Authorization` header.

- **Borrow Book**
  - **Endpoint**: `/books/borrow/:id`
  - **Method**: `GET`
  - **Response**:
    ```json
    {
      "status": "success",
      "code": 200,
      "message": "Book borrowed successfully",
      "data": { ... }
    }
    ```

- **Return Book**
  - **Endpoint**: `/books/return/:id`
  - **Method**: `GET`
  - **Response**:
    ```json
    {
      "status": "success",
      "code": 200,
      "message": "Book returned successfully",
      "data": { ... }
    }
    ```

### Rate Limiting

Rate limiting is enabled on certain endpoints to prevent abuse by limiting the number of requests allowed within a specified timeframe. If the rate limit is exceeded, the following response is returned:

```json
{
  "status": "error",
  "code": 429,
  "message": "Too many requests, please try again later.",
  "data": null
}
```

### API Documentation

You can access the Swagger documentation by navigating to `http://localhost:8080/swagger/index.html` in your web browser.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

Feel free to modify the responses based on your actual implementation. Let me know if you need any further adjustments!