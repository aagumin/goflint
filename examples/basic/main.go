package main

import (
	"context"
	"fmt"
	flint "github.com/aagumin/goflint/flint"
	sc "github.com/aagumin/goflint/flint/sparkconf"
	"os"
	"path"
)

func main() {
	xx := map[string]string{"spark.driver.port": "4031", "spark.driver.host": "localhost"}
	sparkCfg := sc.NewFrozenConf(xx)
	v := os.Getenv("SPARK_HOME")
	scalaExamples := path.Join(v, "examples/jars/spark-examples_2.12-3.5.3.jar")

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
	_, err := app.Submit(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(app.Status(ctx))

}
