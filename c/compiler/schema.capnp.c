#include "schema.capnp.h"
/* AUTO GENERATED DO NOT EDIT*/

void read_Node(struct Node *s, Node_ptr p) {
	s->id = capn_read64(p.p, 0);
	s->displayName = capn_get_text(p.p, 0);
	s->scopeId = capn_read64(p.p, 8);
	s->nestedNodes.p = capn_getp(p.p, 1);
	s->annotations.p = capn_getp(p.p, 2);
	s->body_tag = (enum Node_body) capn_read16(p.p, 16);

	switch (s->body_tag) {
	case Node_fileNode:
	case Node_structNode:
	case Node_enumNode:
	case Node_interfaceNode:
	case Node_constNode:
	case Node_annotationNode:
		s->body.annotationNode.p = capn_getp(p.p, 3);
		break;
	}
}
int write_Node(const struct Node *s, Node_ptr p) {
	int err = 0;
	err = err || capn_write64(p.p, 0, s->id);
	err = err || capn_set_text(p.p, 0, s->displayName);
	err = err || capn_write64(p.p, 8, s->scopeId);
	err = err || capn_setp(p.p, 1, s->nestedNodes.p);
	err = err || capn_setp(p.p, 2, s->annotations.p);
	err = err || capn_write16(p.p, 16, s->body_tag);

	switch (s->body_tag) {
	case Node_fileNode:
	case Node_structNode:
	case Node_enumNode:
	case Node_interfaceNode:
	case Node_constNode:
	case Node_annotationNode:
		err = err || capn_setp(p.p, 3, s->body.annotationNode.p);
		break;
	}
	return err;
}
void get_Node(struct Node *s, Node_list l, int i) {
	Node_ptr p = {capn_getp(l.p, i)};
	read_Node(s, p);
}
int set_Node(const struct Node *s, Node_list l, int i) {
	Node_ptr p = {capn_getp(l.p, i)};
	return write_Node(s, p);
}

void read_Node_NestedNode(struct Node_NestedNode *s, Node_NestedNode_ptr p) {
	s->name = capn_get_text(p.p, 0);
	s->id = capn_read64(p.p, 0);
}
int write_Node_NestedNode(const struct Node_NestedNode *s, Node_NestedNode_ptr p) {
	int err = 0;
	err = err || capn_set_text(p.p, 0, s->name);
	err = err || capn_write64(p.p, 0, s->id);
	return err;
}
void get_Node_NestedNode(struct Node_NestedNode *s, Node_NestedNode_list l, int i) {
	Node_NestedNode_ptr p = {capn_getp(l.p, i)};
	read_Node_NestedNode(s, p);
}
int set_Node_NestedNode(const struct Node_NestedNode *s, Node_NestedNode_list l, int i) {
	Node_NestedNode_ptr p = {capn_getp(l.p, i)};
	return write_Node_NestedNode(s, p);
}

void read_Type(struct Type *s, Type_ptr p) {
	s->body_tag = (enum Type_body) capn_read16(p.p, 0);

	switch (s->body_tag) {
	case Type_voidType:
	case Type_boolType:
	case Type_int8Type:
	case Type_int16Type:
	case Type_int32Type:
	case Type_int64Type:
	case Type_uint8Type:
	case Type_uint16Type:
	case Type_uint32Type:
	case Type_uint64Type:
	case Type_float32Type:
	case Type_float64Type:
	case Type_textType:
	case Type_dataType:
	case Type_objectType:
		break;
	case Type_enumType:
	case Type_structType:
	case Type_interfaceType:
		s->body.interfaceType = capn_read64(p.p, 8);
		break;
	case Type_listType:
		s->body.listType.p = capn_getp(p.p, 0);
		break;
	}
}
int write_Type(const struct Type *s, Type_ptr p) {
	int err = 0;
	err = err || capn_write16(p.p, 0, s->body_tag);

	switch (s->body_tag) {
	case Type_voidType:
	case Type_boolType:
	case Type_int8Type:
	case Type_int16Type:
	case Type_int32Type:
	case Type_int64Type:
	case Type_uint8Type:
	case Type_uint16Type:
	case Type_uint32Type:
	case Type_uint64Type:
	case Type_float32Type:
	case Type_float64Type:
	case Type_textType:
	case Type_dataType:
	case Type_objectType:
		break;
	case Type_enumType:
	case Type_structType:
	case Type_interfaceType:
		err = err || capn_write64(p.p, 8, s->body.interfaceType);
		break;
	case Type_listType:
		err = err || capn_setp(p.p, 0, s->body.listType.p);
		break;
	}
	return err;
}
void get_Type(struct Type *s, Type_list l, int i) {
	Type_ptr p = {capn_getp(l.p, i)};
	read_Type(s, p);
}
int set_Type(const struct Type *s, Type_list l, int i) {
	Type_ptr p = {capn_getp(l.p, i)};
	return write_Type(s, p);
}

void read_Value(struct Value *s, Value_ptr p) {
	s->body_tag = (enum Value_body) capn_read16(p.p, 0);

	switch (s->body_tag) {
	case Value_voidValue:
	case Value_interfaceValue:
		break;
	case Value_boolValue:
		s->body.boolValue = (capn_read8(p.p, 8) & 1) != 0;
		break;
	case Value_int8Value:
	case Value_uint8Value:
		s->body.uint8Value = capn_read8(p.p, 8);
		break;
	case Value_int16Value:
	case Value_uint16Value:
	case Value_enumValue:
		s->body.enumValue = capn_read16(p.p, 8);
		break;
	case Value_int32Value:
	case Value_uint32Value:
	case Value_float32Value:
		s->body.float32Value = capn_read_float(p.p, 8, 0.0f);
		break;
	case Value_int64Value:
	case Value_uint64Value:
	case Value_float64Value:
		s->body.float64Value = capn_read_double(p.p, 8, 0.0);
		break;
	case Value_textValue:
		s->body.textValue = capn_get_text(p.p, 0);
		break;
	case Value_dataValue:
		s->body.dataValue = capn_get_data(p.p, 0);
		break;
	case Value_listValue:
	case Value_structValue:
	case Value_objectValue:
		s->body.objectValue = capn_getp(p.p, 0);
		break;
	}
}
int write_Value(const struct Value *s, Value_ptr p) {
	int err = 0;
	err = err || capn_write16(p.p, 0, s->body_tag);

	switch (s->body_tag) {
	case Value_voidValue:
	case Value_interfaceValue:
		break;
	case Value_boolValue:
		err = err || capn_write1(p.p, 64, s->body.boolValue);
		break;
	case Value_int8Value:
	case Value_uint8Value:
		err = err || capn_write8(p.p, 8, s->body.uint8Value);
		break;
	case Value_int16Value:
	case Value_uint16Value:
	case Value_enumValue:
		err = err || capn_write16(p.p, 8, s->body.enumValue);
		break;
	case Value_int32Value:
	case Value_uint32Value:
	case Value_float32Value:
		err = err || capn_write_float(p.p, 8, s->body.float32Value, 0.0f);
		break;
	case Value_int64Value:
	case Value_uint64Value:
	case Value_float64Value:
		err = err || capn_write_double(p.p, 4, s->body.float64Value, 0.0);
		break;
	case Value_textValue:
		err = err || capn_set_text(p.p, 0, s->body.textValue);
		break;
	case Value_dataValue:
		err = err || capn_set_data(p.p, 0, s->body.dataValue);
		break;
	case Value_listValue:
	case Value_structValue:
	case Value_objectValue:
		err = err || capn_setp(p.p, 0, s->body.objectValue);
		break;
	}
	return err;
}
void get_Value(struct Value *s, Value_list l, int i) {
	Value_ptr p = {capn_getp(l.p, i)};
	read_Value(s, p);
}
int set_Value(const struct Value *s, Value_list l, int i) {
	Value_ptr p = {capn_getp(l.p, i)};
	return write_Value(s, p);
}

void read_Annotation(struct Annotation *s, Annotation_ptr p) {
	s->id = capn_read64(p.p, 0);
	s->value.p = capn_getp(p.p, 0);
}
int write_Annotation(const struct Annotation *s, Annotation_ptr p) {
	int err = 0;
	err = err || capn_write64(p.p, 0, s->id);
	err = err || capn_setp(p.p, 0, s->value.p);
	return err;
}
void get_Annotation(struct Annotation *s, Annotation_list l, int i) {
	Annotation_ptr p = {capn_getp(l.p, i)};
	read_Annotation(s, p);
}
int set_Annotation(const struct Annotation *s, Annotation_list l, int i) {
	Annotation_ptr p = {capn_getp(l.p, i)};
	return write_Annotation(s, p);
}

void read_FileNode(struct FileNode *s, FileNode_ptr p) {
	s->imports.p = capn_getp(p.p, 0);
}
int write_FileNode(const struct FileNode *s, FileNode_ptr p) {
	int err = 0;
	err = err || capn_setp(p.p, 0, s->imports.p);
	return err;
}
void get_FileNode(struct FileNode *s, FileNode_list l, int i) {
	FileNode_ptr p = {capn_getp(l.p, i)};
	read_FileNode(s, p);
}
int set_FileNode(const struct FileNode *s, FileNode_list l, int i) {
	FileNode_ptr p = {capn_getp(l.p, i)};
	return write_FileNode(s, p);
}

void read_FileNode_Import(struct FileNode_Import *s, FileNode_Import_ptr p) {
	s->id = capn_read64(p.p, 0);
	s->name = capn_get_text(p.p, 0);
}
int write_FileNode_Import(const struct FileNode_Import *s, FileNode_Import_ptr p) {
	int err = 0;
	err = err || capn_write64(p.p, 0, s->id);
	err = err || capn_set_text(p.p, 0, s->name);
	return err;
}
void get_FileNode_Import(struct FileNode_Import *s, FileNode_Import_list l, int i) {
	FileNode_Import_ptr p = {capn_getp(l.p, i)};
	read_FileNode_Import(s, p);
}
int set_FileNode_Import(const struct FileNode_Import *s, FileNode_Import_list l, int i) {
	FileNode_Import_ptr p = {capn_getp(l.p, i)};
	return write_FileNode_Import(s, p);
}

void read_StructNode(struct StructNode *s, StructNode_ptr p) {
	s->dataSectionWordSize = capn_read16(p.p, 0);
	s->pointerSectionSize = capn_read16(p.p, 2);
	s->preferredListEncoding = (enum ElementSize) capn_read16(p.p, 4);
	s->members.p = capn_getp(p.p, 0);
}
int write_StructNode(const struct StructNode *s, StructNode_ptr p) {
	int err = 0;
	err = err || capn_write16(p.p, 0, s->dataSectionWordSize);
	err = err || capn_write16(p.p, 2, s->pointerSectionSize);
	err = err || capn_write16(p.p, 4, (uint16_t) s->preferredListEncoding);
	err = err || capn_setp(p.p, 0, s->members.p);
	return err;
}
void get_StructNode(struct StructNode *s, StructNode_list l, int i) {
	StructNode_ptr p = {capn_getp(l.p, i)};
	read_StructNode(s, p);
}
int set_StructNode(const struct StructNode *s, StructNode_list l, int i) {
	StructNode_ptr p = {capn_getp(l.p, i)};
	return write_StructNode(s, p);
}

void read_StructNode_Member(struct StructNode_Member *s, StructNode_Member_ptr p) {
	s->name = capn_get_text(p.p, 0);
	s->ordinal = capn_read16(p.p, 0);
	s->codeOrder = capn_read16(p.p, 2);
	s->annotations.p = capn_getp(p.p, 1);
	s->body_tag = (enum StructNode_Member_body) capn_read16(p.p, 4);

	switch (s->body_tag) {
	case StructNode_Member_fieldMember:
	case StructNode_Member_unionMember:
		s->body.unionMember.p = capn_getp(p.p, 2);
		break;
	}
}
int write_StructNode_Member(const struct StructNode_Member *s, StructNode_Member_ptr p) {
	int err = 0;
	err = err || capn_set_text(p.p, 0, s->name);
	err = err || capn_write16(p.p, 0, s->ordinal);
	err = err || capn_write16(p.p, 2, s->codeOrder);
	err = err || capn_setp(p.p, 1, s->annotations.p);
	err = err || capn_write16(p.p, 4, s->body_tag);

	switch (s->body_tag) {
	case StructNode_Member_fieldMember:
	case StructNode_Member_unionMember:
		err = err || capn_setp(p.p, 2, s->body.unionMember.p);
		break;
	}
	return err;
}
void get_StructNode_Member(struct StructNode_Member *s, StructNode_Member_list l, int i) {
	StructNode_Member_ptr p = {capn_getp(l.p, i)};
	read_StructNode_Member(s, p);
}
int set_StructNode_Member(const struct StructNode_Member *s, StructNode_Member_list l, int i) {
	StructNode_Member_ptr p = {capn_getp(l.p, i)};
	return write_StructNode_Member(s, p);
}

void read_StructNode_Field(struct StructNode_Field *s, StructNode_Field_ptr p) {
	s->offset = capn_read32(p.p, 0);
	s->type.p = capn_getp(p.p, 0);
	s->defaultValue.p = capn_getp(p.p, 1);
}
int write_StructNode_Field(const struct StructNode_Field *s, StructNode_Field_ptr p) {
	int err = 0;
	err = err || capn_write32(p.p, 0, s->offset);
	err = err || capn_setp(p.p, 0, s->type.p);
	err = err || capn_setp(p.p, 1, s->defaultValue.p);
	return err;
}
void get_StructNode_Field(struct StructNode_Field *s, StructNode_Field_list l, int i) {
	StructNode_Field_ptr p = {capn_getp(l.p, i)};
	read_StructNode_Field(s, p);
}
int set_StructNode_Field(const struct StructNode_Field *s, StructNode_Field_list l, int i) {
	StructNode_Field_ptr p = {capn_getp(l.p, i)};
	return write_StructNode_Field(s, p);
}

void read_StructNode_Union(struct StructNode_Union *s, StructNode_Union_ptr p) {
	s->discriminantOffset = capn_read32(p.p, 0);
	s->members.p = capn_getp(p.p, 0);
}
int write_StructNode_Union(const struct StructNode_Union *s, StructNode_Union_ptr p) {
	int err = 0;
	err = err || capn_write32(p.p, 0, s->discriminantOffset);
	err = err || capn_setp(p.p, 0, s->members.p);
	return err;
}
void get_StructNode_Union(struct StructNode_Union *s, StructNode_Union_list l, int i) {
	StructNode_Union_ptr p = {capn_getp(l.p, i)};
	read_StructNode_Union(s, p);
}
int set_StructNode_Union(const struct StructNode_Union *s, StructNode_Union_list l, int i) {
	StructNode_Union_ptr p = {capn_getp(l.p, i)};
	return write_StructNode_Union(s, p);
}

void read_EnumNode(struct EnumNode *s, EnumNode_ptr p) {
	s->enumerants.p = capn_getp(p.p, 0);
}
int write_EnumNode(const struct EnumNode *s, EnumNode_ptr p) {
	int err = 0;
	err = err || capn_setp(p.p, 0, s->enumerants.p);
	return err;
}
void get_EnumNode(struct EnumNode *s, EnumNode_list l, int i) {
	EnumNode_ptr p = {capn_getp(l.p, i)};
	read_EnumNode(s, p);
}
int set_EnumNode(const struct EnumNode *s, EnumNode_list l, int i) {
	EnumNode_ptr p = {capn_getp(l.p, i)};
	return write_EnumNode(s, p);
}

void read_EnumNode_Enumerant(struct EnumNode_Enumerant *s, EnumNode_Enumerant_ptr p) {
	s->name = capn_get_text(p.p, 0);
	s->codeOrder = capn_read16(p.p, 0);
	s->annotations.p = capn_getp(p.p, 1);
}
int write_EnumNode_Enumerant(const struct EnumNode_Enumerant *s, EnumNode_Enumerant_ptr p) {
	int err = 0;
	err = err || capn_set_text(p.p, 0, s->name);
	err = err || capn_write16(p.p, 0, s->codeOrder);
	err = err || capn_setp(p.p, 1, s->annotations.p);
	return err;
}
void get_EnumNode_Enumerant(struct EnumNode_Enumerant *s, EnumNode_Enumerant_list l, int i) {
	EnumNode_Enumerant_ptr p = {capn_getp(l.p, i)};
	read_EnumNode_Enumerant(s, p);
}
int set_EnumNode_Enumerant(const struct EnumNode_Enumerant *s, EnumNode_Enumerant_list l, int i) {
	EnumNode_Enumerant_ptr p = {capn_getp(l.p, i)};
	return write_EnumNode_Enumerant(s, p);
}

void read_InterfaceNode(struct InterfaceNode *s, InterfaceNode_ptr p) {
	s->methods.p = capn_getp(p.p, 0);
}
int write_InterfaceNode(const struct InterfaceNode *s, InterfaceNode_ptr p) {
	int err = 0;
	err = err || capn_setp(p.p, 0, s->methods.p);
	return err;
}
void get_InterfaceNode(struct InterfaceNode *s, InterfaceNode_list l, int i) {
	InterfaceNode_ptr p = {capn_getp(l.p, i)};
	read_InterfaceNode(s, p);
}
int set_InterfaceNode(const struct InterfaceNode *s, InterfaceNode_list l, int i) {
	InterfaceNode_ptr p = {capn_getp(l.p, i)};
	return write_InterfaceNode(s, p);
}

void read_InterfaceNode_Method(struct InterfaceNode_Method *s, InterfaceNode_Method_ptr p) {
	s->name = capn_get_text(p.p, 0);
	s->codeOrder = capn_read16(p.p, 0);
	s->params.p = capn_getp(p.p, 1);
	s->requiredParamCount = capn_read16(p.p, 2);
	s->returnType.p = capn_getp(p.p, 2);
	s->annotations.p = capn_getp(p.p, 3);
}
int write_InterfaceNode_Method(const struct InterfaceNode_Method *s, InterfaceNode_Method_ptr p) {
	int err = 0;
	err = err || capn_set_text(p.p, 0, s->name);
	err = err || capn_write16(p.p, 0, s->codeOrder);
	err = err || capn_setp(p.p, 1, s->params.p);
	err = err || capn_write16(p.p, 2, s->requiredParamCount);
	err = err || capn_setp(p.p, 2, s->returnType.p);
	err = err || capn_setp(p.p, 3, s->annotations.p);
	return err;
}
void get_InterfaceNode_Method(struct InterfaceNode_Method *s, InterfaceNode_Method_list l, int i) {
	InterfaceNode_Method_ptr p = {capn_getp(l.p, i)};
	read_InterfaceNode_Method(s, p);
}
int set_InterfaceNode_Method(const struct InterfaceNode_Method *s, InterfaceNode_Method_list l, int i) {
	InterfaceNode_Method_ptr p = {capn_getp(l.p, i)};
	return write_InterfaceNode_Method(s, p);
}

void read_InterfaceNode_Method_Param(struct InterfaceNode_Method_Param *s, InterfaceNode_Method_Param_ptr p) {
	s->name = capn_get_text(p.p, 0);
	s->type.p = capn_getp(p.p, 1);
	s->defaultValue.p = capn_getp(p.p, 2);
	s->annotations.p = capn_getp(p.p, 3);
}
int write_InterfaceNode_Method_Param(const struct InterfaceNode_Method_Param *s, InterfaceNode_Method_Param_ptr p) {
	int err = 0;
	err = err || capn_set_text(p.p, 0, s->name);
	err = err || capn_setp(p.p, 1, s->type.p);
	err = err || capn_setp(p.p, 2, s->defaultValue.p);
	err = err || capn_setp(p.p, 3, s->annotations.p);
	return err;
}
void get_InterfaceNode_Method_Param(struct InterfaceNode_Method_Param *s, InterfaceNode_Method_Param_list l, int i) {
	InterfaceNode_Method_Param_ptr p = {capn_getp(l.p, i)};
	read_InterfaceNode_Method_Param(s, p);
}
int set_InterfaceNode_Method_Param(const struct InterfaceNode_Method_Param *s, InterfaceNode_Method_Param_list l, int i) {
	InterfaceNode_Method_Param_ptr p = {capn_getp(l.p, i)};
	return write_InterfaceNode_Method_Param(s, p);
}

void read_ConstNode(struct ConstNode *s, ConstNode_ptr p) {
	s->type.p = capn_getp(p.p, 0);
	s->value.p = capn_getp(p.p, 1);
}
int write_ConstNode(const struct ConstNode *s, ConstNode_ptr p) {
	int err = 0;
	err = err || capn_setp(p.p, 0, s->type.p);
	err = err || capn_setp(p.p, 1, s->value.p);
	return err;
}
void get_ConstNode(struct ConstNode *s, ConstNode_list l, int i) {
	ConstNode_ptr p = {capn_getp(l.p, i)};
	read_ConstNode(s, p);
}
int set_ConstNode(const struct ConstNode *s, ConstNode_list l, int i) {
	ConstNode_ptr p = {capn_getp(l.p, i)};
	return write_ConstNode(s, p);
}

void read_AnnotationNode(struct AnnotationNode *s, AnnotationNode_ptr p) {
	s->type.p = capn_getp(p.p, 0);
	s->targetsFile = (capn_read8(p.p, 0) & 1) != 0;
	s->targetsConst = (capn_read8(p.p, 0) & 2) != 0;
	s->targetsEnum = (capn_read8(p.p, 0) & 4) != 0;
	s->targetsEnumerant = (capn_read8(p.p, 0) & 8) != 0;
	s->targetsStruct = (capn_read8(p.p, 0) & 16) != 0;
	s->targetsField = (capn_read8(p.p, 0) & 32) != 0;
	s->targetsUnion = (capn_read8(p.p, 0) & 64) != 0;
	s->targetsInterface = (capn_read8(p.p, 0) & 128) != 0;
	s->targetsMethod = (capn_read8(p.p, 1) & 1) != 0;
	s->targetsParam = (capn_read8(p.p, 1) & 2) != 0;
	s->targetsAnnotation = (capn_read8(p.p, 1) & 4) != 0;
}
int write_AnnotationNode(const struct AnnotationNode *s, AnnotationNode_ptr p) {
	int err = 0;
	err = err || capn_setp(p.p, 0, s->type.p);
	err = err || capn_write1(p.p, 0, s->targetsFile);
	err = err || capn_write1(p.p, 1, s->targetsConst);
	err = err || capn_write1(p.p, 2, s->targetsEnum);
	err = err || capn_write1(p.p, 3, s->targetsEnumerant);
	err = err || capn_write1(p.p, 4, s->targetsStruct);
	err = err || capn_write1(p.p, 5, s->targetsField);
	err = err || capn_write1(p.p, 6, s->targetsUnion);
	err = err || capn_write1(p.p, 7, s->targetsInterface);
	err = err || capn_write1(p.p, 8, s->targetsMethod);
	err = err || capn_write1(p.p, 9, s->targetsParam);
	err = err || capn_write1(p.p, 10, s->targetsAnnotation);
	return err;
}
void get_AnnotationNode(struct AnnotationNode *s, AnnotationNode_list l, int i) {
	AnnotationNode_ptr p = {capn_getp(l.p, i)};
	read_AnnotationNode(s, p);
}
int set_AnnotationNode(const struct AnnotationNode *s, AnnotationNode_list l, int i) {
	AnnotationNode_ptr p = {capn_getp(l.p, i)};
	return write_AnnotationNode(s, p);
}

void read_CodeGeneratorRequest(struct CodeGeneratorRequest *s, CodeGeneratorRequest_ptr p) {
	s->nodes.p = capn_getp(p.p, 0);
	s->requestedFiles.p = capn_getp(p.p, 1);
}
int write_CodeGeneratorRequest(const struct CodeGeneratorRequest *s, CodeGeneratorRequest_ptr p) {
	int err = 0;
	err = err || capn_setp(p.p, 0, s->nodes.p);
	err = err || capn_setp(p.p, 1, s->requestedFiles.p);
	return err;
}
void get_CodeGeneratorRequest(struct CodeGeneratorRequest *s, CodeGeneratorRequest_list l, int i) {
	CodeGeneratorRequest_ptr p = {capn_getp(l.p, i)};
	read_CodeGeneratorRequest(s, p);
}
int set_CodeGeneratorRequest(const struct CodeGeneratorRequest *s, CodeGeneratorRequest_list l, int i) {
	CodeGeneratorRequest_ptr p = {capn_getp(l.p, i)};
	return write_CodeGeneratorRequest(s, p);
}
