personal backup, for ddos protect in opendaylight boron.   
思路：  
>1. 获取ODL网络中一段时间（比如一周，正常运行的网络特征，目前是只有bps，pps和两个速度的加速度。
>2. 在保护中，对这四个特征中有明显异常的，比如加速度过高表示有流量骤增的情况，将该节点列为可疑节点。
>3. 对可疑节点，利用flowtable将流量导入到流量判定网络，初步采用snort作为判定工具。
>4. 对流量进行判断，如果判定正常流量，还原流表，不做操作（或者可以讲其列为高负载节点，降低判断标准或者进行负载均衡。）
>5. 对攻击流量进行其他处理。尝试流量清洗。

程序中假设ODL中，protect network 和 防护系统都在一个odl控制器下。   
而且防护系统在ODL中的node connector ID都应该在配置文件中告知。  

防护网络中，3中的nc id只是nc id， 即认为识别工具可以识别所有的ddos攻击类型。或者在后期加入初步的分类，比如对不同协议类型的ddos进行识别。比如tcp, udp, icmp and the others.  
在5中，流量处理网络可以针对不同的ddos攻击进行不同的处理，这种nc节点应该在配置文件中有相关标识。  

在3中，如何知道判定结果也是一个问题。初步想法是将判定结果分成两类。然后发到不同的网卡上去。也就是说，一个snort有三个网卡。一个in网卡，两个out网卡，分别表示正常流量和攻击流量。   
或者利用流标签法对流量进行标记。  
获取

## 基于流表的流量统计，output action可以试试看是不是output normal来实现。 ##

可能用到的技术：
包标签法，snort， VTN，docker

完成情况：

思路 | 完成度 | 预计时间  
--- | --- | ---
1 |  完成bps, pps, 加速度统计工作。 | 预计2017/2/1 完成，2017/2/5 已完成 2017/2/12 添加获取IP功能。  
2 | 完成简单的基于正常最大流加速度的判断功能。 | 预计2017/2/12 完成，2017/2/7 完成简略的第一版框架。  
3 | 写流表API | 预计2017/2/19 完成，2017/2/14 完成流量转移API，需要测试一下功能。  
4 | 未开工。 | 预计2017/2/24 完成， 未开工。  
5 | 未开工。 | 不定期， 未开工。  

思路来源：[defendse4all](https://wiki.opendaylight.org/view/Defense4All:Tutorial)

FlowEntryMgr 构造restful url。  
 
flow table 的cookie是什么意思。  

## 2017/2/10 ##

写一个炒鸡简化版本的flow entry sender。  



简化版flow entry:  
协议种类有tcp,udp,icmp和其他四种。其中其他这一种不知道该怎么做。  
一定有cookie选项么，不清楚。   
这个entry的目的是针对特定的flow进行筛选。所以需要做的事情是将一个node connector中的流量转移到另外一个网络中去。  
传入的参数有：dest ip, src ip, dest port ,src ip, protocol type等。  
需要根据这些参数构造一个流表 entry。  
其中源目ip固定，根据ip type变更相应的条目。  

## 2017/2/11 ##
读defense4all的源码。

## 2017/2/12 ##
[一种流量迁移的方法。](https://floodlight.atlassian.net/wiki/display/floodlightcontroller/How+to+Perform+Transparent+Packet+Redirection+with+OpenFlow+and+Floodlight)
0. 前提：判断设备是双网卡设备，in和out网卡分别挂在不同的交换机上面。判断设备同一时间只能判断一个nc上的流量。  
1. 按链接所示，将node中的目的ip为nc的ip全部都改为判断设备的ip。  
2. 在判断设备中，进行两件事情，第一，在out网卡所在的交换机，设定将目的ip为入口网卡的ip改为之前的nc，ip。  
3. 在判断完成后，删除所有添加的流表，还原交换机。  
4. 对于判断结果，可以这样，将认为是攻击流量的ip全部转换为flow entries.下发到所有node。
5. 4中描述有待商榷，是属于思路4的内容了。目前暂不考虑。

所以目前需要知道所有node connector的 ip地址。
> 通过inventory api可以拿到node ip。
拿到了。  

写一个API， transfer(src, dst NodeConnectorID, nodeid string),表示在这个node上下发一个流表，把发完src上的流全部转移到dst上去。  
```
<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<flow xmlns="urn:opendaylight:flow:inventory">
    <priority>33000</priority>
    <flow-name>Foo</flow-name>
    <match>
        <ethernet-match>
            <ethernet-type>
                <type>2048</type>
            </ethernet-type>
         </ethernet-match>
         <ip-match>
            <ip-protocol>6</ip-protocol>         
        </ip-match>
       <ipv4-destination>10.0.0.8/32</ipv4-destination>
    </match>
    <id>1</id>
    <table_id>0</table_id>
    <instructions>
        <instruction>
            <order>0</order>
            <apply-actions>
                <action>
                   <order>0</order>
                    <set-nw-dst-action>
                    <ipv4-address>10.0.0.7/32</ipv4-address>
                    </set-nw-dst-action>
                    </action>
                </apply-actions>
            </instruction>
        </instructions>
    </flow>
```

## 2017/2/13 ##
完成初步的流量转移api框架。
需要将之前获得的ncid-ip对照表串联起来。

## 2017/2/14 ##
修改一下。但是感觉api设计上有点问题，可能需要修改修改。  
同时需要测试下功能。  
需要考虑下snort是否需要实现，目前想法是等寒假结束回实验室继续研究。  
## 2017/2/15 ##
功能测试通过，注意Priority需要设大。  