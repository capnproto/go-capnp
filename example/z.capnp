using Go = import "../go.capnp";

# needed to know what package should be used for generated code
$Go.package("main");

# Needed to know how to import types in the capnp file and whether two
# capnp files are in the same package
$Go.import("github.com/jmckaskill/go-capnproto/example");

@0x832bcc6686a26d56;

struct Zdate {
  year  @0   :Int16;
  month @1   :UInt8;
  day   @2   :UInt8;
}
