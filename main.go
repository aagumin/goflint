package main

import (
	"context"
	"fmt"
	goflint "goflint/pkg"
	sconf "goflint/pkg/sparkconf"
	"log"
)

func main() {
	base := goflint.CrateOrUpdate(nil).
		Application("job.jar").
		WithMaster(nil).
		WithName("GoFlint").
		Build()

	// base <- SparkApp{Kill, Status, Submit}

	sparkCfg := sconf.SparkConf{}
	sparkCfg.Set("spark.driver.port", "grpc")

	app := goflint.CrateOrUpdate(&base).
		WithSparkConf(&sparkCfg).
		Build()
	fmt.Println(app)
	ctx := context.Background()

	if err := app.Submit(ctx); err != nil {
		log.Fatal(err)
	}
}
