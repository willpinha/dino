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

## Response format

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

## Organizing errors in your project

Here we define some examples of best practices for you to organize the different types of errors in
your project. Note that these examples are not mandatory; you can come up with your own solutions and
patterns if they make sense for your project

### Functions that return `dino.Error`

You can create functions that return `dino.Error` to standardize the response format for a particular
type of error

!!! example

    We define a generic `NotFound` function, which is used whenever a resource on your server could
    not be found

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

### Grouping error functions into a single package

Your project may have multiple functions that return `dino.Error`. To organize them, you can define
them in a common package, separating them into files within that package

!!! example

    Let's say we have a separate directory structure divided into `cmd`, `internal`, and `pkg`. Inside
    `internal`, we have a `errors` package where we will store all the errors from our project

    ```
    .
    ├── cmd
    ├── internal
    │   └── errors
    │       ├── http.go
    │       └── users.go
    └── pkg
    ```

    The `http.go` file will have generic HTTP errors (not found, bad request, etc.)

    ```go title="http.go"
    package errors

    func NotFound() dino.Error {
    	return dino.NewError(http.StatusNotFound, "Resource not found")
    }

    func BadRequest() dino.Error {
    	return dino.NewError(http.StatusBadRequest, "There is something wrong with your request")
    }
    ```

    And the `users.go` file will have errors specifically related to users

    ```go title="users.go"
    package errors

    func UserIsBlocked(username string) dino.Error {
    	return dino.NewError(
    		http.StatusForbidden,
    		fmt.Sprintf("User %s is blocked", username),
    		dino.WithDetails("Please, contact the HR for more information")
    	)
    }
    ```
