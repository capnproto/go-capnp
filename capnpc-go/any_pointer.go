package main

import "capnproto.org/go/capnp/v3/internal/schema"

type anyPointerRenderStrategy interface {
	StructParams() any
	ListParams() any
	CapabilityParams() any
	PtrParams() any
}

type anyPointer struct {
	G    *generator
	Type schema.Type
}

func (ap anyPointer) Render(s anyPointerRenderStrategy) error {
	switch ap.Type.AnyPointer().Which() {
	case schema.Type_anyPointer_Which_unconstrained:
		return ap.RenderUnconstrained(s)

	case schema.Type_anyPointer_Which_parameter:
		// TODO(soon):  implement parameter
	case schema.Type_anyPointer_Which_implicitMethodParameter:
		// TODO(soon):  implement implicit method parameter
	}

	return ap.G.r.Render(s.PtrParams())
}

func (ap anyPointer) RenderUnconstrained(s anyPointerRenderStrategy) error {
	switch ap.Type.AnyPointer().Unconstrained().Which() {
	case schema.Type_anyPointer_unconstrained_Which_struct:
		return ap.G.r.Render(s.StructParams())

	case schema.Type_anyPointer_unconstrained_Which_list:
		return ap.G.r.Render(s.ListParams())

	case schema.Type_anyPointer_unconstrained_Which_capability:
		return ap.G.r.Render(s.CapabilityParams())
	}

	return ap.G.r.Render(s.PtrParams())
}

type structAnyPointerRenderStrategy struct {
	Params  structFieldParams
	Default staticDataRef
}

func (s structAnyPointerRenderStrategy) StructParams() any {
	return structAnyStructFieldParams{
		structFieldParams: s.Params,
		Default:           s.Default,
	}
}

func (s structAnyPointerRenderStrategy) ListParams() any {
	return structAnyListFieldParams{
		structFieldParams: s.Params,
		Default:           s.Default,
	}
}

func (s structAnyPointerRenderStrategy) CapabilityParams() any {
	return structCapabilityFieldParams(s.Params)
}

func (s structAnyPointerRenderStrategy) PtrParams() any {
	return structPointerFieldParams{
		structFieldParams: s.Params,
		Default:           s.Default,
	}
}

type promiseAnyPointerRenderStrategy struct {
	G     *generator
	Node  *node
	Field field
}

func (s promiseAnyPointerRenderStrategy) StructParams() any {
	return promiseFieldAnyStructParams(s)
}

func (s promiseAnyPointerRenderStrategy) ListParams() any {
	return promiseFieldAnyListParams(s)
}

func (s promiseAnyPointerRenderStrategy) CapabilityParams() any {
	return promiseCapabilityFieldParams(s)
}

func (s promiseAnyPointerRenderStrategy) PtrParams() any {
	return promiseFieldAnyPointerParams(s)
}

func isAnyCap(ap schema.Type_anyPointer) bool {
	if ap.Which() != schema.Type_anyPointer_Which_unconstrained {
		return false
	}

	which := ap.Unconstrained().Which()
	return which == schema.Type_anyPointer_unconstrained_Which_capability
}
