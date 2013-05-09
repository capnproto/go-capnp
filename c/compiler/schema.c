#include "schema.h"

void get_Node(struct Node *s, Node_list l, int i) {
	Node_ptr p = {capn_getp(l.p, i)};
	read_Node(s, p);
}

void get_Node_NestedNode(struct Node_NestedNode *s, Node_NestedNode_list l, int i) {
	Node_NestedNode_ptr p = {capn_getp(l.p, i)};
	read_Node_NestedNode(s, p);
}

void get_Type(struct Type *s, Type_list l, int i) {
	Type_ptr p = {capn_getp(l.p, i)};
	read_Type(s, p);
}

void get_Value(struct Value *s, Value_list l, int i) {
	Value_ptr p = {capn_getp(l.p, i)};
	read_Value(s, p);
}

void get_Annotation(struct Annotation *s, Annotation_list l, int i) {
	Annotation_ptr p = {capn_getp(l.p, i)};
	read_Annotation(s, p);
}

void get_FileNode(struct FileNode *s, FileNode_list l, int i) {
	FileNode_ptr p = {capn_getp(l.p, i)};
	read_FileNode(s, p);
}

void get_FileNode_Import(struct FileNode_Import *s, FileNode_Import_list l, int i) {
	FileNode_Import_ptr p = {capn_getp(l.p, i)};
	read_FileNode_Import(s, p);
}

void get_StructNode(struct StructNode *s, StructNode_list l, int i) {
	StructNode_ptr p = {capn_getp(l.p, i)};
	read_StructNode(s, p);
}

void get_StructNode_Member(struct StructNode_Member *s, StructNode_Member_list l, int i) {
	StructNode_Member_ptr p = {capn_getp(l.p, i)};
	read_StructNode_Member(s, p);
}

void get_StructNode_Field(struct StructNode_Field *s, StructNode_Field_list l, int i) {
	StructNode_Field_ptr p = {capn_getp(l.p, i)};
	read_StructNode_Field(s, p);
}

void get_StructNode_Union(struct StructNode_Union *s, StructNode_Union_list l, int i) {
	StructNode_Union_ptr p = {capn_getp(l.p, i)};
	read_StructNode_Union(s, p);
}

void get_EnumNode(struct EnumNode *s, EnumNode_list l, int i) {
	EnumNode_ptr p = {capn_getp(l.p, i)};
	read_EnumNode(s, p);
}

void get_EnumNode_Enumerant(struct EnumNode_Enumerant *s, EnumNode_Enumerant_list l, int i) {
	EnumNode_Enumerant_ptr p = {capn_getp(l.p, i)};
	read_EnumNode_Enumerant(s, p);
}

void get_InterfaceNode(struct InterfaceNode *s, InterfaceNode_list l, int i) {
	InterfaceNode_ptr p = {capn_getp(l.p, i)};
	read_InterfaceNode(s, p);
}

void get_InterfaceNode_Method(struct InterfaceNode_Method *s, InterfaceNode_Method_list l, int i) {
	InterfaceNode_Method_ptr p = {capn_getp(l.p, i)};
	read_InterfaceNode_Method(s, p);
}

void get_InterfaceNode_Method_Param(struct InterfaceNode_Method_Param *s, InterfaceNode_Method_Param_list l, int i) {
	InterfaceNode_Method_Param_ptr p = {capn_getp(l.p, i)};
	read_InterfaceNode_Method_Param(s, p);
}

void get_ConstNode(struct ConstNode *s, ConstNode_list l, int i) {
	ConstNode_ptr p = {capn_getp(l.p, i)};
	read_ConstNode(s, p);
}

void get_AnnotationNode(struct AnnotationNode *s, AnnotationNode_list l, int i) {
	AnnotationNode_ptr p = {capn_getp(l.p, i)};
	read_AnnotationNode(s, p);
}

void get_CodeGeneratorRequest(struct CodeGeneratorRequest *s, CodeGeneratorRequest_list l, int i) {
	CodeGeneratorRequest_ptr p = {capn_getp(l.p, i)};
	read_CodeGeneratorRequest(s, p);
}



void read_Node(struct Node *s, Node_ptr p) {
	s->id = capn_read64(p.p, 0);
	s->displayName = capn_get_text(p.p, 0);
	s->scopeId = capn_read64(p.p, 8);
	s->nestedNodes.p = capn_getp(p.p, 1);
	s->annotations.p = capn_getp(p.p, 2);
	s->body_tag = capn_read16(p.p, 16);
	switch (s->body_tag) {
	case Node_fileNode:
		s->body.fileNode.p = capn_getp(p.p, 3);
		break;
	case Node_structNode:
		s->body.structNode.p = capn_getp(p.p, 3);
		break;
	case Node_enumNode:
		s->body.enumNode.p = capn_getp(p.p, 3);
		break;
	case Node_interfaceNode:
		s->body.interfaceNode.p = capn_getp(p.p, 3);
		break;
	case Node_constNode:
		s->body.constNode.p = capn_getp(p.p, 3);
		break;
	case Node_annotationNode:
		s->body.annotationNode.p = capn_getp(p.p, 3);
		break;
	default:
		break;
	}
}

void read_Node_NestedNode(struct Node_NestedNode *s, Node_NestedNode_ptr p) {
	s->name = capn_get_text(p.p, 0);
	s->id = capn_read64(p.p, 0);
}

void read_Type(struct Type *s, Type_ptr p) {
	s->body_tag = capn_read16(p.p, 0);
	switch (s->body_tag) {
	case Type_listType:
		s->body.listType.p = capn_getp(p.p, 0);
		break;
	case Type_enumType:
		s->body.enumType = capn_read64(p.p, 0);
		break;
	case Type_structType:
		s->body.structType = capn_read64(p.p, 0);
		break;
	case Type_interfaceType:
		s->body.interfaceType = capn_read64(p.p, 0);
		break;
	default:
		break;
	}
}

void read_Value(struct Value *s, Value_ptr p) {
	s->body_tag = capn_read16(p.p, 0);
	switch (s->body_tag) {
	case Value_boolValue:
		s->body.boolValue = (capn_read8(p.p, 8) & 1) != 0;
		break;
	case Value_int8Value:
		s->body.int8Value = (int8_t) capn_read8(p.p, 8);
		break;
	case Value_int16Value:
		s->body.int16Value = (int16_t) capn_read16(p.p, 8);
		break;
	case Value_int32Value:
		s->body.int32Value = (int32_t) capn_read32(p.p, 8);
		break;
	case Value_int64Value:
		s->body.int64Value = (int64_t) capn_read64(p.p, 8);
		break;
	case Value_uint8Value:
		s->body.uint8Value = capn_read8(p.p, 8);
		break;
	case Value_uint16Value:
		s->body.uint16Value = capn_read16(p.p, 8);
		break;
	case Value_uint32Value:
		s->body.uint32Value = capn_read32(p.p, 8);
		break;
	case Value_uint64Value:
		s->body.uint64Value = capn_read64(p.p, 8);
		break;
	case Value_float32Value:
		s->body.float32Value = capn_read_float(p.p, 8, 0.0f);
		break;
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
		s->body.listValue = capn_getp(p.p, 0);
		break;
	case Value_enumValue:
		s->body.enumValue = capn_read16(p.p, 8);
		break;
	case Value_structValue:
		s->body.structValue = capn_getp(p.p, 0);
		break;
	case Value_objectValue:
		s->body.objectValue = capn_getp(p.p, 0);
		break;
	default:
		break;
	}
}

void read_Annotation(struct Annotation *s, Annotation_ptr p) {
	s->id = capn_read64(p.p, 0);
	s->value.p = capn_getp(p.p, 0);
}

void read_FileNode(struct FileNode *s, FileNode_ptr p) {
	s->imports.p = capn_getp(p.p, 0);
}

void read_FileNode_Import(struct FileNode_Import *s, FileNode_Import_ptr p) {
	s->id = capn_read64(p.p, 0);
	s->name = capn_get_text(p.p, 0);
}

void read_StructNode(struct StructNode *s, StructNode_ptr p) {
	s->dataSectionWordSize = capn_read16(p.p, 0);
	s->pointerSectionSize = capn_read16(p.p, 2);
	s->preferredListEncoding = (enum ElementSize) capn_read16(p.p, 4);
	s->members.p = capn_getp(p.p, 0);
}

void read_StructNode_Member(struct StructNode_Member *s, StructNode_Member_ptr p) {
	s->name = capn_get_text(p.p, 0);
	s->ordinal = capn_read16(p.p, 0);
	s->codeOrder = capn_read16(p.p, 2);
	s->annotations.p = capn_getp(p.p, 1);
	s->body_tag = (enum StructNode_Member_body) capn_read16(p.p, 4);
	switch (s->body_tag) {
	case StructNode_Member_fieldMember:
		s->body.fieldMember.p = capn_getp(p.p, 2);
		break;
	case StructNode_Member_unionMember:
		s->body.unionMember.p = capn_getp(p.p, 2);
		break;
	default:
		break;
	}
}

void read_StructNode_Field(struct StructNode_Field *s, StructNode_Field_ptr p) {
	s->offset = capn_read32(p.p, 0);
	s->type.p = capn_getp(p.p, 0);
	s->defaultValue.p = capn_getp(p.p, 1);
}

void read_StructNode_Union(struct StructNode_Union *s, StructNode_Union_ptr p) {
	s->discriminantOffset = capn_read32(p.p, 0);
	s->members.p = capn_getp(p.p, 0);
}

void read_EnumNode(struct EnumNode *s, EnumNode_ptr p) {
	s->enumerants.p = capn_getp(p.p, 0);
}

void read_EnumNode_Enumerant(struct EnumNode_Enumerant *s, EnumNode_Enumerant_ptr p) {
	s->name = capn_get_text(p.p, 0);
	s->codeOrder = capn_read16(p.p, 0);
	s->annotations.p = capn_getp(p.p, 1);
}

void read_InterfaceNode(struct InterfaceNode *s, InterfaceNode_ptr p) {
	s->methods.p = capn_getp(p.p, 0);
}

void read_InterfaceNode_Method(struct InterfaceNode_Method *s, InterfaceNode_Method_ptr p) {
	s->name = capn_get_text(p.p, 0);
	s->codeOrder = capn_read16(p.p, 0);
	s->params.p = capn_getp(p.p, 1);
	s->requiredParamCount = capn_read16(p.p, 2);
	s->returnType.p = capn_getp(p.p, 2);
	s->annotations.p = capn_getp(p.p, 3);
}

void read_InterfaceNode_Method_Param(struct InterfaceNode_Method_Param *s, InterfaceNode_Method_Param_ptr p) {
	s->name = capn_get_text(p.p, 0);
	s->type.p = capn_getp(p.p, 1);
	s->defaultValue.p = capn_getp(p.p, 2);
	s->annotations.p = capn_getp(p.p, 3);
}

void read_ConstNode(struct ConstNode *s, ConstNode_ptr p) {
	s->type.p = capn_getp(p.p, 0);
	s->value.p = capn_getp(p.p, 1);
}

void read_AnnotationNode(struct AnnotationNode *s, AnnotationNode_ptr p) {
	s->type.p = capn_getp(p.p, 0);
	s->targetsFile = (capn_read8(p.p, 0) & 1) != 0;
	s->targetsConst = (capn_read8(p.p, 0) & 3) != 0;
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

void read_CodeGeneratorRequest(struct CodeGeneratorRequest *s, CodeGeneratorRequest_ptr p) {
	s->nodes.p = capn_getp(p.p, 0);
	s->requestedFiles.p = capn_getp(p.p, 1);
}

