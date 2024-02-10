package tests

import (
	"genericAPI/internal/utils/authentication_utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestBooks(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Authentication")
}

var _ = Describe("Generating password hash", Label("auth"), passwordHashTest)

func passwordHashTest() {
	authentication_utils.HashPassword("1234")
	var password string
	BeforeEach(func() {
		password = strconv.Itoa(rand.Int())
	})
	It("generated a password", func(ctx SpecContext) {
		hashedPw := authentication_utils.HashPassword(password)
		Expect(hashedPw).NotTo(Equal(""))
	}, SpecTimeout(time.Millisecond*5))
}
