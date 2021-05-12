# Test interfaces for RPC tests.

using Go = import "/go.capnp";

@0xef12a34b9807e19c;
$Go.package("testcapnp");
$Go.import("capnproto.org/go/capnp/v3/rpc/internal/testcapnp");

interface PingPong {
  echoNum @0 (n :Int64) -> (n :Int64);
}
