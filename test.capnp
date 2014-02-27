using Go = import "go.capnp";

$Go.package("capn_test");

@0x832bcc6686a26d56;

struct Zdate {
  year  @0   :Int16;
  month @1   :UInt8;
  day   @2   :UInt8;
}

struct Zdata {
  data @0 :Data;
}
