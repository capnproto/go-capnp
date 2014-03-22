using Go = import "go.capnp";

$Go.package("capn_test");
$Go.import("go-capnproto/example");

@0x832bcc6686a26d56;

struct Zdate {
  year  @0   :Int16;
  month @1   :UInt8;
  day   @2   :UInt8;
}

struct Zdata {
  data @0 :Data;
}


enum Airport {
  none @0;
  jfk @1;
  lax @2;
  sfo @3;
  luv @4;
  dfw @5;
  test @6; 
  # test must be last because we use it to count
  # the number of elements in the Airport enum.
}

struct PlaneBase {
  name       @0: Text;
  homes      @1: List(Airport);
  rating     @2: Int64;
  canFly     @3: Bool;
  capacity   @4: Int64;
  maxSpeed   @5: Float64;
}

struct B737 {
  base @0: PlaneBase;
}

struct A320 {
  base @0: PlaneBase;
}

struct F16 {
  base @0: PlaneBase;
}


# need a struct with at least two pointers to catch certain bugs
struct Regression {
  base     @0: PlaneBase;
  b0       @1: Float64; # intercept
  beta     @2: List(Float64);
  planes   @3: List(Aircraft);
  ymu      @4: Float64; # y-mean in original space
  ysd      @5: Float64; # y-standard deviation in original space
}



struct Aircraft {
  #  so we can restrict
  #  and specify a Plane is required in
  #  certain places.

  union {
    void      @0: Void; # @0 will be the default, so always make @0 a Void.
    b737      @1: B737;
    a320      @2: A320;
    f16       @3: F16;
  }
}


struct Z {
  # Z must contain all types, as this is our
  # runtime type identification. It is a thin shim.

  union {
    void              @0: Void; # always first in any union.
    zz                @1: Z;    # any. fyi, this can't be 'z' alone.

    f64               @2: Float64;
    f32               @3: Float32;

    i64               @4: Int64;
    i32               @5: Int32;
    i16               @6: Int16;
    i8                @7: Int8;

    u64               @8:  UInt64;
    u32               @9:  UInt32;
    u16               @10: UInt16;
    u8                @11: UInt8;

    bool              @12: Bool;
    text              @13: Text;
    blob              @14: Data;

    f64vec            @15: List(Float64);
    f32vec            @16: List(Float32);

    i64vec            @17: List(Int64);
    i32vec            @18: List(Int32);
    i16vec            @19: List(Int16);
    i8vec             @20: List(Int8);

    u64vec            @21: List(UInt64);
    u32vec            @22: List(UInt32);
    u16vec            @23: List(UInt16);
    u8vec             @24: List(UInt8);

    zvec              @25: List(Z);
    zvecvec           @26: List(List(Z));

    zdate             @27: Zdate;
    zdata             @28: Zdata;

    aircraftvec       @29: List(Aircraft);
    aircraft          @30: Aircraft;
    regression        @31: Regression;
    planebase         @32: PlaneBase;
  }
}



