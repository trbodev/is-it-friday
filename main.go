package main

import (
	"context"
	"embed"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"text/template"
	"time"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	Host string `env:"HOST,default=127.0.0.1"`
	Port int    `env:"PORT,default=3000"`
}

var config Config

//go:embed templates
var templates embed.FS

type Type string

const (
	TypeExpresive Type = "expressive"
	TypePlain     Type = "plain"
	TypeBoolean   Type = "boolean"
	TypeJSON      Type = "json"
	TypeYAML      Type = "yaml"
	TypeXML       Type = "xml"
	TypeBinary    Type = "binary"
)

func generateText(t Type, friday bool) string {
	var text string
	if friday {
		switch t {
		case TypeExpresive:
			text = "YES!"
		case TypePlain:
			text = "Yes"
		case TypeBoolean:
			text = fmt.Sprint(friday)
		case TypeJSON:
			text = `{"friday":` + fmt.Sprintf("%v", friday) + `}`
		case TypeYAML:
			text = `friday: ` + fmt.Sprintf("%v", friday)
		case TypeXML:
			text = `<?xml version="1.0" encoding="UTF-8"?><friday>` + fmt.Sprintf("%v", friday) + `</friday>`
		case TypeBinary:
			text = "1"
		}
	} else {
		switch t {
		case TypeExpresive:
			text = "no :("
		case TypePlain:
			text = "No"
		case TypeBoolean:
			text = fmt.Sprint(friday)
		case TypeJSON:
			text = `{"friday":` + fmt.Sprintf("%v", !friday) + `}`
		case TypeYAML:
			text = `friday: ` + fmt.Sprintf("%v", !friday)
		case TypeXML:
			text = `<?xml version="1.0" encoding="UTF-8"?><friday>` + fmt.Sprintf("%v", !friday) + `</friday>`
		case TypeBinary:
			text = "0"
		}
	}
	return text
}

func init() {
	ctx := context.Background()
	if err := envconfig.Process(ctx, &config); err != nil {
		panic(err)
	}
}

func main() {
	tmpl, err := template.ParseFS(templates, "templates/index.html")
	if err != nil {
		panic(err)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("%v:%v", config.Host, config.Port))
	if err != nil {
		panic(err)
	}

	fmt.Printf("http://%v\n", listener.Addr())

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		friday := time.Now().Weekday() == time.Friday
		tmpl.Execute(rw, struct {
			Friday string
		}{
			Friday: generateText(TypeExpresive, friday),
		})
	})

	http.HandleFunc("/plain", func(rw http.ResponseWriter, r *http.Request) {
		friday := time.Now().Weekday() == time.Friday
		rw.Header().Set("Content-Type", "text/plain")
		rw.Write([]byte(generateText(TypePlain, friday)))
	})

	http.HandleFunc("/boolean", func(rw http.ResponseWriter, r *http.Request) {
		friday := time.Now().Weekday() == time.Friday
		rw.Header().Set("Content-Type", "text/plain")
		rw.Write([]byte(generateText(TypeBoolean, friday)))
	})

	http.HandleFunc("/json", func(rw http.ResponseWriter, r *http.Request) {
		friday := time.Now().Weekday() == time.Friday
		rw.Header().Set("Content-Type", "application/json")
		rw.Write([]byte(generateText(TypeJSON, friday)))
	})

	http.HandleFunc("/yaml", func(rw http.ResponseWriter, r *http.Request) {
		friday := time.Now().Weekday() == time.Friday
		rw.Header().Set("Content-Type", "application/x-yaml")
		rw.Write([]byte(generateText(TypeYAML, friday)))
	})

	http.HandleFunc("/xml", func(rw http.ResponseWriter, r *http.Request) {
		friday := time.Now().Weekday() == time.Friday
		rw.Header().Set("Content-Type", "application/xml")
		rw.Write([]byte(generateText(TypeXML, friday)))
	})

	http.HandleFunc("/binary", func(rw http.ResponseWriter, r *http.Request) {
		friday := time.Now().Weekday() == time.Friday
		rw.Header().Set("Content-Type", "text/plain")
		rw.Write([]byte(generateText(TypeBinary, friday)))
	})

	http.HandleFunc("/svg", func(rw http.ResponseWriter, r *http.Request) {
		friday := time.Now().Weekday() == time.Friday
		text := generateText(TypeExpresive, friday)
		color := "1CED2A"
		if !friday {
			color = "ED1C1C"
		}
		rw.Header().
			Add("Location",
				"https://img.shields.io/badge/Is%20It%20Friday%3F-"+text+"-%23"+color+"")
		rw.WriteHeader(http.StatusFound)
	})

	http.HandleFunc("/png", func(rw http.ResponseWriter, r *http.Request) {
		friday := time.Now().Weekday() == time.Friday
		text := generateText(TypeExpresive, friday)
		color := "1CED2A"
		if !friday {
			color = "ED1C1C"
		}
		rw.Header().
			Add("Location",
				"https://raster.shields.io/badge/Is%20It%20Friday%3F-"+text+"-%23"+color+"")
		rw.WriteHeader(http.StatusFound)

	})

	if err := http.Serve(listener, nil); err != nil {
		panic(err)
	}

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	<-s

	listener.Close()
}
