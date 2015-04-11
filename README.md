# fever

fever is micro framework.

## Include

fever contains what was modified the following module

* [alice](https://github.com/justinas/alice)
* [stack/mux](https://github.com/nmerouze/stack)

## Usage

```go
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/mix3/fever"
	"github.com/mix3/fever/mux"
	"golang.org/x/net/context"
)

func main() {
	m := mux.New()
	m.Get("/:name").ThenFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
		p := mux.Params(c)
		fmt.Fprintf(w, "Hello %s", p.ByName("name"))
	})
	fever.Run(":19300", 10*time.Second, m)
}
```

```
go run examples/simpl/main.go
```
or
```
go get github.com/lestrrat/go-server-starter/cmd/start_server
start_server --port 19300 -- go run examples/simple/main.go
```

## LICENSE

* fever under MIT license
* [chain/\*](chain) under MIT license
* [mux/\*](mux) under MPL 2.0 license

see [LICENSE](LICENSE) for details.
