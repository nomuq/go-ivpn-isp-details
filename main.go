package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func main() {

	url := "https://api.ivpn.net/v4/servers/stats"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	result, err := UnmarshalResult(body)
	if err != nil {
		fmt.Println(err)
		return
	}

	data := [][]string{}

	for _, server := range result.Servers {
		if server.Hostnames.Openvpn != nil {
			data = append(data, []string{
				server.Country, server.City, server.ISP, server.Gateway, strings.Join(server.Protocols, ", "),
			})
		}
	}
	sort.Slice(data, func(i, j int) bool {
		return data[i][0] < data[j][0]
	})
	sort.Slice(data, func(i, j int) bool {
		switch strings.Compare(data[i][0], data[j][0]) {
		case -1:
			return true
		case 1:
			return false
		}
		return data[i][1] > data[j][1]
	})
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetHeader([]string{"Country", "City", "ISP", "Gateway", "Protocols"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetRowLine(true)
	table.SetAutoMergeCellsByColumnIndex([]int{0, 1})
	table.AppendBulk(data)
	table.Render()
	fmt.Println(tableString.String())
}

func UnmarshalResult(data []byte) (Result, error) {
	var r Result
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Result) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Result struct {
	Servers []Server `json:"servers"`
}

type Server struct {
	Gateway     string    `json:"gateway"`
	Hostnames   Hostnames `json:"hostnames"`
	IsActive    bool      `json:"is_active"`
	CountryCode string    `json:"country_code"`
	Country     string    `json:"country"`
	City        string    `json:"city"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	ISP         string    `json:"isp"`
	Load        float64   `json:"load"`
	Protocols   []string  `json:"protocols"`
	WgPublicKey string    `json:"wg_public_key"`
}

type Hostnames struct {
	Openvpn   *string `json:"openvpn,omitempty"`
	Wireguard *string `json:"wireguard,omitempty"`
}
