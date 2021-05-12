using Go = import "go.capnp";
@0x83c2b5818e83ab19;

$Go.package("template_fix");
$Go.import("capnproto.org/go/capnp/v3/capnpc-go/testdata/group");

struct SomeMisguidedStruct {
  someGroup :group {
    someGroupField @0 :UInt64;
  }
}
