package render

import (
	"context"

	"github.com/ninthsoft/seed"
)

type ResponseRender func(ctx context.Context, data interface{}, err error) seed.Response
