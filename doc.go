package httpbox

/*
HTTP Box is a lightweight library that contains utilities compatible with the net/http package.
If you want to work directly with net/http, without relying on abstractions, but don't want to
implement common functionalities from scratch, HTTP Box is the right library for you

HTTP Box provides features that are not in the standard library, such as:

- Handlers and middlewares that return errors
- Centralized error handling
- Middlewares for various use cases (logging, rate limit, ...)
- Request utilities:
	- Automatic validation of body, query and URL parameters
- Response utilities:
	- Standard error responses
*/
