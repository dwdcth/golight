package golight

import (
	"fmt"

	"go.etcd.io/etcd/clientv3"
	"xdao.top/golight/srv"
)

type Container struct {
	services map[string]interface{}
}

var Di = NewContainer()

func NewContainer() *Container {
	return &Container{services: make(map[string]interface{})}
}

// func (container *Container) GetSrv(srvType interface{}, srvName string) (interface{}, error) {
// 	return nil, nil
// }

// 获取etcd客户端
func (container *Container) GetEtcd(etcdDsn string) (*clientv3.Client, error) {
	val, ok := container.services[fmt.Sprintf("etcd:%s", etcdDsn)]
	if !ok {
		cli, err := srv.NewEtcdClient(etcdDsn)
		if err != nil {
			container.services[fmt.Sprintf("etcd:%s", etcdDsn)] = err
			return nil, err
		}
		container.services[fmt.Sprintf("etcd:%s", etcdDsn)] = cli
		val = cli
	}
	//避免值为error反复调用
	if err, ok := val.(error); ok {
		return nil, err
	}
	return val.(*clientv3.Client), nil
}
