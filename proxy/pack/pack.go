package pack

import (
	"bitcask_master_slave/node/kitex_gen/node"
	prxyservice "bitcask_master_slave/proxy/kitex_gen/prxyService"
)

func RePackOpLogEntryReq(req *prxyservice.LogEntryRequest) *node.LogEntryRequest {
	return &node.LogEntryRequest{
		EntryId:  req.EntryId,
		MasterId: req.MasterId,
		Cmd:      req.Cmd,
		Args_:    req.Args_,
	}
}

func RePackOpLogEntryResp(resp *node.LogEntryResponse) *prxyservice.LogEntryResponse {
	entries := make([]*prxyservice.LogEntry, len(resp.Entries))
	for i, entry := range resp.Entries {
		entries[i] = &prxyservice.LogEntry{
			Key:      entry.Key,
			Value:    entry.Value,
			Score:    entry.Score,
			ExpireAt: entry.ExpireAt,
		}
	}
	return &prxyservice.LogEntryResponse{
		BaseResp: (*prxyservice.BaseResp)(resp.BaseResp),
		Entries:  entries,
		Info:     resp.Info,
	}
}
