# go-locales
[![Build Status](https://github.com/nyaruka/go-locales/workflows/CI/badge.svg)](https://github.com/nyaruka/go-locales/actions?query=workflow%3ACI) 
[![codecov](https://codecov.io/gh/nyaruka/go-locales/branch/main/graph/badge.svg)](https://codecov.io/gh/nyaruka/go-locales) 
[![Go Report Card](https://goreportcard.com/badge/github.com/nyaruka/go-locales)](https://goreportcard.com/report/github.com/nyaruka/go-locales)

Library to make [GNU C Library Locales](https://sourceware.org/glibc/wiki/Locales) accessible in go.

```
go get github.com/nyaruka/go-locales
```

then

```go
locales.Query("es_EC", locales.LC_TIME, "abday")                       // []string{"dom", ..., "sáb"}
locales.QueryString("fr_CA", locales.LC_NAME, "name_fmt")              // "%d%t%g%t%m%t%f"
locales.QueryInteger("zh_CN", locales.LC_MONETARY, "int_frac_digits")  // 2

locales.Codes()  // []string{"POSIX", "aa_DJ", "aa_ER", ... }
```

## localesdump

This command line tool can be used to create JSON dumps of select locale data, e.g.

```
localesdump --pretty days=LC_TIME.day short_days=LC_TIME.abday
```

Extracts the given keywords from each locale and produces JSON like:

```
{
    "aa_DJ": {
        "days": ["Acaada", "Etleeni", "Talaata", "Arbaqa", "Kamiisi", "Gumqata", "Sabti"],
        "short_days": ["Aca", "Etl", "Tal", "Arb", "Kam", "Gum", "Sab"]
    },
    "aa_ER": {
        "days": ["Acaada", "Etleeni", "Talaata", "Arbaqa", "Kamiisi", "Gumqata", "Sabti"],
        "short_days": ["Aca", "Etl", "Tal", "Arb", "Kam", "Gum", "Sab"]
    },
    ...
}
```