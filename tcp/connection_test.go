package tcp

import (
	. "github.com/onsi/gomega"
	"io"
	"net"
	"testing"
)

const address = "127.0.0.1:3355"
const protocol = "tcp"

var data = []byte("test")

func TestNewConnection(t *testing.T) {
	g := NewGomegaWithT(t)
	listener := mockTcpListener(g)
	defer closeConn(g, listener)
	go func() {
		_, err := listener.Accept()
		g.Expect(err).ShouldNot(HaveOccurred())
	}()
	tcpConn := fakeNetConn(g)

	sut := newConn(data, tcpConn)

	g.Expect(sut.connection).Should(BeEquivalentTo(tcpConn))
	g.Expect(sut.message).Should(BeEquivalentTo([]byte("test")))
	g.Expect(sut.Message()).Should(BeEquivalentTo([]byte("test")))
}

func TestWrite(t *testing.T) {
	g := NewGomegaWithT(t)
	listener := mockTcpListener(g)
	defer closeConn(g, listener)
	connConsuming := make(chan net.Conn)
	go func() {
		c, err := listener.Accept()
		g.Expect(err).ShouldNot(HaveOccurred())
		connConsuming <- c
	}()
	tcpConn := fakeNetConn(g)
	sut := newConn(data, tcpConn)
	response := make(chan string)
	go listen(g, <-connConsuming, response)

	err := sut.Write(data)

	g.Expect(err).Should(BeNil())
	g.Eventually(response).Should(Receive(BeEquivalentTo(string(data))))
}

func TestClose(t *testing.T) {
	g := NewGomegaWithT(t)
	listener := mockTcpListener(g)
	go func() {
		_, err := listener.Accept()
		g.Expect(err).ShouldNot(HaveOccurred())
	}()
	tcpConn := fakeNetConn(g)

	sut := newConn(data, tcpConn)

	g.Expect(sut.Close()).ShouldNot(HaveOccurred())
	g.Expect(sut.Close()).Should(HaveOccurred())
}

//
//func TestConnectionType(t *testing.T) {
//	RegisterTestingT(t)
//	address := "127.0.0.1:9054"
//	listener := mockTcpListener(address)
//	defer listener.Close()
//	ok := make(chan bool)
//	go func() {
//		_, err := listener.Accept()
//		Expect(err).ShouldNot(HaveOccurred())
//		ok <- true
//	}()
//	conn := fakeNetConn(address)
//	<-ok
//	addressConnector := connection.NewAddressConnector(conn.LocalAddr(), conn.RemoteAddr())
//	context := NewTcpContext([]byte("teste"), addressConnector, conn)
//	connectionType := context.ConnectionType()
//	Expect(connectionType).Should(BeEquivalentTo("tcp"))
//}

func closeConn(g *GomegaWithT, closer io.Closer) {
	g.Expect(closer.Close()).ShouldNot(HaveOccurred())
}

func mockTcpListener(g *GomegaWithT) *net.TCPListener {
	addr, err := net.ResolveTCPAddr(protocol, address)
	g.Expect(err).Should(BeNil())
	tcpConn, err := net.ListenTCP(protocol, addr)
	g.Expect(err).Should(BeNil())
	return tcpConn
}

func fakeNetConn(g *GomegaWithT) net.Conn {
	netConn, err := net.Dial(protocol, address)
	g.Expect(err).Should(BeNil())
	return netConn
}

func listen(g *GomegaWithT, netConn net.Conn, response chan string) {
	buff := make([]byte, 1024)
	readLen, err := netConn.Read(buff)
	g.Expect(err).ShouldNot(HaveOccurred())
	response <- string(buff[:readLen])
}
