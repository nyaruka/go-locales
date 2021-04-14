# go-locales
[![Build Status](https://github.com/nyaruka/go-locales/workflows/CI/badge.svg)](https://github.com/nyaruka/go-locales/actions?query=workflow%3ACI) 
[![codecov](https://codecov.io/gh/nyaruka/go-locales/branch/main/graph/badge.svg)](https://codecov.io/gh/nyaruka/go-locales) 
[![Go Report Card](https://goreportcard.com/badge/github.com/nyaruka/go-locales)](https://goreportcard.com/report/github.com/nyaruka/go-locales)

Library to make [GNU C Library Locales](https://sourceware.org/glibc/wiki/Locales) accessible in go.

```go
import (
	"github.com/nyaruka/go-locales"
)

func main() {
	locales.Query("es_EC", locales.LC_TIME, "d_fmt")  // "%d/%m/%y"
}
```