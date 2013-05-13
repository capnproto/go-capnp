#ifndef CAPN_B471DF2F45CA32C7
#define CAPN_B471DF2F45CA32C7
/* AUTO GENERATED - DO NOT EDIT */
#include <capn.h>

#ifdef __cplusplus
extern "C" {
#endif

struct Node;
struct Node_NestedNode;
struct Type;
struct Value;
struct Annotation;
struct FileNode;
struct FileNode_Import;
struct StructNode;
struct StructNode_Member;
struct StructNode_Field;
struct StructNode_Union;
struct EnumNode;
struct EnumNode_Enumerant;
struct InterfaceNode;
struct InterfaceNode_Method;
struct InterfaceNode_Method_Param;
struct ConstNode;
struct AnnotationNode;
struct CodeGeneratorRequest;

typedef struct {capn_ptr p;} Node_ptr;
typedef struct {capn_ptr p;} Node_NestedNode_ptr;
typedef struct {capn_ptr p;} Type_ptr;
typedef struct {capn_ptr p;} Value_ptr;
typedef struct {capn_ptr p;} Annotation_ptr;
typedef struct {capn_ptr p;} FileNode_ptr;
typedef struct {capn_ptr p;} FileNode_Import_ptr;
typedef struct {capn_ptr p;} StructNode_ptr;
typedef struct {capn_ptr p;} StructNode_Member_ptr;
typedef struct {capn_ptr p;} StructNode_Field_ptr;
typedef struct {capn_ptr p;} StructNode_Union_ptr;
typedef struct {capn_ptr p;} EnumNode_ptr;
typedef struct {capn_ptr p;} EnumNode_Enumerant_ptr;
typedef struct {capn_ptr p;} InterfaceNode_ptr;
typedef struct {capn_ptr p;} InterfaceNode_Method_ptr;
typedef struct {capn_ptr p;} InterfaceNode_Method_Param_ptr;
typedef struct {capn_ptr p;} ConstNode_ptr;
typedef struct {capn_ptr p;} AnnotationNode_ptr;
typedef struct {capn_ptr p;} CodeGeneratorRequest_ptr;

typedef struct {capn_ptr p;} Node_list;
typedef struct {capn_ptr p;} Node_NestedNode_list;
typedef struct {capn_ptr p;} Type_list;
typedef struct {capn_ptr p;} Value_list;
typedef struct {capn_ptr p;} Annotation_list;
typedef struct {capn_ptr p;} FileNode_list;
typedef struct {capn_ptr p;} FileNode_Import_list;
typedef struct {capn_ptr p;} StructNode_list;
typedef struct {capn_ptr p;} StructNode_Member_list;
typedef struct {capn_ptr p;} StructNode_Field_list;
typedef struct {capn_ptr p;} StructNode_Union_list;
typedef struct {capn_ptr p;} EnumNode_list;
typedef struct {capn_ptr p;} EnumNode_Enumerant_list;
typedef struct {capn_ptr p;} InterfaceNode_list;
typedef struct {capn_ptr p;} InterfaceNode_Method_list;
typedef struct {capn_ptr p;} InterfaceNode_Method_Param_list;
typedef struct {capn_ptr p;} ConstNode_list;
typedef struct {capn_ptr p;} AnnotationNode_list;
typedef struct {capn_ptr p;} CodeGeneratorRequest_list;



enum Node_body {
	Node_fileNode = 0,
	Node_structNode = 1,
	Node_enumNode = 2,
	Node_interfaceNode = 3,
	Node_constNode = 4,
	Node_annotationNode = 5
};

struct Node {
	uint64_t id;
	capn_text displayName;
	uint64_t scopeId;
	Node_NestedNode_list nestedNodes;
	Annotation_list annotations;
	enum Node_body body_tag;
	union {
		FileNode_ptr fileNode;
		StructNode_ptr structNode;
		EnumNode_ptr enumNode;
		InterfaceNode_ptr interfaceNode;
		ConstNode_ptr constNode;
		AnnotationNode_ptr annotationNode;
	} body;
};

struct Node_NestedNode {
	capn_text name;
	uint64_t id;
};

enum Type_body {
	Type_voidType = 0,
	Type_boolType = 1,
	Type_int8Type = 2,
	Type_int16Type = 3,
	Type_int32Type = 4,
	Type_int64Type = 5,
	Type_uint8Type = 6,
	Type_uint16Type = 7,
	Type_uint32Type = 8,
	Type_uint64Type = 9,
	Type_float32Type = 10,
	Type_float64Type = 11,
	Type_textType = 12,
	Type_dataType = 13,
	Type_listType = 14,
	Type_enumType = 15,
	Type_structType = 16,
	Type_interfaceType = 17,
	Type_objectType = 18
};

struct Type {
	enum Type_body body_tag;
	union {
		Type_ptr listType;
		uint64_t enumType;
		uint64_t structType;
		uint64_t interfaceType;
	} body;
};

enum Value_body {
	Value_voidValue = 9,
	Value_boolValue = 1,
	Value_int8Value = 2,
	Value_int16Value = 3,
	Value_int32Value = 4,
	Value_int64Value = 5,
	Value_uint8Value = 6,
	Value_uint16Value = 7,
	Value_uint32Value = 8,
	Value_uint64Value = 0,
	Value_float32Value = 10,
	Value_float64Value = 11,
	Value_textValue = 12,
	Value_dataValue = 13,
	Value_listValue = 14,
	Value_enumValue = 15,
	Value_structValue = 16,
	Value_interfaceValue = 17,
	Value_objectValue = 18
};

struct Value {
	enum Value_body body_tag;
	union {
		unsigned boolValue:1;
		int8_t int8Value;
		int16_t int16Value;
		int32_t int32Value;
		int64_t int64Value;
		uint8_t uint8Value;
		uint16_t uint16Value;
		uint32_t uint32Value;
		uint64_t uint64Value;
		float float32Value;
		double float64Value;
		capn_text textValue;
		capn_data dataValue;
		capn_ptr listValue;
		uint16_t enumValue;
		capn_ptr structValue;
		capn_ptr objectValue;
	} body;
};

struct Annotation {
	uint64_t id;
	Value_ptr value;
};

struct FileNode {
	FileNode_Import_list imports;
};

struct FileNode_Import {
	uint64_t id;
	capn_text name;
};

enum ElementSize {
	ElementSize_empty = 0,
	ElementSize_bit = 1,
	ElementSize_byte = 2,
	ElementSize_twoBytes = 3,
	ElementSize_fourBytes = 4,
	ElementSize_eightBytes = 5,
	ElementSize_pointer = 6,
	ElementSize_inlineComposite = 7
};

struct StructNode {
	uint16_t dataSectionWordSize;
	uint16_t pointerSectionSize;
	enum ElementSize preferredListEncoding;
	StructNode_Member_list members;
};

enum StructNode_Member_body {
	StructNode_Member_fieldMember = 0,
	StructNode_Member_unionMember = 1
};

struct StructNode_Member {
	capn_text name;
	uint16_t ordinal;
	uint16_t codeOrder;
	Annotation_list annotations;
	enum StructNode_Member_body body_tag;
	union {
		StructNode_Field_ptr fieldMember;
		StructNode_Union_ptr unionMember;
	} body;
};

struct StructNode_Field {
	uint32_t offset;
	Type_ptr type;
	Value_ptr defaultValue;
};

struct StructNode_Union {
	uint32_t discriminantOffset;
	StructNode_Member_list members;
};

struct EnumNode {
	EnumNode_Enumerant_list enumerants;
};

struct EnumNode_Enumerant {
	capn_text name;
	uint16_t codeOrder;
	Annotation_list annotations;
};

struct InterfaceNode {
	InterfaceNode_Method_list methods;
};

struct InterfaceNode_Method {
	capn_text name;
	uint16_t codeOrder;
	InterfaceNode_Method_Param_list params;
	uint16_t requiredParamCount;
	Type_ptr returnType;
	Annotation_list annotations;
};

struct InterfaceNode_Method_Param {
	capn_text name;
	Type_ptr type;
	Value_ptr defaultValue;
	Annotation_list annotations;
};

struct ConstNode {
	Type_ptr type;
	Value_ptr value;
};

struct AnnotationNode {
	Type_ptr type;
	unsigned targetsFile:1;
	unsigned targetsConst:1;
	unsigned targetsEnum:1;
	unsigned targetsEnumerant:1;
	unsigned targetsStruct:1;
	unsigned targetsField:1;
	unsigned targetsUnion:1;
	unsigned targetsInterface:1;
	unsigned targetsMethod:1;
	unsigned targetsParam:1;
	unsigned targetsAnnotation:1;
};

struct CodeGeneratorRequest {
	Node_list nodes;
	capn_list64 requestedFiles;
};

Node_ptr new_Node(struct capn_segment*);
Node_NestedNode_ptr new_Node_NestedNode(struct capn_segment*);
Type_ptr new_Type(struct capn_segment*);
Value_ptr new_Value(struct capn_segment*);
Annotation_ptr new_Annotation(struct capn_segment*);
FileNode_ptr new_FileNode(struct capn_segment*);
FileNode_Import_ptr new_FileNode_Import(struct capn_segment*);
StructNode_ptr new_StructNode(struct capn_segment*);
StructNode_Member_ptr new_StructNode_Member(struct capn_segment*);
StructNode_Field_ptr new_StructNode_Field(struct capn_segment*);
StructNode_Union_ptr new_StructNode_Union(struct capn_segment*);
EnumNode_ptr new_EnumNode(struct capn_segment*);
EnumNode_Enumerant_ptr new_EnumNode_Enumerant(struct capn_segment*);
InterfaceNode_ptr new_InterfaceNode(struct capn_segment*);
InterfaceNode_Method_ptr new_InterfaceNode_Method(struct capn_segment*);
InterfaceNode_Method_Param_ptr new_InterfaceNode_Method_Param(struct capn_segment*);
ConstNode_ptr new_ConstNode(struct capn_segment*);
AnnotationNode_ptr new_AnnotationNode(struct capn_segment*);
CodeGeneratorRequest_ptr new_CodeGeneratorRequest(struct capn_segment*);

Node_list new_Node_list(struct capn_segment*, int len);
Node_NestedNode_list new_Node_NestedNode_list(struct capn_segment*, int len);
Type_list new_Type_list(struct capn_segment*, int len);
Value_list new_Value_list(struct capn_segment*, int len);
Annotation_list new_Annotation_list(struct capn_segment*, int len);
FileNode_list new_FileNode_list(struct capn_segment*, int len);
FileNode_Import_list new_FileNode_Import_list(struct capn_segment*, int len);
StructNode_list new_StructNode_list(struct capn_segment*, int len);
StructNode_Member_list new_StructNode_Member_list(struct capn_segment*, int len);
StructNode_Field_list new_StructNode_Field_list(struct capn_segment*, int len);
StructNode_Union_list new_StructNode_Union_list(struct capn_segment*, int len);
EnumNode_list new_EnumNode_list(struct capn_segment*, int len);
EnumNode_Enumerant_list new_EnumNode_Enumerant_list(struct capn_segment*, int len);
InterfaceNode_list new_InterfaceNode_list(struct capn_segment*, int len);
InterfaceNode_Method_list new_InterfaceNode_Method_list(struct capn_segment*, int len);
InterfaceNode_Method_Param_list new_InterfaceNode_Method_Param_list(struct capn_segment*, int len);
ConstNode_list new_ConstNode_list(struct capn_segment*, int len);
AnnotationNode_list new_AnnotationNode_list(struct capn_segment*, int len);
CodeGeneratorRequest_list new_CodeGeneratorRequest_list(struct capn_segment*, int len);

void read_Node(struct Node*, Node_ptr);
void read_Node_NestedNode(struct Node_NestedNode*, Node_NestedNode_ptr);
void read_Type(struct Type*, Type_ptr);
void read_Value(struct Value*, Value_ptr);
void read_Annotation(struct Annotation*, Annotation_ptr);
void read_FileNode(struct FileNode*, FileNode_ptr);
void read_FileNode_Import(struct FileNode_Import*, FileNode_Import_ptr);
void read_StructNode(struct StructNode*, StructNode_ptr);
void read_StructNode_Member(struct StructNode_Member*, StructNode_Member_ptr);
void read_StructNode_Field(struct StructNode_Field*, StructNode_Field_ptr);
void read_StructNode_Union(struct StructNode_Union*, StructNode_Union_ptr);
void read_EnumNode(struct EnumNode*, EnumNode_ptr);
void read_EnumNode_Enumerant(struct EnumNode_Enumerant*, EnumNode_Enumerant_ptr);
void read_InterfaceNode(struct InterfaceNode*, InterfaceNode_ptr);
void read_InterfaceNode_Method(struct InterfaceNode_Method*, InterfaceNode_Method_ptr);
void read_InterfaceNode_Method_Param(struct InterfaceNode_Method_Param*, InterfaceNode_Method_Param_ptr);
void read_ConstNode(struct ConstNode*, ConstNode_ptr);
void read_AnnotationNode(struct AnnotationNode*, AnnotationNode_ptr);
void read_CodeGeneratorRequest(struct CodeGeneratorRequest*, CodeGeneratorRequest_ptr);

int write_Node(const struct Node*, Node_ptr);
int write_Node_NestedNode(const struct Node_NestedNode*, Node_NestedNode_ptr);
int write_Type(const struct Type*, Type_ptr);
int write_Value(const struct Value*, Value_ptr);
int write_Annotation(const struct Annotation*, Annotation_ptr);
int write_FileNode(const struct FileNode*, FileNode_ptr);
int write_FileNode_Import(const struct FileNode_Import*, FileNode_Import_ptr);
int write_StructNode(const struct StructNode*, StructNode_ptr);
int write_StructNode_Member(const struct StructNode_Member*, StructNode_Member_ptr);
int write_StructNode_Field(const struct StructNode_Field*, StructNode_Field_ptr);
int write_StructNode_Union(const struct StructNode_Union*, StructNode_Union_ptr);
int write_EnumNode(const struct EnumNode*, EnumNode_ptr);
int write_EnumNode_Enumerant(const struct EnumNode_Enumerant*, EnumNode_Enumerant_ptr);
int write_InterfaceNode(const struct InterfaceNode*, InterfaceNode_ptr);
int write_InterfaceNode_Method(const struct InterfaceNode_Method*, InterfaceNode_Method_ptr);
int write_InterfaceNode_Method_Param(const struct InterfaceNode_Method_Param*, InterfaceNode_Method_Param_ptr);
int write_ConstNode(const struct ConstNode*, ConstNode_ptr);
int write_AnnotationNode(const struct AnnotationNode*, AnnotationNode_ptr);
int write_CodeGeneratorRequest(const struct CodeGeneratorRequest*, CodeGeneratorRequest_ptr);

void get_Node(struct Node*, Node_list, int i);
void get_Node_NestedNode(struct Node_NestedNode*, Node_NestedNode_list, int i);
void get_Type(struct Type*, Type_list, int i);
void get_Value(struct Value*, Value_list, int i);
void get_Annotation(struct Annotation*, Annotation_list, int i);
void get_FileNode(struct FileNode*, FileNode_list, int i);
void get_FileNode_Import(struct FileNode_Import*, FileNode_Import_list, int i);
void get_StructNode(struct StructNode*, StructNode_list, int i);
void get_StructNode_Member(struct StructNode_Member*, StructNode_Member_list, int i);
void get_StructNode_Field(struct StructNode_Field*, StructNode_Field_list, int i);
void get_StructNode_Union(struct StructNode_Union*, StructNode_Union_list, int i);
void get_EnumNode(struct EnumNode*, EnumNode_list, int i);
void get_EnumNode_Enumerant(struct EnumNode_Enumerant*, EnumNode_Enumerant_list, int i);
void get_InterfaceNode(struct InterfaceNode*, InterfaceNode_list, int i);
void get_InterfaceNode_Method(struct InterfaceNode_Method*, InterfaceNode_Method_list, int i);
void get_InterfaceNode_Method_Param(struct InterfaceNode_Method_Param*, InterfaceNode_Method_Param_list, int i);
void get_ConstNode(struct ConstNode*, ConstNode_list, int i);
void get_AnnotationNode(struct AnnotationNode*, AnnotationNode_list, int i);
void get_CodeGeneratorRequest(struct CodeGeneratorRequest*, CodeGeneratorRequest_list, int i);

int set_Node(const struct Node*, Node_list, int i);
int set_Node_NestedNode(const struct Node_NestedNode*, Node_NestedNode_list, int i);
int set_Type(const struct Type*, Type_list, int i);
int set_Value(const struct Value*, Value_list, int i);
int set_Annotation(const struct Annotation*, Annotation_list, int i);
int set_FileNode(const struct FileNode*, FileNode_list, int i);
int set_FileNode_Import(const struct FileNode_Import*, FileNode_Import_list, int i);
int set_StructNode(const struct StructNode*, StructNode_list, int i);
int set_StructNode_Member(const struct StructNode_Member*, StructNode_Member_list, int i);
int set_StructNode_Field(const struct StructNode_Field*, StructNode_Field_list, int i);
int set_StructNode_Union(const struct StructNode_Union*, StructNode_Union_list, int i);
int set_EnumNode(const struct EnumNode*, EnumNode_list, int i);
int set_EnumNode_Enumerant(const struct EnumNode_Enumerant*, EnumNode_Enumerant_list, int i);
int set_InterfaceNode(const struct InterfaceNode*, InterfaceNode_list, int i);
int set_InterfaceNode_Method(const struct InterfaceNode_Method*, InterfaceNode_Method_list, int i);
int set_InterfaceNode_Method_Param(const struct InterfaceNode_Method_Param*, InterfaceNode_Method_Param_list, int i);
int set_ConstNode(const struct ConstNode*, ConstNode_list, int i);
int set_AnnotationNode(const struct AnnotationNode*, AnnotationNode_list, int i);
int set_CodeGeneratorRequest(const struct CodeGeneratorRequest*, CodeGeneratorRequest_list, int i);

#ifdef __cplusplus
}
#endif
#endif
