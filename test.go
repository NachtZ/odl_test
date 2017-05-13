package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func GetNetworkTopology() {
	baseurl := "http://192.168.32.135:8181/restconf"

	// The URL to get the topology of the default slice
	url := strings.Join([]string{baseurl, "operational/network-topology:network-topology"}, "/")
	log.Println(url)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth("admin", "admin")
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error")
	}
	contents, err := ioutil.ReadAll(resp.Body)
	log.Println(string(contents))
	var t jsonNetworkTopology
	err = json.Unmarshal(contents, &t)
	log.Println(t)
}

func GetOpenflowNodes() []ODLInventoryNode {
	baseurl := "http://10.108.37.153:8181/restconf"
	//url := strings.Join([]string{baseurl,"operational/opendaylight-inventory:nodes"},"/")
	//url := strings.Join([]string{baseurl,"operational/opendaylight-inventory:nodes/node/openflow:1/node-connector/openflow:1:1"},"/")
	url := strings.Join([]string{baseurl, "operational/opendaylight-inventory:nodes"}, "/")
	//log.Println(url)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth("admin", "admin")
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil
	}
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil
	}
	// log.Println(string(contents))
	var t jsonODLInventoryNodes
	err = json.Unmarshal(contents, &t)
	if err != nil {
		log.Println(err)
		return nil
	}
	return t.Nodes.Nodes
}

/*
*v1, 2017/2/3, demo, do not support dynamic topo, and support topo node is less then 1024.
*
 */

func GetBaseRecord(before, now []ODLInventoryNode, beforeRecord []SingleRecord) []SingleRecord {
	totalBefore, totalNow, i := 0, 0, 0
	for _, be := range before {
		totalBefore += len(be.NodeConnectors)
	}
	for _, no := range now {
		totalNow += len(no.NodeConnectors)
	}
	if totalBefore != totalNow || (len(beforeRecord) != 0 && totalBefore != len(beforeRecord)) {
		log.Println("The network topo has changed.", totalBefore, totalNow, len(beforeRecord))
		return nil
	}
	nowRecord := make([]SingleRecord, totalNow)
	for idx, node := range now {
		if idx >= len(before) || before[idx].ID != node.ID {
			log.Println("Network Topo changed, You need to rerun the func now!")
			return nil
		}
		tmp := before[idx]
		for idx1, nc := range node.NodeConnectors {
			if idx1 >= len(tmp.NodeConnectors) || nc.ID != tmp.NodeConnectors[idx1].ID || (len(beforeRecord) > 0 && beforeRecord[i].ID != nc.ID) {
				log.Println("Network Topo changed, Please wait!", idx, len(tmp.NodeConnectors), nc.ID, tmp.NodeConnectors[idx1].ID, len(beforeRecord), beforeRecord[i].ID)
				return nil
			}
			time := nc.OPFstatics.Duration.Second - tmp.NodeConnectors[idx1].OPFstatics.Duration.Second
			if time == 0 { //impossible happend. But sometimes ODL will get error and then eturn the same value every times.
				time = 1
			}
			nowRecord[i].ID = nc.ID
			nowRecord[i].Bytes.Rx = float64((nc.OPFstatics.Bytes.Rx - tmp.NodeConnectors[idx1].OPFstatics.Bytes.Rx) / time)
			nowRecord[i].Pkts.Rx = float64((nc.OPFstatics.Pkts.Rx - tmp.NodeConnectors[idx1].OPFstatics.Pkts.Rx) / time)
			nowRecord[i].Bytes.Tx = float64((nc.OPFstatics.Bytes.Tx - tmp.NodeConnectors[idx1].OPFstatics.Bytes.Tx) / time)
			nowRecord[i].Pkts.Tx = float64((nc.OPFstatics.Pkts.Tx - tmp.NodeConnectors[idx1].OPFstatics.Pkts.Tx) / time)
			if len(beforeRecord) > 0 {
				nowRecord[i].Bytes.AccelerationRx = (nowRecord[i].Bytes.Rx - beforeRecord[i].Bytes.Rx) / float64(time)
				nowRecord[i].Pkts.AccelerationRx = (nowRecord[i].Pkts.Rx - beforeRecord[i].Pkts.Rx) / float64(time)
				nowRecord[i].Bytes.AccelerationTx = (nowRecord[i].Bytes.Tx - beforeRecord[i].Bytes.Tx) / float64(time)
				nowRecord[i].Pkts.AccelerationTx = (nowRecord[i].Pkts.Tx - beforeRecord[i].Pkts.Tx) / float64(time)
			}
			i++
		}
	}
	return nowRecord
}

func testGetBaseRecord() {
	before := GetOpenflowNodes()
	time.Sleep(5 * time.Second)
	now := GetOpenflowNodes()
	beforeRecord := GetBaseRecord(before, now, nil)
	for {
		time.Sleep(5 * time.Second)
		now = GetOpenflowNodes()
		beforeRecord = GetBaseRecord(before, now, beforeRecord)
		log.Println(beforeRecord)
		before = now
	}
}

func (base *Recorder) CheckAttack(rec []SingleRecord) []string {
	var ret []string
	for _, r := range rec {
		t1, ok := (*base).RecordMap[r.ID]
		if ok == false {
			log.Println("Failed to find NC", r.ID)
			return nil // I think all the func may need return an error.
		}
		t := t1.Average
		if r.Bytes.AccelerationRx > t.Bytes.AccelerationRx && r.Pkts.AccelerationRx > t.Pkts.AccelerationRx || r.Bytes.AccelerationTx > t.Bytes.AccelerationTx && r.Pkts.AccelerationTx > r.Pkts.AccelerationTx {
			ret = append(ret, r.ID)
		}
	}
	return ret
}

func (base *Recorder) InitRecord(getStaistic func() []ODLInventoryNode) *Recorder {
	before := getStaistic()
	totalNC := 0
	for _, nc := range before {
		totalNC += len(nc.NodeConnectors)
	}
	if totalNC >= len((*base).RawRecord) {
		log.Println("Too much NodeConnectors in this Topo! Can not init the recorder.")
		return nil
	}
	i := 0
	for _, node := range before {
		for _, nc := range node.NodeConnectors {
			(*base).RawRecord[i].ID = nc.ID
			if len((*base).RawRecord[i].Records) < 10080 {
				log.Println("Error in len(base[i].Records)")
				return nil
			}
			if len(nc.AddressList) > 0 {
				(*base).RawRecord[i].IP = make([]string, len(nc.AddressList))
				(*base).RawRecord[i].MAC = make([]string, len(nc.AddressList))
				for idx, t := range nc.AddressList {
					(*base).RawRecord[i].IP[idx] = t.IP
					(*base).RawRecord[i].MAC[idx] = t.MAC
				}
				log.Println("NC ", (*base).RawRecord[i].ID, "'s address is:")
				for idx := 0; idx < len((*base).RawRecord[i].IP); idx++ {
					log.Println((*base).RawRecord[i].IP[idx], (*base).RawRecord[i].MAC[idx])
				}

			}
			i++
		}
	}
	log.Println("The topo has", len(before), "nodes,", i, "node connectors.")
	for counter := 0; counter < 10080; counter++ { //last for one week
		time.Sleep(60 * time.Second) //wait one min.
		log.Println("\n", counter, "Now", "\n")
		now := getStaistic()
		i = 0
		timenow := time.Now()
		for idx, node := range now {
			if idx >= len(before) || before[idx].ID != node.ID {
				log.Println("Network Topo changed, You need to rerun the func now!")
				return nil
			}
			tmp := before[idx]
			log.Println("Node :", node.ID)
			log.Println("Node Info:", node.Manufacturer, node.Hardware, node.Software, node.SerialNumber)
			log.Println("NC statistic:")
			for idx1, nc := range node.NodeConnectors {
				log.Println(nc.ID, "|", nc.Name)
				if idx1 >= len(tmp.NodeConnectors) || nc.ID != tmp.NodeConnectors[idx1].ID || (*base).RawRecord[i].ID != nc.ID {
					log.Println("Network Topo changed, Please wait!")
					return nil
				}
				(*base).RawRecord[i].Records[counter].Time.Day = timenow.Day()
				(*base).RawRecord[i].Records[counter].Time.Hour = timenow.Hour()
				(*base).RawRecord[i].Records[counter].Time.Min = timenow.Minute()

				time := nc.OPFstatics.Duration.Second - tmp.NodeConnectors[idx1].OPFstatics.Duration.Second
				log.Println("time:", time, nc.OPFstatics.Duration.Second, tmp.NodeConnectors[idx1].OPFstatics.Duration.Second)
				if time == 0 { //impossible happend. But sometimes ODL will get error and then eturn the same value every times.
					time = 1
				}
				(*base).RawRecord[i].Records[counter].Bytes.Rx = float64((nc.OPFstatics.Bytes.Rx - tmp.NodeConnectors[idx1].OPFstatics.Bytes.Rx) / time)
				(*base).RawRecord[i].Records[counter].Pkts.Rx = float64((nc.OPFstatics.Pkts.Rx - tmp.NodeConnectors[idx1].OPFstatics.Pkts.Rx) / time)
				(*base).RawRecord[i].Records[counter].Bytes.Tx = float64((nc.OPFstatics.Bytes.Tx - tmp.NodeConnectors[idx1].OPFstatics.Bytes.Tx) / time)
				(*base).RawRecord[i].Records[counter].Pkts.Tx = float64((nc.OPFstatics.Pkts.Tx - tmp.NodeConnectors[idx1].OPFstatics.Pkts.Tx) / time)
				if counter > 0 {
					(*base).RawRecord[i].Records[counter].Bytes.AccelerationRx = ((*base).RawRecord[i].Records[counter].Bytes.Rx - (*base).RawRecord[i].Records[counter-1].Bytes.Rx) / float64(time)
					(*base).RawRecord[i].Records[counter].Pkts.AccelerationRx = ((*base).RawRecord[i].Records[counter].Pkts.Rx - (*base).RawRecord[i].Records[counter-1].Pkts.Rx) / float64(time)
					(*base).RawRecord[i].Records[counter].Bytes.AccelerationTx = ((*base).RawRecord[i].Records[counter].Bytes.Tx - (*base).RawRecord[i].Records[counter-1].Bytes.Tx) / float64(time)
					(*base).RawRecord[i].Records[counter].Pkts.AccelerationTx = ((*base).RawRecord[i].Records[counter].Pkts.Tx - (*base).RawRecord[i].Records[counter-1].Pkts.Tx) / float64(time)
					if (*base).RawRecord[i].Records[counter].Bytes.AccelerationRx > (*base).RawRecord[i].Average.Bytes.AccelerationRx {
						(*base).RawRecord[i].Average.Bytes.AccelerationRx = (*base).RawRecord[i].Records[counter].Bytes.AccelerationRx
					}
					if (*base).RawRecord[i].Records[counter].Bytes.AccelerationTx > (*base).RawRecord[i].Average.Bytes.AccelerationTx {
						(*base).RawRecord[i].Average.Bytes.AccelerationTx = (*base).RawRecord[i].Records[counter].Bytes.AccelerationTx
					}
					if (*base).RawRecord[i].Records[counter].Pkts.AccelerationRx > (*base).RawRecord[i].Average.Pkts.AccelerationRx {
						(*base).RawRecord[i].Average.Pkts.AccelerationRx = (*base).RawRecord[i].Records[counter].Pkts.AccelerationRx
					}
					if (*base).RawRecord[i].Records[counter].Pkts.AccelerationTx > (*base).RawRecord[i].Average.Pkts.AccelerationTx {
						(*base).RawRecord[i].Average.Pkts.AccelerationTx = (*base).RawRecord[i].Records[counter].Pkts.AccelerationTx
					}
				}
				log.Println("Rx Speed:", (*base).RawRecord[i].Records[counter].Bytes.Rx, "bps", (*base).RawRecord[i].Records[counter].Pkts.Rx, "pps")
				log.Println("Rx Acceleration:", (*base).RawRecord[i].Records[counter].Bytes.AccelerationRx, "bps", (*base).RawRecord[i].Records[counter].Pkts.AccelerationRx, "pps")
				log.Println("Tx Speed:", (*base).RawRecord[i].Records[counter].Bytes.Tx, "bps", (*base).RawRecord[i].Records[counter].Pkts.Tx, "pps")
				log.Println("Tx Acceleration:", (*base).RawRecord[i].Records[counter].Bytes.AccelerationTx, "bps", (*base).RawRecord[i].Records[counter].Pkts.AccelerationTx, "pps")
				i++
			}
		}
		before = now
	}
	(*base).RecordMap = make(map[string]*BaseRecord)
	for _, rec := range (*base).RawRecord {
		(*base).RecordMap[rec.ID] = &rec
	}
	return base
}

func testInitRecord() {
	var recorder Recorder
	recorder.RawRecord = make([]BaseRecord, 1024)
	recorder.InitRecord(GetOpenflowNodes)
}

/*
func printStatistic(before, now []ODLInventoryNode) {
	for idx, node := range now {
		if idx >= len(before) || before[idx].ID != node.ID {
			log.Println("Network Topo changed, Please wait!")
			return
		}
		tmp := before[idx]
		log.Println("Node :", node.ID)
		log.Println("Node Info:", node.Manufacturer, node.Hardware, node.Software, node.SerialNumber)
		log.Println("NC statistic:")
		for idx1, nc := range node.NodeConnectors {
			log.Println(nc.ID, "|", nc.Name)
			if idx1 >= len(tmp.NodeConnectors) || nc.ID != tmp.NodeConnectors[idx1].ID {
				log.Println("Network Topo changed, Please wait!")
				return
			}
			time := nc.OPFstatics.Duration.Second - tmp.NodeConnectors[idx1].OPFstatics.Duration.Second
			log.Println("time:", time, nc.OPFstatics.Duration.Second, tmp.NodeConnectors[idx1].OPFstatics.Duration.Second)
			if time == 0 {
				time = 1
			}
			log.Println("Rx Speed:", (nc.OPFstatics.Bytes.Rx-tmp.NodeConnectors[idx1].OPFstatics.Bytes.Rx)/time, "bps", (nc.OPFstatics.Pkts.Rx-tmp.NodeConnectors[idx1].OPFstatics.Pkts.Rx)/time, "pps")
			log.Println("Tx Speed:", (nc.OPFstatics.Bytes.Tx-tmp.NodeConnectors[idx1].OPFstatics.Bytes.Tx)/time, "bps", (nc.OPFstatics.Pkts.Tx-tmp.NodeConnectors[idx1].OPFstatics.Pkts.Tx)/time, "pps")
		}
	}
}
*/

var data [][][]int64
var name []string
var flag bool

func printStatistic(before, now []ODLInventoryNode) {
	var total int64

	ttmp := [][]int64{}
	for idx, node := range now {
		if idx >= len(before) || before[idx].ID != node.ID {
			//log.Println("Network Topo changed, Please wait!")
			return
		}

		tmp := before[idx]
		//log.Println("Node :", node.ID)
		//log.Println("Node Info:", node.Manufacturer, node.Hardware, node.Software, node.SerialNumber)
		//log.Println("NC statistic:")
		for idx1, nc := range node.NodeConnectors {
			//log.Println(nc.ID, "|", nc.Name)
			if idx1 >= len(tmp.NodeConnectors) || nc.ID != tmp.NodeConnectors[idx1].ID {
				//		log.Println("Network Topo changed, Please wait!")
				return
			}
			ttmp1 := make([]int64, 4)
			time := nc.OPFstatics.Duration.Second - tmp.NodeConnectors[idx1].OPFstatics.Duration.Second
			//log.Println("time:", time, nc.OPFstatics.Duration.Second, tmp.NodeConnectors[idx1].OPFstatics.Duration.Second)

			if time == 0 {
				time = 1
			}
			if flag == false {
				name = append(name, nc.Name)
			}
			ttmp1[0], ttmp1[1], ttmp1[2], ttmp1[3] = (nc.OPFstatics.Bytes.Rx-tmp.NodeConnectors[idx1].OPFstatics.Bytes.Rx)/time, (nc.OPFstatics.Pkts.Rx-tmp.NodeConnectors[idx1].OPFstatics.Pkts.Rx)/time, (nc.OPFstatics.Bytes.Tx-tmp.NodeConnectors[idx1].OPFstatics.Bytes.Tx)/time, (nc.OPFstatics.Pkts.Tx-tmp.NodeConnectors[idx1].OPFstatics.Pkts.Tx)/time
			ttmp = append(ttmp, ttmp1)
			// format: nodename time: rxbps rxpps txbps txpps
			fmt.Println(nc.Name, time, ":", nc.OPFstatics.Duration.Second, (nc.OPFstatics.Bytes.Rx-tmp.NodeConnectors[idx1].OPFstatics.Bytes.Rx)/time, (nc.OPFstatics.Pkts.Rx-tmp.NodeConnectors[idx1].OPFstatics.Pkts.Rx)/time, (nc.OPFstatics.Bytes.Tx-tmp.NodeConnectors[idx1].OPFstatics.Bytes.Tx)/time, (nc.OPFstatics.Pkts.Tx-tmp.NodeConnectors[idx1].OPFstatics.Pkts.Tx)/time)
			//log.Println("Rx Speed:", (nc.OPFstatics.Bytes.Rx-tmp.NodeConnectors[idx1].OPFstatics.Bytes.Rx)/time, "bps", (nc.OPFstatics.Pkts.Rx-tmp.NodeConnectors[idx1].OPFstatics.Pkts.Rx)/time, "pps")
			//log.Println("Tx Speed:", (nc.OPFstatics.Bytes.Tx-tmp.NodeConnectors[idx1].OPFstatics.Bytes.Tx)/time, "bps", (nc.OPFstatics.Pkts.Tx-tmp.NodeConnectors[idx1].OPFstatics.Pkts.Tx)/time, "pps")
		}
		for _, nc := range node.ODLInventoryTables {
			total += nc.Statistic.ActiveFlows
		}

	}
	data = append(data, ttmp)
	flag = true
	//fmt.Println("total flow is ", total)
}

func SpeedMonitor() {
	before := GetOpenflowNodes()
	for i := 0; i < 100; i++ {
		time.Sleep(5 * time.Second)
		now := GetOpenflowNodes()
		printStatistic(before, now)
		before = now
	}

}

func main() {
	SpeedMonitor()
	file, _ := os.Create("d://data" + time.Now().Format("2006_01_02_15_04_05") + ".txt")
	/*
		for i := 0; i < len(data[0]); i++ {
			file.WriteString(name[i] + "\n")
			for j := 0; j < len(data); j++ {
				file.WriteString(strconv.FormatInt(data[j][i][0], 10) + "," + strconv.FormatInt(data[j][i][1], 10) + "," + strconv.FormatInt(data[j][i][2], 10) + "," + strconv.FormatInt(data[j][i][3], 10) + "\n")
			}
			file.WriteString("\n")
		}
	*/
	for j := 0; j < len(data[0]); j++ {
		file.WriteString(name[j] + ",")
	}
	file.WriteString("\n")
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data[0]); j++ {
			file.WriteString(strconv.FormatInt(data[i][j][3], 10) + ",")
		}
		file.WriteString("\n")
	}
}
