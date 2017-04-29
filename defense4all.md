
Defense4All：教程



内容
 [ 隐藏 ] 
1 简介
2 Defense4All设计
3 部署替代
4 Defense4All在ODL环境中
5 框架视图
6 应用视图
7 ODL代表视图
8 基本控制流程
9 配置和设置流程
10 攻击检测流程
11 攻击减轻流程
12 问题和故障排除
13 连续性
14 维护和升级
15 新术语和概念
16 其他信息
16.1 安全和隐私
16.2 兼容性
16.3 性能和可伸缩性信息
16.4 参考资料

## 介绍 ##
将可疑流量从正常网络路径转移到专用攻击缓解设备来进行清洗和威胁检测是一个普遍的DoS攻击环节和威胁检测策略。这些基础架构也被称为安全中心或者清洗中心，主要由第3层到第7层DoS攻击缓解设备组成。这些清洗中心可以以路径外（OOP）方式（不内嵌在本地业务流）部署在网络内的专用远程站点，因此将流量转移到这些中心是必要的。在流量清洗过程期间，攻击缓解基础设施识别并丢弃恶意IP分组，并将合法IP分组转发回其原始目标网络目的地。清洗中心可以位于企业网络，数据中心或云中，也可以作为运营商基础设施的一部分。一般来说，OOP系统中的DDoS防护包括以下主要元素：
* 在正常时期收集被保护网络的流量统计和对流量统计特征的学习。 受保护网络的正常流量基线是由这些收集的统计信息构建的。
* 认为偏离正常流量基线的流量异常可能是DoS攻击。
* 将可疑流量从其正常路径转移到缓解（清洗）中心，用于流量清洗，选择性源阻塞等。从清洗中心出来的干净流量被重新发往到原始目的地中。

## Defense4All设计 ##
Defense4All是一种SDN安全应用程序，用于检测和缓解在不同SDN拓扑中的DoS和DDoS攻击。它实现了基于可编程流的SDN环境的OOP模式下的DoS防护。管理员可以配置Defense4All以保护某些网络和服务器，程序中将之称为受保护网络或受保护对象（POs​​）。Defense4All利用SDN功能来统计指定的流量，并在每个网络位置，为每个配置中的PO，安装针对不同协议的计数流。然后Defense4All监控所有配置的PO的流量，汇总来自所有相关网络位置的读取操作，速率和平均值。如果在特定PO的协议（例如TCP，UDP，ICMP或其余流量）中检测到与正常学习的流量行为的偏差，则Defense4All在目标PO中宣布发现基于该种协议的攻击。Defense4All在安装计数流后至少需要一个星期的学习时间，在这期间Defense4All不检测攻击。为了缓解检测到的攻击，Defense4All执行以下过程：

1. 它确认DefensePro设备是否处于正常状态，并选择和正常工作的设备进行连接（如果DefensePro没有启动或没有心跳活动连接，则不会执行流量转移）。有关详细信息，请参阅“连续性”部分。
2. Defense4All使用安全策略和攻击流量的正常速率来配置DefensePro。后者能提升DefensePro对攻击的缓解。
3. Defense4All开始从DefensePro中监测和记录系统日志。只要它继续从DefensePro接收关于此攻击的攻击通知，Defense4All就会继续转移流量来继续进行攻击缓解，即使Vexternal FlowFilter计数器不提示流量中包含任何攻击。
4. Defense4All通过创建一对Vexternals并将它们映射到与DefensePro连接的所选物理PFS端口对，将所选的物理DefensePro连接映射到相关的VTN。并且会自动学习和保留VLAN标记（如果存在）。如果Defense4All已经在VTN中创建并映射了具有相同VLAN的一对Vexternals，则相同的对也被重新用于新流量的转移（而不是为相同的VTN和VLAN创建新的Vexternals）。
5. Defense4All在每个北向Vexternal中安装较高优先级的流过滤条目，通过该过滤条目，将攻击流量转移重定向到“北向DP-In Vexternal”。它还选择连接到所有Vexternals的Vbr的活北接口之一（可以有一个具有相同VLAN的Vbr）。Defense4All将来自“DP-Out Vexternal”的流量重新注入到选定的Vbr接口。（** 这句话看不懂 **）当Defense4All确定攻击已结束（没有来自PFC FlowFilter计数器或来自DefensePro的指示）时，它会恢复先前的流量：它停止监视有关流量的DefensePro syslog，它会删除流量转移流表项，删除“DP-In和DP-Out Vexternals“（如果这是此VTN和VLAN中的最后一次攻击），并从DefensePro中删除安全配置。Defense4All然后返回到正常监控。

在此版本中，Defense4All作为单个实例（非集群）运行，但它集成了以下主要容错功能：

* 它作为Linux服务运行，如果失败，它会自动重新启动。
* 有状态服务，会将最新的运行状态保存在可靠存储中，并能恢复运行。
* 它带有一个有重启和重置功能的健康跟踪器，以解决某些逻辑和老化的错误。
Defense4All监控DefensePro的状态，连接到DefensePro的交换机，在各种VTN中的相关Vbr(虚拟网桥)，这些Vbr的北向接口和北Vexternals。它依此调整，取消和（重新）发起攻击流量转移。下图说明了任何给定PO的可能状态。Radware的DefensePro（DP）是具体化的AMS的示例。 

![图1：Defense4All攻击减轻的工作流程](https://wiki.opendaylight.org/view/File:Pn_possible_states.jpg)

PN可能状态


## 部署替代 ##
Defense4All支持“短转移”，即AMS（攻击减轻系统）连接到边缘路由器，从而实现只用一跳进行流量重定向。另请参见图3中的PE1和DefensePro 2。“长转向”，即AMS清洗中心位于网络中的任意远程位置，将在未来的Defense4All版本中添加。参见下图中的PE2和DefensePro 1。Defense4All支持自动和手动转移模式。手动模式即在转移前需要向用户进行确认。

![Defense4All部署流量重定向替代](https://wiki.opendaylight.org/images/e/e5/Redirection_alternatives.jpg)

## 基于ODL的Defense4All ##
Defense4All是一种用于检测和缓解DDoS攻击的SDN应用程序。应用程序通过ODC北向REST API与OpenDaylight控制器通信。通过这个API，Defense4All执行两个主要任务：

* 监控受保护流量的行为 - 应用程序在所选网络位置中设置流条目以读取每个PN的流量统计（从多个位置聚集针对给定PN收集的统计）。
* 将攻击流量转移到选定的AMS - 应用程序在选定的网络位置中设置流条目，以将流量转移到所选的AMS。当攻击完成时，应用程序将删除这些流条目，从而返回到正常操作和流量监控。
Defense4All可以与设定好的AMS通信 - 例如，动态配置它们，监视它们，或者收集和操作来自AMS的攻击统计。AMS的API不是标准化的，在任何情况下都超出了OpenDaylight工作设计的范围。Defense4All包含一个参考实现可插拔驱动程序，与Radware的DefensePro AMS通信。应用程序提供其北向REST和CLI API，以允许其管理器：

* 控制和配置应用程序（运行时参数，ODC连接，域中的AMS，PN等等）。
* 从Defense4All和其他来源（如ODC，AMS）统一获取报告数据 - 操作或安全，当前或历史。
![Defense4在OpenDaylight环境中的所有逻辑定位](https://wiki.opendaylight.org/images/thumb/5/5b/D4A_in_odl.jpg/900px-D4A_in_odl.jpg)
Defense4All包括一个SDN应用程序框架和Defense4All应用程序本身，打包为一个单一的实体。应用程序集成到框架中是可插入的，因此任何其他SDN应用程序都可以受益于公共框架服务。这种架构的主要优点是：

* 更快的应用程序开发和更改 - 框架包含用于多个应用程序的公共代码，复杂元素（例如集群和存储库服务）只需要实现一次就能复用。
* 更快，更灵活地部署在不同的环境，形式因素，满足不同的NFR - 框架掩盖SDN应用因素，如所需的数据生存性，规模和弹性，可用性，安全性。
* 更强的鲁棒性 - 复杂的框架代码实现和测试一次，更清晰的关注分离导致更稳定的代码，框架可以主动增加鲁棒性，没有额外的代码在应用程序逻辑（如周期性应用程序循环）。
* 统一管理的共同方面 - 共同的外观和感觉。
框架视图
![从框架的角度看Defense4All结构](https://wiki.opendaylight.org/images/thumb/b/bb/Framework_view.jpg/1350px-Framework_view.jpg)
框架包含以下元素：

**FrameworkMain** - 框架根点包含对所有框架模块和全局存储库的引用，以及部署的SDN应用程序的根（在当前版本中，框架只能容纳一个应用程序）。这也是启动，停止或重置框架（连同其管辖的应用程序）的Web服务器，运行Jersey RESTful Web服务框架的Jetty Web服务器，使用Jackson解析器解析JSON编码参数。REST Web服务器为框架运行一个servlet，同时为每个已部署的应用程序（目前只有一个）运行单独一个servlet。所有REST和CLI API都通过此REST Web服务器提供服务。 

**FrameworkRestService** - 构成响应框架REST请求（获取最新的Flight Recorder记录，执行恢复出厂设置等）的框架servlet的一组类。FrameworkRestService针对FrameworkMgmtPoint调用控制和配置方法，并且报告它直接从相关存储库检索信息。它调用专门的FlightRecorder的方法来获取Flight Recorder。 

**FrameworkMgmtPoint** - 驱动器控制和配置命令（启动，停止，重置，设置主机的地址，等等）的点。FrameworkMgmtPoint反过来按照正确的顺序调用其他相关模块的方法。它将生命周期请求（开始，停止，重置）直接转发到FrameworkMain，以按正确的顺序驱动它们。 

**Defense4All应用程序** - 任何SDN应用程序（在本例中为Defense4All）都应该Implement/extendAppRoot对象。SDN应用程序没有“main”，并且它们的生命周期（启动，停止，重置）由对应用程序根对象操作的框架管理，然后应用程序根对象驱动应用程序中的所有生命周期操作。此模块还包含回框架的引用，允许应用程序使用框架服务（例如创建Repo并记录Flight Record）和常见实用程序。 

**通用类和实用程序** - 一个方便的类和实用程序库，任何框架或SDN应用程序模块都可以从中受益。示例包括包装线程服务（用于异步，定期或后台执行），字符串的短散列和用户确认。 

**存储库服务** - 框架理念中的关键元素之一是将计算状态与计算逻辑分离。所有持久状态都应该存储在一组存储库中，然后可以在不知道计算逻辑（框架或应用程序）的情况下进行复制，缓存，分发。存储库服务包括RepoFactory和Repo或其等价物 - EntityManager。RepoFactory负责与底层存储库插件服务建立连接，实例化新的请求存储库，并返回对现有存储库的引用。所选的底层存储库服务是通过Cassandra NoSQL DB的Hector Client。Repo呈现单个DB表的抽象。它使得能够只用表键就能读取整个表，（表仅由单个主键索引），记录单个单元，以及依照控制要求写入记录。并能够写入子记录（仅具有单元的一部分）。在这种情况下，显示的单元格覆盖存储库中的现有单元格。存储库中的其他单元格保持不变。与关系数据库（其中所有列都必须预先指定（在模式设计中）相反），Repo利用底层Cassandra支持来在相同的表中包含不同列的集合（记录），其中一些列可能甚至没有预先定义。此外，具有新列的单元可以在运行中添加或移除。RepoFactory和Repo（以及它的等价物Entity Manager）构成了一个基于与Cassandra Repository集群通信的Hector客户端库之上，针对框架和SDN应用程序目标的方便的库。Cassandra集群伸缩，跨Cassandra集群成员分发数据，以及配置读/写热切和一致性大部分封装在此层中。 

**日志记录和Flight Recorder服务** - 日志记录服务使用Log4J库记录错误，警告，跟踪或提示消息。这些日志主要用为Defense4All开发人员所用。管理员可以从错误日志中获取有关故障的其他详细信息。FlightRecorder记录由Defense4All模块记录的所有航班记录，包括从外部网络元件（如ODC和AMS）接收的信息。然后，它允许用户或管理员通过REST API或CLI获取该信息。航班记录可以按类别（可以指定零个或多个）和时间范围进行过滤。FlightRecorder将所有航班记录存储在其自己的Repo中（其他repo从存储在其他区域的记录库中使用有效时间范围检索获取这些记录）。由于所有航班记录都存储在Cassandra中，Defense4All可以保留的航班记录数量仅受所有Cassandra服务器的底层持久存储容量的限制，因此即使在单个Cassandra实例上，也可以保留几个月的历史信息。 

**HealthTracker** - 在Defense4All的运行时进行统一健康状况并响应严重恶化的行为。任何模块在感测到其中或任何其他模块中的意外和/或错误行为时，可以在HealthTracker中记录“健康问题”，并提供健康问题的症状。这不会直接触发Defense4All终止。这意味着在短时间内聚合的许多健康问题可能意味着Defense4All存在问题，但是偶发性和/或间歇性操作性“打嗝”可以被忽略，即使Defense4All保持小于100％可操作管理员可以始终将其重置以完全恢复）。结果，每个非永久性健康问题都随着时间的推移逐渐减弱。如果Defense4Al健康恶化低于预定义阈值，HealthTracker会根据健康问题的性质触发响应操作。重启可以治愈瞬态问题，因此HealthTracker触发Defense4All终止（作为Linux服务运行，Defense4All自动重新启动）。要从更永久的问题中恢复，HealthTracker可能还会触发Defense4All重置。如果这没有帮助，下次HealthTracker尝试更深层次的重置。作为最后手段，可以建议管理员执行出厂重置。 

**ClusterMgr** - 当前未实现，**所以不翻译了，保留机翻**。此模块负责管理一个Defense4All集群（与Cassandra或ODC集群分开，建模为单独的层集群）。集群的Defense4All提高了高可用性和可扩展性。Defense4All框架或应用程序中的任何模块都可以向ClusterMgr注册集群操作，指定其功能是由单个还是由多个/所有活动实例（在不同的Defense4All集群成员上运行）执行。当集群成员资格更改时，ClusterMgr会通知每个模块中的每个实例其在该模块的集群操作中的角色。如果存在单个活动实例，则会通知该实例在集群中的角色，而所有其他实例都会通知它们处于待机模式。如果有多个活动实例，则通知每个活动实例关于该范围中的活动实例数及其逻辑枚举。所有状态都存储在一个全局可访问和共享的存储库中，所以模块的任何实例都是无状态的，并且可以在每个成员更改后执行任何角色。例如，在成员资格变化N之后，实例可以被枚举为7中的2个，作为执行工作的相关部分的结果。在成员变化N + 1，相同的实例可以枚举6中的5，并且执行分配给5而不是2的工作部分。跳过对等消息服务，ClusterMgr可以提供用于更协调的跨实例操作。 

Defense4All应用程序是高度可插拔的。它可以适应不同版本的ODC和不同AMS的不同攻击检测机制，不同攻击缓解驱动程序和驱动程序（称为代表）。Defense4All应用程序包括“核心”模块和实现定义明确的Defense4All应用程序API的“可插拔”模块。

应用程序视图
![Defense4All Defense4All应用程序结构](https://wiki.opendaylight.org/images/thumb/3/35/D4a_application_view.jpg/1350px-D4a_application_view.jpg)
以下是对Defense4All应用程序模块的描述：


**DFAppRoot** - Defense4All应用程序的根模块。Defense4All应用程序没有“main”，它的生命周期（启动，停止，重置）由对此模块操作的框架管理，从而驱动Defense4All应用程序中的所有生命周期操作。DFAppRoot还包含对所有Defense4All应用程序模块（核心和可插入），全局存储库和回到框架的引用的引用，从而允许Defense4All应用程序模块使用框架服务（例如创建Repo并记录航班记录）和公用实用程序。


**DFRestService** - 构成Defense4All应用程序servlet的一组类，用于响应Defense4All应用程序REST请求。DFRestService针对DFMgmtPoint调用控制和配置方法，并且报告它直接从相关存储库检索信息。对于飞行记录，它调用针对FlightRecorder的方法。


**DFMgmtPoint** - 驱动器控制和配置命令（如addams和addpn）的点。DFMgmtPoint反过来按照正确的顺序调用其他相关模块的方法。


**ODL Reps** - 用于不同版本的ODC的可插拔模块集。包括两个子模块中的两个功能：相关流量的统计收集和流量转移。这两个子模块遵循StatsCollectionRep DvsnRep API。ODL Reps在图6和其后的描述中详细描述。


**SDNStatsCollector** - 负责在指定的网络位置（物理或逻辑）为每个PN设置“计数器”。计数器是ODC可用的网络交换机和路由器中的一组OpenFlow流条目。SDNStatsCollector定期从这些计数器收集统计信息，并将它们馈送到SDNBasedDetectionMgr（请参阅下面的描述）。模块使用SDNStatsCollectionRep设置计数器并从这些计数器读取最新的统计信息。stat报告由读取时间，计数器规范，PN标签和trafficData信息列表组成，其中每个trafficData元素包含在计数器位置中为`<protocol，port，direction>`配置的流条目的最新字节和分组值。协议可以是`{tcp，udp，icmp，other ip}`，端口是任何第4层端口，方向可以是`{inbound，outbound}`。


**SDNBasedDetectionMgr** - 用于可插拔基于SDN的检测器的容器。它将从SDNStatsCollector接收的统计信息报告给插入的基于SDN的检测器。它还从关于结束攻击的攻击分辨点（参见下面的描述）馈送所有基于SDN的检测器通知（以允许检测机制的重置）。


**RateBasedDetector子模块** - 该检测器针对每个PN了解其随时间的正常流量行为，并在检测到流量异常时通知AttackDecisionPoint（参见下面的描述）。对于每个PN的每个协议{TCP，UDP，ICMP，其他IP}，RateBasedDetector保持字节和分组的最新速率和指数移动平均值（基线）以及最后读取时间。检测器保持每个计数器的这些值以及每个PN的所有计数器的聚合。在两个计算级别（计数器和PN聚合）的组织允许更好的可扩展性（例如使用集群ODC，其中每个实例负责从一部分网络交换机获得统计数据，并绕过ODC单实例映像API）。这样的组织还使得能够进行更精确的统计收集（避免在非常小的时间间隔期间收集所有统计的困难）。统计在计数器级处理，并且在PN级定期聚合。持续检测到流量异时，RateBasedDetector通知AttackDecisionPoint进行攻击检测。然后，在一段时间内没有异常，就通知检测器停止该攻击检测。检测器可以指定检测持续时间，在该检测持续时间内检测有效。之后，检测过期，但可以“延长”另一个关于相同攻击的通知。


**AttackDecisionPoint** - 此模块负责在攻击生命周期内应对攻击。它可以从多个检测器接收攻击检测信息。Defense4All支持RateBasedDetector，外部检测器（未来版本可用）和基于AMS的检测器参考实现（基于Radware的DefensePro）。在当前版本中，AttackDecisionPoint完全尊重每次检测（最大检测器置信度和最大检测置信度）。它为每次检测到新攻击的流量（PN，协议和端口）声明一个新的攻击，并且为现有（已声明的）攻击增加更多的检测。模块定期检查所有攻击的状态。只要至少有一个未到期的检测（每个检测都有一个到期时间），就会继续声明攻击。如果对于给定攻击的所有检测都过期，则AttackDecisionPoint声明攻击已经结束。模块通知MitigationMgr（见下面的描述）开始减轻任何新的声明的攻击。它通知MitigationMgr停止缓解已结束的攻击，并通知detectMgr重置对攻击刚刚结束的流量的统计计算。


**MitigationMgr** - 用于驱动可插拔攻击缓解程序的容器。MitigationMgr维护所有攻击缓解，并负责通知来自AttackDecisionPoint的攻击缓解请求的结果。它包含预先排好序列的MitigationDriver子模块，并尝试按照该顺序满足每个攻击缓解。如果MitigationDriveri向MitigationMgr指示它不缓解一个攻击缓解（由于每个PN偏好，AMS资源的不可用性，网络问题等），MitigationMgr将尝试通过MitigationDriveri + 1进行缓解。如果没有插入的MitigationDrivers处理缓解，则它仍处于“未缓解”状态。


**MitigationDriverLocal** - 此缓解驱动程序负责在其管理范围内使用AMS来驱动攻击缓解。当请求缓解攻击时，此缓和剂执行以下步骤顺序：

1. 它向可插入DvsnRep（见下面的描述）咨询关于从每个相关网络位置的每个管理的AMS的拓扑可行的转移选项。在此版本中，转移始终从安装统计计数器的位置执行。
2. MitigationDriverLocal在所有可行选项中选择一个AMS（在第一个版本中，只选择列表中的第一个）。
3. 在指示将业务转移到每个AMS之前，它随机分配所有AMS（每个转向源可以有与其相关联的不同AMS）。这是通过插入AMSRep完成的。
4. MitigationDriverLocal指示DvsnRep，将来自每个源NetNode（在此版本中，NetNode为SDN交换机）的流量转移到与该NetNode相关联的AMS。转移可以仅用于入站流量，也可以用于入站和出站流量。
5. 攻击缓解驱动程序通知AMSBasedDetector开始监视所有AMS中的攻击状态，并将攻击检测馈送到AttackDecisionPoint。
6. 在未来版本中，MitigationDriverLocal会监视所有AMS和网络拓扑的相关部分的运行状况，可以重新选择AMS（如果该AMS缓解失败），或者依照网络拓扑变动重新选择AMS。
当缓解应该结束时，MitigationDriverLocal通知AMSBasedDetector停止监视已结束攻击的攻击状态，通知DvsnRep停止对所有AMS的流量转移以进行此缓解，最后通知AMSRep可选地清除每个缓冲中的所有缓解相关配置集相关的AMS。


**AMSBasedDetector** - 此可选模块（可以打包为AMSRep的一部分）负责监视/查询AMS的攻击缓解。注册为检测器，此模块然后可以通知AttackDecisionPoint有关攻击持续和结束。它仅监视指定的AMS，并且仅监视指定（攻击的）流量。


**AMSRep** - 用于不同AMS的可插拔模块。该模块遵守AMSRep API。它可以支持所有引入的AMS的配置（永久或在攻击缓解之前/之后）。它还可以接收/查询安全信息（攻击状态），以及操作信息（健康，负载）。AMSRep模块是完全可选的 - AMS可以在外部配置和监控。在许多情况下，攻击可以继续通过SDN计数器单独监控。Defense4All包含与Radware的DefensePro AMS通信的参考实现AMSRep。

## ODL代表视图 ##
![Defense4All Defense4All ODL Reps结构](https://wiki.opendaylight.org/images/thumb/5/54/D4a_odl_reps_view.jpg/1350px-D4a_odl_reps_view.jpg)
上图描述了Defense4All应用程序ODL Reps模块集结构。不同版本的OFC可以用ODL Reps模块集的不同版本实现。ODLReps包括两个功能：相关流量的统计收集和流量转移。这两个功能或任一功能可以在给定部署中使用。因此，它们在与ODC通信，并保存ODC的所有一般信息（见下文）具有相通之处。


ODL Reps支持两种类型的SDN交换机：sdn-hybrid，支持SDN和传统路由; sdn-native，支持仅SDN路由。在传统路由中，通过对预设流表和动作“send to normal”的流条目进行计数来计算sdn混合交换机上的流量。在sdn-native交换机上计数流量需要一个显式路由操作（将流量发送到哪个输出端口）。Defense4All通过要求一个sdn本地交换机避免学习所有路由表，该交换机或多或少是相对于流量路由的有线路径（即，进入端口1的流量通常退出端口2，进入端口3的流量通常退出端口4反之亦然）。**（Defense4All avoids learning all routing tables by requiring an sdn-native switch which is more or less a bump-in the wire with respect to traffic routing (that is, traffic entering port 1 normally exits port 2 and traffic entering port 3 normally exits port 4 and vice versa).）**这种交换机允许简单地对流条目进行编程，以便计数流量或者将流量转移到附接的AMS /从附接的AMS转移流量。当Defense4All编写具有包括端口1的选择标准的流量计数流条目时，其操作被输出到端口2，并且类似地被输出从3到4.在未来的版本中，该限制被排除。


以下是子模块的描述：


**StatsCollectionRep** - 该模块遵守StatsCollectionRep API。其主要任务是：

* 在网络中提供计数器展示位置NetNodes。提供的NetNodes是为PN定义的所有NetNode。这基本上映射了哪些SDN交换给定PN流的业务。
* 在选定的NetNodes中添加一个平时计数器，以收集给定PN的统计信息。StatsCollectionRep为每个PN中的NetNode创建一个单独的计数器。（总体上，NetNode可以具有用于不同PN的多个计数器;并且PN可以在为给定PN指定的NetNode中具有多个计数器）。StatsCollectionRep将NetNode中的计数器的安装转换为为NetNode端口中的每个“北向接口”编程四个流条目（对于TCP，UDP，ICMP和其余IP），从客户端到受保护的PN进入SDN交换机。例如，StatsCollectionRep对具有三个端口的SDN交换机中的给定PN 12流条目添加，其中PN的入站业务进入OFS。并且，如果指定另一个NetNode（SDN交换机）使该PN的入站流量通过两个端口进入，则StatsCollectionRep在该第二个NetNode中为该PN添加八个流条目。
* 删除一个平时计数器。
* 读取指定计数器的最新计数器值。StatsCollectionRep返回在每个方向上（每个方向上当前只支持“从北向南”）为每个协议端口计数的最新字节和分组的向量，以及从ODC接收到读取的时间。

**DvsnRep** - 该模块遵守DvsnRep API。其主要任务是：

* 将指定NetNode的转移属性返回给定的AMS。在这个版本中，如果这样的转移是拓扑可行的（AMS直接附加到指定的NetNode被建模的SDN交换机上，则返回空属性，否则不返回任何属性，这在将来的版本中留有用于远程转移的空间，每个远程AMS的拓扑成本，例如延迟，带宽预留和成本）。
* 通过AMS转发（攻击）来自指定NetNode的流量。因此，新的流条目优先于和平时段的流条目。DvsnRep编程流条目，以将每个“北向接口”的入站攻击流量（如果PN需要的话，也可以为所有流量）转移到AMS“北”端口。如果已经为该PN指定了“对称转移”（对于入站和返回，出站业务），DvsnRep编程另一组流条目以将来自每个“南业务端口”的受攻击（或所有）流量转移到AMS“南”港口。在sdn混合交换机部署中，DvsnRep为从AMS南端口返回的入站流量添加流条目，并将操作发送到正常，并且类似地，它为来自AMS北端口的出站返回流量添加流条目，也送正常。在SDN本地交换机中，操作是发送到正确的输出端口，但是如果这种情况下，过程对于确定正确的端口更复杂。北端口MAC学习用于从分组中的源/目的MAC确定正确的输出端口。这种流条目的方案适用于TCP，UDP和ICMP攻击。对于“其他IP”攻击，流条目编程更复杂，这里为了清楚起见而被抑制。被编程为转移（但仍然计数）业务的流条目的集合包括“攻击业务量楼层”。可能有许多攻击流量楼层，所有这些都优先于和平时间统计收集楼层（通过编程较高优先级流条目）。额外的攻击（除了“其他IP”攻击，这是一种特殊情况，在此被抑制）创建具有高于先前设置的攻击流量楼层的优先级流量楼层。攻击可以完全或部分地“早期”攻击（例如，通过TCP的TCP端口80，反之亦然），或者不连接（例如TCP和UDP）。统计收集来自所有交通楼层，包括平时和攻击。基于SDN的检测器将所有统计数据聚合为总速率，从而确定攻击是否仍在进行中。（请注意，黯淡的和平时间计算的交通可能显示零利率，并且计数由较高优先权楼层计数器补充）。**在细看**
* 结束转移。DvsnRep删除相关的攻击流量层（即从NetNode中删除其所有流条目）。注意，这既不影响已移除楼层“上方”的流量层，也不影响“下面的流量层”。此外，基于SDN的检测器从剩余楼层的计数器接收相同的累计费率，因此其操作也不受影响。

**ODLCommon** - 此模块包含编程ODC中流条目所需的所有常见元素。这允许通过StatsCollectionRep和DvsnRep对配置的ODC（在该版本中，至多一个）进行相干编程。例如ODLCommon实例化与ODC的连接，维护分配给每个ODC的编程流条目和cookie的列表。它还维护对DFAppRoot和FrameworkMain的引用。为每个受保护链路添加sdn本地NetNode时ODLCommon（输入到输出端口对）添加2个流条目，以在两个端口之间传输流量（进入北端口的流量路由到南端口，反之亦然）。ODLCommon为连接到AMS的每个端口添加两个流条目，以阻止返回ARP流量（以便在未配置AMS的情况下避免ARP洪泛）。该“公共业务层”流条目被设置为具有最低优先级。他们的柜台既不是统计收集也不是交通分流。当删除NetNode时，ODLCommon删除此公共流量底层流条目。


**FlowEntryMgr** - 此模块提供了一个API，用于对ODC管理的SDN交换机中的流条目执行操作，并检索有关ODC管理的所有节点的信息。流条目操作包括在指定的NetNode（SDN交换机/路由器）中添加指定的流条目，移除流条目，切换流条目，获得流条目的详细信息以及读取由流条目收集的统计信息。FlowEntryMgr使用连接器模块与ODC通信。


**连接器** - 此模块提供了包括REST通信的与ODC通信的基本API调用。在使用指定的ODC初始化连接详细信息后，连接器允许从ODC获取或删除数据，以及将数据发布或放入ODC。


**ODL REST Pojos** - 此组Java类是ODC REST API的一部分，指定参数的Java类和与ODC交互的结果。

## 基本控制流程 ##
控制流在逻辑上根据模块运行时依赖性排序，因此如果模块A依赖于模块B，则应在模块A之前初始化模块B，并在模块A之后终止。Defense4All应用程序模块依赖于大多数Framework模块，WebServer除外。

* 启动 - Defense4All初始化所有模块并重新应用以前配置的基础架构和安全设置，从持久性存储库获取它们。在启动过程结束时，Defense4All恢复其先前的操作。
* 终止 - 重新启动 - Defense4All将任何相关数据保存到稳定的存储库中，并自行终止。如果终止是重新启动，则自动重新启动机制重新启动Defense4All。否则（如升级）Defense4All不会自动重启。
* 复位 - 在此流程中，所有模块都复位为出厂设置。这意味着将删除所有动态获取的数据以及用户配置。
## 配置和设置流程 ##
* OFC（OpenFlowController = ODC） - 当DFMgmtPoint从DFRestService接收到添加OFC的请求时，它首先在OFC的Repo中添加OFC记录，然后通知ODLStatsCollectionRep和ODLDvsnRep，它们又通知ODL发起到添加的OFC（ODC）。ODL实例化用于与ODC通信的REST客户端。


* NetNode - 可以添加多个NetNodes。每个NetNode对交换机或类似的网络设备及其流量端口，受保护的链路和与AMS的连接建模。当DFMgmtPoint从DFRestService接收到添加NetNode的请求时，它首先在NetNodes Repo中记录添加的NetNode，然后通知ODLStatsCollectionRep和ODLDvsnRep，然后通知MitigationMgr。ODLStatsCollectionRep和ODLDvsnRep然后通知ODL，并且ODL安装低优先级流条目以在受保护链路的端口对之间传递流量。MitigationMgr通知MitigationDriverLocal，它更新其NetNode-AMS连接组以从给定的NetNodes中转移AMS的一致性分配。


* AMS - 可以添加多个AMS。当DFMgmtPoint从DFRestService接收到添加AMS的请求时，它首先在AMS的Repo中记录添加的AMS，然后通知AMSRep。AMSRep可以选择在添加的AMS中预配置保护功能，并开始监控其运行状况。


* PN - 可以添加多个PN。当DFMgmtPoint从DFRestService接收到添加PN的请求时，它首先在PN的Repo中记录添加的PN，通知MitigationMgr，然后最终通知DetectionMgr。MitigationMgr通知MitigationDriverLocal，然后通知AMSRep。AMSRep可以预配置此PN的AMS以及其EventMgr以接受与该PN的业务相关的事件。DetectionMgr通知RateBasedDetector，然后通知StatsCollector。StatsCollector查询ODLStatsCollectionRep可能放置此PN的统计信息收集计数器。ODLStatsCollectionRep返回为此PN配置的所有NetNodes（如果没有配置，则返回当前已知为Defense4All的所有NetNodes）。StatsCollector“选择”计数器位置选项（此版本中唯一可用的选项）。对于每个NetNode，它然后询问ODLStatsCollectionRep以创建用于主题PN的计数器。计数器本质上是在每个北通信端口上针对感兴趣的协议（TCP，UDP，ICMP和其余IP）设置的一组流条目。计数器被给予优先级，并且这构成了平时通信量楼层（通过周期性地读取所有计数器流条目通信量计数值来监视通信量）。因为PN可以在重新启动时被重新引入，或者网络拓扑的改变可能需要重新计算计数器位置，所以一些/所有计数器可能已经就位。仅添加新计数器。不再删除的计数器。ODLStatsCollectionRep根据NetNode类型配置流条目。对于混合NetNodes，流条目操作是“发送到正常”（继续到传统路由），而对于本地NetNodes，操作是匹配输出端口（在每个受保护的链路中）。OdlStatsCollectionRep调用ODL来创建每个指定的流条目。后者调用FlowEntryMgr和Connector将请求发送到ODC。

## 攻击检测流程 ##
周期性地，StatsCollector请求ODL StatsCollectionRep向ODC查询每个配置的PN的每个集合计数器的最新统计。ODLStatsCollectionRep调用FlowEntryMgr以获取计数器中每个流条目的统计信息。后者调用连接器从ODC获取所需的统计信息。

ODLStatsCollectionRep将获得的结果聚合在stats（每个协议的最新字节和包读取）向量中，并返回该向量。StatsCollector将每个计数器stats向量提供给DetectionMgr，然后它将stats向量转发到RateBasedDetector。RateBasedDetector维护每个计数器的统计信息以及每个PN的聚合计数器统计信息。统计信息包括先前读取的时间，以及对于每个协议的最新速率和指数平均值。

RateBasedDetector检查与平均值的显着和延长的最新速率偏差，并且如果在PN聚合等级中发现这样的偏差，则它通知攻击判断点关于攻击检测。只要偏差继续，RateBasedDetector继续通知AttackDecisionPoint关于检测。它为每个检测通知设置到期时间，并且可重复的通知基本上延长检测到期。

AttackDecisionPoint尊重所有检测。如果它已经声明对该协议端口的攻击，则AttackDecisionPoint将附加检测与该现有攻击相关联。否则，它会创建一个新的攻击，并通知MitigationMgr缓解该攻击（如下所述）。AttackDecisionPoint定期检查每个实时攻击的所有检测的状态。如果所有检测已过期，AttackDecisionPoint声明攻击结束，并通知MitigationMgr停止缓解攻击。

## 攻击缓解流程 ##
MitigationMgr在接收到来自AttackDecisionPoint的缓解通知时，尝试找到插入的MitigationDriver来处理缓解。目前，它只请求插入MitigationDriverLocal。

MitigationDriverLocal检查是否存在已知，实时和可用的AMS，受攻击（或所有）流量可以从受攻击流量流经的NetNodes转向。它选择一个合适的AMS并在将攻击流量转移到所选择的AMS之前对其进行配置。例如，MitigationDriverLocal从Repo检索相关的协议平均值，并通过AMSRep在AMS中配置它们。

MitigationDriverLocal然后请求ODLDvsnRep将来自PN流量流经的每个NetNode的受攻击的PN协议端口（或所有PN）流量转移到所选择的AMS。

ODLDvsnRep创建新的最高优先级流量层（其包含优先级高于先前设置的流量层中的任何流条目的流条目）。流量层包含所有流条目，以转移和计数从每个入口/北行流量端口到AMS的流量，并且从AMS返回到相关输出（南向）端口。可选地，转向可以是“对称的”（在两个方向上），在这种情况下，流条目被添加以将流量从南向端口转移到AMS中，并且从AMS返回到北向端口。注意，StatsCollector将此添加的流量地板视为任何其他**（？）**，并将获取的统计信息从此层传递到DetectionMgr / RateBasedDetector。因为对于给定的PN，流量底板被聚合（在相同的NetNode中以及在NetNodes中），所以合并的速率保持与转移之前相同。与ODLStatsCollectionRep一样，ODLDvsnRep也使用较低级别的模块在所需的NetNodes中安装流条目。

最后，MitigationDriverLocal通知AMSRep可选地开始监视此攻击，并通知AttackDecisionPoint如果攻击继续或新的攻击发展。AMSRep可以通过AMSBasedDetector模块来实现。

如果MitigationDriverLocal找不到合适的AMS，或无法配置其任何缓解步骤，则它将中止缓解尝试，异步通知MitigationMgr。减缓措施然后保持在“无资源”状态。

当MitigationMgr接收到停止缓解攻击的通知时，它会将此通知转发到相关（且目前是唯一的）MitigationDriver，MitigationDriverLocal。MitigationDriverLocal在缓解开始时反转动作。它通知AMSRep停止监视此攻击，取消被攻击流量的转移，最后通知AMSRep可选地删除预缓解配置。

## 问题和故障排除 ##
(defense4all已经停止开发了，所以这部分无所谓了)
除了下面列出的过程，还参考Defense4All日志有关您遇到的具体问题的信息。Defense4All日志位于/var/log/Defense4All/server.log。

Defense4All无法启动 - 如果WebServer无法启动，请检查是否与端口号存在冲突。Defense4All使用端口8086.如果RepoFactory无法初始化，请检查Cassandra服务是否正在运行（sudo服务Cassandra停止/启动/重新启动）。如果Defense4All无法初始化，则问题可能在于系统资源（如线程或内存）。尝试重新启动机器。另一个问题可能会损坏Cassandra DF DB（键空间）。如果是，请尝试执行还原或重置（请参阅准则）。

Defense4All无法终止 - 可能是其WebServer崩溃。在这种情况下停止Defense4All的唯一方法是使用以下命令杀死其JVM Linux进程：kill -9 <defense4all-JVM-process-number>

Defense4All无法重置 - 启动Defense4All的一些问题也可能适用于重置。具体来说，Cassandra服务应该重置。如果重置失败，您可能需要手动清除Cassandra，如下所示：

卡桑德拉 - 克利
drop keyspace DF;
放弃;
Defense4All无法添加OFC - 检查故障原因（REST或CLI）。除了不正确的参数之外，问题可能是添加的PFC不存在或在API中指定的可寻址性（地址+端口）或安全性（用户+密码）。此外，Cassandra服务应该启动。

Defense4All无法添加NetNode - 检查故障原因（REST或API）。除了不正确的参数，Cassandra服务可能已关闭。

Defense4All无法添加AMS - 检查故障原因（REST或API）。除了不正确的参数，Cassandra服务可能已关闭。另外，问题可能在于AMS不存活或未连接。

Defense4All无法删除AMS - 检查故障原因（REST或API）。除了不正确的参数，Cassandra服务可能已关闭。另外，问题可能在于AMS不存活或未连接。

Defense4All无法添加PO - 检查失败原因（REST或API）。除了不正确的参数，Cassandra服务可能已关闭。此外，问题可能是PFC可能不存在并连接。

Defense4All无法删除PO - 检查失败原因（REST或API）。除了不正确的参数外，Cassandra服务可能已关闭。此外，问题可能是DF_GLOBAL_PNS表（列族）已损坏，因此应该完全删除，如下所示：

Cassandra-cli
使用DF;
截断列族DF_GLOBAL_PNS;
另一个问题可能是控制器没有启动，所以Defense4All无法删除它已经设置。

Defense4All无法检索/转储/清理飞行记录器日志记录 - 检查故障原因（REST或API）。除了不正确的参数外，Cassandra服务可能已关闭。此外，问题可能是FWORK_FLIGHT_RECORDER_EVENTS或FWORK_FLIGHT_RECORDER_SLICES表（列族）已损坏，因此两者都应完全删除，如下所示：

Cassandra-cli“
使用DF;
截断列族FWORK_FLIGHT_RECORDER_EVENTS;
截断列族FWORK_FLIGHT_RECORDER_SLICES;
警告：截断这些列族将导致当前存储在Cassandra中的所有航班记录的丢失。

Defense4All无法检索攻击或缓解 - 检查失败原因（REST或API）。除了不正确的参数外，Cassandra服务可能已关闭。此外，问题可能是DF_GLOBAL_ATTACKS / DF_GLOBAL_MITIGATIONS表（列族）已损坏，因此应该完全删除，如下所示：

Cassandra-cli
使用DF;
截断列族DF_GLOBAL_ATTACKS或DF_GLOBAL_MITIGATIONS;
警告：截断这些列族将导致丢失当前攻击和缓解的所有记录。在这种情况下，流量重定向流条目可能需要从所有相关的NetNodes手动删除（寻找优先级5X，7X，9X，等等，与action重定向）。

Defense4All操作错误和故障 - 根据错误/故障的性质，外部实体（Cassandra，PFC，NetNodes，AMS）的活性可能需要检查/重新启动。如果是内部Defense4All错误，应按此顺序尝试以下恢复步骤：

Defense4All重启
重新启动Defense4All宿主机
Defense4All重置并可能还原到较早的状态（以及手动清理由Defense4All设置的流条目和AMS配置）。
未检测到攻击 - 检查Defense4All日志，以查看在stats收集和检测机制中是否记录了任何错误。检查PFC和相关NetNodes是否存活。检查与平均值相比的最新速率。平均值可能偏斜，但重置Defense4All在攻击期间没有帮助，因为Defense4All将获得歪斜的平均值。等待攻击结束，然后重置Defense4All，并重新添加受攻击的PO。

缓解状态NO_RESOURCES - 这意味着Defense4All mitigationDriverLocal无法驱动此缓解，因为内部Defense4All错误或缺少AMS资源。如果确实没有AMS资源，则不需要恢复。否则检查相关AMS的活性。还要检查Cassandra是否正在运行，如果PFC已启动，并且如果相关的NetNode正在运行。如果存在Defense4All内部错误（根据Defense4All日志），则可能会在DF_GLOBAL_ATTACKS / DF_GLOBAL_MITIGATIONS表（列系列）中出现损坏，因此应全部删除，如下所示：

Cassandra-cli
使用DF;
截断列族DF_GLOBAL_ATTACKS或DF_GLOBAL_MITIGATIONS;
后续检测将重新创建相关的攻击和缓解记录。如果这没有帮助，请重新启动Defense4All。

缓解未终止 - 清除此缓解中的外部元素，如下所示：

手动删除相关攻击。
重新启动Defense4All。
重置Defense4All。
连续性
服务连续性（而不是高可用性）在这里定义为在存在中断事件的情况下以可承受的成本提供所需级别的服务的能力，其中

中断事件可以是加载，更改，逻辑错误，故障和灾难，管理操作（如升级），外部攻击等。
服务级别可以包括响应时间，吞吐量，数据/操作的生存性，安全/隐私等。对于每种类型的事件，在不同的事件处理阶段，所需的服务级别可以针对每个服务功能而不同。
成本可以包括人（数量，专业知识），设备（硬件，软件），设施（空间，电力）。
集群和容错 - 集群有助于解决可伸缩性和高可用性。如果其中一个集群成员失败，另一个集群成员可以快速承担其责任。这克服了成员故障，成员托管机器故障和成员网络连接故障。Defense4All集群安排在未来版本。在版本1.0中，Defense4All作为Linux可重新启动服务运行，因此如果它失败，托管Linux操作系统恢复Defense4All。这可以克服间歇性/偶发性Defense4All故障。Defense4All宿主机的故障意味着更长的时间和温和的额外的人类努力恢复机器及其托管的Defense4All。如果机器无法启动，Defense4All可以在网络中的另一台机器上启动。为了确保Defense4All恢复其操作（而不是从头开始重新启动），必须在该计算机上预加载Defense4All（最新或更早）状态快照。非集群环境会影响从机器故障中恢复的时间和人力。时间因素不那么重要，因为Defense4All运行路径不足，所以其较长的非可用性期间意味着较长的时间来检测和减轻新的攻击。

状态持久性 - Defense4All保持在同一台机器上运行的Cassandra DB中的状态。在版本1.0中，仅配置一个Cassandra实例集群。只要本地稳定存储不会崩溃，Linux4重新启动Defense4All服务就能使Defense4All从Cassandra快速检索其最新状态并恢复其最新操作。在发生故障和重新启动承载Defense4All的机器时也会发生同样的情况。采取Defense4All状态备份，并在另一台机器上恢复允许恢复该机器上的Defense4All操作。多节点Cassandra集群（计划用于未来版本）将增加状态持久性，同时减少恢复时间和工作量。

重新启动过程 - 当Defense4All（re）启动时，它首先检查保存的配置数据，并对所有相关模块重新执行配置步骤，驱动任何相关的外部编程和/或配置操作（例如针对PFC或AMS设备），例如，重新添加PO。此配置重放和原始配置之间的唯一区别是，保留任何动态获取的数据，例如所有PO统计信息。这允许轻松达到内部一致性，特别是在Defense4All或其主机计算机崩溃的情况下。当针对外部实体重放配置动作派生时，例如添加缺失的PO统计计数器，并且去除不再需要的计数器，也达到与外部实体的一致性。Defense4All变得可操作（启动其Web服务器），让您或一些其他组件根据可能的更改完成Defense4All缺少的配置，而Defense4All关闭。这导致达到端到端的一致性。

Reset - Defense4All允许您重置其动态获取的数据和配置信息（恢复出厂设置）。这使您能够克服许多逻辑错误和错误配置。注意，Defense4All重新启动或故障转移不会克服这样的问题。因此，这种机制是对重新启动 - 故障转移机制的补充，通常应该作为最后手段应用。

故障隔离和健康跟踪器 - 在Defense4All中，故障隔离以立即恢复或补偿（尽可能多）失败的形式发生，并在称为健康跟踪器的特殊模块中记录故障。除了少数实质性故障（例如无法启动框架），任何模块中的任何故障都不会立即导致Defense4All停止。相反，每个模块在其范围内记录每个故障，提供严重性规范和故障持久性的指示。如果所有故障的组合严重性（永久或临时）超过全局设置的阈值，HealthTracker会触发Defense4All关闭（并由Linux进行复原）。以后，永久或重复的临时故障将导致HealthTracker触发Defense4All软动态和动态重置（动态获取的数据）或建议管理员执行恢复出厂设置（也包括配置信息）。

状态备份和恢复 - 管理员可以对Defense4All状态进行快照，将备份保存在其他位置，然后恢复到原始或新的Defense4All位置。这允许克服某些逻辑错误和错误配置，以及承载Defense4All的机器的永久故障。要快照Defense4All状态，请执行以下操作：

停顿（关闭）Defense4All，导致当前状态刷新到稳定的存储）。避免在恢复时执行任何配置更改，避免新的状态更改。
获取Defense4All DB - “DF”的Cassandra快照：有关备份还原准则，请参阅 http://www.datastax.com/docs/1.0/operations/backup_restore。
将快照文件复制到所需的存储归档。
要将Defense4All备份恢复到目标计算机，请执行以下操作：

在目标计算机中恢复所需的已保存快照（与备份相同或不同）。有关Cassandra备份 - 恢复指南，请参阅http://www.datastax.com/docs/1.0/operations/ backup_restore。
在该机器上启动Cassandra。
在该机器上启动Defense4All。
维护和升级
关于升级的一个关键问题是数据的格式是否改变。未来版本的Defense4All将必须通过以下两种方式之一处理更改的数据格式：

作为升级过程的一部分，自动升级存储库中的状态格式
需要Defense4All重置，然后删除所有现有存储库表，然后创建新格式的表。
有关升级的另一个关键问题是与外部实体的兼容性：OFC，NetNodes和AMS。升级版本中的StatsCollectionRep，DvsnRep和AmsRep必须能够处理其外部实体，无论是从头开始还是从先前设置的配置和在运行时获得的数据。

Defense4All升级过程包括：

备份它的状态
或者出厂重置它
停止它
升级任何外部实体
升级它
（Re）启动它
要降级Defense4All：

出厂复位
停下来
降级任何外部实体
降级它
在升级之前恢复其备份状态
启动它
因为在这个版本Defense4All没有集群，集群滚动升级不适用这里。

新术语和概念
AMS - 攻击减轻系统，用于检测，减轻和报告网络攻击。例如，Radware的DefensePro是一种能够检测，减轻和报告大量网络攻击的AMS。
攻击 - 怀疑或检测到PN上的DDoS或其他网络攻击。攻击可以是网络链路，目的地址，协议或第4层端口的任何组合。Defense4All维持攻击生命周期，其中它试图根据每个主题PN的规范来减轻攻击。
检测 - 检测到监控的交通异常的指示。检测具有到期时间，并且可以更新。
缓解 - 正在执行以缓解给定攻击的活动。在本版本中，所有攻击都通过将攻击流量转移到DefensePro来缓解，并将干净的流量重新注入VTN。
NetNode - 模拟交换机或类似的网络设备，以及其流量端口，受保护的链路和与AMS的连接。NetNode指定一个或多个PN的业务通常流过（如果不被重定向）和/或AMS连接到的感兴趣的网络位置。在和平时间Defense4All通过该PN的业务流在每个网络中设置PN的计数器。攻击Defense4All选择一个或多个连接到引入的NetNodes的AMS，并将攻击/所有流量重定向到连接到这些NetNode的AMS。
OFC - 支持OpenFlow网络编程的SDN控制器（OFC代表OpenFlow控制器）。OpenDaylight控制器为支持OpenFlow的网络设备和其他网络设备提供这种风格。
受保护网络（PN） - 这是具有给定保护规范的用户定义的受保护网络元素。网络由两部分的任意组合指定：1）网络地址范围（并且可选地仅协议L4端口），以及2）指定业务流经过的网络链路。两个部分（但不是两者）中的任一个可以是未指定的。在版本1.0中，仅实现了第二部分。保护规范指示与对主题PN的攻击的检测和缓解相关的属性的范围。
受保护对象（PO​​） - 在Defense4All中，受保护网络（PN）称为受保护对象（PO​​）。
保护链路 - 交换机中的一组入口 - 出口端口。Defense4All通过要求在网络拓扑中设置为“bump-in-the-wire”的sdn-native开关，避免学习所有路由表。Defense4All程序使未攻击的流量进入其中一个入口 - 出口对端口以退出另一个。
sdn-hybrid和sdn- native - 有两种类型的SDN交换机：sdn-hybrid，它支持SDN和传统路由，sdn-native支持仅SDN路由。在sdn混合交换机上计数流量可以通过对具有期望的流量选择标准的流条目进行编程来完成，并且动作被发送到正常，意味着继续传统路由。对sdn本机交换机上的流量计数需要显式路由操作（将流量发送到哪个输出端口）。Defense4All通过需要一个sdn本地交换机来避免学习所有路由表，该交换机或多或少是相对于流量路由的线路，意味着进入端口1的业务通常退出端口2，进入端口3的业务通常退出端口4，反之亦然。这样的开关允许简单地对流条目进行编程，以便计数流量或将流量转移到附接的AMS /从附接的AMS转移流量。当Defense4All对包含端口1的选择标准的流量计数流条目进行编程时，其操作将输出到端口2，类似于3到4.此限制预计在将来的版本中取消。另请参阅受保护的链接项。
流量楼层 - 在给定的NetNode上Defense4All程序的一组流条目。不同的PN具有它们自己的业务层用于和平时间攻击检测，以及用于攻击缓解的业务重定向。攻击流量楼层包含优先级高于该PN的所有先前设置的攻击缓解流量楼层的流条目，以及平时流量楼层（其包含仅用于计数PN流量以便学习行为和异常的流条目）。攻击可以完全或部分地删除早期的攻击（例如TCP上的TCP端口80通过TCP或反之亦然）或不连接（例如TCP和UDP）。统计收集来自所有交通楼层：平时和攻击。基于SDN的检测器将所有统计信息汇总到总速率，确定攻击是否仍然开启。被削减的和平时间计数的交通可以显示零利率，并且计数由更高优先权的楼层计数器补充。
流量端口 - 北Vexternal，流量通过该端口进入PN的VTN。ProtectedLink是NetNode中的Vbridge的名称。
其他信息
安全和隐私
Defense4All REST API目前不检查凭据。它也不根据允许或限制某些REST API的使用来定义用户角色。

兼容性
这个Defense4All版本（1.0.7）与ODC 1.0兼容。AmsRep的参考实现是通过Radware的DefensePro版本：硬件版本ODS-VL，软件版本6.03,6.07和6.09。

性能和可伸缩性信息
TBD。

参考资料
[ODC]


返回Defense4All用户指南页面

隐私政策 关于OpenDaylight项目 免责声明
 的页面