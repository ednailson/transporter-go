package udp

import (
	. "github.com/onsi/gomega"
	"net"
	"testing"
)

const address = "127.0.0.1:1099"
const protocol = "udp"

var message = []byte("test")

func TestNewConnection(t *testing.T) {
	g := NewGomegaWithT(t)
	udpConn := mockUdpConn(g)
	defer closeConn(g, udpConn)

	sut := newConn(message, udpConn, fakeRemoteAddr(g), protocol)

	g.Expect(sut.udpConn).Should(BeEquivalentTo(udpConn))
	g.Expect(sut.message).Should(BeEquivalentTo(message))
	g.Expect(sut.Message()).Should(BeEquivalentTo(message))
	g.Expect(sut.remoteAddr).Should(BeEquivalentTo(fakeRemoteAddr(g)))
	g.Expect(sut.protocol).Should(BeEquivalentTo(protocol))
}

func TestWrite(t *testing.T) {
	g := NewGomegaWithT(t)
	udpConn := mockUdpConn(g)
	defer closeConn(g, udpConn)
	response := make(chan string)
	go listen(g, udpConn, response)
	sut := newConn(message, udpConn, fakeRemoteAddr(g), protocol)

	err := sut.Write(message)

	g.Expect(err).ShouldNot(HaveOccurred())
	g.Eventually(response).Should(Receive(BeEquivalentTo(string(message))))
}

func TestClose(t *testing.T) {
	g := NewGomegaWithT(t)
	udpConn := mockUdpConn(g)
	sut := newConn(message, udpConn, fakeRemoteAddr(g), protocol)

	err := sut.Close()

	g.Expect(err).ShouldNot(HaveOccurred())
}

func fakeRemoteAddr(g *GomegaWithT) *net.UDPAddr {
	addr, err := net.ResolveUDPAddr("udp", address)
	g.Expect(err).Should(BeNil())
	return addr
}

func mockUdpConn(g *GomegaWithT) *net.UDPConn {
	addr, err := net.ResolveUDPAddr("udp", address)
	g.Expect(err).Should(BeNil())
	udpConn, err := net.ListenUDP("udp", addr)
	g.Expect(err).Should(BeNil())
	return udpConn
}

func listen(g *GomegaWithT, ServerConn *net.UDPConn, response chan string) {
	buff := make([]byte, 1024)
	readLen, _, err := ServerConn.ReadFromUDP(buff)
	g.Expect(err).ShouldNot(HaveOccurred())
	response <- string(buff[:readLen])
}
