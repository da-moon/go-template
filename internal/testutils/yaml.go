package testutils

import (
	"bytes"
	"html/template"
	"strconv"

	"github.com/brianvoe/gofakeit/v6"
)

var ymlTemplateWithList = `
foo:
  bar: {{.foobar}}
fizz: {{.fizz}}
dict:
  list:
    - "{{.dictlist0}}"
    - "{{.dictlist1}}"
    - "{{.dictlist2}}"
    - "4": {{.dictlist3}}
      "5": {{.dictlist4}}
int: {{.int}}
bool: {{.bool}}
`
var jsonTemplateWithList = `
`

// RandomMap uses json codec to return a randomized map
func RandomMap() (map[string]interface{}, error) {
	m := map[string]string{
		"foobar":    gofakeit.HipsterWord(),
		"fizz":      gofakeit.HipsterWord(),
		"dictlist0": strconv.Itoa(gofakeit.Number(1, 100)),
		"dictlist1": strconv.Itoa(gofakeit.Number(1, 100)),
		"dictlist2": strconv.Itoa(gofakeit.Number(1, 100)),
		"dictlist3": gofakeit.HipsterWord(),
		"dictlist4": gofakeit.HipsterWord(),
		"int":       strconv.Itoa(gofakeit.Number(1, 100)),
		"bool":      strconv.FormatBool(gofakeit.Bool()),
	}
	t := template.Must(template.New("yaml").Parse(ymlTemplateWithList))
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "yaml", m); err != nil {
		panic(err)
	}

	return nil, nil
}
func RandomYaml() bytes.Buffer {
	// ymlTemplateWithList = strings.Trim(ymlTemplateWithList, "\n")
	m := map[string]string{
		"foobar":    gofakeit.HipsterWord(),
		"fizz":      gofakeit.HipsterWord(),
		"dictlist0": strconv.Itoa(gofakeit.Number(1, 100)),
		"dictlist1": strconv.Itoa(gofakeit.Number(1, 100)),
		"dictlist2": strconv.Itoa(gofakeit.Number(1, 100)),
		"dictlist3": gofakeit.HipsterWord(),
		"dictlist4": gofakeit.HipsterWord(),
		"int":       strconv.Itoa(gofakeit.Number(1, 100)),
		"bool":      strconv.FormatBool(gofakeit.Bool()),
	}
	t := template.Must(template.New("yaml").Parse(ymlTemplateWithList))
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "yaml", m); err != nil {
		panic(err)
	}
	return buf
}

// [ TODO ] => implement
func RandomJSON() bytes.Buffer {
	panic("not implemented!")
	// ymlTemplateWithList = strings.Trim(ymlTemplateWithList, "\n")
	m := map[string]string{
		"foobar":    gofakeit.HipsterWord(),
		"fizz":      gofakeit.HipsterWord(),
		"dictlist0": strconv.Itoa(gofakeit.Number(1, 100)),
		"dictlist1": strconv.Itoa(gofakeit.Number(1, 100)),
		"dictlist2": strconv.Itoa(gofakeit.Number(1, 100)),
		"dictlist3": gofakeit.HipsterWord(),
		"dictlist4": gofakeit.HipsterWord(),
		"int":       strconv.Itoa(gofakeit.Number(1, 100)),
		"bool":      strconv.FormatBool(gofakeit.Bool()),
	}
	t := template.Must(template.New("yaml").Parse(ymlTemplateWithList))
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "yaml", m); err != nil {
		panic(err)
	}
	return buf
}

// [ TODO ] => implement
func RandomProperties() bytes.Buffer {
	panic("not implemented!")
	// ymlTemplateWithList = strings.Trim(ymlTemplateWithList, "\n")
	m := map[string]string{
		"foobar":    gofakeit.HipsterWord(),
		"fizz":      gofakeit.HipsterWord(),
		"dictlist0": strconv.Itoa(gofakeit.Number(1, 100)),
		"dictlist1": strconv.Itoa(gofakeit.Number(1, 100)),
		"dictlist2": strconv.Itoa(gofakeit.Number(1, 100)),
		"dictlist3": gofakeit.HipsterWord(),
		"dictlist4": gofakeit.HipsterWord(),
		"int":       strconv.Itoa(gofakeit.Number(1, 100)),
		"bool":      strconv.FormatBool(gofakeit.Bool()),
	}
	t := template.Must(template.New("yaml").Parse(ymlTemplateWithList))
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "yaml", m); err != nil {
		panic(err)
	}
	return buf
}
