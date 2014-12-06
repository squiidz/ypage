package ypage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type Data []byte

type Probe struct {
	Id    int
	Key   string
	Start int
	Stop  int
	Total int
	Info  *Fetch
	Data  []byte
	Errs  []string
}

func NewProbe(id int, key string, start int, stop int, tot int, in *Fetch) *Probe {
	return &Probe{Id: id, Key: key, Start: start, Stop: stop, Total: tot, Info: in, Errs: []string{}}
}

func MakeProbe(n int, dataSize int, f *Fetch, keys []string) ([]*Probe, error) {
	probes := []*Probe{}
	splitW := dataSize / n
	if len(keys) != n {
		return nil, errors.New("Not enough key for probs")
	}
	for i, k := range keys {
		workstart := i * splitW
		probes = append(probes, NewProbe(i+1, k, workstart, workstart+splitW, splitW, f))
	}
	return probes, nil
}

func (p *Probe) Work(wg *sync.WaitGroup) {
	fmt.Printf("Probe #%d start working at %d end at %d \n", p.Id, p.Start, p.Stop)
	for i := p.Start; i < p.Stop; i++ {
		rx, err := p.scraper(i)
		if err != nil {
			p.Errs = append(p.Errs, err.Error())
			continue
		}
		raw, err := ioutil.ReadAll(rx.Body)
		if err != nil {
			p.Errs = append(p.Errs, err.Error())
			continue
		}
		for _, b := range raw {
			p.Data = append(p.Data, b)
		}

		at := (100 / p.Total) * (i - p.Start)
		// API Limitation 1req/sec
		fmt.Printf("Probe #%d  at %d/100 \n", p.Id, at)
		time.Sleep(time.Second * 1)
	}
	fmt.Printf("Probe #%d Finish\n", p.Id)
	wg.Done()
}

func (p *Probe) Extract() Data {
	return p.Data
}

func (d Data) ToJson() []byte {
	result, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	return result
}

func (p *Probe) scraper(pageNbr int) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s/?what=%s&where=%s&pgLen=%d&pg=%d&dist=%d&fmt=%s&lang=%s&UID=%s&apikey=%s",
		p.Info.BaseUrl,
		p.Info.Service,
		p.Info.What,
		p.Info.Where,
		p.Info.PageLen,
		pageNbr,
		p.Info.Dist,
		p.Info.Format,
		p.Info.Lang,
		p.Info.UID,
		p.Key,
	)
	cli := http.Client{}
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := cli.Do(r)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
