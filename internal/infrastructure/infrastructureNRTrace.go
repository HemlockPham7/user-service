package infrastructure

import (
	"github.com/HemlockPham7/common-libs/pkg/common"
	"github.com/HemlockPham7/common-libs/pkg/nrtrace"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func CreateNRClient(envPrefix string) *newrelic.Application {
	nrClient, err := nrtrace.NewClient(envPrefix)
	common.HandleError(err)
	return nrClient
}
