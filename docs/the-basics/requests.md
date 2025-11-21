# Requests

Dino provides several utilities to help you easily extract and validate data from your requests.
This request data can come from the URL (path & query parameters), the body, or the headers. Below
you will see how to obtain and validate the data from each of these sources

## Reading the URL (path & query parameters)

The [`dino.Param`](https://pkg.go.dev/github.com/willpinha/dino#Param) type represents parameters
that come from the URL path (e.g. `/users/{id}`) or query string (e.g. `?name=will`).

### Path parameters

### Query parameters

```go
func MyHandler() dino.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
        id, err := dino.NewPathParam(r, "id").Int()
        if err != nil {
            return err
        }

        name := dino.NewPathParam(r, "name").String()

        return nil
    }
}
```

| Function                                                                                      | Description |
| --------------------------------------------------------------------------------------------- | ----------- |
| [`NewPathParam`](https://pkg.go.dev/github.com/willpinha/dino#NewPathParam)                   |             |
| [`NewQueryParam`](https://pkg.go.dev/github.com/willpinha/dino#NewQueryParam)                 |             |
| [`NewDefaultQueryParam`](https://pkg.go.dev/github.com/willpinha/dino#NewDefaultQueryParam)   |             |
| [`NewRequiredQueryParam`](https://pkg.go.dev/github.com/willpinha/dino#NewRequiredQueryParam) |             |

## Reading the body

## Reading headers
