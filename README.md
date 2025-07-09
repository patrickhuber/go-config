# Go Config

configuration library in go

# Get it

```bash
go get github.com/patrickhuber/go-config
```

# Example

> config.yml
```yaml
yaml: yes
```

> config.json
```json
{"json": "yes"}
```

> config.toml
```toml
toml="yes"
```

> .env
```
dotenv="yes"
```

```go
package main
import (
    "fmt"
    "log"
    "github.com/patrickhuber/go-config"
    "github.com/patrickhuber/go-cross"
    "github.com/patrickhuber/go-cross/arch"
    "github.com/patrickhuber/go-cross/platform"
)
func main(){
    args := []string{"--hello", "world"}

    // Create a target for cross-platform abstractions
    target := cross.NewTest(platform.Linux, arch.AMD64)
    filesystem := target.FS()

    // Use environment provider from target and set environment variable
    osEnv := target.Env()
    osEnv.Set("env", "yes")

    builder := config.NewBuilder(
        config.NewYaml(filesystem, "config.yml"),
        config.NewJson(filesystem, "config.json"),
        config.NewToml(filesystem, "config.toml"),
        config.NewEnv(osEnv, config.EnvOption{Prefix: "env"}),
        config.NewDotEnv(filesystem, ".env"),
        config.NewFlag([]config.Flag{
            &config.StringFlag{
                Name: "hello",
            },
        }, args),
    )
    root := builder.Build()
    cfg, err := root.Get(&config.GetContext{})
    if err != nil {
        log.Fatal(err)
    } else {
        log.Printf("%v", cfg)
    }
}
```

> output
```
yaml: yes
json: yes
toml: yes
dotenv: yes
env: yes
```