package Zookeeper

import (
	"errors"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
)

// Registry 注册中心接口
type Registry interface {
	Register(key string, data []byte) error
	Unregister(key string) error
}

// ZkRegistry 实现了 Registry
type ZkRegistry struct {
	RootPath string
	Conn     *zk.Conn
}

func NewZkRegistry(rootPath string, conn *zk.Conn) *ZkRegistry {
	return &ZkRegistry{
		RootPath: fmt.Sprintf("/%s", rootPath),
		Conn:     conn,
	}
}

func (r *ZkRegistry) getAcl() []zk.ACL {
	return zk.WorldACL(zk.PermAll)
}

func (r *ZkRegistry) GetAbsPath(path string) string {
	return fmt.Sprintf("%s/%s", r.RootPath, path)
}

// Register 向 zookeeper 创建一个临时节点
func (r *ZkRegistry) Register(path string, data []byte) error {
	_, err := r.Conn.Create(r.GetAbsPath(path), data, 3, r.getAcl())
	return err
}

// Unregister 显式删除一个临时节点
func (r *ZkRegistry) Unregister(path string) error {
	path = r.GetAbsPath(path)
	_, stat, err := r.Conn.Get(path)
	if err != nil {
		if errors.Is(err, zk.ErrNodeExists) {
			return nil
		}
		return err
	}
	return r.Conn.Delete(path, stat.Version)
}
