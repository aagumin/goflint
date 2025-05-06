package main

import (
	"context"
	"fmt"
	flint "github.com/aagumin/goflint/flint"
	sc "github.com/aagumin/goflint/flint/sparkconf"
	"log"
	"os"
)

type XSparkApp struct {
	app flint.SparkApp
}

func NewXSparkApp(app flint.SparkApp) *XSparkApp {
	return &XSparkApp{app}
}

func (s *XSparkApp) Submit(ctx context.Context, c chan []byte) error {
	out, err := s.app.Submit(ctx)
	if err != nil {
		fmt.Println(err)
	}
	c <- out
	return nil
}

func main() {
	xx := map[string]string{"spark.driver.port": "4031", "spark.driver.host": "localhost"}
	sparkCfg := sc.NewFrozenConf(xx)
	err := os.Setenv("SPARK_HOME", "/Users/21370371/trash/spark-3.5.3-bin-hadoop3")
	if err != nil {
		log.Fatal(err)
	}

	submit := flint.NewSparkApp(
		flint.WithApplication("examples/python/freeze.py"),
		flint.WithSparkConf(sparkCfg),
		flint.WithName("GoFlint"),
	)

	base := submit.Build()

	updatedSubmit := flint.ExtendSparkApp(
		&base,
		flint.WithMaster(""),
		// Other options...

	)

	app := updatedSubmit.Build()
	//shortDuration := 5
	//d := time.Now().Add(time.Duration(shortDuration) * time.Second)
	//ctx, _ := context.WithDeadline(context.Background(), d)
	ctx := context.Background()

	xapp := NewXSparkApp(app)
	cc := make(chan []byte)
	err = xapp.Submit(context.Background(), cc)
	if err != nil {
		return
	}
	zzz := <-cc
	println(string(zzz))
	_, err = app.Submit(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(app.Status(ctx))

}
