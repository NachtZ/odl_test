package main

type Metadata struct {
	Metadata int64 `json:"metadata,omitempty"`
	Mask     int64 `json:"metadata-mask,omitempty"`
}
type Tunnel struct {
	ID   int64 `json:"tunnel-id,omitempty"`
	Mask int64 `json:"tunnel-mask,omitempty"`
}
type EthernetAddr struct {
	Address string `json:"address,omitempty"`
	Mask    string `json:"mask,omitempty"`
}
type EthernetType struct {
	Type uint32 `json:"type,omitempty"`
}
type EthernetMatch struct {
	EthernetSource      *EthernetAddr `json:"ethernet-source,omitempty"`
	EthernetDestination *EthernetAddr `json:"ethernet-destination,omitempty"`
	EthernetType        *EthernetType `json:"ethernet-type,omitempty"`
}
type VlanMatch struct {
	VlanID struct {
		VlanIDPresent bool `json:"vlan-id-present,omitempty"`
		VlanID        int  `json:"vlan-id,omitempty"`
	} `json:"vlan-id,omitempty"`
	VlanPcp int `json:"vlan-pcp,omitempty"`
}
type IpMatch struct {
	IpProtocol uint8 `json:"ip-protocol,omitempty"`
	IpDscp     uint8 `json:"ip-dscp,omitempty"`
	IpEcn      uint8 `json:"ip-ecn,omitempty"`
	IpProto    uint8 `json:"ip-proto,omitempty"`
}
type IcmpV4Match struct {
	Type uint8 `json:"icmpv4-type,omitempty"`
	Code uint8 `json:"icmpv4-code,omitempty"`
}
type IcmpV6Match struct {
	Type uint8 `json:"icmpv6-type,omitempty"`
	Code uint8 `json:"icmpv6-code,omitempty"`
}
type ProtocolMatchFields struct {
	MplsLabel uint32 `json:"mpls-label,omitempty"`
	MplsTc    uint8  `json:"mpls-tc,omitempty"`
	MplsBos   uint8  `json:"mpls-bos,omitempty"`
	Pbb       struct {
		ISID uint32 `json:"pbb-isid,omitempty"`
		Mask uint32 `json:"pbb-mask,omitempty"`
	} `json:"pdd,omitempty"`
}
type TcpFlagsMatch struct {
	TcpFlags     uint16 `json:"tcp-flags,omitempty"`
	TcpFlagsMask uint16 `json:"tcp-flags-mask,omitempty"`
}

type Layer3Match struct {
}

//unused.
type FlowEntryMatchFiled struct {
	InPort        string         `json:"in-port,omitempty"`
	InPhyPort     string         `json:"in-phy-port,omitempty"`
	Metadata      *Metadata      `json:"metadata,omitempty"`
	Tunnel        *Tunnel        `json:"tunnel,omitempty"`
	EthernetMatch *EthernetMatch `json:"ethernet-match,omitempty"`
	VlanMatch     *VlanMatch     `json:"vlan-match,omitempty"`
	IpMatch       *IpMatch       `json:"ip-match,omitempty"`
	Ipv4Src       string         `json:"ipv4-source,omitempty"`
	Ipv4Dest      string         `json:"ipv4-destination,omitempty"`
	//	Layer3Match         string               `json:"layer-3-match,omitempty"`
	//	Layer4Match         string               `json:"layer-4-match,omitempty"`
	IcmpV4Match         *IcmpV4Match         `json:"icmpv4-match,omitempty"`
	IcmpV6Match         *IcmpV6Match         `json:"icmpv6-match,omitempty"`
	ProtocolMatchFields *ProtocolMatchFields `json:"-"` //`json:"protocol-match-fields,omitempty"`
	TcpFlagsMatch       *TcpFlagsMatch       `json:"tcp-flag-match,omitempty"`
	ExtensionLists      []struct {
		ExtensionKey string `json:"extension-key,omitempty"`
		Extension    struct {
			DosEkis string `json:"dos-ekis ,omitempty"`
		} `json:"extension,omitempty"`
	} `json:"extensionLists,omitempty"`
}

type FLowEntryApplyAction struct { //Output Action
	Order       int32 `json:"order"`
	OutputAcion struct {
		OutputNodeConnector string `json:"output-node-connector,omitempty"`
		MaxLength           uint16 `json:"max-length,omitempty"`
	} `json:"output-action,omitempty"`
}

//unused
type FlowEntryInstruction struct { //Only support Output Action
	Order       int32 `json:"order"`
	ApplyAction struct {
		Actions []FLowEntryApplyAction `json:"action,omitempty"`
	} `json:"apply-actions,omitempty"`
}

//unused
type FlowEntry struct {
	ID           string              `json:"id"`
	FlowName     string              `json:"flow-name,omitempty"`
	Priority     int                 `json:"priority,omitempty"`
	TableID      int64               `json:"table_id"`
	Match        FlowEntryMatchFiled `json:"match,omitempty"`
	Instructions struct {
		List []FlowEntryInstruction `json:"instruction,omitempty"`
	} `json:"instructions,omitempty"`
	Cookie int64 `json:"cookie,omitempty"`
}

type FlowNodeInventoryFlow struct {
	Flow []FlowEntry `json:"flow,omitempty"`
}

type jsonFlowNodeInventoryFlow struct {
	Flow []FlowEntry `json:"flow-node-inventory:flow,omitempty"`
}

type FlowConfig struct {
	Name       string //base info of a flow config
	Cookie     int64
	Node       string
	ID         string
	Priority   int
	TableId    int64
	Outputnode string // output destination of this flow config
	EtherType  int    //Ether config
	EthDst     string
	EthSrc     string
	IpType     uint8 // default v4. means if IpType !=6, this Type is v4
	IpConfig   struct {
		Dst      string
		Src      string
		Protocol uint8 // use Protocol to judge What protocal match.
	}
	TcpConfig struct { //every element can be zero here.
		SrcPort uint32
		DstPort uint32
	}
	UdpConfig struct { //every element can be zero here.
		SrcPort uint32
		DstPort uint32
	}
	IcmpConfig struct { //every element can be zero here. if IpType is 6, it's ICMPv6, other is ICMPv4.
		Type uint8
		Code uint8
	}
} //configure of a flow entry, for construct a FlowEntry, just for construct tcp, udp, icmp, or others.
