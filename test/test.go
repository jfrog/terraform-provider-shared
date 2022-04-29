package test

import (
	"bytes"
	"math/rand"
	"text/template"
	"time"
)

func RandomInt() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(10000000)
}

func RandBool() bool {
	return RandomInt()%2 == 0
}

func RandSelect(items ...interface{}) interface{} {
	return items[RandomInt()%len(items)]
}

func ExecuteTemplate(name, temp string, fields interface{}) string {
	var tpl bytes.Buffer
	if err := template.Must(template.New(name).Parse(temp)).Execute(&tpl, fields); err != nil {
		panic(err)
	}

	return tpl.String()
}
