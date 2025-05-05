# GoFlint 🔥

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=aagumin_goflint&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=aagumin_goflint)
[![Go Report Card](https://goreportcard.com/badge/github.com/aagumin/goflint)](https://goreportcard.com/report/github.com/aagumin/goflint)
[![Maintainability](https://api.codeclimate.com/v1/badges/aabaf6a2d52511b3d581/maintainability)](https://codeclimate.com/github/aagumin/goflint/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/aabaf6a2d52511b3d581/test_coverage)](https://codeclimate.com/github/aagumin/goflint/test_coverage)
[![Known Vulnerabilities](https://snyk.io/org/github/aagumin/goflint/badge.svg)](https://snyk.io/org/github/aagumin/goflint)

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
go get github.com/aagumin/flint
```

## Example usage

```go
package main

import (
   "context"
   "fmt"
   "os"
   "path"
   "github.com/aagumin/goflint/flint"
   sc "github.com/aagumin/goflint/flint/sparkconf"
)

func main() {
   xx := map[string]string{"spark.driver.port": "4031", "spark.driver.host": "localhost"}
   sparkCfg := sc.NewFrozenConf(xx)

   scalaExamples := path.Join(os.Getenv("SPARK_HOME"), "examples/jars/spark-examples_2.12-3.5.3.jar")

   submit := flint.NewSparkApp(
      flint.WithApplication(scalaExamples),
      flint.WithSparkConf(sparkCfg),
      flint.WithName("GoFlint"),
      flint.WithMainClass("org.apache.spark.examples.parkPi"),
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
