package main

type TerminationPoint struct {
	TpID string `json:"tp-id"`
	Ref  string `jsonL"opendaylight-topology-inventory:inventory-node-connector-ref"`
}

type TopoNode struct {
	Id     string             `json:"node-id"`
	Ref    string             `json:"opendaylight-topology-inventory:inventory-node-ref"`
	Points []TerminationPoint `json:"termination-point"`
}

type Source struct {
	SourceTp   string `json:"source-tp"`
	SourceNode string `json:"source-node"`
}

type Dest struct {
	DestNode string `json:"dest-node"`
	DestTP   string `json:"dest-tp"`
}

type Link struct {
	Id     string `json:"link-id"`
	Source Source `json:"source"`
	Dest   Dest   `json:"destination"`
}

type State struct {
	Linkdown bool
	Blocked  bool
	Live     bool
}

type NodeConnector struct {
	ID              string `json:"id"`
	Supported       string `json:"flow-node-inventory:supported"`
	PeerFeatures    string `json:"flow-node-inventory:peer-features"`
	AdvFeatures     string `json:"flow-node-inventory:advertised-features"`
	PortNumber      string `json:"flow-node-inventory:port-number"`
	HardwareAddress string `json:"flow-node-inventory:hardware-address"`
	CurrentSpeed    int64  `json:"flow-node-inventory:current-speed"`
	CurrnetFeature  string `json:"flow-node-inventory:current-feature"`
	MaxSpeed        int64  `json:"flow-node-inventory:maximum-speed"`
	Name            string `json:"flow-node-inventory:name"`
	State           State  `json:"flow-node-inventory:state"`
	Configure       string `json:"flow-node-inventory:configuration"`
	Status          string `json:"stp-status-aware-node-connector:status"`
	OPFstatics      struct {
		Count      int64 `json:"collision-count"`
		Drop       int64 `json:"transmit-drops"`
		RxFrameErr int64 `json:"receive-frame-error"`
		TxErr      int64 `json:"transmit-errors"`
		Bytes      struct {
			Rx int64 `json:"received"`
			Tx int64 `json:"transmitted"`
		} `json:"bytes"`
		RxCrcErr int64 `json:"receive-crc-error"`
		Duration struct {
			Second     int64 `json:"second"`
			NanoSecond int64 `json:"nanosecond"`
		} `json:"duration"`
		RxErr        int64 `json:"receive-errors"`
		RxDrop       int64 `json:"receive-drops"`
		RxoverRunErr int64 `json:"receive-over-run-error"`
		Pkts         struct {
			Rx int64 `json:"received"`
			Tx int64 `json:"transmitted"`
		} `json:"packets"`
	} `json:"opendaylight-port-statistics:flow-capable-node-connector-statistics"`
	AddressList []struct {
		ID        int    `json:"id"`
		MAC       string `json:"mac"`
		IP        string `json:"ip"`
		Firstseen uint64 `json:"first-seen"`
		Lastseen  uint64 `json:"last-seen"`
	} `json:"address-tracker:addresses"`
}

type jsonNodeConnector struct {
	NodeConnectors []NodeConnector `json:"node-connector"`
}

type Edge struct {
	// for decoding it seems to be case insensitive
	TailNodeConnector NodeConnector //`json:"tailNodeConnector"`
	HeadNodeConnector NodeConnector //`json:"headNodeConnector"`
}

type Properties struct {
	TimeStamp TimeStamp
	Name      ValueString
	State     ValueInt
	Config    ValueInt
	Bandwidth ValueInt
}

type TimeStamp struct {
	Value int
	Name  string
}

type UserLinks struct {
	UserLinks []UserLink
}

type UserLink struct {
	Status           string
	Name             string
	SrcNodeConnector string
	DstNodeConnector string
}

type ValueString struct {
	Value string
}

type ValueInt struct {
	Value int
}

type EdgeProperty struct {
	Edge       Edge       //`json:"edge"`
	Properties Properties //`json:"properties"`
}

type EdgeProperties struct {
	EdgeProperties []EdgeProperty //`json:"edgeProperties"`
}

type jsonNetworkTopology struct {
	NetworkTp NetworkTopology `json:"network-topology"`
}

type NetworkTopology struct {
	Topologies []Topology `json:"topology"`
}
type Topology struct {
	TopologyId string     `json:"topology-id"`
	Nodes      []TopoNode `json:"node"`
	Links      []Link     `json:"link"`
}

type ODLInventoryTable struct {
	ID        int64 `json:"id"`
	Statistic struct {
		PktMatched  int64 `json:"packets-matched"`
		ActiveFlows int64 `json:"active-flows"`
		PktLookedup int64 `json:"packets-looked-up"`
	} `json:"opendaylight-flow-table-statistics:flow-table-statistics"`
}

type jsonODLInventoryTable struct {
	Tables []ODLInventoryTable `json:"flow-node-inventory:table"`
}
type jsonODLInventoryNode struct {
	Nodes []ODLInventoryNode `json:"node"`
}

type jsonODLInventoryNodes struct {
	Nodes jsonODLInventoryNode `json:"nodes"`
}

type ODLInventoryNode struct {
	ID             string          `json:"id"`
	NodeConnectors []NodeConnector `json:"node-connector"`
	MeterFeatures  struct {
		MaxMeter uint64 `json:"max_meter"`
		MaxBands uint64 `json:"max_bands"`
		MaxColor uint64 `json:"max_color"`
	} `json:"opendaylight-meter-statistics:meter-features"`
	Manufacturer  string `json:"flow-node-inventory:manufacturer"`
	Hardware      string `json:"flow-node-inventory:hardware"`
	Software      string `json:"flow-node-inventory:software"`
	SerialNumber  string `json:"flow-node-inventory:serial-number"`
	Description   string `json:"flow-node-inventory:description"`
	SwitchFeature struct {
		MaxBuffers   uint64   `json:"max_buffers"`
		MaxTables    uint64   `json:"max_tables"`
		Capabilities []string `json:"capabilities"`
	} `json:"flow-node-inventory:switch-features"`
	ODLInventoryTables []ODLInventoryTable `json:"flow-node-inventory:table"`
	PortNumber         uint64              `json:"flow-node-inventory:port-number"`
	IPAddress          string              `json:"flow-node-inventory:ip-address"`
	TableFeatures      []struct {
		TableID         uint64 `json:"table-id"`
		Name            string `json:"name"`
		Config          string `json:"config"`
		MetadataMatch   uint64 `json:"metadata-match"`
		MetadataWrite   uint64 `json:"metadata-write"`
		MaxEntries      uint64 `json:"max-entries"`
		TableProperties struct {
			featureProperties struct {
				Type struct {
					SetFieldMatch []struct {
						MatchType string `json:"match-type"`
					} `json:"set-field-match"`
				} `json:"wildcard-setfield"`
				Order uint64 `json:"order"`
			} `json:"table-feature-properties"`
		} `json:"table-properties"`
	} `json:"flow-node-inventory:table-features"`
	GroupFeatures struct {
		TypesSupported        []string `json:"group-types-supported"`
		CapabilitiesSupported []string `json:"group-capabilities-supported"`
		MaxGroups             []uint64 `json:"max-groups"`
		Actions               []uint64 `json:"actions"`
	} `json:"opendaylight-group-statistics:group-features"`
}

type SingleRecord struct {
	ID    string // node ID
	Bytes struct {
		Rx             float64
		Tx             float64
		AccelerationRx float64
		AccelerationTx float64
	} //save the rate of bps
	Pkts struct {
		Rx             float64
		Tx             float64
		AccelerationRx float64
		AccelerationTx float64
	} //save the rate of pps
}

type BaseRecord struct {
	ID      string // node connector ID
	Records [10800]struct {
		Time struct {
			Day  int
			Hour int
			Min  int
		} //time for save record.
		Bytes struct {
			Rx             float64
			Tx             float64
			AccelerationRx float64
			AccelerationTx float64
		} //save the rate of bps
		Pkts struct {
			Rx             float64
			Tx             float64
			AccelerationRx float64
			AccelerationTx float64
		} //save the rate of pps
	}
	Average struct { // here in fact just record the max Acceleration not the max one.
		Bytes struct {
			AccelerationRx float64
			AccelerationTx float64
		} //save the rate of bps
		Pkts struct {
			AccelerationRx float64
			AccelerationTx float64
		} //save the rate of pps
	}
	IP  []string // ip address list for traffic transfer
	MAC []string // mac address list for traffic transfer
}

type StaticRecord interface {
	UpdateRecord(rec BaseRecord)                                 // add a record.
	AddNodes(recs []BaseRecord)                                  // add a set of record
	CheckAttack(rec []SingleRecord) []string                     //Judge the network's flow is normal or not. Return a set of NodeID or NodeConnector ID.
	InitRecord(getStatistic func() []ODLInventoryNode) *Recorder // Study the hole network's traffic statistic.
}

type Recorder struct {
	RawRecord []BaseRecord
	RecordMap map[string]*BaseRecord // RecordMap[NodeID]BaseRecordPtr

}
