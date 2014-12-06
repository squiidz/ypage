package ypage

import (
	"encoding/json"
	"net/http"
)

type Payload struct {
	Summary  Summary  `json:"summary"`
	Listings []Entity `json:"listings"`
}

type Summary struct {
	What            string `json:"what"`
	Where           string `json:"where"`
	Latitude        string `json:"latitude"`
	Longitude       string `json:"longitude"`
	FirstListing    int    `json:"firstListing"`
	LastListing     int    `json:"lastListing"`
	TotalListings   int    `json:"totalListings"`
	PageCount       int    `json:"pageCount"`
	CurrentPage     int    `json:"currentPage"`
	ListingsPerPage int    `json:"listingsPerPage"`
	Prov            string `json:"prov"`
}

type Entity struct {
	Id       string  `json:"id"`
	Name     string  `json:"name"`
	Address  Address `json:"address"`
	GeoCode  GeoCode `json:"geocode"`
	Distance string  `json:"distance"`
	ParentId string  `json:"parentId"`
	IsParent bool    `json:"isParent"`
	Content  Content `json:"content"`
}

type Address struct {
	Street string `json:"street"`
	City   string `json:"city"`
	Prov   string `json:"prov"`
	PCode  string `json:"pCode"`
}

type GeoCode struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type Content struct {
	Video   Watcher `json:"video"`
	Photo   Watcher `json:"photo"`
	Profile Watcher `json:"profile"`
	DspAd   Watcher `json:"dspAd"`
	Logo    Watcher `json:"logo"`
	Url     Watcher `json:"url"`
}

type Watcher struct {
	Avail bool `json:"avail"`
	InMkt bool `json:"inMkt"`
}

func NewLoad(resp *http.Response) *Payload {
	pay := &Payload{Listings: []Entity{}}
	err := json.NewDecoder(resp.Body).Decode(pay)
	if err != nil {
		return nil
	}
	return pay
}

func BuildList(b []byte) *Entity {
	ent := &Entity{}
	err := json.Unmarshal(b, ent)
	if err != nil {
		return nil
	}
	return ent
}

func (p *Payload) Readable() (string, error) {
	stack, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return "", err
	}
	return string(stack), nil
}

func (p *Payload) Writable() ([]byte, error) {
	stack, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return nil, err
	}
	return stack, nil
}

func (e *Entity) Writable() ([]byte, error) {
	stack, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return nil, err
	}
	return stack, nil
}

func (e *Entity) GetMoreInfo() {

}
