// Code generated by Kitex v0.4.4. DO NOT EDIT.

package proxyservice

import (
	prxyservice "bitcaskDB/internal/bitcask_master_slaves/proxy/kitex_gen/prxyService"
	"context"
	client "github.com/cloudwego/kitex/client"
	callopt "github.com/cloudwego/kitex/client/callopt"
)

// Client is designed to provide IDL-compatible methods with call-option parameter for kitex framework.
type Client interface {
	Ping(ctx context.Context, req bool, callOptions ...callopt.Option) (r bool, err error)
	OpLogEntry(ctx context.Context, req *prxyservice.LogEntryRequest, callOptions ...callopt.Option) (r *prxyservice.LogEntryResponse, err error)
	Proxy(ctx context.Context, masterAddr string, callOptions ...callopt.Option) (r bool, err error)
}

// NewClient creates a client for the service defined in IDL.
func NewClient(destService string, opts ...client.Option) (Client, error) {
	var options []client.Option
	options = append(options, client.WithDestService(destService))

	options = append(options, opts...)

	kc, err := client.NewClient(serviceInfo(), options...)
	if err != nil {
		return nil, err
	}
	return &kProxyServiceClient{
		kClient: newServiceClient(kc),
	}, nil
}

// MustNewClient creates a client for the service defined in IDL. It panics if any error occurs.
func MustNewClient(destService string, opts ...client.Option) Client {
	kc, err := NewClient(destService, opts...)
	if err != nil {
		panic(err)
	}
	return kc
}

type kProxyServiceClient struct {
	*kClient
}

func (p *kProxyServiceClient) Ping(ctx context.Context, req bool, callOptions ...callopt.Option) (r bool, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.Ping(ctx, req)
}

func (p *kProxyServiceClient) OpLogEntry(ctx context.Context, req *prxyservice.LogEntryRequest, callOptions ...callopt.Option) (r *prxyservice.LogEntryResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.OpLogEntry(ctx, req)
}

func (p *kProxyServiceClient) Proxy(ctx context.Context, masterAddr string, callOptions ...callopt.Option) (r bool, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.Proxy(ctx, masterAddr)
}
