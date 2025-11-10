package httpbox

/*
# Handlers

Similar to the `net/http` handlers, the [Handler] type represents the core mechanism for processing
HTTP requests and generating responses. The difference is that in httpbox, handlers return errors

We can define a handler in the following way:

	func MyHandler() httpbox.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			// Handler logic here
			return nil
		}
	}

These handlers implement the http.Handler interface, and therefore can be used anywhere that interface
is used. An example of this is in http.ServeMux:

	mux := http.NewServeMux()

	mux.Handle("GET /path", MyHandler())

In fact, http.ServeMux is the ideal multiplexer (router) to use with httpbox. This is because many of
httpbox's features are built upon the functionalities that http.ServeMux provides. Besides, it's part
of the standard library

# Error handling

httpbox has a centralized error handling logic. Errors returned by the handlers are automatically
processed, which greatly simplifies error handling

There are two possibilities for returning an error in a handler: an error of type [Error], or any
other error that is not of that type

	func KnownErrorHandler() httpbox.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			return httpbox.NewError(http.StatusBadRequest, "some message")
		}
	}

	func UnknownErrorHandler() httpbox.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			return errors.New("an unknown error occurred")
		}
	}

When we return an error of type [Error], it is automatically serialized to JSON in the response body.
This response contains the following fields:

- `code` (required): The status code of the response
- `message` (required): A descriptive message about the error
- `details` (optional): Additional information about the error (any serializable type)

When we return an error that is not of type [Error], httpbox generates a generic internal server
error response with status code 500. This prevents unhandled errors from being sent in the response,
which could lead to leaks of sensitive data (e.g. database info). This means that ideally all errors
returned by a handler are of type [Error]

We use [NewError] to create new errors of type [Error]. In addition to passing the status code and
message, we can also pass the following options:

	func MyHandler() httpbox.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			return httpbox.NewError(
				http.StatusNotFound,
				"resource not found",
				httpbox.WithDetails("Some additional details"),
				httpbox.WithInternalError(errors.New("database: record not found")
				httpbox.WithLog(),
			)
		}
	}

- [WithDetails]: Additional details. Ideal for when the message is insufficient to describe the error.
	It can be of any type that is serializable to JSON
- [WithInternalError]: An internal error that caused the error to be returned. This internal error is
	not serialized in the response, but can be used for logging or debugging purposes
- [WithLog]: Indicates that this error should be logged. By default, errors are not logged. It must
	be used in conjunction with the [AccessLogMiddleware] middleware

# Middlewares

## RecoverMiddleware

## AccessLogMiddleware

# Request utilities

# Response utilities

# Logging
*/
