# Kitex DNS Resolver

Some application runtime use DNS as service discovery, e.g. Kubernetes.

## How to use with Kitex client?

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

Use Kitex `client.WithResolver` function optional, we can set DNS resolver with our client.
