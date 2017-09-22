# go-reachsms
A Go lang wrapper for the [Reach SMS API](https://www.reach-interactive.com/developers/api/)

## Install
```shell
go get github.com/brentmountford/go-reachsms
```

## Usage

Import the library

```go
import "github.com/brentmountford/go-reachsms"
```

## Example
```go
package main

import (
    "fmt"
    "github.com/brentmountford/go-reachsms"
)

func main() {
    reachSmsApi, _ := reachsms.Create("username", "password")

    fmt.Println(reachSmsApi.GetBalance())
    fmt.Println(reachSmsApi.GetMessage("f1dcd3fe-4cca-46a6-b3c5-fd4dd0xxx"))

    message := reachsms.NewMessage(
        "07376959xxx,07867440xxx",
        "Test",
        "This is a test message",
    )
    fmt.Println(message)
    fmt.Println(reachSmsApi.SendMessage(message))
}

```

## Documentation
https://godoc.org/github.com/brentmountford/go-reachsms


