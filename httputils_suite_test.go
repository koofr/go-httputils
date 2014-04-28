package httputils_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestHttputils(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Httputils Suite")
}
