package dns

import (
	"context"
	"fmt"
	"testing"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/stretchr/testify/assert"
)

func TestDNSResolver(t *testing.T) {
	cr := NewDNSResolver()

	target := cr.Target(context.Background(), rpcinfo.NewEndpointInfo("example.com:8888", "", nil, nil))
	r, err := cr.Resolve(context.Background(), target)
	assert.Nil(t, err)
	for _, ins := range r.Instances {
		fmt.Printf("result %s\n", ins.Address())
	}
}
