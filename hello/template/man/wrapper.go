package man

import "context"

type Wrapper func(Man) Man

type WorkFunc func(ctx context.Context, req string, msg string, opt WorkOptions) string

type WorkWrapper func(WorkFunc) WorkFunc
