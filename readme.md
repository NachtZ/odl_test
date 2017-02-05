personal backup.   
思路：
>1. 获取ODL网络中一段时间（比如一周，正常运行的网络特征，目前是只有bps，pps和两个速度的加速度。
>2. 在保护中，对这四个特征中有明显异常的，比如加速度过高表示有流量骤增的情况，将该节点列为可疑节点。
>3. 对可疑节点，利用flowtable将流量导入到流量判定网络，初步采用snort作为判定工具。
>4. 对流量进行判断，如果判定正常流量，还原流表，不做操作（或者可以讲其列为高负载节点，降低判断标准或者进行负载均衡。）
>5. 对攻击流量进行其他处理。尝试流量清洗。

完成情况：

思路 | 完成度 | 预计时间
--- | --- | ---
1 |  完成bps, pps, 加速度统计工作 | 2017/2/1 完成，2017/2/5 已完成
2 | 未开工。 | 2017/2/12 完成， 未开工。
3 | 未开工。 | 2017/2/19 完成， 未开工。
4 | 未开工。 | 2017/2/24 完成， 未开工。
5 | 未开工。 | 不定期， 未开工。