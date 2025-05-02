# GoFlint 🔥

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=aagumin_goflint&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=aagumin_goflint)

> **[WIP]** Project in active dev state.

**GoFlint** is an idiomatic Go library for submitting Apache Spark jobs via `spark-submit`.  
Designed for clarity and type safety, it wraps Spark’s CLI with a fluent Go API.


## Features

- **Fluent API** for all `spark-submit` options (masters, deploy modes, args, etc.)
- **Context support** for cancellation/timeouts
- **Extensible** with custom logging and monitoring
- **Zero dependencies** (except Go’s standard library)

## Installation

```bash
go get github.com/aagumin/goflint
```

## Example usage

```go
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

   scalaExamples := path.Join(os.Getenv("SPARK_HOME"), "spark-examples.jar")

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
      // Other options...
   )

   app := updatedSubmit.Build()
   ctx := context.TODO()
   if err := app.Submit(ctx); err != nil {
      log.Fatal(err)
   }
}
```

## Design Principles

1. **Idiomatic Go**
    - Errors as `error`, not panics
    - `context.Context` support
    - Interfaces for extensibility

2. **Spark Compatibility**
    - Maps 1:1 with `spark-submit` CLI
    - No hidden magic – transparent command generation

3. **Batteries Included**
    - Defaults for quick starts
    - Extensible for complex cases

---

## Roadmap

- [ ] Async job monitoring
- [ ] YARN/K8s auth helpers
- [ ] Prometheus metrics integration

---

## Contributing

PRs welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

MIT © Arsen Gumin
