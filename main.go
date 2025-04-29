package main

import (
	"fmt"
	goflint "goflint/pkg"
)

func main() {
	// base <- SparkApp{Kill, Status, Submit}
	//xx := map[string]string{"spark.driver.port": "grpc"}
	//sparkCfg := sconf.NewFrozenConf(xx)

	base := goflint.CrateOrUpdate().
		Application("job.jar").
		//WithSparkConf(sparkCfg).
		WithMaster(nil).
		WithName("GoFlint").
		Build()

	//sparkCfg.Get("spark.driver.port")
	//
	//app := goflint.CrateOrUpdate(&base).
	//	WithSparkConf(sparkCfg).
	//	Build()
	fmt.Println(base)
	//ctx := context.Background()
	//
	//if err := base.Submit(ctx); err != nil {
	//	log.Fatal(err)
	//}
}
