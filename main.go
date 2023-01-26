package main

import (
    "fmt"
    "io"
    "log"
    "net/http"
    "net/url"
    "regexp"
    "github.com/milad-abbasi/gonfig"
)

type Config struct {
    listeningPort string `default:"8084"`
    destinationURL string `default:"http://127.0.0.1:8080"`
}

func main() {

    var c Config

    err := gonfig.Load().FromFile("config.json").Into(&c)
    if err != nil {
        fmt.Println(err)
    }

    // define destination server URL
    destServerURL, err := url.Parse("http://127.0.0.1:8080")
    if err != nil {
        log.Fatal("invalid origin server URL")
    }

    fmt.Println("Starting server on port " + c.listeningPort)

    waf := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
        //fmt.Printf("[SQLWall] received request at: %s\n", time.Now())

        // set req Host, URL and Request URI to forward a request to the origin server
        req.Host = destServerURL.Host
        req.URL.Host = destServerURL.Host
        req.URL.Scheme = destServerURL.Scheme
        req.RequestURI = ""

        /* QUERY STRING PARSING */

        // get the query string from the request URL
        query := req.URL.RawQuery

        // QUERY Regex

        sqliRegex   := `(?i)(\%3D)|(\%27)|(\')|(\-\-)|(\%3B)|(;)|(\%23)|(\#)|(\%2A)|(\*)`
        orRegex     := `(?i)(\%27)|\'((\%6F)|o|(\%4F))((\%72)|r|(\%52))`
        selectRegex := `(?i)(\%27)|\'((\%73)|s|(\%53))((\%65)|e|(\%45))((\%6c)|l|(\%4c))((\%65)|e|(\%45))((\%63)|c|(\%43))((\%74)|t|(\%54))`
        insertRegex := `(?i)(\%27)|\'((\%69)|i|(\%49))((\%6e)|n|(\%4e))((\%73)|s|(\%53))((\%65)|e|(\%45))((\%72)|r|(\%52))((\%74)|t|(\%54))`
        updateRegex := `(?i)(\%27)|\'((\%75)|u|(\%55))((\%70)|p|(\%50))((\%64)|d|(\%44))((\%61)|a|(\%41))((\%74)|t|(\%54))((\%65)|e|(\%45))`
        deleteRegex := `(?i)(\%27)|\'((\%64)|d|(\%44))((\%65)|e|(\%45))((\%6c)|l|(\%4c))((\%65)|e|(\%45))((\%74)|t|(\%54))((\%65)|e|(\%45))`
        dropRegex   := `(?i)(\%27)|\'((\%64)|d|(\%44))((\%72)|r|(\%52))((\%6f)|o|(\%4f))((\%70)|p|(\%50))`
        escapeCharRegex := `(?i)(exec.*\(.*\))|(char.*\(.*\))|(ASCII.*\(.*\))|(BIN.*\(.*\))|(HEX.*\(.*\))|(UNHEX.*\(.*\))|(BASE64.*\(.*\))|(DEC.*\(.*\))|(ROT13.*\(.*\))`

        queryMatch,  _ := regexp.MatchString(sqliRegex, query)
        orMatch,     _ := regexp.MatchString(orRegex, query)
        selectMatch, _ := regexp.MatchString(selectRegex, query)
        insertMatch, _ := regexp.MatchString(insertRegex, query)
        updateMatch, _ := regexp.MatchString(updateRegex, query)
        deleteMatch, _ := regexp.MatchString(deleteRegex, query)
        dropMatch,   _ := regexp.MatchString(dropRegex, query)
        escapeMatch, _ := regexp.MatchString(escapeCharRegex, query)


        if (query != "" && ( queryMatch || orMatch || selectMatch || insertMatch || updateMatch || deleteMatch || dropMatch || escapeMatch )) {
            rw.WriteHeader(http.StatusBadRequest)
            rw.Write([]byte("Possible SQL Injection detected"))
            return
        } 

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

    log.Fatal(http.ListenAndServe(":8084", waf))
}

