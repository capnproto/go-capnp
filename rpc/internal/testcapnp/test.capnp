# Test interfaces for RPC tests.

using Go = import "../../../go.capnp";

@0xef12a34b9807e19c;
$Go.package("testcapnp");
$Go.import("zombiezen.com/go/capnproto/rpc/internal/testcapnp");

interface Handle {}

interface HandleFactory {
  newHandle @0 () -> (handle :Handle);
}
