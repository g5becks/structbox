package structbox

import (
	"fmt"

	bolt "go.etcd.io/bbolt"
)

type StructBox struct {
	db *bolt.DB
}

func New(db *bolt.DB) *StructBox {
	return &StructBox{db: db}
}

type Collection[T any] struct {
	Name string
	sb   *StructBox
}

func bucketName[T interface{}]() []byte {
	var ref T
	return []byte(fmt.Sprintf("structbox:%T", &ref))
}

var (
	ErrCollectionNotFound = fmt.Errorf("collection not found")
	ErrCreatingCollection = fmt.Errorf("error creating collection")
)

func NewCollection[T interface{}](sb *StructBox) (*Collection[T], error) {
	path := bucketName[T]()
	err := sb.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(path)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrCreatingCollection, err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCreatingCollection, err)
	}
	return &Collection[T]{sb: sb, Name: string(path)}, nil
}
