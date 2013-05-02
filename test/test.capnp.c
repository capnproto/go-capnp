#include "../test/test.capnp.h"

static const uint8_t capnbuf[];
const int boolConst = 1;
const int8_t int8Const = INT8_C(-123);
const int16_t int16Const = INT16_C(-12345);
const int32_t int32Const = INT32_C(-12345678);
const int64_t int64Const = INT64_C(-123456789012345);
const uint8_t uInt8Const = UINT8_C(234);
const uint16_t uInt16Const = UINT16_C(45678);
const uint32_t uInt32Const = UINT32_C(3456789012);
const uint64_t uInt64Const = UINT64_C(12345678901234567890);
const float float32Const = 1234.5f;
const double float64Const = -1.23e+47;
const struct capn_text textConst = {3,"foo"};
const struct capn_data dataConst = {3,(uint8_t*)"bar"};
const struct TestAllTypes_ptr structConst = {{CAPN_STRUCT,0,(char*)capnbuf+0,0,72,22}};
const enum TestEnum enumConst = ((enum TestEnum) 1);
const struct capn_ptr voidListConst = {CAPN_LIST,6,(char*)capnbuf+1896,0,0,0};
const struct capn_list1 boolListConst = {{CAPN_BIT_LIST,4,(char*)capnbuf+1904}};
const struct capn_list8 int8ListConst = {{CAPN_LIST,2,(char*)capnbuf+1920,0,1}};
const struct capn_list16 int16ListConst = {{CAPN_LIST,2,(char*)capnbuf+1936,0,2}};
const struct capn_list32 int32ListConst = {{CAPN_LIST,2,(char*)capnbuf+1952,0,4}};
const struct capn_list64 int64ListConst = {{CAPN_LIST,2,(char*)capnbuf+1968,0,8}};
const struct capn_list8 uInt8ListConst = {{CAPN_LIST,2,(char*)capnbuf+1992,0,1}};
const struct capn_list16 uInt16ListConst = {{CAPN_LIST,2,(char*)capnbuf+2008,0,2}};
const struct capn_list32 uInt32ListConst = {{CAPN_LIST,1,(char*)capnbuf+2024,0,4}};
const struct capn_list64 uInt64ListConst = {{CAPN_LIST,1,(char*)capnbuf+2040,0,8}};
const struct capn_list32 float32ListConst = {{CAPN_LIST,2,(char*)capnbuf+2056,0,4}};
const struct capn_list64 float64ListConst = {{CAPN_LIST,2,(char*)capnbuf+2072,0,8}};
const struct capn_ptr textListConst = {CAPN_PTR_LIST,3,(char*)capnbuf+2096};
const struct capn_ptr dataListConst = {CAPN_PTR_LIST,3,(char*)capnbuf+2152};
const struct capn_ptr structListConst = {CAPN_LIST,3,(char*)capnbuf+2216,0,72,22};
const struct capn_list16 enumListConst = {{CAPN_LIST,2,(char*)capnbuf+3016,0,2}};
static const struct capn_text val_0 = {3,"abc"};
static const int val_1 = 1;
static const int8_t val_2 = INT8_C(-123);
static const uint8_t val_3 = UINT8_C(124);
static const int16_t val_4 = INT16_C(-12345);
static const uint16_t val_5 = UINT16_C(12456);
static const int32_t val_6 = INT32_C(-125678);
static const uint32_t val_7 = UINT32_C(345786);
static const int64_t val_8 = INT64_C(-123567379234);
static const uint64_t val_9 = UINT64_C(1235768497284);
static const float val_10 = 33.3f;
static const double val_11 = 340000;
static const struct capn_text val_12 = {3,"foo"};
static const struct capn_data val_13 = {3,(uint8_t*)"bar"};
static const struct TestAllTypes_ptr val_14 = {{CAPN_STRUCT,0,(char*)capnbuf+3032,0,72,22}};
static const enum TestEnum val_15 = ((enum TestEnum) 1);
static const int val_17 = 1;
static const int8_t val_18 = INT8_C(-123);
static const int16_t val_19 = INT16_C(-12345);
static const int32_t val_20 = INT32_C(-12345678);
static const int64_t val_21 = INT64_C(-123456789012345);
static const uint8_t val_22 = UINT8_C(234);
static const uint16_t val_23 = UINT16_C(45678);
static const uint32_t val_24 = UINT32_C(3456789012);
static const uint64_t val_25 = UINT64_C(12345678901234567890);
static const float val_26 = 1234.5f;
static const double val_27 = -1.23e+47;
static const struct capn_text val_28 = {3,"foo"};
static const struct capn_data val_29 = {3,(uint8_t*)"bar"};
static const struct TestAllTypes_ptr val_30 = {{CAPN_STRUCT,0,(char*)capnbuf+3288,0,72,22}};
static const enum TestEnum val_31 = ((enum TestEnum) 1);
static const struct capn_ptr val_32 = {CAPN_LIST,6,(char*)capnbuf+5184,0,0,0};
static const struct capn_list1 val_33 = {{CAPN_BIT_LIST,4,(char*)capnbuf+5192}};
static const struct capn_list8 val_34 = {{CAPN_LIST,2,(char*)capnbuf+5208,0,1}};
static const struct capn_list16 val_35 = {{CAPN_LIST,2,(char*)capnbuf+5224,0,2}};
static const struct capn_list32 val_36 = {{CAPN_LIST,2,(char*)capnbuf+5240,0,4}};
static const struct capn_list64 val_37 = {{CAPN_LIST,2,(char*)capnbuf+5256,0,8}};
static const struct capn_list8 val_38 = {{CAPN_LIST,2,(char*)capnbuf+5280,0,1}};
static const struct capn_list16 val_39 = {{CAPN_LIST,2,(char*)capnbuf+5296,0,2}};
static const struct capn_list32 val_40 = {{CAPN_LIST,1,(char*)capnbuf+5312,0,4}};
static const struct capn_list64 val_41 = {{CAPN_LIST,1,(char*)capnbuf+5328,0,8}};
static const struct capn_list32 val_42 = {{CAPN_LIST,2,(char*)capnbuf+5344,0,4}};
static const struct capn_list64 val_43 = {{CAPN_LIST,2,(char*)capnbuf+5360,0,8}};
static const struct capn_ptr val_44 = {CAPN_PTR_LIST,3,(char*)capnbuf+5384};
static const struct capn_ptr val_45 = {CAPN_PTR_LIST,3,(char*)capnbuf+5440};
static const struct capn_ptr val_46 = {CAPN_LIST,3,(char*)capnbuf+5504,0,72,22};
static const struct capn_list16 val_47 = {{CAPN_LIST,2,(char*)capnbuf+6304,0,2}};
static const int val_48 = 1;
static const int8_t val_49 = INT8_C(-123);
static const uint8_t val_50 = UINT8_C(124);
static const int16_t val_51 = INT16_C(-12345);
static const uint16_t val_52 = UINT16_C(12456);
static const int32_t val_53 = INT32_C(-125678);
static const uint32_t val_54 = UINT32_C(345786);
static const int64_t val_55 = INT64_C(-123567379234);
static const uint64_t val_56 = UINT64_C(1235768497284);
static const float val_57 = 33.3f;
static const double val_58 = 340000;
static const enum TestEnum val_59 = ((enum TestEnum) 1);
static const struct TestInterface_vt TestInterface_remote_vt;

struct TestAllTypes_ptr new_TestAllTypes(struct capn_segment* seg) {
	struct TestAllTypes_ptr ret = {capn_new_struct(seg, 72, 22)};
	return ret;
}

struct capn_ptr new_TestAllTypes_list(struct capn_segment *seg, int sz) {
	return capn_new_list(seg, sz, 72, 22);
}

int read_TestAllTypes(const struct TestAllTypes_ptr *p, struct TestAllTypes *s) {
	if (p->p.type != CAPN_STRUCT)
		return -1;
	s->boolField = (capn_get8(&p->p, 0) & 1) != 0;
	s->int8Field = ((int8_t) capn_get8(&p->p, 1));
	s->int16Field = ((int16_t) capn_get16(&p->p, 2));
	s->int32Field = ((int32_t) capn_get32(&p->p, 4));
	s->int64Field = ((int64_t) capn_get64(&p->p, 8));
	s->uInt8Field = capn_get8(&p->p, 16);
	s->uInt16Field = capn_get16(&p->p, 18);
	s->uInt32Field = capn_get32(&p->p, 20);
	s->uInt64Field = capn_get64(&p->p, 24);
	s->float32Field = capn_get_float(&p->p, 32, 0.0f);
	s->float64Field = capn_get_double(&p->p, 40, 0.0);
	s->textField = capn_read_text(&p->p, 0);
	s->dataField = capn_read_data(&p->p, 1);
	s->structField.p = capn_read_ptr(&p->p, 2);
	s->enumField = (enum TestEnum) capn_get16(&p->p, 48);
	s->interfaceField.vt = &TestInterface_remote_vt;
	s->interfaceField.p = capn_read_ptr(&p->p, 3);
	s->voidList = capn_read_ptr(&p->p, 4);
	s->boolList.p = capn_read_ptr(&p->p, 5);
	s->int8List.p = capn_read_ptr(&p->p, 6);
	s->int16List.p = capn_read_ptr(&p->p, 7);
	s->int32List.p = capn_read_ptr(&p->p, 8);
	s->int64List.p = capn_read_ptr(&p->p, 9);
	s->uInt8List.p = capn_read_ptr(&p->p, 10);
	s->uInt16List.p = capn_read_ptr(&p->p, 11);
	s->uInt32List.p = capn_read_ptr(&p->p, 12);
	s->uInt64List.p = capn_read_ptr(&p->p, 13);
	s->float32List.p = capn_read_ptr(&p->p, 14);
	s->float64List.p = capn_read_ptr(&p->p, 15);
	s->textList = capn_read_ptr(&p->p, 16);
	s->dataList = capn_read_ptr(&p->p, 17);
	s->structList = capn_read_ptr(&p->p, 18);
	s->enumList.p = capn_read_ptr(&p->p, 19);
	s->interfaceList = capn_read_ptr(&p->p, 20);
	s->unionField_tag = (enum TestAllTypes_unionField) capn_get16(&p->p, 50);
	switch (s->unionField_tag) {
	case 35:
		break;
	case 36:
		s->unionField.boolUnion = (capn_get8(&p->p, 52) & 1) != 0;
		break;
	case 37:
		s->unionField.int8Union = ((int8_t) capn_get8(&p->p, 53));
		break;
	case 38:
		s->unionField.uint8Union = capn_get8(&p->p, 53);
		break;
	case 39:
		s->unionField.int16Union = ((int16_t) capn_get16(&p->p, 54));
		break;
	case 40:
		s->unionField.uint16Union = capn_get16(&p->p, 54);
		break;
	case 41:
		s->unionField.int32Union = ((int32_t) capn_get32(&p->p, 56));
		break;
	case 42:
		s->unionField.uint32Union = capn_get32(&p->p, 56);
		break;
	case 43:
		s->unionField.int64Union = ((int64_t) capn_get64(&p->p, 64));
		break;
	case 44:
		s->unionField.uint64Union = capn_get64(&p->p, 64);
		break;
	case 45:
		s->unionField.float32Union = capn_get_float(&p->p, 64, 0.0f);
		break;
	case 46:
		s->unionField.float64Union = capn_get_double(&p->p, 64, 0.0);
		break;
	case 47:
		s->unionField.textUnion = capn_read_text(&p->p, 21);
		break;
	case 48:
		s->unionField.dataUnion = capn_read_data(&p->p, 21);
		break;
	case 49:
		s->unionField.structUnion.p = capn_read_ptr(&p->p, 21);
		break;
	case 50:
		s->unionField.enumUnion = (enum TestEnum) capn_get16(&p->p, 64);
		break;
	case 51:
		s->unionField.interfaceUnion.vt = &TestInterface_remote_vt;
		s->unionField.interfaceUnion.p = capn_read_ptr(&p->p, 21);
		break;
	}
	return 0;
}

int write_TestAllTypes(struct TestAllTypes_ptr *p, const struct TestAllTypes *s) {
	if (p->p.type != CAPN_STRUCT)
		return -1;
	capn_set8(&p->p, 0, (capn_get8(&p->p, 0) & ~1) | (s->boolField << 0));
	capn_set8(&p->p, 1, (uint8_t) (s->int8Field));
	capn_set16(&p->p, 2, (uint16_t) (s->int16Field));
	capn_set32(&p->p, 4, (uint32_t) (s->int32Field));
	capn_set64(&p->p, 8, (uint64_t) (s->int64Field));
	capn_set8(&p->p, 16, s->uInt8Field);
	capn_set16(&p->p, 18, s->uInt16Field);
	capn_set32(&p->p, 20, s->uInt32Field);
	capn_set64(&p->p, 24, s->uInt64Field);
	capn_set_float(&p->p, 32, s->float32Field, 0.0f);
	capn_set_double(&p->p, 40, s->float64Field, 0.0);
	capn_write_text(&p->p, 0, s->textField);
	capn_write_data(&p->p, 1, s->dataField);
	capn_write_ptr(&p->p, 2, &s->structField.p);
	capn_set16(&p->p, 48, (uint16_t) (s->enumField));
	s->interfaceField.vt->marshal(&s->interfaceField, &p->p, 3);
	capn_write_ptr(&p->p, 4, &s->voidList);
	capn_write_ptr(&p->p, 5, &s->boolList.p);
	capn_write_ptr(&p->p, 6, &s->int8List.p);
	capn_write_ptr(&p->p, 7, &s->int16List.p);
	capn_write_ptr(&p->p, 8, &s->int32List.p);
	capn_write_ptr(&p->p, 9, &s->int64List.p);
	capn_write_ptr(&p->p, 10, &s->uInt8List.p);
	capn_write_ptr(&p->p, 11, &s->uInt16List.p);
	capn_write_ptr(&p->p, 12, &s->uInt32List.p);
	capn_write_ptr(&p->p, 13, &s->uInt64List.p);
	capn_write_ptr(&p->p, 14, &s->float32List.p);
	capn_write_ptr(&p->p, 15, &s->float64List.p);
	capn_write_ptr(&p->p, 16, &s->textList);
	capn_write_ptr(&p->p, 17, &s->dataList);
	capn_write_ptr(&p->p, 18, &s->structList);
	capn_write_ptr(&p->p, 19, &s->enumList.p);
	capn_write_ptr(&p->p, 20, &s->interfaceList);
	capn_set16(&p->p, 50, (uint16_t) s->unionField_tag);
	switch (s->unionField_tag) {
	case 35:
		break;
	case 36:
		capn_set8(&p->p, 52, (capn_get8(&p->p, 52) & ~1) | (s->unionField.boolUnion << 0));
		break;
	case 37:
		capn_set8(&p->p, 53, (uint8_t) (s->unionField.int8Union));
		break;
	case 38:
		capn_set8(&p->p, 53, s->unionField.uint8Union);
		break;
	case 39:
		capn_set16(&p->p, 54, (uint16_t) (s->unionField.int16Union));
		break;
	case 40:
		capn_set16(&p->p, 54, s->unionField.uint16Union);
		break;
	case 41:
		capn_set32(&p->p, 56, (uint32_t) (s->unionField.int32Union));
		break;
	case 42:
		capn_set32(&p->p, 56, s->unionField.uint32Union);
		break;
	case 43:
		capn_set64(&p->p, 64, (uint64_t) (s->unionField.int64Union));
		break;
	case 44:
		capn_set64(&p->p, 64, s->unionField.uint64Union);
		break;
	case 45:
		capn_set_float(&p->p, 64, s->unionField.float32Union, 0.0f);
		break;
	case 46:
		capn_set_double(&p->p, 64, s->unionField.float64Union, 0.0);
		break;
	case 47:
		capn_write_text(&p->p, 21, s->unionField.textUnion);
		break;
	case 48:
		capn_write_data(&p->p, 21, s->unionField.dataUnion);
		break;
	case 49:
		capn_write_ptr(&p->p, 21, &s->unionField.structUnion.p);
		break;
	case 50:
		capn_set16(&p->p, 64, (uint16_t) (s->unionField.enumUnion));
		break;
	case 51:
		s->unionField.interfaceUnion.vt->marshal(&s->unionField.interfaceUnion, &p->p, 21);
		break;
	}
	return 0;
}

struct TestDefaults_ptr new_TestDefaults(struct capn_segment* seg) {
	struct TestDefaults_ptr ret = {capn_new_struct(seg, 72, 22)};
	return ret;
}

struct capn_ptr new_TestDefaults_list(struct capn_segment *seg, int sz) {
	return capn_new_list(seg, sz, 72, 22);
}

int read_TestDefaults(const struct TestDefaults_ptr *p, struct TestDefaults *s) {
	if (p->p.type != CAPN_STRUCT)
		return -1;
	s->boolField = (capn_get8(&p->p, 0) & 1) == 0;
	s->int8Field = ((int8_t) capn_get8(&p->p, 1)) ^ INT8_C(-123);
	s->int16Field = ((int16_t) capn_get16(&p->p, 2)) ^ INT16_C(-12345);
	s->int32Field = ((int32_t) capn_get32(&p->p, 4)) ^ INT32_C(-12345678);
	s->int64Field = ((int64_t) capn_get64(&p->p, 8)) ^ INT64_C(-123456789012345);
	s->uInt8Field = capn_get8(&p->p, 16) ^ UINT8_C(234);
	s->uInt16Field = capn_get16(&p->p, 18) ^ UINT16_C(45678);
	s->uInt32Field = capn_get32(&p->p, 20) ^ UINT32_C(3456789012);
	s->uInt64Field = capn_get64(&p->p, 24) ^ UINT64_C(12345678901234567890);
	s->float32Field = capn_get_float(&p->p, 32, 1234.5f);
	s->float64Field = capn_get_double(&p->p, 40, -1.23e+47);
	s->textField = capn_read_text(&p->p, 0);
	if (!s->textField.str) {
		s->textField = val_28;
	}
	s->dataField = capn_read_data(&p->p, 1);
	if (!s->dataField.data) {
		s->dataField = val_29;
	}
	s->structField.p = capn_read_ptr(&p->p, 2);
	if (s->structField.p.type == CAPN_NULL) {
		s->structField = val_30;
	}
	s->enumField = (enum TestEnum) capn_get16(&p->p, 48) ^ ((enum TestEnum) 1);
	s->interfaceField.vt = &TestInterface_remote_vt;
	s->interfaceField.p = capn_read_ptr(&p->p, 3);
	s->voidList = capn_read_ptr(&p->p, 4);
	if (s->voidList.type == CAPN_NULL) {
		s->voidList = val_32;
	}
	s->boolList.p = capn_read_ptr(&p->p, 5);
	if (s->boolList.p.type == CAPN_NULL) {
		s->boolList = val_33;
	}
	s->int8List.p = capn_read_ptr(&p->p, 6);
	if (s->int8List.p.type == CAPN_NULL) {
		s->int8List = val_34;
	}
	s->int16List.p = capn_read_ptr(&p->p, 7);
	if (s->int16List.p.type == CAPN_NULL) {
		s->int16List = val_35;
	}
	s->int32List.p = capn_read_ptr(&p->p, 8);
	if (s->int32List.p.type == CAPN_NULL) {
		s->int32List = val_36;
	}
	s->int64List.p = capn_read_ptr(&p->p, 9);
	if (s->int64List.p.type == CAPN_NULL) {
		s->int64List = val_37;
	}
	s->uInt8List.p = capn_read_ptr(&p->p, 10);
	if (s->uInt8List.p.type == CAPN_NULL) {
		s->uInt8List = val_38;
	}
	s->uInt16List.p = capn_read_ptr(&p->p, 11);
	if (s->uInt16List.p.type == CAPN_NULL) {
		s->uInt16List = val_39;
	}
	s->uInt32List.p = capn_read_ptr(&p->p, 12);
	if (s->uInt32List.p.type == CAPN_NULL) {
		s->uInt32List = val_40;
	}
	s->uInt64List.p = capn_read_ptr(&p->p, 13);
	if (s->uInt64List.p.type == CAPN_NULL) {
		s->uInt64List = val_41;
	}
	s->float32List.p = capn_read_ptr(&p->p, 14);
	if (s->float32List.p.type == CAPN_NULL) {
		s->float32List = val_42;
	}
	s->float64List.p = capn_read_ptr(&p->p, 15);
	if (s->float64List.p.type == CAPN_NULL) {
		s->float64List = val_43;
	}
	s->textList = capn_read_ptr(&p->p, 16);
	if (s->textList.type == CAPN_NULL) {
		s->textList = val_44;
	}
	s->dataList = capn_read_ptr(&p->p, 17);
	if (s->dataList.type == CAPN_NULL) {
		s->dataList = val_45;
	}
	s->structList = capn_read_ptr(&p->p, 18);
	if (s->structList.type == CAPN_NULL) {
		s->structList = val_46;
	}
	s->enumList.p = capn_read_ptr(&p->p, 19);
	if (s->enumList.p.type == CAPN_NULL) {
		s->enumList = val_47;
	}
	s->interfaceList = capn_read_ptr(&p->p, 20);
	s->unionField_tag = (enum TestDefaults_unionField) capn_get16(&p->p, 50);
	switch (s->unionField_tag) {
	case 35:
		break;
	case 36:
		s->unionField.boolUnion = (capn_get8(&p->p, 52) & 1) == 0;
		break;
	case 37:
		s->unionField.int8Union = ((int8_t) capn_get8(&p->p, 53)) ^ INT8_C(-123);
		break;
	case 38:
		s->unionField.uint8Union = capn_get8(&p->p, 53) ^ UINT8_C(124);
		break;
	case 39:
		s->unionField.int16Union = ((int16_t) capn_get16(&p->p, 54)) ^ INT16_C(-12345);
		break;
	case 40:
		s->unionField.uint16Union = capn_get16(&p->p, 54) ^ UINT16_C(12456);
		break;
	case 41:
		s->unionField.int32Union = ((int32_t) capn_get32(&p->p, 56)) ^ INT32_C(-125678);
		break;
	case 42:
		s->unionField.uint32Union = capn_get32(&p->p, 56) ^ UINT32_C(345786);
		break;
	case 43:
		s->unionField.int64Union = ((int64_t) capn_get64(&p->p, 64)) ^ INT64_C(-123567379234);
		break;
	case 44:
		s->unionField.uint64Union = capn_get64(&p->p, 64) ^ UINT64_C(1235768497284);
		break;
	case 45:
		s->unionField.float32Union = capn_get_float(&p->p, 64, 33.3f);
		break;
	case 46:
		s->unionField.float64Union = capn_get_double(&p->p, 64, 340000);
		break;
	case 47:
		s->unionField.textUnion = capn_read_text(&p->p, 21);
		if (!s->unionField.textUnion.str) {
			s->unionField.textUnion = val_12;
		}
		break;
	case 48:
		s->unionField.dataUnion = capn_read_data(&p->p, 21);
		if (!s->unionField.dataUnion.data) {
			s->unionField.dataUnion = val_13;
		}
		break;
	case 49:
		s->unionField.structUnion.p = capn_read_ptr(&p->p, 21);
		if (s->unionField.structUnion.p.type == CAPN_NULL) {
			s->unionField.structUnion = val_14;
		}
		break;
	case 50:
		s->unionField.enumUnion = (enum TestEnum) capn_get16(&p->p, 64) ^ ((enum TestEnum) 1);
		break;
	case 51:
		s->unionField.interfaceUnion.vt = &TestInterface_remote_vt;
		s->unionField.interfaceUnion.p = capn_read_ptr(&p->p, 21);
		break;
	}
	return 0;
}

int write_TestDefaults(struct TestDefaults_ptr *p, const struct TestDefaults *s) {
	if (p->p.type != CAPN_STRUCT)
		return -1;
	capn_set8(&p->p, 0, (capn_get8(&p->p, 0) & ~1) & ~(s->boolField << 0));
	capn_set8(&p->p, 1, (uint8_t) (s->int8Field ^ INT8_C(-123)));
	capn_set16(&p->p, 2, (uint16_t) (s->int16Field ^ INT16_C(-12345)));
	capn_set32(&p->p, 4, (uint32_t) (s->int32Field ^ INT32_C(-12345678)));
	capn_set64(&p->p, 8, (uint64_t) (s->int64Field ^ INT64_C(-123456789012345)));
	capn_set8(&p->p, 16, s->uInt8Field ^ UINT8_C(234));
	capn_set16(&p->p, 18, s->uInt16Field ^ UINT16_C(45678));
	capn_set32(&p->p, 20, s->uInt32Field ^ UINT32_C(3456789012));
	capn_set64(&p->p, 24, s->uInt64Field ^ UINT64_C(12345678901234567890));
	capn_set_float(&p->p, 32, s->float32Field, 1234.5f);
	capn_set_double(&p->p, 40, s->float64Field, -1.23e+47);
	if (s->textField.str != val_28.str || s->textField.size != val_28.size) {
		capn_write_text(&p->p, 0, s->textField);
	} else {
		capn_write_ptr(&p->p, 0, 0);
	}
	if (s->dataField.data != val_29.data || s->dataField.size != val_29.size) {
		capn_write_data(&p->p, 1, s->dataField);
	} else {
		capn_write_ptr(&p->p, 1, 0);
	}
	capn_write_ptr(&p->p, 2, s->structField.p.data != val_30.p.data ? &s->structField.p : 0);
	capn_set16(&p->p, 48, (uint16_t) (s->enumField ^ ((enum TestEnum) 1)));
	s->interfaceField.vt->marshal(&s->interfaceField, &p->p, 3);
	capn_write_ptr(&p->p, 4, (s->voidList.data != val_32.data || s->voidList.size != val_32.size) ? &s->voidList : 0);
	capn_write_ptr(&p->p, 5, (s->boolList.p.data != val_33.p.data || s->boolList.p.size != val_33.p.size) ? &s->boolList.p : 0);
	capn_write_ptr(&p->p, 6, (s->int8List.p.data != val_34.p.data || s->int8List.p.size != val_34.p.size) ? &s->int8List.p : 0);
	capn_write_ptr(&p->p, 7, (s->int16List.p.data != val_35.p.data || s->int16List.p.size != val_35.p.size) ? &s->int16List.p : 0);
	capn_write_ptr(&p->p, 8, (s->int32List.p.data != val_36.p.data || s->int32List.p.size != val_36.p.size) ? &s->int32List.p : 0);
	capn_write_ptr(&p->p, 9, (s->int64List.p.data != val_37.p.data || s->int64List.p.size != val_37.p.size) ? &s->int64List.p : 0);
	capn_write_ptr(&p->p, 10, (s->uInt8List.p.data != val_38.p.data || s->uInt8List.p.size != val_38.p.size) ? &s->uInt8List.p : 0);
	capn_write_ptr(&p->p, 11, (s->uInt16List.p.data != val_39.p.data || s->uInt16List.p.size != val_39.p.size) ? &s->uInt16List.p : 0);
	capn_write_ptr(&p->p, 12, (s->uInt32List.p.data != val_40.p.data || s->uInt32List.p.size != val_40.p.size) ? &s->uInt32List.p : 0);
	capn_write_ptr(&p->p, 13, (s->uInt64List.p.data != val_41.p.data || s->uInt64List.p.size != val_41.p.size) ? &s->uInt64List.p : 0);
	capn_write_ptr(&p->p, 14, (s->float32List.p.data != val_42.p.data || s->float32List.p.size != val_42.p.size) ? &s->float32List.p : 0);
	capn_write_ptr(&p->p, 15, (s->float64List.p.data != val_43.p.data || s->float64List.p.size != val_43.p.size) ? &s->float64List.p : 0);
	capn_write_ptr(&p->p, 16, (s->textList.data != val_44.data || s->textList.size != val_44.size) ? &s->textList : 0);
	capn_write_ptr(&p->p, 17, (s->dataList.data != val_45.data || s->dataList.size != val_45.size) ? &s->dataList : 0);
	capn_write_ptr(&p->p, 18, (s->structList.data != val_46.data || s->structList.size != val_46.size) ? &s->structList : 0);
	capn_write_ptr(&p->p, 19, (s->enumList.p.data != val_47.p.data || s->enumList.p.size != val_47.p.size) ? &s->enumList.p : 0);
	capn_write_ptr(&p->p, 20, &s->interfaceList);
	capn_set16(&p->p, 50, (uint16_t) s->unionField_tag);
	switch (s->unionField_tag) {
	case 35:
		break;
	case 36:
		capn_set8(&p->p, 52, (capn_get8(&p->p, 52) & ~1) & ~(s->unionField.boolUnion << 0));
		break;
	case 37:
		capn_set8(&p->p, 53, (uint8_t) (s->unionField.int8Union ^ INT8_C(-123)));
		break;
	case 38:
		capn_set8(&p->p, 53, s->unionField.uint8Union ^ UINT8_C(124));
		break;
	case 39:
		capn_set16(&p->p, 54, (uint16_t) (s->unionField.int16Union ^ INT16_C(-12345)));
		break;
	case 40:
		capn_set16(&p->p, 54, s->unionField.uint16Union ^ UINT16_C(12456));
		break;
	case 41:
		capn_set32(&p->p, 56, (uint32_t) (s->unionField.int32Union ^ INT32_C(-125678)));
		break;
	case 42:
		capn_set32(&p->p, 56, s->unionField.uint32Union ^ UINT32_C(345786));
		break;
	case 43:
		capn_set64(&p->p, 64, (uint64_t) (s->unionField.int64Union ^ INT64_C(-123567379234)));
		break;
	case 44:
		capn_set64(&p->p, 64, s->unionField.uint64Union ^ UINT64_C(1235768497284));
		break;
	case 45:
		capn_set_float(&p->p, 64, s->unionField.float32Union, 33.3f);
		break;
	case 46:
		capn_set_double(&p->p, 64, s->unionField.float64Union, 340000);
		break;
	case 47:
		if (s->unionField.textUnion.str != val_12.str || s->unionField.textUnion.size != val_12.size) {
			capn_write_text(&p->p, 21, s->unionField.textUnion);
		} else {
			capn_write_ptr(&p->p, 21, 0);
		}
		break;
	case 48:
		if (s->unionField.dataUnion.data != val_13.data || s->unionField.dataUnion.size != val_13.size) {
			capn_write_data(&p->p, 21, s->unionField.dataUnion);
		} else {
			capn_write_ptr(&p->p, 21, 0);
		}
		break;
	case 49:
		capn_write_ptr(&p->p, 21, s->unionField.structUnion.p.data != val_14.p.data ? &s->unionField.structUnion.p : 0);
		break;
	case 50:
		capn_set16(&p->p, 64, (uint16_t) (s->unionField.enumUnion ^ ((enum TestEnum) 1)));
		break;
	case 51:
		s->unionField.interfaceUnion.vt->marshal(&s->unionField.interfaceUnion, &p->p, 21);
		break;
	}
	return 0;
}

static const uint8_t capnbuf[] = {
 0, 0, 0, 0, 9, 0, 22, 0,
 1, 244, 128, 13, 14, 16, 76, 251,
 78, 115, 232, 56, 166, 51, 0, 0,
 90, 0, 210, 4, 20, 136, 98, 3,
 210, 10, 111, 18, 33, 25, 204, 4,
 95, 112, 9, 175, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 144, 117, 64,
 1, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 85, 0, 0, 0, 34, 0, 0, 0,
 85, 0, 0, 0, 26, 0, 0, 0,
 84, 0, 0, 0, 9, 0, 22, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 81, 1, 0, 0, 24, 0, 0, 0,
 77, 1, 0, 0, 41, 0, 0, 0,
 77, 1, 0, 0, 34, 0, 0, 0,
 77, 1, 0, 0, 35, 0, 0, 0,
 77, 1, 0, 0, 36, 0, 0, 0,
 81, 1, 0, 0, 37, 0, 0, 0,
 93, 1, 0, 0, 34, 0, 0, 0,
 93, 1, 0, 0, 35, 0, 0, 0,
 93, 1, 0, 0, 36, 0, 0, 0,
 97, 1, 0, 0, 37, 0, 0, 0,
 109, 1, 0, 0, 52, 0, 0, 0,
 117, 1, 0, 0, 53, 0, 0, 0,
 137, 1, 0, 0, 30, 0, 0, 0,
 157, 1, 0, 0, 30, 0, 0, 0,
 177, 1, 0, 0, 31, 0, 0, 0,
 57, 3, 0, 0, 19, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 98, 97, 122, 0, 0, 0, 0, 0,
 113, 117, 120, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 85, 0, 0, 0, 58, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 80, 0, 0, 0, 9, 0, 22, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 110, 101, 115, 116, 101, 100, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 85, 0, 0, 0, 114, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 114, 101, 97, 108, 108, 121, 32, 110,
 101, 115, 116, 101, 100, 0, 0, 0,
 26, 0, 0, 0, 0, 0, 0, 0,
 12, 222, 128, 127, 0, 0, 0, 0,
 210, 4, 210, 233, 0, 128, 255, 127,
 78, 97, 188, 0, 64, 211, 160, 250,
 0, 0, 0, 248, 255, 255, 255, 7,
 121, 223, 13, 134, 72, 112, 0, 0,
 46, 117, 19, 253, 138, 150, 253, 255,
 0, 0, 0, 0, 0, 0, 0, 128,
 255, 255, 255, 255, 255, 255, 255, 127,
 12, 34, 0, 255, 0, 0, 0, 0,
 210, 4, 46, 22, 0, 0, 255, 255,
 78, 97, 188, 0, 192, 44, 95, 5,
 0, 0, 0, 0, 255, 255, 255, 255,
 121, 223, 13, 134, 72, 112, 0, 0,
 210, 138, 236, 2, 117, 105, 2, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 255, 255, 255, 255, 255, 255, 255, 255,
 0, 0, 0, 0, 56, 180, 150, 73,
 194, 189, 240, 124, 194, 189, 240, 252,
 234, 28, 8, 2, 234, 28, 8, 130,
 0, 0, 0, 0, 0, 0, 0, 0,
 64, 222, 119, 131, 33, 18, 220, 66,
 41, 144, 35, 202, 229, 200, 118, 127,
 41, 144, 35, 202, 229, 200, 118, 255,
 145, 247, 80, 55, 158, 120, 102, 0,
 145, 247, 80, 55, 158, 120, 102, 128,
 9, 0, 0, 0, 42, 0, 0, 0,
 9, 0, 0, 0, 50, 0, 0, 0,
 9, 0, 0, 0, 58, 0, 0, 0,
 113, 117, 117, 120, 0, 0, 0, 0,
 99, 111, 114, 103, 101, 0, 0, 0,
 103, 114, 97, 117, 108, 116, 0, 0,
 9, 0, 0, 0, 50, 0, 0, 0,
 9, 0, 0, 0, 42, 0, 0, 0,
 9, 0, 0, 0, 34, 0, 0, 0,
 103, 97, 114, 112, 108, 121, 0, 0,
 119, 97, 108, 100, 111, 0, 0, 0,
 102, 114, 101, 100, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 77, 1, 0, 0, 122, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 217, 0, 0, 0, 122, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 101, 0, 0, 0, 122, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 120, 32, 115, 116, 114, 117, 99, 116,
 108, 105, 115, 116, 32, 49, 0, 0,
 120, 32, 115, 116, 114, 117, 99, 116,
 108, 105, 115, 116, 32, 50, 0, 0,
 120, 32, 115, 116, 114, 117, 99, 116,
 108, 105, 115, 116, 32, 51, 0, 0,
 2, 0, 1, 0, 0, 0, 0, 0,
 1, 0, 0, 0, 48, 0, 0, 0,
 1, 0, 0, 0, 33, 0, 0, 0,
 9, 0, 0, 0, 0, 0, 0, 0,
 1, 0, 0, 0, 18, 0, 0, 0,
 111, 145, 0, 0, 0, 0, 0, 0,
 1, 0, 0, 0, 19, 0, 0, 0,
 103, 43, 153, 212, 0, 0, 0, 0,
 1, 0, 0, 0, 20, 0, 0, 0,
 199, 107, 159, 6, 57, 148, 96, 249,
 1, 0, 0, 0, 21, 0, 0, 0,
 199, 113, 196, 43, 171, 117, 107, 15,
 57, 142, 59, 212, 84, 138, 148, 240,
 1, 0, 0, 0, 18, 0, 0, 0,
 111, 222, 0, 0, 0, 0, 0, 0,
 1, 0, 0, 0, 19, 0, 0, 0,
 53, 130, 156, 173, 0, 0, 0, 0,
 1, 0, 0, 0, 12, 0, 0, 0,
 85, 161, 174, 198, 0, 0, 0, 0,
 1, 0, 0, 0, 13, 0, 0, 0,
 199, 113, 172, 181, 175, 152, 50, 154,
 1, 0, 0, 0, 20, 0, 0, 0,
 0, 156, 173, 69, 0, 228, 10, 69,
 1, 0, 0, 0, 21, 0, 0, 0,
 0, 0, 0, 0, 192, 97, 190, 64,
 0, 0, 0, 0, 128, 92, 145, 64,
 1, 0, 0, 0, 30, 0, 0, 0,
 9, 0, 0, 0, 50, 0, 0, 0,
 9, 0, 0, 0, 50, 0, 0, 0,
 9, 0, 0, 0, 42, 0, 0, 0,
 112, 108, 117, 103, 104, 0, 0, 0,
 120, 121, 122, 122, 121, 0, 0, 0,
 116, 104, 117, 100, 0, 0, 0, 0,
 1, 0, 0, 0, 30, 0, 0, 0,
 9, 0, 0, 0, 34, 0, 0, 0,
 9, 0, 0, 0, 74, 0, 0, 0,
 13, 0, 0, 0, 58, 0, 0, 0,
 111, 111, 112, 115, 0, 0, 0, 0,
 101, 120, 104, 97, 117, 115, 116, 101,
 100, 0, 0, 0, 0, 0, 0, 0,
 114, 102, 99, 51, 48, 57, 50, 0,
 1, 0, 0, 0, 31, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 77, 1, 0, 0, 106, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 217, 0, 0, 0, 106, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 101, 0, 0, 0, 106, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 115, 116, 114, 117, 99, 116, 108, 105,
 115, 116, 32, 49, 0, 0, 0, 0,
 115, 116, 114, 117, 99, 116, 108, 105,
 115, 116, 32, 50, 0, 0, 0, 0,
 115, 116, 114, 117, 99, 116, 108, 105,
 115, 116, 32, 51, 0, 0, 0, 0,
 1, 0, 0, 0, 19, 0, 0, 0,
 1, 0, 2, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 9, 0, 22, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 133, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 9, 0, 22, 0,
 1, 244, 128, 13, 14, 16, 76, 251,
 78, 115, 232, 56, 166, 51, 0, 0,
 90, 0, 210, 4, 20, 136, 98, 3,
 210, 10, 111, 18, 33, 25, 204, 4,
 95, 112, 9, 175, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 144, 117, 64,
 1, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 85, 0, 0, 0, 34, 0, 0, 0,
 85, 0, 0, 0, 26, 0, 0, 0,
 84, 0, 0, 0, 9, 0, 22, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 81, 1, 0, 0, 24, 0, 0, 0,
 77, 1, 0, 0, 41, 0, 0, 0,
 77, 1, 0, 0, 34, 0, 0, 0,
 77, 1, 0, 0, 35, 0, 0, 0,
 77, 1, 0, 0, 36, 0, 0, 0,
 81, 1, 0, 0, 37, 0, 0, 0,
 93, 1, 0, 0, 34, 0, 0, 0,
 93, 1, 0, 0, 35, 0, 0, 0,
 93, 1, 0, 0, 36, 0, 0, 0,
 97, 1, 0, 0, 37, 0, 0, 0,
 109, 1, 0, 0, 52, 0, 0, 0,
 117, 1, 0, 0, 53, 0, 0, 0,
 137, 1, 0, 0, 30, 0, 0, 0,
 157, 1, 0, 0, 30, 0, 0, 0,
 177, 1, 0, 0, 31, 0, 0, 0,
 57, 3, 0, 0, 19, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 98, 97, 122, 0, 0, 0, 0, 0,
 113, 117, 120, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 85, 0, 0, 0, 58, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 80, 0, 0, 0, 9, 0, 22, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 110, 101, 115, 116, 101, 100, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 85, 0, 0, 0, 114, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 114, 101, 97, 108, 108, 121, 32, 110,
 101, 115, 116, 101, 100, 0, 0, 0,
 26, 0, 0, 0, 0, 0, 0, 0,
 12, 222, 128, 127, 0, 0, 0, 0,
 210, 4, 210, 233, 0, 128, 255, 127,
 78, 97, 188, 0, 64, 211, 160, 250,
 0, 0, 0, 248, 255, 255, 255, 7,
 121, 223, 13, 134, 72, 112, 0, 0,
 46, 117, 19, 253, 138, 150, 253, 255,
 0, 0, 0, 0, 0, 0, 0, 128,
 255, 255, 255, 255, 255, 255, 255, 127,
 12, 34, 0, 255, 0, 0, 0, 0,
 210, 4, 46, 22, 0, 0, 255, 255,
 78, 97, 188, 0, 192, 44, 95, 5,
 0, 0, 0, 0, 255, 255, 255, 255,
 121, 223, 13, 134, 72, 112, 0, 0,
 210, 138, 236, 2, 117, 105, 2, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 255, 255, 255, 255, 255, 255, 255, 255,
 0, 0, 0, 0, 56, 180, 150, 73,
 194, 189, 240, 124, 194, 189, 240, 252,
 234, 28, 8, 2, 234, 28, 8, 130,
 0, 0, 0, 0, 0, 0, 0, 0,
 64, 222, 119, 131, 33, 18, 220, 66,
 41, 144, 35, 202, 229, 200, 118, 127,
 41, 144, 35, 202, 229, 200, 118, 255,
 145, 247, 80, 55, 158, 120, 102, 0,
 145, 247, 80, 55, 158, 120, 102, 128,
 9, 0, 0, 0, 42, 0, 0, 0,
 9, 0, 0, 0, 50, 0, 0, 0,
 9, 0, 0, 0, 58, 0, 0, 0,
 113, 117, 117, 120, 0, 0, 0, 0,
 99, 111, 114, 103, 101, 0, 0, 0,
 103, 114, 97, 117, 108, 116, 0, 0,
 9, 0, 0, 0, 50, 0, 0, 0,
 9, 0, 0, 0, 42, 0, 0, 0,
 9, 0, 0, 0, 34, 0, 0, 0,
 103, 97, 114, 112, 108, 121, 0, 0,
 119, 97, 108, 100, 111, 0, 0, 0,
 102, 114, 101, 100, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 77, 1, 0, 0, 122, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 217, 0, 0, 0, 122, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 101, 0, 0, 0, 122, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 120, 32, 115, 116, 114, 117, 99, 116,
 108, 105, 115, 116, 32, 49, 0, 0,
 120, 32, 115, 116, 114, 117, 99, 116,
 108, 105, 115, 116, 32, 50, 0, 0,
 120, 32, 115, 116, 114, 117, 99, 116,
 108, 105, 115, 116, 32, 51, 0, 0,
 2, 0, 1, 0, 0, 0, 0, 0,
 1, 0, 0, 0, 48, 0, 0, 0,
 1, 0, 0, 0, 33, 0, 0, 0,
 9, 0, 0, 0, 0, 0, 0, 0,
 1, 0, 0, 0, 18, 0, 0, 0,
 111, 145, 0, 0, 0, 0, 0, 0,
 1, 0, 0, 0, 19, 0, 0, 0,
 103, 43, 153, 212, 0, 0, 0, 0,
 1, 0, 0, 0, 20, 0, 0, 0,
 199, 107, 159, 6, 57, 148, 96, 249,
 1, 0, 0, 0, 21, 0, 0, 0,
 199, 113, 196, 43, 171, 117, 107, 15,
 57, 142, 59, 212, 84, 138, 148, 240,
 1, 0, 0, 0, 18, 0, 0, 0,
 111, 222, 0, 0, 0, 0, 0, 0,
 1, 0, 0, 0, 19, 0, 0, 0,
 53, 130, 156, 173, 0, 0, 0, 0,
 1, 0, 0, 0, 12, 0, 0, 0,
 85, 161, 174, 198, 0, 0, 0, 0,
 1, 0, 0, 0, 13, 0, 0, 0,
 199, 113, 172, 181, 175, 152, 50, 154,
 1, 0, 0, 0, 20, 0, 0, 0,
 0, 156, 173, 69, 0, 228, 10, 69,
 1, 0, 0, 0, 21, 0, 0, 0,
 0, 0, 0, 0, 192, 97, 190, 64,
 0, 0, 0, 0, 128, 92, 145, 64,
 1, 0, 0, 0, 30, 0, 0, 0,
 9, 0, 0, 0, 50, 0, 0, 0,
 9, 0, 0, 0, 50, 0, 0, 0,
 9, 0, 0, 0, 42, 0, 0, 0,
 112, 108, 117, 103, 104, 0, 0, 0,
 120, 121, 122, 122, 121, 0, 0, 0,
 116, 104, 117, 100, 0, 0, 0, 0,
 1, 0, 0, 0, 30, 0, 0, 0,
 9, 0, 0, 0, 34, 0, 0, 0,
 9, 0, 0, 0, 74, 0, 0, 0,
 13, 0, 0, 0, 58, 0, 0, 0,
 111, 111, 112, 115, 0, 0, 0, 0,
 101, 120, 104, 97, 117, 115, 116, 101,
 100, 0, 0, 0, 0, 0, 0, 0,
 114, 102, 99, 51, 48, 57, 50, 0,
 1, 0, 0, 0, 31, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 77, 1, 0, 0, 106, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 217, 0, 0, 0, 106, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 101, 0, 0, 0, 106, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 0, 0, 0, 0, 0, 0, 0, 0,
 115, 116, 114, 117, 99, 116, 108, 105,
 115, 116, 32, 49, 0, 0, 0, 0,
 115, 116, 114, 117, 99, 116, 108, 105,
 115, 116, 32, 50, 0, 0, 0, 0,
 115, 116, 114, 117, 99, 116, 108, 105,
 115, 116, 32, 51, 0, 0, 0, 0,
 1, 0, 0, 0, 19, 0, 0, 0,
 1, 0, 2, 0, 0, 0, 0, 0,
};
