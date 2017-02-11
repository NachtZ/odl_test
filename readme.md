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

可能用到的技术：
包标签法，snort， VTN，docker

完成情况：

思路 | 完成度 | 预计时间
--- | --- | ---
1 |  完成bps, pps, 加速度统计工作。 | 2017/2/1 完成，2017/2/5 已完成
2 | 完成简单的基于正常最大流加速度的判断功能。 | 2017/2/12 完成，2017/2/7 完成简略的第一版框架。
3 | 写流表API | 2017/2/19 完成， 完成简单的添加流表功能 for json。
4 | 未开工。 | 2017/2/24 完成， 未开工。
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
