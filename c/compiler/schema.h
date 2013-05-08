/* vim: set sw=8 ts=8 sts=8 noet: */
#include <capn.h>

struct Node_ptr {struct capn_ptr p;};
struct Node_NestedNode_ptr {struct capn_ptr p;};
struct Type_ptr {struct capn_ptr p;};
struct Value_ptr {struct capn_ptr p;};
struct Annotation_ptr {struct capn_ptr p;};
struct FileNode_ptr {struct capn_ptr p;};
struct FileNode_Import_ptr {struct capn_ptr p;};
struct StructNode_ptr {struct capn_ptr p;};
struct StructNode_Member_ptr {struct capn_ptr p;};
struct StructNode_Field_ptr {struct capn_ptr p;};
struct StructNode_Union_ptr {struct capn_ptr p;};
struct EnumNode_ptr {struct capn_ptr p;};
struct EnumNode_Enumerant_ptr {struct capn_ptr p;};
struct InterfaceNode_ptr {struct capn_ptr p;};
struct InterfaceNode_Method_ptr {struct capn_ptr p;};
struct InterfaceNode_Method_Param_ptr {struct capn_ptr p;};
struct ConstNode_ptr {struct capn_ptr p;};
struct AnnotationNode_ptr {struct capn_ptr p;};
struct CodeGeneratorRequest_ptr {struct capn_ptr p;};

enum Node_body {
	Node_fileNode = 0,
	Node_structNode = 1,
	Node_enumNode = 2,
	Node_interfaceNode = 3,
	Node_constNode = 4,
	Node_annotationNode = 5,
};

struct Node {
	uint64_t id;
	struct capn_text displayName;
	uint64_t scopeId;
	struct capn_ptr nestedNodes; /* List(Node_NestedNode) */
	struct capn_ptr annotations; /* List(Annotation) */
	enum Node_body body_tag;
	union {
		struct FileNode_ptr fileNode;
		struct StructNode_ptr structNode;
		struct EnumNode_ptr enumNode;
		struct InterfaceNode_ptr interfaceNode;
		struct ConstNode_ptr constNode;
		struct AnnotationNode_ptr annotationNode;
	} body;
};

struct Node_NestedNode {
	struct capn_text name;
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
	Type_objectType = 18,
};

struct Type {
	enum Type_body body_tag;
	union {
		struct Type_ptr listType;
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
	Value_objectValue = 18,
};

struct Value {
	enum Value_body body_tag;
	union {
		unsigned int boolValue : 1;
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
		struct capn_text textValue;
		struct capn_data dataValue;
		struct capn_ptr listValue;
		uint16_t enumValue;
		struct capn_ptr structValue;
		struct capn_ptr objectValue;
	} body;
};

struct Annotation {
	uint64_t id;
	struct Value_ptr value;
};

struct FileNode {
	struct capn_ptr imports; /* List(FileNode_Import) */
};

struct FileNode_Import {
	uint64_t id;
	struct capn_text name;
};

enum ElementSize {
	ElementSize_empty = 0,
	ElementSize_bit = 1,
	ElementSize_byte = 2,
	ElementSize_twoBytes = 3,
	ElementSize_fourBytes = 4,
	ElementSize_eightBytes = 5,
	ElementSize_pointer = 6,
	ElementSize_inlineComposite = 7,
};

struct StructNode {
	uint16_t dataSectionWordSize;
	uint16_t pointerSectionSize;
	enum ElementSize preferredListEncoding;
	struct capn_ptr members; /* List(StructNode_Member) */
};

enum StructNode_Member_body {
	StructNode_Member_fieldMember = 0,
	StructNode_Member_unionMember = 1,
};

struct StructNode_Member {
	struct capn_text name;
	uint16_t ordinal;
	uint16_t codeOrder;
	struct capn_ptr annotations; /* List(Annotation) */
	enum StructNode_Member_body body_tag;
	union {
		struct StructNode_Field_ptr fieldMember;
		struct StructNode_Field_ptr unionMember;
	} body;
};

struct StructNode_Field {
	uint32_t offset;
	struct Type_ptr type;
	struct Value_ptr defaultValue;
};

struct StructNode_Union {
	uint32_t discriminantOffset;
	struct capn_ptr members; /* List(StructNode_Member) */
};

struct EnumNode {
	struct capn_ptr enumerants; /* List(EnumNode_Enumerant) */
};

struct EnumNode_Enumerant {
	struct capn_text name;
	uint16_t codeOrder;
	struct capn_ptr annotations; /* List(Annotation) */
};

struct InterfaceNode {
	struct capn_ptr methods; /* List(InterfaceNode_Method) */
};

struct InterfaceNode_Method {
	struct capn_text name;
	uint16_t codeOrder;
	struct capn_ptr params; /* List(InterfaceNode_Method_Param) */
	uint16_t requiredParamCount;
	struct Type_ptr returnType;
	struct capn_ptr annotations; /* List(Annotation) */
};

struct InterfaceNode_Method_Param {
	struct capn_text name;
	struct Type_ptr type;
	struct Value_ptr defaultValue;
	struct capn_ptr annotations; /* List(Annotation) */
};

struct ConstNode {
	struct Type_ptr type;
	struct Value_ptr value;
};

struct AnnotationNode {
	struct Type_ptr type;
	unsigned int targetsFile : 1;
	unsigned int targetsConst : 1;
	unsigned int targetsEnum : 1;
	unsigned int targetsEnumerant : 1;
	unsigned int targetsStruct : 1;
	unsigned int targetsField : 1;
	unsigned int targetsUnion : 1;
	unsigned int targetsInterface : 1;
	unsigned int targetsMethod : 1;
	unsigned int targetsParam : 1;
	unsigned int targetsAnnotation : 1;
};

struct CodeGeneratorRequest {
	struct capn_ptr nodes; /* List(Node) */
	struct capn_ptr requestedFiles; /* List(uint64_t) */
};

void read_Node(const struct Node_ptr*, struct Node*);
void read_Node_NestedNode(const struct Node_NestedNode_ptr*, struct Node_NestedNode*);
void read_Type(const struct Type_ptr*, struct Type*);
void read_Value(const struct Value_ptr*, struct Value*);
void read_Annotation(const struct Annotation_ptr*, struct Annotation*);
void read_FileNode(const struct FileNode_ptr*, struct FileNode*);
void read_FileNode_Import(const struct FileNode_Import_ptr*, struct FileNode_Import*);
void read_StructNode(const struct StructNode_ptr*, struct StructNode*);
void read_StructNode_Member(const struct StructNode_Member_ptr*, struct StructNode_Member*);
void read_StructNode_Field(const struct StructNode_Field_ptr*, struct StructNode_Field*);
void read_StructNode_Union(const struct StructNode_Union_ptr*, struct StructNode_Union*);
void read_EnumNode(const struct EnumNode_ptr*, struct EnumNode*);
void read_EnumNode_Enumerant(const struct EnumNode_Enumerant_ptr*, struct EnumNode_Enumerant*);
void read_InterfaceNode(const struct InterfaceNode_ptr*, struct InterfaceNode*);
void read_InterfaceNode_Method(const struct InterfaceNode_Method_ptr*, struct InterfaceNode_Method*);
void read_InterfaceNode_Method_Param(const struct InterfaceNode_Method_Param_ptr*, struct InterfaceNode_Method_Param*);
void read_ConstNode(const struct ConstNode_ptr*, struct ConstNode*);
void read_AnnotationNode(const struct AnnotationNode_ptr*, struct AnnotationNode*);
void read_CodeGeneratorRequest(const struct CodeGeneratorRequest_ptr*, struct CodeGeneratorRequest*);

