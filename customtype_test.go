package capnp_test

import (
	"fmt"
	"net"
	"testing"

	cv "github.com/smartystreets/goconvey/convey"
	"zombiezen.com/go/capnproto"
	air "zombiezen.com/go/capnproto/internal/aircraftlib"
)

func Example_customType() (*capnp.Segment, []byte) {
	msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		panic(err)
	}
	e, err := air.NewRootEndpoint(seg)
	if err != nil {
		panic(err)
	}
	e.SetIp(net.ParseIP("1.2.3.4").To4())
	e.SetPort(56)
	e.SetHostname("test.com")

	ip, err := e.Ip()
	if err != nil {
		panic(err)
	}
	fmt.Printf("ip: %s\n", ip.String())
	fmt.Printf("port: %d\n", e.Port())
	hostname, err := e.Hostname()
	if err != nil {
		panic(err)
	}
	fmt.Printf("hostname: %s\n", hostname)

	buf, err := msg.Marshal()
	if err != nil {
		panic(err)
	}
	return seg, buf
}

func TestCreationOfEndpoint(t *testing.T) {
	seg, _ := Example_customType()
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
			endpoint, err := air.ReadRootEndpoint(seg.Message())
			cv.So(err, cv.ShouldEqual, nil)
			cv.Convey(fmt.Sprintf("Then we should get the expected ip '%s'", expectedIP), func() {
				ip, err := endpoint.Ip()
				cv.So(err, cv.ShouldEqual, nil)
				cv.So(ip, cv.ShouldResemble, expectedIP)
			})
			cv.Convey(fmt.Sprintf("Then we should get the expected port '%d'", expectedPort), func() {
				cv.So(endpoint.Port(), cv.ShouldEqual, expectedPort)
			})
			cv.Convey(fmt.Sprintf("Then we should get the expected hostname '%s'", expectedHostname), func() {
				hostname, err := endpoint.Hostname()
				cv.So(err, cv.ShouldEqual, nil)
				cv.So(hostname, cv.ShouldEqual, expectedHostname)
			})
		})
	})
}
