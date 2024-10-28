package public

import (
	"context"
	"github.com/loveyu233/gu/ctx"
)

type UseConfig struct {
	Ctx   context.Context
	Debug bool
}

var DefaultUseConfig = UseConfig{
	Ctx:   ctx.Timeout(),
	Debug: true,
}

const (
	JWTSECRET = "yFv034etfGlidKMwQzK9-Y8Y3ajaBs1Qu0L2SQzVj"
)
