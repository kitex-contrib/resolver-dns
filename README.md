# KiteX DNS Resolver

For some application runtime, maybe we need use DNS as service discovery.
Here is a simple DNS resolver inspired by gRPC and Kubernetes.

## How to use with KiteX client?

```go
import (
    ...
    dns "github.com/kitex-contrib/dns-resolver"
    "github.com/cloudwego/kitex/client"
    ...
)

func main() {
    ...
    client, err := echo.NewClient("echo", client.WithResolver(dns.NewDNSResolver()))
	if err != nil {
		log.Fatal(err)
	}
    ...
}
```

Use KiteX `client.WithResolver` function optional, we can set DNS resolver with our client.

That's all, easy to use.
