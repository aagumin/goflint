package main

import (
	"context"
	"fmt"
	flint "github.com/aagumin/goflint/flint"
	sc "github.com/aagumin/goflint/flint/sparkconf"
	"log"
	"os"
	"path"
)

func main() {
	xx := map[string]string{"spark.driver.port": "4031", "spark.driver.host": "localhost"}
	sparkCfg := sc.NewFrozenConf(xx)
	err := os.Setenv("SPARK_HOME", "/Users/21370371/trash/spark-3.5.3-bin-hadoop3")
	if err != nil {
		log.Fatal(err)
	}
	scalaExamples := path.Join(os.Getenv("SPARK_HOME"), "examples/jars/spark-examples_2.12-3.5.3.jar")

	submit := flint.NewSparkApp(
		flint.WithApplication(scalaExamples),
		flint.WithSparkConf(sparkCfg),
		flint.WithName("GoFlint"),
		flint.WithMainClass("org.apache.spark.examples.SparkPi"),
	)

	base := submit.Build()

	updatedSubmit := flint.ExtendSparkApp(
		&base,
		flint.WithMaster(""),
		// Other options...

	)

	app := updatedSubmit.Build()
	ctx := context.Background()
	_, err = app.Submit(ctx)
	if err != nil {
		fmt.Println(err)
	}

}
