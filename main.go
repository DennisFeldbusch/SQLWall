package main

import (
    "fmt"
    "io"
    "log"
    "net/http"
    "net/url"
    "time"
    "regexp"
)

func main() {

    // define origin server URL
    originServerURL, err := url.Parse("http://127.0.0.1:8081")
    if err != nil {
        log.Fatal("invalid origin server URL")
    }

    reverseProxy := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
        fmt.Printf("[reverse proxy server] received request at: %s\n", time.Now())

        // set req Host, URL and Request URI to forward a request to the origin server
        req.Host = originServerURL.Host
        req.URL.Host = originServerURL.Host
        req.URL.Scheme = originServerURL.Scheme
        req.RequestURI = ""

        // get the path from the request url
        path := req.URL.Path

        // get the query string from the request URL
        query := req.URL.RawQuery

        fmt.Println("[PATH] ", path)
        fmt.Println("[QUERY] ", query)

        /* PATH STRING PARSING */

        pathRegex := `^(/\w|\d)*(/(\w|\d)*\.(\w|\d)*)`

        /* QUERY STRING PARSING */
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

    log.Fatal(http.ListenAndServe(":8080", reverseProxy))
}

