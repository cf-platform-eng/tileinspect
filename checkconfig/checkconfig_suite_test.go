package checkconfig_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCheckConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Check Config Suite")
}
