package main

//unused.
type FlowEntryMatchFiled struct {
	InPort    string `json:"in-port,omitempty"`
	InPhyPort string `json:"in-phy-port,omitempty"`
	Metadata  struct {
		Metadata int64 `json:"metadata,omitempty"`
		Mask     int64 `json:"metadata-mask,omitempty"`
	} `json:"metadata,omitempty"`
	Tunnel struct {
		ID   int64 `json:"tunnel-id,omitempty"`
		Mask int64 `json:"tunnel-mask,omitempty"`
	} `json:"tunnel,omitempty"`
	EthernetMatch struct {
		EthernetSource struct {
			Address string `json:"address,omitempty"`
			Mask    string `json:"mask,omitempty"`
		} `json:"ethernet-source,omitempty"`
		EthernetDestination struct {
			Address string `json:"address,omitempty"`
			Mask    string `json:"mask,omitempty"`
		} `json:"ethernet-destination,omitempty"`
		EthernetType struct {
			Type uint32 `json:"type,omitempty"`
		} `json:"ethernet-type,omitempty"`
	} `json:"ethernet-match,omitempty"`
	VlanMatch struct {
		VlanID struct {
			VlanIDPresent bool `json:"vlan-id-present,omitempty"`
			VlanID        int  `json:"vlan-id,omitempty"`
		} `json:"vlan-id,omitempty"`
		VlanPcp int `json:"vlan-pcp,omitempty"`
	} `json:"vlan-match,omitempty"`
	IpMatch struct {
		IpProtocol uint8 `json:"ip-protocol,omitempty"`
		IpDscp     uint8 `json:"ip-dscp,omitempty"`
		IpEcn      uint8 `json:"ip-ecn,omitempty"`
		IpProto    uint8 `json:"ip-proto,omitempty"`
	} `json:"ip-match,omitempty"`
	Layer3Match string `json:"layer-3-match,omitempty"`
	Layer4Match string `json:"layer-4-match,omitempty"`
	IcmpV4Match struct {
		Type uint8 `json:"icmpv4-type,omitempty"`
		Code uint8 `json:"icmpv4-code,omitempty"`
	} `json:"icmpv4-match,omitempty"`
	IcmpV6Match struct {
		Type uint8 `json:"icmpv6-type,omitempty"`
		Code uint8 `json:"icmpv6-code,omitempty"`
	} `json:"icmpv6-match,omitempty"`
	ProtocolMatchFields struct {
		MplsLabel uint32 `json:"mpls-label,omitempty"`
		MplsTc    uint8  `json:"mpls-tc,omitempty"`
		MplsBos   uint8  `json:"mpls-bos,omitempty"`
		Pbb       struct {
			ISID uint32 `json:"pbb-isid,omitempty"`
			Mask uint32 `json:"pbb-mask,omitempty"`
		} `json:"pdd,omitempty"`
	} `json:"-"` //`json:"protocol-match-fields,omitempty"`
	TcpFlagsMatch struct {
		TcpFlags     uint16 `json:"tcp-flags,omitempty"`
		TcpFlagsMask uint16 `json:"tcp-flags-mask,omitempty"`
	} `json:"tcp-flag-match,omitempty"`
	ExtensionLists []struct {
		ExtensionKey string `json:"extension-key,omitempty"`
		Extension    struct {
			DosEkis string `json:"dos-ekis ,omitempty"`
		} `json:"extension,omitempty"`
	} `json:"extensionLists,omitempty"`
}

type FLowEntryApplyAction struct { //Output Action
	Order       int32 `json:"order,omitempty"`
	OutputAcion struct {
		OutputNodeConnector string `json:"output-node-connector,omitempty"`
		MaxLength           uint16 `json:"max-length,omitempty"`
	} `json:"output-action,omitempty"`
}

//unused
type FlowEntryInstruction struct { //Only support Output Action
	Order       int32 `json:"order,omitempty"`
	ApplyAction struct {
		Actions []FLowEntryApplyAction `json:"action,omitempty"`
	} `json:"apply-actions,omitempty"`
}

//unused
type FlowEntry struct {
	ID           string               `json:"id"`
	FlowName     string               `json:"flow-name,omitempty"`
	Priority     int                  `json:"priority,omitempty"`
	TableID      int64                `json:"table_id"`
	Match        *FlowEntryMatchFiled `json:"match,omitempty"`
	Instructions *struct {
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
