Dino provides several utilities to help you easily extract and validate data from your requests.
This request data can come from the URL (path & query parameters), the body, or the headers. Below
you will see how to obtain and validate the data from each of these sources

## Reading the URL (path & query parameters)

The [`dino.Param`]() type represents parameters that come from the URL path or query string. It
contains validations that you can use to convert the data received from the request into a specific
type in Go. To obtain these parameters, we can use 4 different functions:

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

| Function                    | Description |
| --------------------------- | ----------- |
| [`NewPathParam`]()          |             |
| [`NewQueryParam`]()         |             |
| [`NewDefaultQueryParam`]()  |             |
| [`NewRequiredQueryParam`]() |             |

## Reading the body

## Reading headers
