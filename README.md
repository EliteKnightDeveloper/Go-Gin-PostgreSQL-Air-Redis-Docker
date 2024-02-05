# Go-Gin-Air-PostgreSQL-Redis-Docker-Starter

This project is written purely in Go Language. Gin (Http Web Frame Work) is used in this project. PostgreSQL Database is used to manage the data.

## Framework Used

Gin-Gonic: This whole project is built on Gin frame work. Its is a popular http web frame work.

```
go get -u github.com/gin-gonic/gin
```

## Database used:

PostgreSQL: PostgreSQL is a powerful, open source object-relational database. The data managment in this project is done using PostgreSQL. ORM tool named GORM is also been used to simplify the forms of queries for better understanding.

```
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
```

## External Packages Used

#### Validator

Package validator implements value validations for structs and individual fields based on tags.

```
github.com/go-playground/validator/v10
```

#### JWT

JSON Web Tokens are an open, industry standard RFC 7519 method for representing claims securely between two parties.

```
github.com/golang-jwt/jwt/v4
```

#### Commands to run project:

```
air
```

## 👉 Documentation

```
http://localhost:8000/api/v1/docs/index.html
```

## 👉 Signup as user

### Endpoint :

```
http://localhost:8000/api/v1/auth/register
```

### Method:

`POST`

### Request Body:

| Parameter  | Type   | Description          |
| ---------- | ------ | -------------------- |
| `email`    | string | Email ID of the user |
| `password` | string | Password of the user |

### Example Request:

```
 POST   http://localhost:8000/api/v1/auth/register

-H "Content-Type: application/json"
-d '{
    "email" : "tony@yopmail.com",
    "password" : "12345",
}'
```

### Success Response:

HTTP Code: `200 OK`

```
{
  "message": "Go to /register/validate"
}
```

## 👉 To verify the otp

### Endpoint :

```
http://localhost:8000/user/register/validate
```

### Method:

`POST`

### Request Body:

### Request Body:

| Parameter | Type    | Description       |
| --------- | ------- | ----------------- |
| `otp`     | Intiger | Otp of the user   |
| `email`   | string  | Email of the user |

### Example Request:

```
 POST  http://localhost:8000/user/register/validate
-H "Content-Type: application/json" \
-d '{
       "otp" : "1904",
       "email" : "tony@yopmail.com"
    }'
```

### Success Response:

HTTP Code: `200 OK`

```
{
  "Message": "New User Successfully Registered"
}
```

## 👉 To login as a user

### Endpoint :

```
  http://localhost:8000/api/v1/auth/login
```

### Method:

`POST`

### Request Body:

| Parameter  | Type   | Description          |
| ---------- | ------ | -------------------- |
| `email`    | string | Email ID of the user |
| `password` | string | Password of the user |

### Success Response:

HTTP Code: `200 OK`

```
{
  "message": "User login successfully"
}
```
