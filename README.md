# Go Config

configuration library in go

# Get it

```bash
go get github.com/patrickhuber/go-config
```

# Example

> config.yml
```
yaml: yes
```

> config.json
```
json: yes
```

> config.toml
```
toml="yes"
```

> .env
```
EXAMPLE_ENV="yes"
```

```go
package main
import (
    "log"
    "github.com/patrickhuber/go-config"
)
func main(){
    args := []string{"--hello", "world"}
    builder := config.NewBuilder(
        config.NewYaml("config.yml"),
        config.NewJson("config.json"),
        config.NewToml("config.toml"),
        config.NewEnv("EXAMPLE_"),
        config.NewFlag([]config.Flag{
            config.StringFlag{
                Name: "hello",
            },
        }, args),
    )    
    cfg, err := builder.Build()
    if err != nil{
        log.Fatal(err)
    }else{
        fmt.Println("%v", cfg)
    }
}
```

> output
```

```