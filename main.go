package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"time"
	"github.com/milad-abbasi/gonfig"
)

type Config struct {
    Port string `default:"8080"`
    Destination string `default:"http://127.0.0.1:8081"`
}

func main() {

    var c Config

    err := gonfig.Load().FromFile("config.json").Into(&c)
	if err != nil {
		fmt.Println(err)
	}

    // define destination server URL
    destServerURL, err := url.Parse(c.Destination)
    if err != nil {
        log.Fatal("invalid origin server URL")
    }

    fmt.Println("Starting server on port " + c.Port)

    waf := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
        fmt.Printf("[SQLWall] received request at: %s\n", time.Now())

        // set req Host, URL and Request URI to forward a request to the origin server
        req.Host = destServerURL.Host
        req.URL.Host = destServerURL.Host
        req.URL.Scheme = destServerURL.Scheme
        req.RequestURI = ""

        /* PATH STRING PARSING */

        // get the path from the request url
        path := req.URL.Path
        fmt.Println("[PATH] ", path)

        // PATH Regex
        pathRegex := `^(/\w|\d)*(/(\w|\d)*\.(\w|\d)*)`

        pathMatch, _ := regexp.MatchString(pathRegex, path)

        if (path != "" && pathMatch == false) {
            rw.WriteHeader(http.StatusBadRequest)
            rw.Write([]byte("Bad Request"))
            return
        }

        /* QUERY STRING PARSING */

        // get the query string from the request URL
        query := req.URL.RawQuery
        fmt.Println("[QUERY] ", query)

        // QUERY Regex
        queryRegex := `^(\w|\d)+=(\w|\d)+(&(\w|\d)+=(\w|\d)+)*$` 

        queryMatch, _ := regexp.MatchString(queryRegex, query)

        if (query != "" && queryMatch == false) {
            rw.WriteHeader(http.StatusBadRequest)
            rw.Write([]byte("Bad Request"))
            return
        }

        // print out the http Request
        // EXPECTED: GET / HTTP2

        //requestDump, err := httputil.DumpRequest(req, true)
        //fmt.Println("[WHOLE] "+ string(requestDump))

        // save the response from the origin server
        originServerResponse, err := http.DefaultClient.Do(req)
        if err != nil {
            rw.WriteHeader(http.StatusInternalServerError)
            _, _ = fmt.Fprint(rw, err)
            return
        }

        // return response to the client
        rw.WriteHeader(http.StatusOK)
        io.Copy(rw, originServerResponse.Body)
    })

    log.Fatal(http.ListenAndServe(":"+c.Port, waf))
}

