package store

import (
	"context"
	"ws/model"
)

type DataStore interface {
	StoreFileMetaData(ctx context.Context, metadata model.Metadata) error
}
