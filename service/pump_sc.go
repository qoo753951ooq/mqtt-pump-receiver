package service

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"mqtt-pump-receiver/model"
	"mqtt-pump-receiver/model/api_md"
	"mqtt-pump-receiver/util"

	"strings"
)

func GetPumpStMap(org_id, display, etype string) map[string][]*api_md.Pump {

	m := make(map[string][]*api_md.Pump, 0)

	api := api_md.PumpStation{Org_id: org_id, Display: display, Etype: etype}
	pumpSTs := api.GetPumpStation()

	for _, st := range pumpSTs {
		api := api_md.Pump{Pump_st_ids: st.Pump_st_id, Display: display}
		m[st.Pump_st_id] = api.GetPumpList()
	}

	return m
}

func GetSendPumpDatas(rawData string, basePStnMap map[string][]*api_md.Pump) []*api_md.Pump {

	rawPump := getRawPump(rawData)

	if rawPump.PumpStn == "" {
		return make([]*api_md.Pump, 0)
	}

	return getSendPumpDatas(rawPump, basePStnMap[rawPump.PumpStn])
}

func PostPumpDataToCOVM(sendDatas []*api_md.Pump) string {

	pump := api_md.Pump{
		RequestBody: util.ConverToJsonString(sendDatas),
	}

	postResult := pump.PostPumpValue()
	fmt.Printf("%s 更新 %d臺 抽水機 %s \n", util.GetTimeNow(), len(sendDatas), postResult)
	return postResult
}

//解析原始資料至struct
func getRawPump(rawData string) model.RawPump {

	var rawPump model.RawPump

	body := []byte(rawData)

	if err := json.Unmarshal(body, &rawPump); err != nil {
		return rawPump
	}

	return rawPump
}

//取得待傳送資料
func getSendPumpDatas(rawPump model.RawPump, basePumps []*api_md.Pump) []*api_md.Pump {

	var sendResult []*api_md.Pump

	dataType := util.SubString(rawPump.PumpStn, 2, 2)

	//fmt.Printf("dataType: %s \n", dataType)

	switch dataType {
	case "PS", "PV":
		sendResult = getPumpOnlyTachDatas(rawPump, basePumps)
	case "PU":
		sendResult = getPumpOnlyRunStatusDatas(rawPump, basePumps)
	}

	return sendResult
}

//取得只有轉速表資料的抽水機
func getPumpOnlyTachDatas(rawPump model.RawPump, basePumps []*api_md.Pump) []*api_md.Pump {

	sendDatas := make([]*api_md.Pump, 0)

	if rawPump.Tachometer == "" {
		return sendDatas
	}

	//rawPump.Tachometer = "00F80029000000A800F8" //test

	//fmt.Printf("len(rawPump.Tachometer): %d \n", len(rawPump.Tachometer))
	//fmt.Printf("len(basePumps): %d \n", len(basePumps)*4)

	if len(rawPump.Tachometer) != len(basePumps)*4 {
		return sendDatas
	}

	for index, base := range basePumps {

		var send api_md.Pump

		send.St_no = base.St_no

		if tach := getTachometerFormat(base.St_no, rawPump.Tachometer, index); tach != nil {
			send.Tachometer = tach
		}

		send.Datatime = strings.ReplaceAll(rawPump.Datatime, model.Slash, model.Hyphen)
		send.Trust = true

		sendDatas = append(sendDatas, &send)
	}

	return sendDatas
}

//取得只有運轉狀態資料的抽水機
func getPumpOnlyRunStatusDatas(rawPump model.RawPump, basePumps []*api_md.Pump) []*api_md.Pump {
	sendDatas := make([]*api_md.Pump, 0)

	if rawPump.RunStatus == "" {
		return sendDatas
	}

	//rawPump.RunStatus = "01" //test

	binaryStr := util.HexToBinaryStr(rawPump.RunStatus, len(basePumps))

	if binaryStr == "" {
		return sendDatas
	}

	rawPump.RunStatus = binaryStr

	for index, base := range basePumps {

		var send api_md.Pump

		send.St_no = base.St_no

		if rStatus := getRunStatusFormat(rawPump.RunStatus, index); rStatus != "" {
			send.Run_status = rStatus
		}

		send.Datatime = strings.ReplaceAll(rawPump.Datatime, model.Slash, model.Hyphen)
		send.Trust = true

		sendDatas = append(sendDatas, &send)
	}

	return sendDatas
}

//取得格式化後之轉速
func getTachometerFormat(st_no, tachStr string, index int) *float64 {

	s := util.SubString(tachStr, index*4, 4)

	//fmt.Printf("16進制: %s \n", s)

	b, err := hex.DecodeString(s)

	if err != nil {
		return nil
	}

	num := binary.BigEndian.Uint16(b[:2])
	result := float64(num)

	//fmt.Printf("10進制: %v \n", result)

	return &result
}

//取得格式化後之啟閉
func getRunStatusFormat(rStatusStr string, index int) string {

	var result string

	s := util.SubString(rStatusStr, index*1, 1)

	switch s {
	case "1":
		result = model.Pump_run_status_off
	case "0":
		result = model.Pump_run_status_on
	}

	return result
}
