package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"example.com/m/utils/envget"
	"github.com/milad-abbasi/gonfig"
)

type Config struct {
	Port        string `default:"8080"`
	Destination string `default:"http://127.0.0.1:8081"`
}

func main() {

	var c Config

	err := gonfig.Load().FromFile("config.json").Into(&c)
	if err != nil {
		fmt.Println(err)
	}
	c.Destination = envget.GetEnv("DESTINATION", c.Destination) // get the destination server URL from the environment variable (else use the default value)
	c.Port = envget.GetEnv("PORT", c.Port)                      // get the port from the environment variable (else use the default value)

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

		/* QUERY STRING PARSING */

		// get the query string from the request URL
		query := req.URL.RawQuery
		fmt.Println("[QUERY] ", query)

		// QUERY Regex
		queryRegex := `^(\w|\-|(%[0-9A-Fa-f]{2}))+=(\w|\-|(%[0-9A-Fa-f]{2}))+(&(\w|\-|(%[0-9A-Fa-f]{2}))+=(\w|\-|(%[0-9A-Fa-f]{2}))+)*$`

		/*

		   REGEX EXPLANATION

		   ^ asserts position at start of a line

		   \w matches any digit or letter A-z or a underscore sign _
		   \- matches the - sign
		   (%[0-9A-Fa-f]{2}) matches any url encoded char like %20, %22, $a4, ... does NOT match %2x, %4-4, %xx, ...

		   (pattern)+ matches one or more occurences of the pattern
		   (pattern)* matches zero or more occurences of the pattern

		   simplified the used regex can be expressed like the following

		   ^(pattern)+=(pattern)+(&(pattern)+=(pattern)+)*$

		   where pattern is one of the matches described above

		   this is the definition of the query part in a url = key=value&key2=value2&key3=value3

		   $ asserts position at the end of a line

		*/

		queryMatch, _ := regexp.MatchString(queryRegex, query)

		if query != "" && !queryMatch {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("Bad Request"))
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

	log.Fatal(http.ListenAndServe(":"+c.Port, waf))
}
