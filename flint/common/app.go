package common

import (
	"context"
)

//
//┌────────────────────────┐
//│    interface Submitter │◄────────┐
//└────────────────────────┘         │
//┌────────────────────────┐         │
//│    interface Monitor   │◄──┐     │
//└────────────────────────┘   │     │
//							   │     │
//							   ▼     ▼
//						 ┌──────────────┐
//						 │  SparkApp    │
//						 │(immutable)   │
//						 └──────────────┘
//							 ▲
//                .Build()   │
//					┌───────────────────┐
//					│ SparkSubmit    	│
//					│ (builder)         │
//					└───────────────────┘

type Submitter interface {
	Submit(ctx context.Context) ([]byte, error)
}

type Monitor interface {
	Status(ctx context.Context) (string, error)
	Kill(ctx context.Context) error
}

type SparkConf interface {
	Get(key string) string
	GetAll() map[string]string
	ToCommandLineArgs() []string
	Merge(other SparkConf) SparkConf
}

type SparkConfHelpers interface {
	SparkConf
	IsEmpty() bool
	Repr() string
}
