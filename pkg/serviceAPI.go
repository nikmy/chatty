package chatty

import detail "github.com/nikmy/chatty/internal"

func Init() error {
	if err := detail.WithRedis().Init(); err != nil {
		return err
	}
	return detail.WithKafka().Init()
}

func Finalize() error {
	if err := detail.WithRedis().Finalize(); err != nil {
		return err
	}
	return detail.WithKafka().Finalize()
}
