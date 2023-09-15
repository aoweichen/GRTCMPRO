package NSQQueue

import (
	"github.com/nsqio/go-nsq"
	"go.uber.org/zap"
)

// NSQProducer 定义了一个全局变量NSQProducer，类型为*nsq.Producer
var (
	NSQProducer *nsq.Producer
)

// NewNSQProducer 定义了一个名为NewNSQProducer的函数，用于创建一个NSQ的生产者
// 函数的参数是一个字符串addr，表示NSQ的地址
// 函数返回一个错误类型，用于指示函数执行过程中是否发生了错误
func NewNSQProducer(addr string) error {
	// 创建一个NSQ的配置对象NSQConfig
	NSQConfig := nsq.NewConfig()
	var err error
	// 使用给定的地址和配置对象创建一个NSQ的生产者，并将其赋值给全局变量NSQProducer
	NSQProducer, err = nsq.NewProducer(addr, NSQConfig)
	if err != nil {
		// 如果创建生产者时发生错误，使用zap.S().Error打印错误信息，并返回错误
		zap.S().Errorln(err)
		return err
	} else {
		// 如果创建生产者成功，则返回nil表示没有错误发生
		return nil
	}
}
