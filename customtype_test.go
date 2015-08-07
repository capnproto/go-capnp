package capnp_test

import (
	"bytes"
	"fmt"
	"net"
	"testing"

	cv "github.com/smartystreets/goconvey/convey"
	"zombiezen.com/go/capnproto"
	air "zombiezen.com/go/capnproto/internal/aircraftlib"
)

func Example_createEndpoint() (*capnp.Segment, []byte) {
	seg := capnp.NewBuffer(nil)
	e := air.NewRootEndpoint(seg)
	e.SetIp(net.ParseIP("1.2.3.4").To4())
	e.SetPort(56)
	e.SetHostname("test.com")

	fmt.Printf("ip: %s\n", e.Ip().String())
	fmt.Printf("port: %d\n", e.Port())
	fmt.Printf("hostname: %s\n", e.Hostname())

	buf := bytes.Buffer{}
	seg.WriteTo(&buf)

	return seg, buf.Bytes()
}

func TestCreationOfEndpoint(t *testing.T) {
	seg, _ := Example_createEndpoint()
	text := CapnpDecodeSegment(seg, "", schemaPath, "Endpoint")

	expectedText := `(ip = "\x01\x02\x03\x04", port = 56, hostname = "test.com")`
	expectedIP := net.IP([]byte{1, 2, 3, 4})
	const expectedPort = 56
	expectedHostname := "test.com"

	cv.Convey("Given a go-capnproto created Endpoint", t, func() {
		cv.Convey("When we decode it with capnp", func() {
			cv.Convey(fmt.Sprintf("Then we should get the expected text '%s'", expectedText), func() {
				cv.So(text, cv.ShouldEqual, expectedText)
			})
		})
		cv.Convey("When we decode it", func() {
			endpoint := air.ReadRootEndpoint(seg)
			cv.Convey(fmt.Sprintf("Then we should get the expected ip '%s'", expectedIP), func() {
				cv.So(endpoint.Ip(), cv.ShouldResemble, expectedIP)
			})
			cv.Convey(fmt.Sprintf("Then we should get the expected port '%d'", expectedPort), func() {
				cv.So(endpoint.Port(), cv.ShouldEqual, expectedPort)
			})
			cv.Convey(fmt.Sprintf("Then we should get the expected hostname '%s'", expectedHostname), func() {
				cv.So(endpoint.Hostname(), cv.ShouldEqual, expectedHostname)
			})
		})
	})
}
