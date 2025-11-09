<div align="center">

# HTTP Box

🧰 Lightweight HTTP library compatible with net/http

[Documentation]() ⋅ [License]()

</div>

## About

HTTP Box is a lightweight library that contains utilities compatible with the net/http package.
If you want to work directly with net/http, without relying on abstractions, but don't want to
implement common functionalities from scratch, HTTP Box is the right library for you

## Features

- Handlers that return errors
- Centralized error handling
- Middlewares for various use cases (access logging, rate limit, ...)
- Request utilities:
  - Manual body validation
  - Automatic validation of query and URL parameters
- Response utilities:
  - Responses for various content types (JSON, XML, ...)
  - Standard error responses

## Documentation

All the documentation, including more information about the project, its philosophy, and multiple
usage examples, can be found at [pkg.go.dev/github.com/willpinha/httpbox]()

## License

This library is under the [MIT license]()
