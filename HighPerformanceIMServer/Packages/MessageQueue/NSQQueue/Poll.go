package NSQQueue

import (
	"HighPerformanceIMServer/Configs"
	"github.com/nsqio/go-nsq"
	"github.com/pkg/errors"
	"github.com/silenceper/pool"
	"go.uber.org/zap"
	"time"
)

var (
	NSQProducerPool pool.Pool
)

// InitNewNSQProducerPoll 定义了一个名为InitNewNSQProducerPoll的函数，用于初始化一个生产者池
// 函数返回一个错误类型，用于指示函数执行过程中是否发生了错误
func InitNewNSQProducerPoll() error {
	// 定义一个函数factory，用于创建一个NSQ生产者并返回
	factory := func() (interface{}, error) {
		// 创建一个NSQ生产者，使用配置文件中的NsqHost作为地址，配置对象使用nsq.NewConfig()
		producer, err := nsq.NewProducer(Configs.ConfigData.Nsq.NsqHost, nsq.NewConfig())
		if err != nil {
			return nil, err
		} else {
			return producer, nil
		}
	}

	// 定义一个函数closeError，用于关闭生产者
	closeError := func(v interface{}) error {
		// 将v转换为*nsq.Producer类型，并调用Stop()方法关闭生产者
		v.(*nsq.Producer).Stop()
		return nil
	}

	// 创建一个生产者池的配置对象poolConfig
	poolConfig := &pool.Config{
		InitialCap:  20,              // 初始容量为20
		MaxIdle:     40,              // 最大空闲数为40
		MaxCap:      50,              // 最大容量为50
		Factory:     factory,         // 使用factory函数创建对象
		Close:       closeError,      // 使用closeError函数关闭对象
		IdleTimeout: 0 * time.Second, // 空闲超时时间为0秒
	}
	var err error
	// 创建一个通道池NSQProducerPool，使用配置对象poolConfig
	NSQProducerPool, err = pool.NewChannelPool(poolConfig)
	if err != nil {
		return errors.New("NewChannelPool init failed")
	} else {
		return err
	}
}

// PublishMessage 定义了一个名为PublishMessage的函数，用于发布消息到指定主题
// 函数接收两个参数：topic为要发布消息的主题，content为消息内容
// 函数返回一个错误类型，用于指示函数执行过程中是否发生了错误
func PublishMessage(topic string, content []byte) error {
	// 从NSQ生产者池中获取一个生产者对象
	NSQProducer, err := NSQProducerPool.Get()
	if err != nil {
		return err
	}
	// 在函数结束时将生产者对象放回NSQ生产者池中
	defer func(NSQProducerPool pool.Pool, i interface{}) {
		err := NSQProducerPool.Put(i)
		if err != nil {
			zap.S().Infoln(err)
			panic(err)
		}
	}(NSQProducerPool, NSQProducer)
	// 将生产者对象转换为*nsq.Producer类型，并调用Publish方法发布消息到指定主题
	err = NSQProducer.(*nsq.Producer).Publish(topic, content)
	if err != nil {
		return err
	} else {
		return nil
	}
}

// Exit 定义了一个名为Exit的函数，用于退出程序
func Exit() {
	// 释放NSQ生产者池
	NSQProducerPool.Release()
}
