package observer_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestObserver(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Observer Suite")
}
