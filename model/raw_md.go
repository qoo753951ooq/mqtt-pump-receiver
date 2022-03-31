package model

type RawPump struct {
	PumpStn    string `json:"id"`
	Tachometer string `json:"tach"`
	RunStatus  string `json:"r_status"`
	Datatime   string `json:"time"`
}
