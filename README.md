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
json: yes
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
    "log"
    "github.com/patrickhuber/go-config"
)
func main(){
    args := []string{"--hello", "world"}
    os.SetEnv("env", "yes")
    builder := config.NewBuilder(
        config.NewYaml("config.yml"),
        config.NewJson("config.json"),
        config.NewToml("config.toml"),
        config.NewEnv("env"),
        config.NewDotEnv(".env"),
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
yaml: yes
json: yes
toml: yes
dotenv: yes
env: yes
```