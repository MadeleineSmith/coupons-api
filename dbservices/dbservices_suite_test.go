package dbservices_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDbservices(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dbservices Suite")
}
