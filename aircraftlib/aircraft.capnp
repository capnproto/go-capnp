using Go = import "go.capnp";

$Go.package("aircraftlib");
$Go.import("go-capnproto/aircraftlib");

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
    airport           @33: Airport;
    b737              @34: B737;
    a320              @35: A320;
    f16               @36: F16;
    zdatevec          @37: List(Zdate);
    zdatavec          @38: List(Zdata);

    boolvec           @39: List(Bool);
  }
}

# tests for Text/List(Text) recusion handling

struct Counter {
  size  @0: Int64;
  words @1: Text;
  wordlist @2: List(Text);
}

struct Bag {
  counter  @0: Counter;
}

struct Zserver {
   waitingjobs       @0: List(Zjob);
}

struct Zjob {
    cmd        @0: Text;
    args       @1: List(Text);
}

# versioning test structs

struct VerEmpty {
}

struct VerOneData {
    val @0: Int16;
}

struct VerTwoData {
    val @0: Int16;
    duo @1: Int64;
}

struct VerOnePtr {
    ptr @0: VerOneData;
}

struct VerTwoPtr {
       ptr1 @0: VerOneData;
       ptr2 @1: VerOneData;
}

struct VerTwoDataTwoPtr {
    val @0: Int16;
    duo @1: Int64;
    ptr1 @2: VerOneData;
    ptr2 @3: VerOneData;
}

struct HoldsVerEmptyList {
  mylist @0: List(VerEmpty);
}

struct HoldsVerOneDataList {
  mylist @0: List(VerOneData);
}

struct HoldsVerTwoDataList {
  mylist @0: List(VerTwoData);
}

struct HoldsVerOnePtrList {
  mylist @0: List(VerOnePtr);
}

struct HoldsVerTwoPtrList {
  mylist @0: List(VerTwoPtr);
}

struct HoldsVerTwoTwoList {
  mylist @0: List(VerTwoDataTwoPtr);
}

struct HoldsVerTwoTwoPlus {
  mylist @0: List(VerTwoTwoPlus);
}

struct VerTwoTwoPlus {
    val @0: Int16;
    duo @1: Int64;
    ptr1 @2: VerTwoDataTwoPtr;
    ptr2 @3: VerTwoDataTwoPtr;
    tre  @4: Int64;
    lst3 @5: List(Int64);
}

# text handling

struct HoldsText {
       txt @0: Text;
       lst @1: List(Text);
       lstlst @2: List(List(Text));
}

# test that we avoid unnecessary truncation

struct WrapEmpty {
   mightNotBeReallyEmpty @0: VerEmpty;
}

struct Wrap2x2 {
   mightNotBeReallyEmpty @0: VerTwoDataTwoPtr;
}

struct Wrap2x2plus {
   mightNotBeReallyEmpty @0: VerTwoTwoPlus;
}

# test customtype annotation for Data

struct Endpoint {
   ip   @0: Data $Go.customtype("net.IP");
   port @1: Int16;
   hostname @2: Text;
}

# test voids in a union

struct VoidUnion {
  union {
    a @0 :Void;
    b @1 :Void;
  }
}
