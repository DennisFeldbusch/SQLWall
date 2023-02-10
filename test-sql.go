package main

import (
	"fmt"
	"net/http"
	"time"
    "os"
    "bufio"
)

func main() {
	startTime := time.Now()

	counter := 0
    
    // read flag
    flag := os.Args[1]

    file := openFile(flag)

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		url := "http://127.0.0.1:8084/?id=" + line
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == 400 {
			counter++
		}
	}

	endTime := time.Now()
	fmt.Println("Time: ", endTime.Sub(startTime))
	fmt.Println("Counter: ", counter)
}

func openFile(flag string) http.File {
    if flag == "s" {
        file, _ := os.Open("./SQL_3.txt")
        return file
    } else if flag == "v" {
        file, _ := os.Open("./httparchive_parameters_top_1m_2022_12_28.txt")
        return file
    } else {
        return nil
    }
}
