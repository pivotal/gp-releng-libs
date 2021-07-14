package retry_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/pivotal/gp-releng-libs/retry"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRetry(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "retry tests")
}

var _ = Describe("Retry", func() {
	It("Test Retry", func() {
		errors := make([]error, 0)
		f := func() error {
			e := fmt.Errorf("test error")
			errors = append(errors, e)
			return e
		}

		err := retry.Retry(f, 5, 100*time.Nanosecond)
		Expect(err).To(HaveOccurred())
		Expect(len(errors)).To(Equal(5))
	})

	It("be able to retry even under panic", func() {
		funcCalled := 0
		f := func() error {
			funcCalled++
			if funcCalled == 1 {
				panic("function panic!")
			}
			return nil
		}

		err := retry.Retry(f, 3, 10*time.Nanosecond)
		Expect(err).ToNot(HaveOccurred())
		// expect the function called twice, first time hit panic, second time passed successfully
		Expect(funcCalled).To(Equal(2))
	})
})
