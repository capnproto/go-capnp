# Test interfaces for RPC tests.

using Go = import "/go.capnp";

@0xef12a34b9807e19c;
$Go.package("testcapnp");
$Go.import("capnproto.org/go/capnp/v3/rpc/internal/testcapnp");

interface PingPong {
  echoNum @0 (n :Int64) -> (n :Int64);
}

interface StreamTest {
  push @0 (data :Data) -> stream;
}

interface CapArgsTest {
  call @0 (cap :Capability);
  self @1 () -> (self :CapArgsTest);
}

interface PingPongProvider {
  pingPong @0 () -> (pingPong :PingPong);
}
