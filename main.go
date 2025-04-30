package main

import (
	"context"
	"fmt"
	goflint "goflint/pkg"
	sconf "goflint/pkg/sparkconf"
	"log"
)

func main() {
	//base <- SparkApp{Kill, Status, Submit}
	xx := map[string]string{"spark.driver.port": "grpc",
		"spark.driver.host": "localhost"}
	sparkCfg := sconf.NewFrozenConf(xx)

	base := goflint.CrateOrUpdate(nil).
		Application("job.jar").
		WithSparkConf(sparkCfg).
		WithMaster(nil).
		WithName("GoFlint").
		Build()

	fmt.Println(sparkCfg.Get("spark.driver.port"))
	fmt.Println(base)

	app := goflint.CrateOrUpdate(&base).
		WithSparkConf(sparkCfg).
		Build()

	ctx := context.Background()

	if err := app.Submit(ctx); err != nil {
		log.Fatal(err)
	}
}
