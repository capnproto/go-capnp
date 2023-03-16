# Test interfaces for RPC tests.

using Go = import "/go.capnp";

@0xef12a34b9807e19c;
$Go.package("testcapnp");
$Go.import("capnproto.org/go/capnp/v3/rpc/internal/testcapnp");

interface Empty {
  # Empty interface, handy for testing shutdown hooks and stuff that just
  # needs an arbitrary capability.
}

interface EmptyProvider {
  getEmpty @0 () -> (empty :Empty);
}

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
