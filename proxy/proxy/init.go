package proxy

import (
	"bitcask_master_slave/log"
	"bitcask_master_slave/node/kitex_gen/node"
	"bitcask_master_slave/node/kitex_gen/node/nodeservice"
	"bitcask_master_slave/pkg/consts"
	"context"

	"github.com/cloudwego/kitex/client"
)

func (proxy *Proxy) HandleProxyReq(masterAddr string) (bool, error) {
	proxy.mu.Lock()
	defer proxy.mu.Unlock()

	// 初始化MasterRPC
	masterRpc, err := nodeservice.NewClient(
		consts.NodeServiceName,
		client.WithHostPorts(masterAddr),
	)
	if err != nil {
		log.Errorf("Init master rpc err [%v]", err)
		return false, err
	}
	proxy.masterRpc = masterRpc

	// 获取所有从节点信息
	resp, err := masterRpc.GetAllNodesInfo(context.Background(), &node.GetAllNodesInfoReq{})
	if err != nil {
		log.Errorf("Get all nodes into err [%v]", err)
		return false, err
	}

	// 初始化所有SlaveRPC
	for _, info := range resp.Infos {

		// for i := 0; i < len(resp.SlaveAddress); i++ {
		tmpRpc, err := nodeservice.NewClient(
			consts.NodeServiceName,
			client.WithHostPorts(info.Addr),
		)
		if err != nil {
			log.Errorf("Init slave rpc err [%v]", err)
			continue
		}
		proxy.slaveRpcs[info.Id] = tmpRpc
		proxy.node = append(proxy.node, Node{
			addr:   info.Addr,
			id:     info.Addr,
			weight: int(info.Weight),
		})
	}

	// 返回结果
	return true, nil
}
