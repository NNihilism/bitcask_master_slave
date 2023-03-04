// bitcask-cluster的节点
package main

import (
	"bitcask_master_slave/log"
	"bitcask_master_slave/node/config"
	nodeCore "bitcask_master_slave/node/node_core"
	"bitcask_master_slave/pkg/consts"
	"strings"

	"github.com/NNihilism/bitcaskdb/options"
)

func init() {
	parts := strings.Split(consts.NodeAddr, ":") // []string{"ip", "port"}
	nodeConfig := &config.NodeConfig{
		Role:      config.Master,
		Addr:      consts.NodeAddr,
		Path:      config.BaseDBPath + parts[1],
		ID:        consts.NodeAddr,
		RemakeDir: config.RemakeDir,
		Weight:    consts.Weight,
	}

	var err error
	bitcaskNode, err = nodeCore.NewBitcaskNode(nodeConfig, options.Options{})
	if err != nil {
		log.Errorf("create bitcasknode err : %v", err)
	}

}

var bitcaskNode *nodeCore.BitcaskNode
