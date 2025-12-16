package config

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type BaseURLer interface {
	BaseURL() string
}

type baseURLer struct {
	getter kv.Getter
	once   comfig.Once
}

type baseURLConfig struct {
	BaseURL string
}

func NewBaseURLer(getter kv.Getter) BaseURLer {
	return &baseURLer{
		getter: getter,
	}
}

func (b *baseURLer) BaseURL() string {
	return b.once.Do(func() interface{} {
		var config baseURLConfig

		raw := kv.MustGetStringMap(b.getter, "short_url")

		err := figure.Out(&config).From(raw).Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to get base url"))
		}

		return config.BaseURL
	}).(string)
}
