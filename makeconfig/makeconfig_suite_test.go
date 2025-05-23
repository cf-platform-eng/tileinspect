package makeconfig_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMakeConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Make Config Suite")
}
