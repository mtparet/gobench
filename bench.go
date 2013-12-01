package main

import (
    "net/http"
    "runtime"
    "sync"
    "fmt"
    "time"
    "flag"
    "os"
)

var wg sync.WaitGroup

func get(nb int, url string){
  tr := &http.Transport{
    DisableKeepAlives:  true,
    DisableCompression: true,
  }

  client := &http.Client{Transport: tr}

  for c := 0; c < nb; c++ {
	  response, err := client.Get(url)
    if err != nil {
      fmt.Printf("%s", err)
    } else {
      response.Body.Close()
    }
  }
  wg.Done()
}

func main() {
    concurrencyPtr := flag.Int("c", 4, "Number of multiple requests")
    processorNbPtr := flag.Int("proc", 1, "Number of processor to use. 1 should be sufficient")
    numberRequestPtr := flag.Int("n", 8000, "Number of requests to perform")

    flag.Parse()

    args := flag.Args()

    if len(args) != 1 {
	    fmt.Println("We must specify an url and that's it")
	    os.Exit(1)
    }

    url := args[0] // url to which we want to perform the benchmark

    runtime.GOMAXPROCS(*processorNbPtr)
    numberRequestGoRun := *numberRequestPtr / *concurrencyPtr

    fmt.Println("Number of multiple requests:", *concurrencyPtr)
    fmt.Println("Number of request by go runtime", numberRequestGoRun)
    fmt.Println("Number of processor used:", *processorNbPtr)
    fmt.Println("Number of total requests:", *numberRequestPtr)
    fmt.Println("Url to benchmark:", url)

    timestart := time.Now()

    wg.Add(*concurrencyPtr)

    fmt.Println("Starting");

    for c := 0; c < *concurrencyPtr; c++ {
	    go get(numberRequestGoRun, url)
    }

    fmt.Println("Performing");
    wg.Wait()
    fmt.Println("Finished")

    timeend := time.Now()
    fmt.Println("Total time:", timeend.Sub(timestart))
    fmt.Println("Requests by second", float64(*numberRequestPtr) / timeend.Sub(timestart).Seconds())
}
