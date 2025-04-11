# GoFlint 🔥

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

import github.com/aagumin/goflint

func main() {
   base := goflint.NewSparkSubmit().
      Application("job.jar").
      WithMaster("k8s://http://k8s-master:443").
      WithName("GoFlint").
      Build()
   
   // base <- SparkApp{Kill, Status, Submit}
   
   sparkCfg := goflint.SparkConf{}
   sparkCfg.Set("spark.driver.port", "grpc")
   
   app := goflint.NewSparkSubmit(base).
      WithSparkConf(sparkCfg)
      NumExecutors(5).
      DriverMemory("16Gi").
      Build()

   ctx := context.Background()

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