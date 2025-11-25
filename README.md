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

## 🦖 License

Dino is under the [MIT license](LICENSE)
