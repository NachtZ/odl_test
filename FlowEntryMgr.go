package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	//"fmt"
)

type ODLBasic struct {
	BaseUrl  string // ODL contorller address, like http://10.108.20.110:8181
	User     string // admin
	Password string // admin
}

//return json of the dest flow entry
//light-inventory:nodes/node/openflow:1/table/0/flow/1
func (base ODLBasic) GetFlowEntry(nodeid string, tableid int, flowid string) string {
	url := base.BaseUrl + "/restconf/config/opendaylight-inventory:nodes/node/" + nodeid + "/table/" + strconv.Itoa(tableid) + "/flow/" + flowid
	log.Println(url)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(base.User, base.Password)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error")
	}
	contents, err := ioutil.ReadAll(resp.Body)
	var t jsonFlowNodeInventoryFlow
	err = json.Unmarshal(contents, &t)
	log.Println(string(contents))
	if err != nil {
		log.Println(err)
		return ""
	}

	log.Printf("%+v\n", t)
	log.Printf("%+v\n", *t.Flow[0].Match)
	log.Printf("%+v\n", *t.Flow[0].Instructions)

	return string(contents)
}

func (base ODLBasic) PutFlowEntry(nodeid string, tableid int, flow FlowEntry) {
	url := base.BaseUrl + "/restconf/config/opendaylight-inventory:nodes/node/" + nodeid + "/table/" + strconv.Itoa(tableid) + "/flow/" + flow.ID
	log.Println(url)
	client := &http.Client{}
	ctx, err := json.Marshal(flow)
	str := "{\"flow-node-inventory:flow\":[" + string(ctx) + "]}"
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(str)
	req, err := http.NewRequest("PUT", url, strings.NewReader(str))
	req.SetBasicAuth(base.User, base.Password)
	req.Header.Set("Content-type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error")
	}
	contents, err := ioutil.ReadAll(resp.Body)
	log.Println(string(contents))
}

func testPutFlowEntry() {
	base := ODLBasic{
		"http://10.108.20.110:8181",
		"admin",
		"admin",
	}
	var flow FlowEntry
	flow.ID = "1"
	flow.FlowName = "test"
	flow.Match.EthernetMatch.EthernetType.Type = 2048
	flow.Instructions.List = make([]FlowEntryInstruction, 1)
	flow.Instructions.List[0].ApplyAction.Actions = make([]FLowEntryApplyAction, 1)
	flow.Instructions.List[0].ApplyAction.Actions[0].Order = 0
	flow.Instructions.List[0].ApplyAction.Actions[0].OutputAcion.OutputNodeConnector = "1"
	flow.Instructions.List[0].ApplyAction.Actions[0].OutputAcion.MaxLength = 60
	base.PutFlowEntry("openflow:1", 0, flow)
}

func testFlowEntryMgr() {
	base := ODLBasic{
		"http://10.108.20.110:8181",
		"admin",
		"admin",
	}
	base.GetFlowEntry("openflow:1", 0, "1")
}
