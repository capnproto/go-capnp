#include "schema.h"

void read_Node(const struct Node_ptr *p, struct Node *s) {
	s->id = capn_read64(&p->p, 0);
	s->displayName = capn_get_text(&p->p, 0);
	s->scopeId = capn_read64(&p->p, 8);
	s->nestedNodes = capn_getp(&p->p, 1);
	s->annotations = capn_getp(&p->p, 2);
	s->body_tag = capn_read16(&p->p, 16);
	switch (s->body_tag) {
	case Node_fileNode:
		s->body.fileNode.p = capn_getp(&p->p, 3);
		break;
	case Node_structNode:
		s->body.structNode.p = capn_getp(&p->p, 3);
		break;
	case Node_enumNode:
		s->body.enumNode.p = capn_getp(&p->p, 3);
		break;
	case Node_interfaceNode:
		s->body.interfaceNode.p = capn_getp(&p->p, 3);
		break;
	case Node_constNode:
		s->body.constNode.p = capn_getp(&p->p, 3);
		break;
	case Node_annotationNode:
		s->body.annotationNode.p = capn_getp(&p->p, 3);
		break;
	default:
		break;
	}
}

void read_Node_NestedNode(const struct Node_NestedNode_ptr *p, struct Node_NestedNode *s) {
	s->name = capn_get_text(&p->p, 0);
	s->id = capn_read64(&p->p, 0);
}

void read_Type(const struct Type_ptr *p, struct Type *s) {
	s->body_tag = capn_read16(&p->p, 0);
	switch (s->body_tag) {
	case Type_listType:
		s->body.listType.p = capn_getp(&p->p, 0);
		break;
	case Type_enumType:
		s->body.enumType = capn_read64(&p->p, 0);
		break;
	case Type_structType:
		s->body.structType = capn_read64(&p->p, 0);
		break;
	case Type_interfaceType:
		s->body.interfaceType = capn_read64(&p->p, 0);
		break;
	default:
		break;
	}
}

void read_Value(const struct Value_ptr *p, struct Value *s) {
	s->body_tag = capn_read16(&p->p, 0);
	switch (s->body_tag) {
	case Value_boolValue:
		s->body.boolValue = (capn_read8(&p->p, 8) & 1) != 0;
		break;
	case Value_int8Value:
		s->body.int8Value = (int8_t) capn_read8(&p->p, 8);
		break;
	case Value_int16Value:
		s->body.int16Value = (int16_t) capn_read16(&p->p, 8);
		break;
	case Value_int32Value:
		s->body.int32Value = (int32_t) capn_read32(&p->p, 8);
		break;
	case Value_int64Value:
		s->body.int64Value = (int64_t) capn_read64(&p->p, 8);
		break;
	case Value_uint8Value:
		s->body.uint8Value = capn_read8(&p->p, 8);
		break;
	case Value_uint16Value:
		s->body.uint16Value = capn_read16(&p->p, 8);
		break;
	case Value_uint32Value:
		s->body.uint32Value = capn_read32(&p->p, 8);
		break;
	case Value_uint64Value:
		s->body.uint64Value = capn_read64(&p->p, 8);
		break;
	case Value_float32Value:
		s->body.float32Value = capn_read_float(&p->p, 8, 0.0f);
		break;
	case Value_float64Value:
		s->body.float64Value = capn_read_double(&p->p, 8, 0.0);
		break;
	case Value_textValue:
		s->body.textValue = capn_get_text(&p->p, 0);
		break;
	case Value_dataValue:
		s->body.dataValue = capn_get_data(&p->p, 0);
		break;
	case Value_listValue:
		s->body.listValue = capn_getp(&p->p, 0);
		break;
	case Value_enumValue:
		s->body.enumValue = capn_read16(&p->p, 8);
		break;
	case Value_structValue:
		s->body.structValue = capn_getp(&p->p, 0);
		break;
	case Value_objectValue:
		s->body.objectValue = capn_getp(&p->p, 0);
		break;
	default:
		break;
	}
}

void read_Annotation(const struct Annotation_ptr *p, struct Annotation *s) {
	s->id = capn_read64(&p->p, 0);
	s->value.p = capn_getp(&p->p, 0);
}

void read_FileNode(const struct FileNode_ptr *p, struct FileNode *s) {
	s->imports = capn_getp(&p->p, 0);
}

void read_FileNode_Import(const struct FileNode_Import_ptr *p, struct FileNode_Import *s) {
	s->id = capn_read64(&p->p, 0);
	s->name = capn_get_text(&p->p, 0);
}

void read_StructNode(const struct StructNode_ptr *p, struct StructNode *s) {
	s->dataSectionWordSize = capn_read16(&p->p, 0);
	s->pointerSectionSize = capn_read16(&p->p, 2);
	s->preferredListEncoding = (enum ElementSize) capn_read16(&p->p, 4);
	s->members = capn_getp(&p->p, 0);
}

void read_StructNode_Member(const struct StructNode_Member_ptr *p, struct StructNode_Member *s) {
	s->name = capn_get_text(&p->p, 0);
	s->ordinal = capn_read16(&p->p, 0);
	s->codeOrder = capn_read16(&p->p, 2);
	s->annotations = capn_getp(&p->p, 1);
	s->body_tag = (enum StructNode_Member_body) capn_read16(&p->p, 4);
	switch (s->body_tag) {
	case StructNode_Member_fieldMember:
		s->body.fieldMember.p = capn_getp(&p->p, 2);
		break;
	case StructNode_Member_unionMember:
		s->body.unionMember.p = capn_getp(&p->p, 2);
		break;
	default:
		break;
	}
}

void read_StructNode_Field(const struct StructNode_Field_ptr *p, struct StructNode_Field *s) {
	s->offset = capn_read32(&p->p, 0);
	s->type.p = capn_getp(&p->p, 0);
	s->defaultValue.p = capn_getp(&p->p, 1);
}

void read_StructNode_Union(const struct StructNode_Union_ptr *p, struct StructNode_Union *s) {
	s->discriminantOffset = capn_read32(&p->p, 0);
	s->members = capn_getp(&p->p, 0);
}

void read_EnumNode(const struct EnumNode_ptr *p, struct EnumNode *s) {
	s->enumerants = capn_getp(&p->p, 0);
}

void read_EnumNode_Enumerant(const struct EnumNode_Enumerant_ptr *p, struct EnumNode_Enumerant *s) {
	s->name = capn_get_text(&p->p, 0);
	s->codeOrder = capn_read16(&p->p, 0);
	s->annotations = capn_getp(&p->p, 1);
}

void read_InterfaceNode(const struct InterfaceNode_ptr *p, struct InterfaceNode *s) {
	s->methods = capn_getp(&p->p, 0);
}

void read_InterfaceNode_Method(const struct InterfaceNode_Method_ptr *p, struct InterfaceNode_Method *s) {
	s->name = capn_get_text(&p->p, 0);
	s->codeOrder = capn_read16(&p->p, 0);
	s->params = capn_getp(&p->p, 1);
	s->requiredParamCount = capn_read16(&p->p, 2);
	s->returnType.p = capn_getp(&p->p, 2);
	s->annotations = capn_getp(&p->p, 3);
}

void read_InterfaceNode_Method_Param(const struct InterfaceNode_Method_Param_ptr *p, struct InterfaceNode_Method_Param *s) {
	s->name = capn_get_text(&p->p, 0);
	s->type.p = capn_getp(&p->p, 1);
	s->defaultValue.p = capn_getp(&p->p, 2);
	s->annotations = capn_getp(&p->p, 3);
}

void read_ConstNode(const struct ConstNode_ptr *p, struct ConstNode *s) {
	s->type.p = capn_getp(&p->p, 0);
	s->value.p = capn_getp(&p->p, 1);
}

void read_AnnotationNode(const struct AnnotationNode_ptr *p, struct AnnotationNode *s) {
	s->type.p = capn_getp(&p->p, 0);
	s->targetsFile = (capn_read8(&p->p, 0) & 1) != 0;
	s->targetsConst = (capn_read8(&p->p, 0) & 3) != 0;
	s->targetsEnum = (capn_read8(&p->p, 0) & 4) != 0;
	s->targetsEnumerant = (capn_read8(&p->p, 0) & 8) != 0;
	s->targetsStruct = (capn_read8(&p->p, 0) & 16) != 0;
	s->targetsField = (capn_read8(&p->p, 0) & 32) != 0;
	s->targetsUnion = (capn_read8(&p->p, 0) & 64) != 0;
	s->targetsInterface = (capn_read8(&p->p, 0) & 128) != 0;
	s->targetsMethod = (capn_read8(&p->p, 1) & 1) != 0;
	s->targetsParam = (capn_read8(&p->p, 1) & 2) != 0;
	s->targetsAnnotation = (capn_read8(&p->p, 1) & 4) != 0;
}

void read_CodeGeneratorRequest(const struct CodeGeneratorRequest_ptr *p, struct CodeGeneratorRequest *s) {
	s->nodes = capn_getp(&p->p, 0);
	s->requestedFiles = capn_getp(&p->p, 1);
}

