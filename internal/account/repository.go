package account

import "context"

type AccountRepository interface {
	Register(ctx context.Context)
}
