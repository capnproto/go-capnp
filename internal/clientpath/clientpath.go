package clientpath

import "strconv"

// A PipelineOp describes a step in transforming a pipeline.
// It maps closely with the PromisedAnswer.Op struct in rpc.capnp.
type PipelineOp struct {
	Field        uint16
	DefaultValue []byte
}

// Path is an encoded version of a list of pipeline operations.
// It is suitable as a map key.
//
// It specifically ignores default values, because a capability can't have a
// default value other than null.
type Path string

// FromTransform converts a list of pipeline operations to a Path.
func FromTransform(ops []PipelineOp) Path {
	buf := make([]byte, 0, len(ops)*2)
	for i := range ops {
		f := ops[i].Field
		buf = append(buf, byte(f&0x00ff), byte(f&0xff00>>8))
	}
	return Path(buf)
}

// Transform converst a Path into a list of PipelineOps.
func (cp Path) Transform() []PipelineOp {
	ops := make([]PipelineOp, len(cp)/2)
	for i := range ops {
		ops[i].Field = uint16(cp[i*2]) | uint16(cp[i*2+1])<<8
	}
	return ops
}

// String returns a human-readable description of op.
func (op PipelineOp) String() string {
	s := make([]byte, 0, 32)
	s = append(s, "get field "...)
	s = strconv.AppendInt(s, int64(op.Field), 10)
	if op.DefaultValue == nil {
		return string(s)
	}
	s = append(s, " with default"...)
	return string(s)
}
