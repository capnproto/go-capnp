# Generate scopes.capnp.out with:
# capnp compile -o- scopes.capnp > scopes.capnp.out
# Must run inside this directory to preserve paths.

using Go = import "go.capnp";
using Other = import "otherscopes.capnp";

@0xd68755941d99d05e;

$Go.package("scopes");
$Go.import("zombiezen.com/go/capnproto2/capnpc-go/testdata/scopes");

struct Foo @0xc8d7b3b4e07f8bd9 {
}

const fooVar @0x84efedc75e99768d :Foo = ();
const otherFooVar @0x836faf1834d91729 :Other.Foo = ();
