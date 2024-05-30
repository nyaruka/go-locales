package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/nyaruka/go-locales"
	"github.com/nyaruka/go-locales/localedata"
)

// matches locale codes in the form xx_YY or xxx_YY
var bcp47Regex = regexp.MustCompile(`^[a-z]{2,3}_[A-Z]{2}$`)

// data dumped for a single locale
type localeDump map[string][]string

func main() {
	var merge, pretty bool

	flags := flag.NewFlagSet("", flag.ExitOnError)
	flags.BoolVar(&merge, "merge", false, "Merge identical data of same language")
	flags.BoolVar(&pretty, "pretty", false, "Pretty format output")
	flags.Parse(os.Args[1:])

	// read rest of arguments as mappings
	mappings := make(map[string]string, len(flags.Args()))
	for _, m := range flags.Args() {
		v := strings.Split(m, "=")
		if len(v) != 2 || strings.Count(v[1], ".") != 1 {
			fmt.Printf("invalid mapping: %s", m)
			os.Exit(1)
		}
		mappings[v[0]] = v[1]
	}

	if err := localesDump(mappings, merge, pretty); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func localesDump(mappings map[string]string, merge, pretty bool) error {
	data, err := extractData(mappings)
	if err != nil {
		return err
	}

	if merge {
		data = mergeLocales(data)
	}

	var marshaled []byte

	if pretty {
		marshaled, err = json.MarshalIndent(data, "  ", "  ")
	} else {
		marshaled, err = json.Marshal(data)
	}
	if err != nil {
		return err
	}

	fmt.Println(string(marshaled))

	return nil
}

// extract all mappings from all locales
func extractData(mappings map[string]string) (map[string]localeDump, error) {
	data := make(map[string]localeDump)

	for _, code := range locales.Codes() {
		// only interesting in real languages
		if !bcp47Regex.MatchString(code) {
			continue
		}

		locale := make(localeDump, len(mappings))
		for key, path := range mappings {
			parts := strings.Split(path, ".")
			category := parts[0]
			keyword := parts[1]

			ops, err := locales.Query(code, localedata.LC(category), keyword)
			if err != nil {
				return nil, fmt.Errorf("error querying path %s: %w", path, err)
			}

			locale[key] = ops
		}

		data[code] = locale
	}

	return data, nil
}

// merge locales with same language and same extracted data
func mergeLocales(data map[string]localeDump) map[string]localeDump {
	distinctByLang := make(map[string][]localeDump, len(data))

	for code, dump := range data {
		lang := strings.SplitN(code, "_", 2)[0]
		distinct := true

		for _, existing := range distinctByLang[lang] {
			if reflect.DeepEqual(existing, dump) {
				distinct = false
				break
			}
		}

		if distinct {
			distinctByLang[lang] = append(distinctByLang[lang], dump)
		}
	}

	merged := make(map[string]localeDump, len(data))

	for code, dump := range data {
		lang := strings.SplitN(code, "_", 2)[0]
		variesByCountry := len(distinctByLang[lang]) > 1

		if variesByCountry {
			merged[code] = dump
		} else {
			merged[lang] = dump
		}
	}

	return merged
}
