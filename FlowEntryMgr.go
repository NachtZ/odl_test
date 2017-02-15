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
	//log.Printf("%+v\n", *t.Flow[0].Match)
	//log.Printf("%+v\n", *t.Flow[0].Instructions)

	return string(contents)
}

func (base ODLBasic) PutFlowEntry(nodeid string, flow FlowEntry) {
	url := base.BaseUrl + "/restconf/config/opendaylight-inventory:nodes/node/" + nodeid + "/table/" + strconv.Itoa(int(flow.TableID)) + "/flow/" + flow.ID //here may have bug when convent int64 to int
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
		return
	}
	contents, err := ioutil.ReadAll(resp.Body)
	log.Println(string(contents))
}

func (base ODLBasic) SentFlowConfig(nodeid string, flowconfig FlowConfig) {
	var flow FlowEntry
	flow.Cookie = flowconfig.Cookie
	flow.FlowName = flowconfig.Name
	flow.ID = flowconfig.ID
	flow.Priority = flowconfig.Priority
	flow.TableID = flowconfig.TableId
	flow.Instructions.List = make([]FlowEntryInstruction, 1)
	flow.Instructions.List[0].ApplyAction.Actions = make([]FLowEntryApplyAction, 1)
	flow.Instructions.List[0].ApplyAction.Actions[0].Order = 0 //default order
	if flowconfig.Outputnode != "" {
		flow.Instructions.List[0].ApplyAction.Actions[0].OutputAction = &OutputAction{
			OutputNodeConnector: flowconfig.Outputnode,
		}
	}
	if flowconfig.NwDstActionIP != "" {
		if flowconfig.IpType == 6 {
			flow.Instructions.List[0].ApplyAction.Actions[0].SetNwDstAction = &SetNwDstAction{
				Ipv6Address: flowconfig.NwDstActionIP,
			}
		} else {
			flow.Instructions.List[0].ApplyAction.Actions[0].SetNwDstAction = &SetNwDstAction{
				Ipv4Address: flowconfig.NwDstActionIP,
			}
		}
	}
	//flow.Instructions.List[0].ApplyAction.Actions[0].OutputAcion.MaxLength =
	if flowconfig.EtherType != 0 || flowconfig.EthDst != "" || flowconfig.EthSrc != "" { //ether config
		var t EthernetMatch
		flow.Match.EthernetMatch = &t
		if flowconfig.EtherType != 0 {
			t.EthernetType = &EthernetType{
				uint32(flowconfig.EtherType),
			}
		}
		if flowconfig.EthDst != "" {
			t.EthernetDestination = &EthernetAddr{
				flowconfig.EthDst,
				"", //Should add EthMask support.
			}
		}
		if flowconfig.EthSrc != "" {
			t.EthernetSource = &EthernetAddr{
				flowconfig.EthSrc,
				"", //Should add EthMask support.
			}
		}
	}
	if flowconfig.IpConfig.Dst != "" || flowconfig.IpConfig.Src != "" || flowconfig.IpConfig.Protocol != 0 { //here just support v4 config
		if flowconfig.IpConfig.Dst != "" {
			flow.Match.Ipv4Dest = flowconfig.IpConfig.Dst
		}
		if flowconfig.IpConfig.Src != "" {
			flow.Match.Ipv4Src = flowconfig.IpConfig.Src
		}
		if flowconfig.IpConfig.Protocol != 0 {
			flow.Match.IpMatch = &IpMatch{
				flowconfig.IpConfig.Protocol,
				0,
				0,
				flowconfig.IpType,
			}
		}
	}
	if flowconfig.TcpConfig.SrcPort != 0 || flowconfig.TcpConfig.DstPort != 0 { // now ignort layer4Match
	}
	if flowconfig.IcmpConfig.Code != 0 || flowconfig.IcmpConfig.Type != 0 {
		if flowconfig.IpType == 6 {
			flow.Match.IcmpV6Match = &IcmpV6Match{
				flowconfig.IcmpConfig.Type,
				flowconfig.IcmpConfig.Code,
			}
		} else {
			flow.Match.IcmpV4Match = &IcmpV4Match{
				flowconfig.IcmpConfig.Type,
				flowconfig.IcmpConfig.Code,
			}
		}
	}
	base.PutFlowEntry(nodeid, flow)
}

func (base ODLBasic) DeleteFlowEntry(nodeid string, tableid int, flowid string) {
	url := base.BaseUrl + "/restconf/config/opendaylight-inventory:nodes/node/" + nodeid + "/table/" + strconv.Itoa(tableid) + "/flow/" + flowid
	log.Println(url)
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", url, nil)
	req.SetBasicAuth(base.User, base.Password)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error")
	}
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(string(contents))
}

func (base ODLBasic) TransferFlow(nodeid string, from, to string, rec *Recorder) { // set a flowentry in Node "nodeid", aim to tranfer flow from "from" to "to"
	flowid := "1" //todo : need to add some
	url := base.BaseUrl + "/restconf/config/opendaylight-inventory:nodes/node/" + nodeid + "/table/" + strconv.Itoa(0) + "/flow/" + flowid
	log.Println(url)
	_, ok1 := (*rec).RecordMap[from]
	_, ok2 := (*rec).RecordMap[to]
	if !ok1 || !ok2 {
		log.Println("Node ", from, "exsit:", ok1)
		log.Println("Node ", to, "exsit:", ok2)
		return
	}
	if len((*rec).RecordMap[from].IP) == 0 || len((*rec).RecordMap[to].IP) == 0 {
		log.Println("IP list of from or to is nil", "from", len((*rec).RecordMap[from].IP), "to", len((*rec).RecordMap[to].IP))
		return
	}
	flowconfig := FlowConfig{
		Name:          "Tran" + from + "to" + to + "in" + nodeid,
		Node:          nodeid,
		ID:            flowid,
		Priority:      33000,
		EtherType:     2048,
		TableId:       0,                          //todo
		NwDstActionIP: (*rec).RecordMap[to].IP[0], //todo: I think one nc can only have one ip address.
	}
	flowconfig.IpConfig.Dst = (*rec).RecordMap[from].IP[0] //todo: I think one nc can only have one ip address.
	base.SentFlowConfig(nodeid, flowconfig)
}

func testDeleteFlowEntry() {
	base := ODLBasic{
		"http://10.108.20.110:8181",
		"admin",
		"admin",
	}
	base.DeleteFlowEntry("openflow:1", 0, "1")
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
	flow.TableID = 0
	flow.Priority = 2
	flow.Match.EthernetMatch = &EthernetMatch{
		nil,
		nil,
		&EthernetType{
			2048,
		},
	}
	flow.Instructions.List = make([]FlowEntryInstruction, 1)
	flow.Instructions.List[0].ApplyAction.Actions = make([]FLowEntryApplyAction, 1)
	flow.Instructions.List[0].ApplyAction.Actions[0].Order = 0
	flow.Instructions.List[0].ApplyAction.Actions[0].OutputAction = &OutputAction{
		"1",
		60,
	}
	base.PutFlowEntry("openflow:1", flow)
}

func testSentFlowConfig() {
	fc := FlowConfig{
		Name:       "testSentFlowConfig",
		Node:       "openflow:1",
		ID:         "2",
		Priority:   33000,
		TableId:    0,
		Outputnode: "openflow:1:2",
		EtherType:  2048,
	}
	fc.IpConfig.Protocol = 6
	base := ODLBasic{
		"http://10.108.20.110:8181",
		"admin",
		"admin",
	}
	base.SentFlowConfig("openflow:1", fc)
}
func testFlowEntryMgr() {
	base := ODLBasic{
		"http://10.108.20.110:8181",
		"admin",
		"admin",
	}
	base.GetFlowEntry("openflow:1", 0, "1")
}

func testTransferFlow() {
	base := ODLBasic{
		"http://10.108.20.110:8181",
		"admin",
		"admin",
	}
	from := &BaseRecord{
		IP: []string{"10.0.0.2/32"},
	}
	to := &BaseRecord{
		IP: []string{"10.0.0.3/32"},
	}
	var rec Recorder
	rec.RecordMap = make(map[string]*BaseRecord)
	rec.RecordMap["from"] = from
	rec.RecordMap["to"] = to
	base.TransferFlow("openflow:2", "from", "to", &rec)
}
