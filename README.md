# GO_SATORI
This is the go driver for Satori, it has all of the operations implemented.

# Example

```go
package main

import (
	"fmt"

	satori "github.com/satori-db/go_satori"
)

func main() {

	st := satori.Satori{Host: "127.0.0.1", Port: "2310", Username: "john wick", Token: "example"}
	fmt.Print(st.Heartbeat())

}

```