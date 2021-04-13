# go-locales

Library to make glibc locale data accessible in go.

```go
import (
	"github.com/nyaruka/go-locales"
)

func main() {
	locales.Query("es_EC", locales.LC_TIME, "d_fmt")  // "%d/%m/%y"
}
```