package render

import (
	"context"

	"github.com/hexthink/seed"
)

type ResponseRender func(ctx context.Context, data interface{}, err error) seed.Response
