put, url : `http://10.108.20.110:8181/restconf/config/opendaylight-inventory:nodes/node/openflow:1/table/0/flow/1`
```
<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<flow xmlns="urn:opendaylight:flow:inventory">
    <priority>2</priority>
    <flow-name>TCP</flow-name>
    <match>
        <ethernet-match>
            <ethernet-type>
                <type>2048</type>
            </ethernet-type>
        </ethernet-match>
        <ip-match>
            <ip-protocol>6</ip-protocol>
        </ip-match>
    </match>
    <id>1</id>
    <table_id>0</table_id>
    <instructions>
        <instruction>
            <order>0</order>
            <apply-actions>
                <action>
                   <order>0</order>
                   <dec-nw-ttl/>
                </action>
            </apply-actions>
        </instruction>
    </instructions>
</flow>
```
put, url : `http://10.108.20.110:8181/restconf/config/opendaylight-inventory:nodes/node/openflow:1/table/0/flow/2`
```
<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<flow xmlns="urn:opendaylight:flow:inventory">
    <priority>2</priority>
    <flow-name>UDP</flow-name>
    <match>
        <ethernet-match>
            <ethernet-type>
                <type>2048</type>
            </ethernet-type>
        </ethernet-match>
        <ip-match>
            <ip-protocol>17</ip-protocol>
        </ip-match>
    </match>
    <id>2</id>
    <table_id>0</table_id>
    <instructions>
        <instruction>
            <order>0</order>
            <apply-actions>
                <action>
                   <order>0</order>
                   <dec-nw-ttl/>
                </action>
            </apply-actions>
        </instruction>
    </instructions>
</flow>
```

put, url : `http://10.108.20.110:8181/restconf/config/opendaylight-inventory:nodes/node/openflow:1/table/0/flow/3`
```
<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<flow xmlns="urn:opendaylight:flow:inventory">
    <priority>2</priority>
    <flow-name>ICMP</flow-name>
    <match>
        <ethernet-match>
            <ethernet-type>
                <type>2048</type>
            </ethernet-type>
        </ethernet-match>
        <ip-match>
            <ip-protocol>1</ip-protocol>
        </ip-match>
    </match>
    <id>3</id>
    <table_id>0</table_id>
    <instructions>
        <instruction>
            <order>0</order>
            <apply-actions>
                <action>
                   <order>0</order>
                   <dec-nw-ttl/>
                </action>
            </apply-actions>
        </instruction>
    </instructions>
</flow>
```

put, url : `http://10.108.20.110:8181/restconf/config/opendaylight-inventory:nodes/node/openflow:1/table/0/flow/4`
```
<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<flow xmlns="urn:opendaylight:flow:inventory">
    <priority>2</priority>
    <flow-name>Other</flow-name>
    <match>
        <ethernet-match>
            <ethernet-type>
                <type>2048</type>
            </ethernet-type>
        </ethernet-match>
    </match>
    <id>4</id>
    <table_id>0</table_id>
    <instructions>
        <instruction>
            <order>0</order>
            <apply-actions>
                <action>
                   <order>0</order>
                   <dec-nw-ttl/>
                </action>
            </apply-actions>
        </instruction>
    </instructions>
</flow>
```


