package storage

import "context"

type BannerStorage interface {
	Close(ctx context.Context) error
}
