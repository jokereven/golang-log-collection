package model

type TransferConfig struct {
	KafkaConfig `ini:"kafka"`
	EsConfig    `ini:"es"`
}

type KafkaConfig struct {
	Address string `ini:"address"`
	Topic   string `ini:"topic"`
}

type EsConfig struct {
	Address      string `ini:"address"`
	Index        string `ini:"index"`
	MaxChanSize  int64  `ini:"max_chan_size"`
	GoroutineNum int    `ini:"goroutine_num"`
}
