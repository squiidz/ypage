package main

import (
	"fmt"
	"github.com/squiidz/ypage"
	"sync"
)

type KeyStack []Key

type Key string

// First Key "qs9x872kthgk4aur4u6x2xr9"
// Second Key "je2dq59qgv5b7nrx3a6swpfj"
// Third Key "eh2vk49jvdgmm66dymcre2xy"

func main() {
	// Create a Array of API Key
	ks := []string{"eh2vk49jvdgmm66dymcre2xy"}

	// Create a First request to see the data lenght
	f := ypage.NewFetch("http://api.sandbox.yellowapi.com", "Dev", ks[0])

	// Source file to put incoming data
	fd := ypage.NewFile(".", "data_1.json", 0600)

	// Make the request with the provided infos
	res, err := f.FindBusiness("Montreal", "Restaurant")
	if err != nil {
		fmt.Println(err)
	}

	// Make a Payload from the init request
	p := ypage.NewLoad(res.Resp)
	fmt.Println(p.Summary.PageCount)

	// Create and start probes
	var wg sync.WaitGroup
	probes, err := ypage.MakeProbe(1, p.Summary.PageCount, f, ks)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, pr := range probes {
			wg.Add(1)
			go pr.Work(&wg)
		}
		wg.Wait()
	}

	// For all pages build a payload with the response,
	// make the result writable to the source
	// and write the content data in the Source file
	for _, p := range probes {
		// Loop over the response to make a payload with it
		ex := ypage.BuildLoad(p.Extract())
		// Loop over the Payloads
		// Insert Data to the source
		fd.Insert(ex)
		// Add a comma to make the json format valid in the source
		fd.Insert([]byte(","))

	}

}
