#ifndef CAPN_DDAA5E267A490CC9
#define CAPN_DDAA5E267A490CC9
/* AUTO GENERATED - DO NOT EDIT */
#include <capn.h>

#ifdef __cplusplus
extern "C" {
#endif

struct TestAllTypes;
struct TestDefaults;

typedef struct {capn_ptr p;} TestAllTypes_ptr;
typedef struct {capn_ptr p;} TestDefaults_ptr;

typedef struct {capn_ptr p;} TestAllTypes_list;
typedef struct {capn_ptr p;} TestDefaults_list;

typedef struct {capn_ptr p;} TestInterface_ptr;

typedef struct {capn_ptr p;} TestInterface_list;

enum TestEnum {
	TestEnum_foo = 0,
	TestEnum_bar = 1
};
extern unsigned int boolConst;
extern int8_t int8Const;
extern int16_t int16Const;
extern int32_t int32Const;
extern int64_t int64Const;
extern uint8_t uInt8Const;
extern uint16_t uInt16Const;
extern uint32_t uInt32Const;
extern uint64_t uInt64Const;
extern union capn_conv_f32 float32Const;
extern union capn_conv_f64 float64Const;
extern capn_text textConst;
extern capn_data dataConst;
extern TestAllTypes_ptr structConst;
extern enum TestEnum enumConst;
extern capn_ptr voidListConst;
extern capn_list1 boolListConst;
extern capn_list8 int8ListConst;
extern capn_list16 int16ListConst;
extern capn_list32 int32ListConst;
extern capn_list64 int64ListConst;
extern capn_list8 uInt8ListConst;
extern capn_list16 uInt16ListConst;
extern capn_list32 uInt32ListConst;
extern capn_list64 uInt64ListConst;
extern capn_list32 float32ListConst;
extern capn_list64 float64ListConst;
extern capn_ptr textListConst;
extern capn_ptr dataListConst;
extern TestAllTypes_list structListConst;
extern capn_list16 enumListConst;

enum TestAllTypes_unionField {
	TestAllTypes_voidUnion = 0,
	TestAllTypes_boolUnion = 1,
	TestAllTypes_int8Union = 2,
	TestAllTypes_uint8Union = 3,
	TestAllTypes_int16Union = 4,
	TestAllTypes_uint16Union = 5,
	TestAllTypes_int32Union = 6,
	TestAllTypes_uint32Union = 7,
	TestAllTypes_int64Union = 8,
	TestAllTypes_uint64Union = 9,
	TestAllTypes_float32Union = 10,
	TestAllTypes_float64Union = 11,
	TestAllTypes_textUnion = 12,
	TestAllTypes_dataUnion = 13,
	TestAllTypes_structUnion = 14,
	TestAllTypes_enumUnion = 15,
	TestAllTypes_interfaceUnion = 16
};

struct TestAllTypes {
	unsigned int boolField:1;
	int8_t int8Field;
	int16_t int16Field;
	int32_t int32Field;
	int64_t int64Field;
	uint8_t uInt8Field;
	uint16_t uInt16Field;
	uint32_t uInt32Field;
	uint64_t uInt64Field;
	float float32Field;
	double float64Field;
	capn_text textField;
	capn_data dataField;
	TestAllTypes_ptr structField;
	enum TestEnum enumField;
	TestInterface_ptr interfaceField;
	capn_ptr voidList;
	capn_list1 boolList;
	capn_list8 int8List;
	capn_list16 int16List;
	capn_list32 int32List;
	capn_list64 int64List;
	capn_list8 uInt8List;
	capn_list16 uInt16List;
	capn_list32 uInt32List;
	capn_list64 uInt64List;
	capn_list32 float32List;
	capn_list64 float64List;
	capn_ptr textList;
	capn_ptr dataList;
	TestAllTypes_list structList;
	capn_list16 enumList;
	TestInterface_list interfaceList;
	enum TestAllTypes_unionField unionField_tag;
	union {
		unsigned int boolUnion:1;
		int8_t int8Union;
		uint8_t uint8Union;
		int16_t int16Union;
		uint16_t uint16Union;
		int32_t int32Union;
		uint32_t uint32Union;
		int64_t int64Union;
		uint64_t uint64Union;
		float float32Union;
		double float64Union;
		capn_text textUnion;
		capn_data dataUnion;
		TestAllTypes_ptr structUnion;
		enum TestEnum enumUnion;
		TestInterface_ptr interfaceUnion;
	} unionField;
};

enum TestDefaults_unionField {
	TestDefaults_voidUnion = 0,
	TestDefaults_boolUnion = 1,
	TestDefaults_int8Union = 2,
	TestDefaults_uint8Union = 3,
	TestDefaults_int16Union = 4,
	TestDefaults_uint16Union = 5,
	TestDefaults_int32Union = 6,
	TestDefaults_uint32Union = 7,
	TestDefaults_int64Union = 8,
	TestDefaults_uint64Union = 9,
	TestDefaults_float32Union = 10,
	TestDefaults_float64Union = 11,
	TestDefaults_textUnion = 12,
	TestDefaults_dataUnion = 13,
	TestDefaults_structUnion = 14,
	TestDefaults_enumUnion = 15,
	TestDefaults_interfaceUnion = 16
};

struct TestDefaults {
	unsigned int boolField:1;
	int8_t int8Field;
	int16_t int16Field;
	int32_t int32Field;
	int64_t int64Field;
	uint8_t uInt8Field;
	uint16_t uInt16Field;
	uint32_t uInt32Field;
	uint64_t uInt64Field;
	float float32Field;
	double float64Field;
	capn_text textField;
	capn_data dataField;
	TestAllTypes_ptr structField;
	enum TestEnum enumField;
	TestInterface_ptr interfaceField;
	capn_ptr voidList;
	capn_list1 boolList;
	capn_list8 int8List;
	capn_list16 int16List;
	capn_list32 int32List;
	capn_list64 int64List;
	capn_list8 uInt8List;
	capn_list16 uInt16List;
	capn_list32 uInt32List;
	capn_list64 uInt64List;
	capn_list32 float32List;
	capn_list64 float64List;
	capn_ptr textList;
	capn_ptr dataList;
	TestAllTypes_list structList;
	capn_list16 enumList;
	TestInterface_list interfaceList;
	enum TestDefaults_unionField unionField_tag;
	union {
		unsigned int boolUnion:1;
		int8_t int8Union;
		uint8_t uint8Union;
		int16_t int16Union;
		uint16_t uint16Union;
		int32_t int32Union;
		uint32_t uint32Union;
		int64_t int64Union;
		uint64_t uint64Union;
		float float32Union;
		double float64Union;
		capn_text textUnion;
		capn_data dataUnion;
		TestAllTypes_ptr structUnion;
		enum TestEnum enumUnion;
		TestInterface_ptr interfaceUnion;
	} unionField;
};

TestAllTypes_ptr new_TestAllTypes(struct capn_segment*);
TestDefaults_ptr new_TestDefaults(struct capn_segment*);

TestAllTypes_list new_TestAllTypes_list(struct capn_segment*, int len);
TestDefaults_list new_TestDefaults_list(struct capn_segment*, int len);

void read_TestAllTypes(struct TestAllTypes*, TestAllTypes_ptr);
void read_TestDefaults(struct TestDefaults*, TestDefaults_ptr);

int write_TestAllTypes(const struct TestAllTypes*, TestAllTypes_ptr);
int write_TestDefaults(const struct TestDefaults*, TestDefaults_ptr);

void get_TestAllTypes(struct TestAllTypes*, TestAllTypes_list, int i);
void get_TestDefaults(struct TestDefaults*, TestDefaults_list, int i);

int set_TestAllTypes(const struct TestAllTypes*, TestAllTypes_list, int i);
int set_TestDefaults(const struct TestDefaults*, TestDefaults_list, int i);

#ifdef __cplusplus
}
#endif
#endif
