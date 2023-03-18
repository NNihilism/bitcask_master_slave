# 主从模式、哨兵模式、集群模式系列 `Part 1`

## bitcask的主从模式

### 整体流程：
- 初始化
- 用户与代理proxy建立连接
- 用户发起数据操作命令
- proxy将读写请求进行分离，转发至不同的节点
- 客户端进行读写请求分离【没有选择用】
  - 读请求：根据负载均衡策略转发至从节点
  - 写请求：转发至主节点
    - 主节点主动选择择同步\半同步\异步进行数据同步
    - 从节点主动请求数据更新
- proxy收到命令结果并返回给客户
        
### 主从节点之间的结构：
#### 星型复制【√】
- 优点：
  - 从节点之间相互独立，单一节点异常不会影响其他节点
- 缺点：
  - 从节点越多，主节点压力越大
  - 扩展难题，整体性能受限于mastet
![星型复制](https://github.com/NNihilism/bitcask_master_slave/blob/master/resource/%E6%98%9F%E5%9E%8B%E5%A4%8D%E5%88%B6.png)
#### 链式复制【x】
- 优点：
  - 解决了星型结构的扩展难题
- 缺点：
  - 复制链越长，末端同步延迟越大
  - 某个节点异常，会影响到下游所有节点
![星型复制](https://github.com/NNihilism/bitcask_master_slave/blob/master/resource/%E9%93%BE%E5%BC%8F%E5%A4%8D%E5%88%B6.png)

### 通讯模型
- 节点之间使用rpc通讯【√】
  - 优点
    - 功能完备，不用担心tcp命令解析等一系列问题
    - 可以借此重构下通讯模块部分
  - 缺点
    - 还没用过grpc，不知上手情况
- 节点之间直接使用tcp连接【x】
  - 缺点
      - 现成代码只有客户端主动向服务端发出请求的部分
    - 实现起来繁琐且功能不足，如重连次数，超时机制都得自己实现
  - 优点
    - 简历上有东西可写
    - 有现成的代码，bitcask单机模式下已经实现了c/s通讯

### 主从模式优点：
    读写分离，分担主节点的压力
    负载均衡，在读多写少的场景中，可以增加从节点来分担redis-server读操作的负载能力，从而大大提高redis-server的并发量
    保证高可用，容灾快速恢复：如果某台从节点挂了，客户端会切换到其他节点进行读操作，如果主节点出现故障后，可以切换到从节点继续工作，保证redisserver的高可用。
    数据冗余，主从复制实现了数据的热备，是除了持久化机制之外的另外一种数据冗余方式，数据不易丢失。

### 主从模式缺点：
    主从节点之间的数据一致问题
    容量有限

### 使用场景：
    读写分离主要使用在对一致性要求不高的场景下

### 高可用性：
    HA(高可用模块)对所有节点健康状态进行监听，master宕机时自动切换到新主节点
    proxy对所有节点健康状态进行监听，异常从节点的权重会被下调，直至下线
    redis-proxy和HA一起做到尽量减少业务对后端异常的感知，提高服务可用性。

### 高性能：
    对于读多写少的业务场景，直接使用集群版本往往不是最合适的方案，现在读写分离提供了更多的选择，业务可以根据场景选择最适合的规格，充分利用每一个read-only replica的资源。
  
### 启动流程
    1. 填写配置文件 xxx.config，包括node,proxy的地址
    2. 开启多个bitcask_node实例，且监听不同的地址/端口
    3. 设置主从关系，在从节点连接窗口内输入命令 slaveof xxx
    4. 开启代理实例，代理会根据配置文件中的地址去与bitcask_node节点进行连接，并查看是否只有一个主节点，若是，则代理启动成功
    5. 开启HA模块（可能有）,HA与代理相同，与节点逐个进行连接

### 讨论
#### 1. 是否需要代理
  在go-redis中是没有代理这一层的，由client直接配置集群信息，代码如下
  ```
  rdb := redis.NewFailoverClient(&redis.FailoverOptions{
    MasterName:    "master-name",
    SentinelAddrs: []string{":9126", ":9127", ":9128"},
})
  ```
  比起代理模式，这样做的好处的会减少一次网络请求的转发，然而，这样做的不足就是，不容易对节点的变更及时察觉。无论是选择让客户端不断地询问主节点信息是否变更，还是让服务端主动的通知客户端变更信息，当客户端数量上来后，都会给服务端造成不小的压力，原因就是客户端与主节点直接相连。  
  而使用代理模式，则能缓解上述问题，无论节点信息如何改变，主节点都只需要与代理节点进行通信。

#### 2. 可改进的点
  - 目前master发往同一个slave的数据是"串行化的"，因为slave的接受逻辑是只接受序号递增的写请求，对于序号不符条件的都直接丢弃掉，因此可以添加一个缓冲区，用于接受序号大于当前offset但不是想要的序号的请求，同时，master的锁的范围也可以因此减小，对于序号有大小关系的请求，不必限制其发送顺序。



redis主从分析 ：https://help.aliyun.com/document_detail/65001.html
redis主从使用 ：https://blog.csdn.net/weixin_40980639/article/details/125569460



#### benchmark
```
一主零从
goos: linux
goarch: amd64
pkg: bitcask_master_slave/benchmark
cpu: AMD Ryzen 5 2500U with Radeon Vega Mobile Gfx  
BenchmarkMS_Get-8   	    1906	    584207 ns/op	     528 B/op	      14 allocs/op
PASS
ok  	bitcask_master_slave/benchmark	1.325s
```

```
一主一从
goos: linux
goarch: amd64
pkg: bitcask_master_slave/benchmark
cpu: AMD Ryzen 5 2500U with Radeon Vega Mobile Gfx  
BenchmarkMS_Get-8   	    1815	    581828 ns/op	     527 B/op	      14 allocs/op
PASS
ok  	bitcask_master_slave/benchmark	11.273s
```

```
一主两从
goos: linux
goarch: amd64
pkg: bitcask_master_slave/benchmark
cpu: AMD Ryzen 5 2500U with Radeon Vega Mobile Gfx  
BenchmarkMS_Get-8   	    1770	    619130 ns/op	     526 B/op	      14 allocs/op
PASS
ok  	bitcask_master_slave/benchmark	11.308s
```