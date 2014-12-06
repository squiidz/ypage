package ypage

import (
	"fmt"
	"net/http"
	"time"
)

type Fetch struct {
	BaseUrl string
	Service string
	What    string
	Where   string
	PageLen int
	NbrPage int
	Dist    int
	Format  string
	Lang    string
	UID     string
	Key     string
}

type exist struct {
	Fetched *Fetch
	Resp    *http.Response
}

func NewFetch(base string, uid string, key string, option ...func(*Fetch)) *Fetch {
	fetch := &Fetch{
		BaseUrl: base,
		PageLen: 700,
		NbrPage: 1,
		Dist:    3,
		Format:  "JSON",
		Lang:    "en",
		UID:     uid,
		Key:     key,
	}

	for _, op := range option {
		op(fetch)
	}

	return fetch
}

func (f *Fetch) FindBusiness(where string, what string) (*exist, error) {
	c := http.Client{}
	f.Service = "FindBusiness"
	f.Where = where
	f.What = what

	req := f.builder()
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
		fmt.Println("ERROR AT FINDBUSINESS FUNC")
	}
	return &exist{f, resp}, nil
}

func (f *Fetch) builder() *http.Request {
	url := fmt.Sprintf("%s/%s/?what=%s&where=%s&pgLen=%d&pg=%d&dist=%d&fmt=%s&lang=%s&UID=%s&apikey=%s",
		f.BaseUrl,
		f.Service,
		f.What,
		f.Where,
		f.PageLen,
		f.NbrPage,
		f.Dist,
		f.Format,
		f.Lang,
		f.UID,
		f.Key,
	)
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	return r
}

func (e *exist) FetchAll(size int) []*http.Response {
	data := []*http.Response{}
	fmt.Println("Size :", size)
	for i := 0; i < size; i++ {
		rx, err := e.scraper(i)
		if err != nil {
			continue
		}
		data = append(data, rx)
		fmt.Println("Page :", i)
		// API Limitation 1req/sec
		time.Sleep(time.Second * 1)
	}
	return data
}

func (e *exist) scraper(pageNbr int) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s/?what=%s&where=%s&pgLen=%d&pg=%d&dist=%d&fmt=%s&lang=%s&UID=%s&apikey=%s",
		e.Fetched.BaseUrl,
		e.Fetched.Service,
		e.Fetched.What,
		e.Fetched.Where,
		e.Fetched.PageLen,
		pageNbr,
		e.Fetched.Dist,
		e.Fetched.Format,
		e.Fetched.Lang,
		e.Fetched.UID,
		e.Fetched.Key,
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
