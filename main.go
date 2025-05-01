package main

import (
	"context"
	goflint "goflint/pkg"
	sconf "goflint/pkg/sparkconf"
	"log"
	"os"
	"path"
)

func main() {
	xx := map[string]string{"spark.driver.port": "4031", "spark.driver.host": "localhost"}
	sparkCfg := sconf.NewFrozenConf(xx)
	err := os.Setenv("SPARK_HOME", "/Users/21370371/trash/spark-3.5.3-bin-hadoop3")
	if err != nil {
		log.Fatal(err)
	}
	scalaExamples := path.Join(os.Getenv("SPARK_HOME"), "examples/jars/spark-examples_2.12-3.5.3.jar")

	submit := goflint.NewSparkApp(
		goflint.WithApplication(scalaExamples),
		goflint.WithSparkConf(sparkCfg),
		goflint.WithName("GoFlint"),
		goflint.WithMainClass("org.apache.spark.examples.SparkPi"),
	)

	base := submit.Build()

	updatedSubmit := goflint.ExtendSparkApp(
		&base,
		goflint.WithMaster(""),
		// Другие опции...
	)

	app := updatedSubmit.Build()
	ctx := context.TODO()
	if err := app.Submit(ctx); err != nil {
		log.Fatal(err)
	}
}
