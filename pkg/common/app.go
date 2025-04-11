package common

import "context"

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
//					│ SparkSubmitCmd    │
//					│ (builder)         │
//					└───────────────────┘

type Submitter interface {
	Submit(ctx context.Context) error
}

type Monitor interface {
	Status(ctx context.Context) (string, error)
	Kill(ctx context.Context) error
}

type SparkConf interface {
	Set(key, value string) *SparkConf
	Get(key string) *SparkConf
	GetAll() map[string]*SparkConf
	ToCommandLineArgs() []string
	SetIfMissing(key, value string) *SparkConf
}
