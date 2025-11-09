package httpbox

/*
HTTP Box is a lightweight library that contains utilities compatible with the net/http package.
If you want to work directly with net/http, without relying on abstractions, but don't want to
implement common functionalities from scratch, HTTP Box is the right library for you

HTTP Box provides the following features:

- Handlers and middlewares that return errors
- Centralized error handling
- Middlewares for various use cases (logging, rate limit, ...)
- Request utilities:
	- Body validation
	- Automatic validation of query and URL parameters
- Response utilities:
	- Responses for various content types (JSON, XML, ...)
	- Standard error responses
*/
