package api_md

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
)

type Pump struct {
	St_no         string   `json:"st_no"`
	Datatime      string   `json:"datatime"`
	Tachometer    *float64 `json:"tachometer,omitempty"`
	Run_status    string   `json:"run_status,omitempty"`
	Device_status string   `json:"device_status,omitempty"`
	Trust         bool     `json:"trust"`
	Pump_st_ids   string   `json:"-"`
	Display       string   `json:"-"`
	RequestBody   []byte   `json:"-"`
}

func (p *Pump) GetPumpList() []*Pump {

	req, err := http.NewRequest(http.MethodGet, PumpStation_api_url+"getPS_PumpList", nil)

	if err != nil {
		fmt.Printf("requsetError: %s\n", err)
		return make([]*Pump, 0)
	}

	q := req.URL.Query()

	q.Add("pump_st_ids", p.Pump_st_ids)
	q.Add("display", p.Display)

	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Printf("clientDoError: %s\n", err)
		return make([]*Pump, 0)
	}

	defer resp.Body.Close()

	bodys, err := ioutil.ReadAll(resp.Body)

	pumps := make([]*Pump, 0)

	err = json.Unmarshal([]byte(bodys), &pumps)

	if err != nil {
		fmt.Printf("jsonError: %s\n", err)
		return make([]*Pump, 0)
	}

	SortPump(pumps, func(p, q *Pump) bool {
		return p.St_no < q.St_no
	})

	return pumps
}

func (p *Pump) PostPumpValue() string {
	req, err := http.NewRequest(http.MethodPost, PumpStation_api_url+"postPumpValue", bytes.NewBuffer(p.RequestBody))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("post Pump %s\n", err)
		return "Not OK"
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	return string(body)
}

type PumpWrapper struct {
	pump []*Pump
	by   func(p, q *Pump) bool
}

type SortPumpBy func(p, q *Pump) bool

func (p PumpWrapper) Len() int { // 覆寫Len()方法
	return len(p.pump)
}

func (p PumpWrapper) Swap(i, j int) { // 覆寫 Swap() 方法
	p.pump[i], p.pump[j] = p.pump[j], p.pump[i]
}

func (p PumpWrapper) Less(i, j int) bool { // 覆寫 Less() 方法
	return p.by(p.pump[i], p.pump[j])
}

// 封裝成 SortPump 方法
func SortPump(pump []*Pump, by SortPumpBy) {
	sort.Sort(PumpWrapper{pump, by})
}
