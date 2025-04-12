package domain

import "gorm.io/gorm"

type BucketType string

const (
	PublicBucketType  BucketType = "PUBLIC"
	PrivateBucketType BucketType = "PRIVATE"
)

type File struct {
	gorm.Model
	Name       string
	Key        string
	BucketType BucketType
}
