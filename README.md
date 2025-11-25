<div align="center">
<img alt="Dino" width="120" src="https://raw.githubusercontent.com/willpinha/dino/refs/heads/master/logo.svg" />

# Dino

🦖 Lightweight HTTP library compatible with `net/http`

[About](#about) • [Philosophy](#philosophy) • [Installation](#installation) • [Documentation](#documentation) • [License](#license)

</div>

## 🦖 About

**Dino** is a lightweight HTTP library that contains utilities compatible with the `net/http` package

If you want to work directly with `net/http` and the standard library, without relying on abstractions (third-party routers, frameworks, etc.), but don't want to implement common functionalities from scratch, Dino is the right library for you

## 🦖 Philosophy

Dino doesn't try to be the magic solution that solves all your problems. Instead, it provides a thin layer built on top of `net/http` for functionalities commonly needed when building applications

It also follows the Go philosophy, which is to maintain simplicity and not apply breaking changes or major versions all the time. This makes Dino an easy-to-use and stable library in the long term

## 🦖 Installation

Want to try Dino? Go get the package below and start reading the next sections of the documentation

```
go get github.com/willpinha/dino
```

## 🦖 Documentation

1. [Handlers]()
2. [Routing]()
3. [Error handling]()
4. [Middlewares]()
5. [Requests]()
6. [Responses]()
7. [Logging]()

### 1. Handlers

Handlers (a.k.a. Controllers) represents the core mechanism for processing HTTP requests and returning HTTP responses to the clients

Similar to the [`http.HandlerFunc`](https://pkg.go.dev/net/http#HandlerFunc) type of `net/http`, a handler in Dino is a function, with the difference that this function returns an error. See below for an example of a handler:

```go
func MyHandler() dino.Handler {
    return func(w http.ResponseWriter, r *http.Request) error {
        // Handler logic here
        return nil
    }
}
```

Defining a custom handler that returns an error is a very common pattern because it simplifies error handling, as we will see in the following sections. This pattern is used, for example, on [Fiber](https://gofiber.io/) and [Echo](https://echo.labstack.com/) handlers

### 2. Routing

Dino's handlers implement the [`http.Handler`](https://pkg.go.dev/net/http#Handler) interface, and therefore can be used anywhere that interface is used. An example of this is in [`http.ServeMux`](https://pkg.go.dev/net/http#ServeMux):

```go
mux := http.NewServeMux()

mux.Handle("GET /my/path", MyHandler())
mux.Handle("POST /another/path", AnotherHandler())
```

In fact, `http.ServeMux` is the ideal multiplexer (router) to use with Dino. This is because many of Dino's features are built upon the functionalities that `http.ServeMux` provides. Besides, it's part of the standard library and therefore a very stable router

### 3. Error handling

Dino has a centralized error handling logic. Errors returned by the handlers are automatically processed, which greatly simplifies error handling

There are two possibilities for returning an error in a handler: an error of type [`dino.Error`](https://pkg.go.dev/github.com/willpinha/dino#Error),
or any other error that is not of that type. Below are two examples of handlers that demonstrate these two possibilities:

```go
func KnownErrorHandler() dino.Handler {
    return func(w http.ResponseWriter, r *http.Request) error {
        return dino.NewError(http.StatusBadRequest, "invalid input")
    }
}

func UnknownErrorHandler() dino.Handler {
    return func(w http.ResponseWriter, r *http.Request) error {
        return errors.New("database error")
    }
}
```

#### Response format

When we return an error of type `dino.Error`, it is automatically serialized to a JSON object in the response body. This standardizes the
response format for all errors that your application returns

This error response contains the following fields:

|Field|Required?|Description|
|---|---|---|
|`code`|Yes|The status code of the response|
|`message`|Yes|A descriptive message about the error|
|`details`|No|Additional information about the error (any JSON serializable type)|

When we return an error that is not of type `dino.Error`, Dino generates a generic internal server error response with status code 500.
This prevents unhandled errors from being sent in the response, which could lead to leaks of sensitive data (e.g. database information)

> [!TIP]
> Ideally, all errors returned by a handler are of type `dino.Error`
> 
> This will force you to validate all types of errors that may occur in your handler, which results in greater control and transparency
> regarding errors and data that are returned to clients

#### Creating instances of `dino.Error`

The correct way to create new instances of dino.Error is through the `dino.NewError` function. In addition to passing the status code and
message, we can also pass any [`dino.ErrorOption`](https://pkg.go.dev/github.com/willpinha/dino#ErrorOption):

|Option|Description|
|---|---|
|[`dino.WithDetails`](https://pkg.go.dev/github.com/willpinha/dino#WithDetails)|Additional details. Ideal for when the message is insufficient to describe the error. It can be of any type that is serializable to JSON|
|[`dino.WithInternalError`](https://pkg.go.dev/github.com/willpinha/dino#WithInternalError)|An internal error that caused the error to be returned. This internal error is not serialized in the response, but can be used for logging or debugging purposes|
|[`dino.WithoutLog`](https://pkg.go.dev/github.com/willpinha/dino#WithoutLog)|Indicates that this error should not be logged. By default, all errors are logged|

#### Organizing errors in your project

Here we define some examples of best practices for you to organize the different types of errors in your project. Note that
these examples are not mandatory; you can come up with your own solutions and patterns if they make sense for your project

##### Functions that return `dino.Error`

You can create functions that return `dino.Error` to standardize the response format for a particular type of error

We define a generic `NotFound` function, which is used whenever a resource on your server could not be found

```go
func NotFound() dino.Error {
    return dino.NewError(http.StatusNotFound, "Resource not found")
}
```

We can also pass an integer as a parameter, representing the ID of the resource that was not found

```go
func IdentifierNotFound(id int) dino.Error {
    return dino.NewError(http.StatusNotFound, fmt.Sprintf("ID %d not found", id))
}
```

We can then call this function in our handlers, instead of directly calling `dino.NewError`

```go
func GetUserHandler() dino.Handler {
    return func(w http.ResponseWriter, r *http.Request) error {
        userID := 123

        return IdentifierNotFound(userID)
    }
}
```

## 🦖 License

Dino is under the [MIT license](LICENSE)
