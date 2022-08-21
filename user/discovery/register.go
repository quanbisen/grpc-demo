package discovery

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"os"
	"strings"
	"time"
)

type Register struct {
	EtcdAddresses []string
	User          string
	Password      string
	DialTimeout   int
	closeCh       chan struct{}
	leaseID       clientv3.LeaseID
	keepAliveCh   <-chan *clientv3.LeaseKeepAliveResponse
	srvInfo       Server
	srvTTL        int
	cli           *clientv3.Client
	logger        *logrus.Logger
}

func NewRegister(etcdAddresses []string, user, password string, logger *logrus.Logger) *Register {
	return &Register{EtcdAddresses: etcdAddresses, User: user, Password: password, logger: logger, DialTimeout: 3}
}

func (r *Register) Register(serInfo Server, ttl int) (chan<- struct{}, error) {
	var err error
	if strings.Split(serInfo.Address, ":")[0] == "" {
		return nil, errors.New("invalid ip address")
	}
	if r.cli, err = clientv3.New(clientv3.Config{
		Endpoints:   r.EtcdAddresses,
		Username:    r.User,
		Password:    r.Password,
		DialTimeout: time.Duration(r.DialTimeout) * time.Second,
	}); err != nil {
		return nil, err
	}
	r.srvInfo = serInfo
	r.srvTTL = ttl
	if err = r.register(); err != nil {
		return nil, err
	}
	r.closeCh = make(chan struct{})
	go r.keepAlice()
	return r.closeCh, nil
}

func (r *Register) register() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.DialTimeout)*time.Second)
	defer cancel()
	leaseResp, err := r.cli.Grant(ctx, int64(r.srvTTL))
	if err != nil {
		return err
	}
	r.leaseID = leaseResp.ID
	if r.keepAliveCh, err = r.cli.KeepAlive(context.Background(), r.leaseID); err != nil {
		return err
	}
	data, err := json.Marshal(r.srvInfo)
	if err != nil {
		return err
	}
	_, err = r.cli.Put(context.Background(), BuildRegisterPath(r.srvInfo), string(data), clientv3.WithLease(r.leaseID))
	fmt.Println(BuildRegisterPath(r.srvInfo))
	return err
}

func (r *Register) keepAlice() error {
	ticker := time.NewTicker(time.Duration(r.srvTTL) * time.Second)
	for {
		select {
		case <-r.closeCh:
			if _, err := r.cli.Revoke(context.Background(), r.leaseID); err != nil {
				r.logger.Fatalf("revoke failed, err: %s", err)
			}
			if err := r.unregister(); err != nil {
				r.logger.Fatalf("unregister failed, err: %s", err)
			}
			os.Exit(0)
		case res := <-r.keepAliveCh:
			if res == nil {
				if err := r.register(); err != nil {
					r.logger.Warnf("register failed, err: %s", err)
				}
			}
		case <-ticker.C:
			if r.keepAliveCh == nil {
				if err := r.register(); err != nil {
					r.logger.Warnf("register failed, err: %s", err)
				}
			}
		}
	}
}

func (r *Register) unregister() error {
	_, err := r.cli.Delete(context.Background(), BuildRegisterPath(r.srvInfo))
	return err
}
