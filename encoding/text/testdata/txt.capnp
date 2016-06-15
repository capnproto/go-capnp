@0x8ae03d633330d781;

struct KeyValue @0x8df8bc5abdc060a6 {
  key @0 :Text;
  value @1 :Value;
}

struct Value @0xd3602730c572a43b {
  union {
    void @0 :Void;
    bool @1 :Bool;
    int8 @2 :Int8;
    int16 @3 :Int16;
    int32 @4 :Int32;
    int64 @5 :Int64;
    uint8 @6 :UInt8;
    uint16 @7 :UInt16;
    uint32 @8 :UInt32;
    uint64 @9 :UInt64;
    float32 @10 :Float32;
    float64 @11 :Float64;
    text @12 :Text;
    textDefault @15 :Text = "hello";
    data @13 :Data;
    map @14 :List(KeyValue);
  }
}

const kv @0xc0b634e19e5a9a4e :KeyValue = (key = "42", value = (int32 = -123));
const floatKv @0x967c8fe21790b0fb :KeyValue = (key = "float", value = (float64 = 3.14));
const boolKv @0xdf35cb2e1f5ea087 :KeyValue = (key = "bool", value = (bool = false));
const mapVal @0xb167974479102805 :Value = (map = [
  (key = "foo", value = (void = void)),
  (key = "bar", value = (void = void)),
]);
const emptyMap @0x81fdbfdc91779421 :Value = (map = []);
