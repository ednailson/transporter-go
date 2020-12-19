package tcp

import (
	"fmt"
	"github.com/ednailson/transporter-go/connection"
	. "github.com/onsi/gomega"
	"net"
	"testing"
)

const ip = "127.0.0.1"
const port = 4499

var host = fmt.Sprintf("%s:%d", ip, port)

func TestNewServer(t *testing.T) {
	g := NewGomegaWithT(t)
	sut, err := New(protocol, ip, port)
	g.Expect(err).ShouldNot(HaveOccurred())
	defer closeConn(g, sut)
	g.Expect(sut.listener).ShouldNot(BeNil())
}

func TestNewServerInvalidIP(t *testing.T) {
	g := NewGomegaWithT(t)
	sut, err := New(protocol, "-@", port)
	g.Expect(err).Should(HaveOccurred())
	g.Expect(sut).Should(BeNil())
}

func TestNewServerAtABusyPort(t *testing.T) {
	g := NewGomegaWithT(t)
	sut, err := New(protocol, ip, port)
	g.Expect(err).ShouldNot(HaveOccurred())
	defer closeConn(g, sut)
	g.Expect(sut.listener).ShouldNot(BeNil())
	sut.Listen(func(conn connection.Connection) {})

	tcpServer, err := New(protocol, ip, port)
	g.Expect(err).Should(HaveOccurred())
	g.Expect(tcpServer).Should(BeNil())
}

func TestListen(t *testing.T) {
	g := NewGomegaWithT(t)
	sut, err := New(protocol, ip, port)
	g.Expect(err).Should(BeNil())
	defer closeConn(g, sut)
	msg := make(chan []byte)
	chErr := sut.Listen(func(conn connection.Connection) {
		msg <- conn.Message()
	})
	netConn := fakeMockConn(g)
	defer closeConn(g, netConn)
	content := "data test"

	dataSender(g, content, netConn)

	g.Eventually(msg).Should(Receive(BeEquivalentTo(content)))
	g.Expect(chErr).ShouldNot(Receive())
}

func TestCloseServer(t *testing.T) {
	g := NewGomegaWithT(t)
	sut, err := New(protocol, ip, port)
	g.Expect(err).ShouldNot(HaveOccurred())

	err = sut.Close()
	g.Expect(err).ShouldNot(HaveOccurred())

	err = sut.Close()
	g.Expect(err).Should(HaveOccurred())
}

func fakeMockConn(g *GomegaWithT) net.Conn {
	netConn, err := net.Dial(protocol, host)
	g.Expect(err).Should(BeNil())
	return netConn
}

func dataSender(g *GomegaWithT, content string, conn net.Conn) {
	sendLen, sendErr := conn.Write([]byte(content))
	g.Expect(sendErr).ShouldNot(HaveOccurred())
	g.Expect(sendLen).To(BeEquivalentTo(len(content)))
}
