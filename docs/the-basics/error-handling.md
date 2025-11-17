Dino has a centralized error handling logic. Errors returned by the handlers are automatically
processed, which greatly simplifies error handling

There are two possibilities for returning an error in a handler: an error of type
[`dino.Error`](https://pkg.go.dev/github.com/willpinha/dino#Error), or any other error that is not
of that type. Below are two examples of handlers that demonstrate these two possibilities:

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

When we return an error of type `dino.Error`, it is automatically serialized to a JSON object in the
response body. This standardizes the response format for all errors that your application returns

This error response contains the following fields:

| Field     | Required? | Description                                                         |
| --------- | --------- | ------------------------------------------------------------------- |
| `code`    | Yes       | The status code of the response                                     |
| `message` | Yes       | A descriptive message about the error                               |
| `details` | No        | Additional information about the error (any JSON serializable type) |

When we return an error that is not of type `dino.Error`, Dino generates a generic internal server
error response with status code 500. This prevents unhandled errors from being sent in the response,
which could lead to leaks of sensitive data (e.g. database information)

!!! tip

    Ideally, all errors returned by a handler are of type `dino.Error`

    This will force you to validate all types of errors that may occur in your handler, which results
    in greater control and transparency regarding errors and data that are returned to clients

## Creating instances of `dino.Error`

The correct way to create new instances of `dino.Error` is through the
[`dino.NewError`](https://pkg.go.dev/github.com/willpinha/dino#NewError) function. In addition to
passing the status code and message, we can also pass any
[`dino.ErrorOption`](https://pkg.go.dev/github.com/willpinha/httpbox#ErrorOption):

| Option                       | Description                                                                                                                                                      |
| ---------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [`dino.WithDetails`]()       | Additional details. Ideal for when the message is insufficient to describe the error. It can be of any type that is serializable to JSON                         |
| [`dino.WithInternalError`]() | An internal error that caused the error to be returned. This internal error is not serialized in the response, but can be used for logging or debugging purposes |
| [`dino.WithLog`]()           | Indicates that this error should be logged. By default, errors are not logged                                                                                    |
