package transporter

import (
	. "github.com/onsi/gomega"
	"testing"
)

const ip = "127.0.0.1"
const port = 9066

func TestNewTransporter(t *testing.T) {
	g := NewGomegaWithT(t)
	server, err := New("tcp", ip, port)
	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(server).ShouldNot(BeNil())
	server, err = New("udp", ip, port)
	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(server).ShouldNot(BeNil())
	server, err = New("invalid-protocol", ip, port)
	g.Expect(err).Should(HaveOccurred())
	g.Expect(server).Should(BeNil())
}
