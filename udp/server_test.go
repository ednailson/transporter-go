package udp

import (
	"fmt"
	"github.com/ednailson/transporter-go/connection"
	. "github.com/onsi/gomega"
	"io"
	"net"
	"testing"
)

const port = 1155
const ip = "127.0.0.1"

var host = fmt.Sprintf("%s:%d", ip, port)

func TestNewServer(t *testing.T) {
	g := NewGomegaWithT(t)

	sut, err := New(protocol, ip, port)

	g.Expect(err).ShouldNot(HaveOccurred())
	defer closeConn(g, sut)
	g.Expect(sut.conn.LocalAddr().String()).Should(BeEquivalentTo(host))
	g.Expect(sut.protocol).Should(BeEquivalentTo(protocol))
}

func TestNewServerAtABusyPort(t *testing.T) {
	g := NewGomegaWithT(t)
	udpServer, err := New(protocol, ip, port)
	g.Expect(err).ShouldNot(HaveOccurred())
	defer closeConn(g, udpServer)

	sut, err := New(protocol, ip, port)

	g.Expect(err).Should(HaveOccurred())
	g.Expect(sut).Should(BeNil())

}

func TestNewServerInvalidIp(t *testing.T) {
	g := NewGomegaWithT(t)

	sut, err := New(protocol, "-@", port)

	g.Expect(err).Should(HaveOccurred())
	g.Expect(sut).Should(BeNil())
}

func TestListen(t *testing.T) {
	g := NewGomegaWithT(t)
	sut, err := New(protocol, ip, port)
	g.Expect(err).ShouldNot(HaveOccurred())
	defer closeConn(g, sut)
	msg := make(chan []byte)
	chErr := sut.Listen(func(conn connection.Connection) {
		msg <- conn.Message()
	})
	udpConn := mockNetConn(g, host)
	defer closeConn(g, udpConn)
	data := "data test"

	udpDataSender(g, data, udpConn)

	g.Eventually(msg).Should(Receive(BeEquivalentTo([]byte(data))))
	g.Expect(chErr).ShouldNot(Receive())
}

func TestCloseConnection(t *testing.T) {
	g := NewGomegaWithT(t)
	sut, err := New(protocol, ip, port)
	g.Expect(err).ShouldNot(HaveOccurred())

	err = sut.Close()
	g.Expect(err).ShouldNot(HaveOccurred())

	err = sut.Close()
	g.Expect(err).Should(HaveOccurred())
}

func closeConn(g *GomegaWithT, closer io.Closer) {
	g.Expect(closer.Close()).ShouldNot(HaveOccurred())
}

func mockNetConn(g *GomegaWithT, address string) net.Conn {
	udpConn, err := net.Dial("udp", address)
	g.Expect(err).Should(BeNil())
	return udpConn
}

func udpDataSender(g *GomegaWithT, content string, conn net.Conn) {
	sendLen, sendErr := conn.Write([]byte(content))
	g.Expect(sendErr).ShouldNot(HaveOccurred())
	g.Expect(sendLen).To(BeEquivalentTo(len(content)))
}
