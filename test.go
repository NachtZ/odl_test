package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func GetNetworkTopology() {
	baseurl := "http://192.168.32.135:8181/restconf"

	// The URL to get the topology of the default slice
	url := strings.Join([]string{baseurl, "operational/network-topology:network-topology"}, "/")
	fmt.Println(url)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth("admin", "admin")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error")
	}
	contents, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(contents))
	var t jsonNetworkTopology
	err = json.Unmarshal(contents, &t)
	fmt.Println(t)
}

func GetOpenflowNodes() []ODLInventoryNode {
	baseurl := "http://10.108.20.110:8181/restconf"
	//url := strings.Join([]string{baseurl,"operational/opendaylight-inventory:nodes"},"/")
	//url := strings.Join([]string{baseurl,"operational/opendaylight-inventory:nodes/node/openflow:1/node-connector/openflow:1:1"},"/")
	url := strings.Join([]string{baseurl, "operational/opendaylight-inventory:nodes"}, "/")
	fmt.Println(url)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth("admin", "admin")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	// fmt.Println(string(contents))
	var t jsonODLInventoryNodes
	err = json.Unmarshal(contents, &t)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return t.Nodes.Nodes
}

/*
*v1, 2017/2/3, demo, do not support dynamic topo, and support topo node is less then 1024.
*
 */
func (base *Recorder) InitRecord(getStaistic func() []ODLInventoryNode) []BaseRecord {
	before := getStaistic()
	totalNC := 0
	for _, nc := range before {
		totalNC += len(nc.NodeConnectors)
	}
	if totalNC >= len(*base) {
		fmt.Println("Too much NodeConnectors in this Topo! Can not init the recorder.")
		return nil
	}
	i := 0
	for _, node := range before {
		for _, nc := range node.NodeConnectors {
			(*base)[i].ID = nc.ID
			if len((*base)[i].Records) < 10080 {
				fmt.Println("Error in len(base[i].Records)")
				return nil
			}
			i++
		}
	}
	fmt.Println("The topo has", len(before), "nodes,", i, "node connectors.")
	for counter := 0; counter < 10080; counter++ { //last for one week
		time.Sleep(60 * time.Second) //wait one min.
		fmt.Println("\n", counter, "Now", "\n")
		now := getStaistic()
		i = 0
		timenow := time.Now()
		for idx, node := range now {
			if idx >= len(before) || before[idx].ID != node.ID {
				fmt.Println("Network Topo changed, You need to rerun the func now!")
				return nil
			}
			tmp := before[idx]
			fmt.Println("Node :", node.ID)
			fmt.Println("Node Info:", node.Manufacturer, node.Hardware, node.Software, node.SerialNumber)
			fmt.Println("NC statistic:")
			for idx1, nc := range node.NodeConnectors {
				fmt.Println(nc.ID, "|", nc.Name)
				if idx1 >= len(tmp.NodeConnectors) || nc.ID != tmp.NodeConnectors[idx1].ID || (*base)[i].ID != nc.ID {
					fmt.Println("Network Topo changed, Please wait!")
					return nil
				}
				(*base)[i].Records[counter].Time.Day = timenow.Day()
				(*base)[i].Records[counter].Time.Hour = timenow.Hour()
				(*base)[i].Records[counter].Time.Min = timenow.Minute()

				time := nc.OPFstatics.Duration.Second - tmp.NodeConnectors[idx1].OPFstatics.Duration.Second
				fmt.Println("time:", time, nc.OPFstatics.Duration.Second, tmp.NodeConnectors[idx1].OPFstatics.Duration.Second)
				if time == 0 { //impossible happend. But sometimes ODL will get error and return the same value every times.
					time = 1
				}
				(*base)[i].Records[counter].Bytes.Rx = float64((nc.OPFstatics.Bytes.Rx - tmp.NodeConnectors[idx1].OPFstatics.Bytes.Rx) / time)
				(*base)[i].Records[counter].Pkts.Rx = float64((nc.OPFstatics.Pkts.Rx - tmp.NodeConnectors[idx1].OPFstatics.Pkts.Rx) / time)
				(*base)[i].Records[counter].Bytes.Tx = float64((nc.OPFstatics.Bytes.Tx - tmp.NodeConnectors[idx1].OPFstatics.Bytes.Tx) / time)
				(*base)[i].Records[counter].Pkts.Tx = float64((nc.OPFstatics.Pkts.Tx - tmp.NodeConnectors[idx1].OPFstatics.Pkts.Tx) / time)
				if counter > 0 {
					(*base)[i].Records[counter].Bytes.AccelerationRx = ((*base)[i].Records[counter].Bytes.Rx - (*base)[i].Records[counter-1].Bytes.Rx) / float64(time)
					(*base)[i].Records[counter].Pkts.AccelerationRx = ((*base)[i].Records[counter].Pkts.Rx - (*base)[i].Records[counter-1].Pkts.Rx) / float64(time)
					(*base)[i].Records[counter].Bytes.AccelerationTx = ((*base)[i].Records[counter].Bytes.Tx - (*base)[i].Records[counter-1].Bytes.Tx) / float64(time)
					(*base)[i].Records[counter].Pkts.AccelerationTx = ((*base)[i].Records[counter].Pkts.Tx - (*base)[i].Records[counter-1].Pkts.Tx) / float64(time)
				}
				fmt.Println("Rx Speed:", (*base)[i].Records[counter].Bytes.Rx, "bps", (*base)[i].Records[counter].Pkts.Rx, "pps")
				fmt.Println("Rx Acceleration:", (*base)[i].Records[counter].Bytes.AccelerationRx, "bps", (*base)[i].Records[counter].Pkts.AccelerationRx, "pps")
				fmt.Println("Tx Speed:", (*base)[i].Records[counter].Bytes.Tx, "bps", (*base)[i].Records[counter].Pkts.Tx, "pps")
				fmt.Println("Tx Acceleration:", (*base)[i].Records[counter].Bytes.AccelerationTx, "bps", (*base)[i].Records[counter].Pkts.AccelerationTx, "pps")
				i++
			}
		}
		before = now
	}
	return *base
}

func testInitRecord() {
	recorder := make(Recorder, 1024)
	recorder.InitRecord(GetOpenflowNodes)
}

func printStatistic(before, now []ODLInventoryNode) {
	for idx, node := range now {
		if idx >= len(before) || before[idx].ID != node.ID {
			fmt.Println("Network Topo changed, Please wait!")
			return
		}
		tmp := before[idx]
		fmt.Println("Node :", node.ID)
		fmt.Println("Node Info:", node.Manufacturer, node.Hardware, node.Software, node.SerialNumber)
		fmt.Println("NC statistic:")
		for idx1, nc := range node.NodeConnectors {
			fmt.Println(nc.ID, "|", nc.Name)
			if idx1 >= len(tmp.NodeConnectors) || nc.ID != tmp.NodeConnectors[idx1].ID {
				fmt.Println("Network Topo changed, Please wait!")
				return
			}
			time := nc.OPFstatics.Duration.Second - tmp.NodeConnectors[idx1].OPFstatics.Duration.Second
			fmt.Println("time:", time, nc.OPFstatics.Duration.Second, tmp.NodeConnectors[idx1].OPFstatics.Duration.Second)
			if time == 0 {
				time = 1
			}
			fmt.Println("Rx Speed:", (nc.OPFstatics.Bytes.Rx-tmp.NodeConnectors[idx1].OPFstatics.Bytes.Rx)/time, "bps", (nc.OPFstatics.Pkts.Rx-tmp.NodeConnectors[idx1].OPFstatics.Pkts.Rx)/time, "pps")
			fmt.Println("Tx Speed:", (nc.OPFstatics.Bytes.Tx-tmp.NodeConnectors[idx1].OPFstatics.Bytes.Tx)/time, "bps", (nc.OPFstatics.Pkts.Tx-tmp.NodeConnectors[idx1].OPFstatics.Pkts.Tx)/time, "pps")
		}
	}
}

func SpeedMonitor() {
	before := GetOpenflowNodes()
	for {
		time.Sleep(5 * time.Second)
		now := GetOpenflowNodes()
		printStatistic(before, now)
		before = now
	}

}

func main() {
	testInitRecord()
}
