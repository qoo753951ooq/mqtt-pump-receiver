package api_md

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
)

type PumpStation struct {
	Pump_st_id   string `json:"pump_st_id"`
	Pump_st_name string `json:"Pump_st_name"`
	Display      string `json:"-"`
	Org_id       string `json:"-"`
	Etype        string `json:"-"`
}

func (p *PumpStation) GetPumpStation() []*PumpStation {

	req, err := http.NewRequest(http.MethodGet, PumpStation_api_url+"getPumpStation", nil)

	if err != nil {
		fmt.Printf("requsetError: %s\n", err)
		return make([]*PumpStation, 0)
	}

	q := req.URL.Query()

	q.Add("org_id", p.Org_id)
	q.Add("display", p.Display)
	q.Add("etype", p.Etype)

	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Printf("clientDoError: %s\n", err)
		return make([]*PumpStation, 0)
	}

	defer resp.Body.Close()

	bodys, err := ioutil.ReadAll(resp.Body)

	pumpStations := make([]*PumpStation, 0)

	err = json.Unmarshal([]byte(bodys), &pumpStations)

	if err != nil {
		fmt.Printf("jsonError: %s\n", err)
		return make([]*PumpStation, 0)
	}

	SortPumpStn(pumpStations, func(p, q *PumpStation) bool {
		return p.Pump_st_id < q.Pump_st_id
	})

	return pumpStations
}

type PumpStnWrapper struct {
	pumpStn []*PumpStation
	by      func(p, q *PumpStation) bool
}

type SortPumpStnBy func(p, q *PumpStation) bool

func (p PumpStnWrapper) Len() int { // 覆寫Len()方法
	return len(p.pumpStn)
}

func (p PumpStnWrapper) Swap(i, j int) { // 覆寫 Swap() 方法
	p.pumpStn[i], p.pumpStn[j] = p.pumpStn[j], p.pumpStn[i]
}

func (p PumpStnWrapper) Less(i, j int) bool { // 覆寫 Less() 方法
	return p.by(p.pumpStn[i], p.pumpStn[j])
}

// 封裝成 SortPumpStation 方法
func SortPumpStn(pumpStn []*PumpStation, by SortPumpStnBy) {
	sort.Sort(PumpStnWrapper{pumpStn, by})
}
