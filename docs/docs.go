package docs

import "github.com/swaggo/swag"

const docTemplate = `{
  "swagger": "2.0",
  "info": {
    "description": "Movie booking backend intern assignment implemented in Go.",
    "title": "Movie Ticket Booking API (Go)",
    "version": "1.0"
  },
  "basePath": "/api",
  "schemes": [
    "http"
  ],
  "paths": {
    "/signup": {
      "post": {
        "summary": "User signup",
        "description": "Register a new user",
        "consumes": ["application/json"],
        "produces": ["application/json"],
        "parameters": [
          {
            "in": "body",
            "name": "body",
            "description": "Signup payload",
            "required": true,
            "schema": { "$ref": "#/definitions/SignupRequest" }
          }
        ],
        "responses": {
          "201": {
            "description": "User created"
          },
          "400": {
            "description": "Bad request"
          }
        }
      }
    },
    "/login": {
      "post": {
        "summary": "User login",
        "description": "Login and obtain a JWT token",
        "consumes": ["application/json"],
        "produces": ["application/json"],
        "parameters": [
          {
            "in": "body",
            "name": "body",
            "description": "Login payload",
            "required": true,
            "schema": { "$ref": "#/definitions/LoginRequest" }
          }
        ],
        "responses": {
          "200": {
            "description": "OK (token returned)",
            "schema": { "$ref": "#/definitions/LoginResponse" }
          },
          "401": {
            "description": "Invalid credentials"
          }
        }
      }
    },
    "/movies": {
      "get": {
        "summary": "List movies",
        "description": "Get all movies",
        "produces": ["application/json"],
        "responses": {
          "200": {
            "description": "List of movies"
          }
        }
      }
    },
    "/movies/{id}/shows": {
      "get": {
        "summary": "List shows for a movie",
        "description": "Get all shows for a given movie",
        "produces": ["application/json"],
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "description": "Movie ID",
            "required": true,
            "type": "integer",
            "format": "int64"
          }
        ],
        "responses": {
          "200": {
            "description": "List of shows"
          }
        }
      }
    },
    "/shows/{id}/book": {
      "post": {
        "summary": "Book a seat",
        "description": "Book a seat for a specific show (JWT required)",
        "consumes": ["application/json"],
        "produces": ["application/json"],
        "parameters": [
          {
            "name": "Authorization",
            "in": "header",
            "description": "Bearer JWT token",
            "required": true,
            "type": "string"
          },
          {
            "name": "id",
            "in": "path",
            "description": "Show ID",
            "required": true,
            "type": "integer",
            "format": "int64"
          },
          {
            "in": "body",
            "name": "body",
            "description": "Seat number payload",
            "required": true,
            "schema": { "$ref": "#/definitions/BookSeatRequest" }
          }
        ],
        "responses": {
          "201": {
            "description": "Seat booked"
          },
          "400": {
            "description": "Validation / business rule error"
          },
          "401": {
            "description": "Unauthorized"
          }
        }
      }
    },
    "/bookings/{id}/cancel": {
      "post": {
        "summary": "Cancel booking",
        "description": "Cancel an existing booking of the logged-in user",
        "produces": ["application/json"],
        "parameters": [
          {
            "name": "Authorization",
            "in": "header",
            "description": "Bearer JWT token",
            "required": true,
            "type": "string"
          },
          {
            "name": "id",
            "in": "path",
            "description": "Booking ID",
            "required": true,
            "type": "integer",
            "format": "int64"
          }
        ],
        "responses": {
          "200": {
            "description": "Booking cancelled"
          },
          "400": {
            "description": "Invalid status"
          },
          "401": {
            "description": "Unauthorized"
          },
          "403": {
            "description": "Forbidden (not your booking)"
          }
        }
      }
    },
    "/my-bookings": {
      "get": {
        "summary": "Get my bookings",
        "description": "List all bookings of the logged-in user",
        "produces": ["application/json"],
        "parameters": [
          {
            "name": "Authorization",
            "in": "header",
            "description": "Bearer JWT token",
            "required": true,
            "type": "string"
          }
        ],
        "responses": {
          "200": {
            "description": "List of bookings"
          },
          "401": {
            "description": "Unauthorized"
          }
        }
      }
    }
  },
  "definitions": {
    "SignupRequest": {
      "type": "object",
      "properties": {
        "name": { "type": "string" },
        "email": { "type": "string" },
        "password": { "type": "string" }
      }
    },
    "LoginRequest": {
      "type": "object",
      "properties": {
        "email": { "type": "string" },
        "password": { "type": "string" }
      }
    },
    "LoginResponse": {
      "type": "object",
      "properties": {
        "token": { "type": "string" }
      }
    },
    "BookSeatRequest": {
      "type": "object",
      "properties": {
        "seat_number": { "type": "integer", "format": "int64" }
      }
    }
  }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "Movie Ticket Booking API (Go)",
	Description:      "Movie booking backend intern assignment implemented in Go.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
