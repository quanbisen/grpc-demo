package discovery

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
	"strings"
	"time"
)

const schema = "etcd"

type Resolver struct {
	schema         string
	EtcdAddresses  []string
	user           string
	password       string
	DialTimeout    int
	closeCh        chan struct{}
	watchCh        clientv3.WatchChan
	cli            *clientv3.Client
	keyPrefix      string
	srvAddressList []resolver.Address
	cc             resolver.ClientConn
	logger         *logrus.Logger
}

func NewResolver(etcdAddresses []string, user, password string, logger *logrus.Logger) *Resolver {
	return &Resolver{
		schema:        schema,
		EtcdAddresses: etcdAddresses,
		user:          user,
		password:      password,
		DialTimeout:   3,
		logger:        logger,
	}
}

func (r *Resolver) Scheme() string {
	return r.schema
}

func (r *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r.cc = cc
	r.keyPrefix = BuildPrefix(Server{Name: strings.TrimLeft(target.URL.Path, "/")})
	if _, err := r.start(); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Resolver) ResolveNow(o resolver.ResolveNowOptions) {
	fmt.Println("ResolveNow calling")
}

func (r *Resolver) Close() {
	r.closeCh <- struct{}{}
}

func (r *Resolver) start() (chan<- struct{}, error) {
	var err error
	r.cli, err = clientv3.New(clientv3.Config{
		Endpoints:   r.EtcdAddresses,
		Username:    r.user,
		Password:    r.password,
		DialTimeout: time.Duration(r.DialTimeout) * time.Second,
	})
	if err != nil {
		return nil, err
	}
	r.closeCh = make(chan struct{})
	if err = r.sync(); err != nil {
		return nil, err
	}
	go r.watch()
	return r.closeCh, err
}

func (r *Resolver) sync() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	res, err := r.cli.Get(ctx, r.keyPrefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	r.srvAddressList = []resolver.Address{}
	for _, v := range res.Kvs {
		info, err := ParseValue(v.Value)
		if err != nil {
			continue
		}
		var addr = resolver.Address{}
		if info.Address == "127.0.0.1:50001" {
			addr = resolver.Address{ServerName: info.Name, Addr: info.Address, BalancerAttributes: attributes.New("xx", 3)}
		} else {
			addr = resolver.Address{ServerName: info.Name, Addr: info.Address, BalancerAttributes: attributes.New("xx", 1)}
		}

		r.srvAddressList = append(r.srvAddressList, addr)
	}
	r.cc.UpdateState(resolver.State{Addresses: r.srvAddressList})
	return nil
}

func (r *Resolver) watch() {
	ticker := time.NewTicker(time.Minute)
	r.watchCh = r.cli.Watch(context.Background(), r.keyPrefix, clientv3.WithPrefix())
	for {
		select {
		case <-r.closeCh:
			return
		case res, ok := <-r.watchCh:
			if ok {
				r.update(res.Events)
			}
		case <-ticker.C:
			if err := r.sync(); err != nil {
				r.logger.Error("sync failed", err)
			}
		}
	}
}

// update
func (r *Resolver) update(events []*clientv3.Event) {
	for _, ev := range events {
		var info Server
		var err error

		switch ev.Type {
		case clientv3.EventTypePut:
			info, err = ParseValue(ev.Kv.Value)
			if err != nil {
				continue
			}
			addr := resolver.Address{Addr: info.Address}
			if !Exist(r.srvAddressList, addr) {
				r.srvAddressList = append(r.srvAddressList, addr)
				r.cc.UpdateState(resolver.State{Addresses: r.srvAddressList})
			}
		case clientv3.EventTypeDelete:
			info, err = SplitPath(string(ev.Kv.Key))
			if err != nil {
				continue
			}
			addr := resolver.Address{Addr: info.Address}
			if s, ok := Remove(r.srvAddressList, addr); ok {
				r.srvAddressList = s
				r.cc.UpdateState(resolver.State{Addresses: r.srvAddressList})
			}
		}
	}
}
