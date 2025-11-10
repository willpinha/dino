# httpbox

> Lightweight HTTP library compatible with `net/http`

[Documentation](https://pkg.go.dev/github.com/willpinha/httpbox)

## About

httpbox is a lightweight library that contains utilities compatible with the `net/http` package.
If you want to work directly with `net/http`, without relying on abstractions (routers, frameworks,
...), but don't want to implement common functionalities from scratch, HTTP Box is the right library
for you

## Philosophy

httpbox doesn't try to be the magic solution that solves all your problems. Instead, it provides a
thin layer built on top of `net/http` for functionalities commonly needed when building applications

It also follows the Go philosophy, which is to maintain simplicity and not apply breaking changes or
major versions all the time. That's why httpbox is still in version 1 (just like Go)

## Features

- Handlers that return errors
- Centralized error handling
- Middlewares for various use cases (access logging, rate limit, ...)
- Request utilities:
  - Manual body validation
  - Automatic validation of query and URL parameters
- Response utilities:
  - Responses for all content types (JSON, XML, ...)
  - Standard error responses
