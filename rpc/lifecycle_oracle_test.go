package rpc_test

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"
)

type lifecycleSide uint8

const (
	sideInvalid lifecycleSide = iota
	sideGo
	sidePeer
)

func (s lifecycleSide) String() string {
	switch s {
	case sideGo:
		return "go"
	case sidePeer:
		return "peer"
	default:
		return "invalid"
	}
}

type lifecycleIDSpace uint8

const (
	spaceInvalid lifecycleIDSpace = iota
	spaceQuestion
	spaceExport
	spaceEmbargo
)

func (s lifecycleIDSpace) String() string {
	switch s {
	case spaceQuestion:
		return "question"
	case spaceExport:
		return "export"
	case spaceEmbargo:
		return "embargo"
	default:
		return "invalid"
	}
}

// incarnation names a protocol identity without baking the implementation's
// concrete allocator choice into the oracle. Owner and space keep equal wire
// integers in opposite directions distinct. Generation changes only after a
// live identity has been retired and the concrete ID is observed again.
type incarnation struct {
	Owner      lifecycleSide
	Space      lifecycleIDSpace
	Ordinal    uint32
	Generation uint32
}

func (id incarnation) String() string {
	return fmt.Sprintf("%s/%s/%d@g%d", id.Owner, id.Space, id.Ordinal, id.Generation)
}

type questionRef struct {
	ID incarnation
}

func (q questionRef) String() string { return q.ID.String() }

type capRef struct {
	Export  incarnation
	Promise bool
}

func (c capRef) String() string {
	kind := "cap"
	if c.Promise {
		kind = "promise"
	}
	return kind + "(" + c.Export.String() + ")"
}

type rawIDKey struct {
	Owner lifecycleSide
	Space lifecycleIDSpace
	Wire  uint32
}

type normalizedBinding struct {
	ID     incarnation
	Active bool
}

type ordinalKey struct {
	Owner lifecycleSide
	Space lifecycleIDSpace
}

type idNormalizer struct {
	bindings map[rawIDKey]normalizedBinding
	next     map[ordinalKey]uint32
}

func newIDNormalizer() *idNormalizer {
	return &idNormalizer{
		bindings: make(map[rawIDKey]normalizedBinding),
		next:     make(map[ordinalKey]uint32),
	}
}

func (n *idNormalizer) normalize(owner lifecycleSide, space lifecycleIDSpace, wire uint32) incarnation {
	if owner == sideInvalid || space == spaceInvalid {
		panic("lifecycle oracle: invalid ID namespace")
	}
	key := rawIDKey{Owner: owner, Space: space, Wire: wire}
	if binding, ok := n.bindings[key]; ok {
		if !binding.Active {
			binding.ID.Generation++
			binding.Active = true
			n.bindings[key] = binding
		}
		return binding.ID
	}
	ordinalKey := ordinalKey{Owner: owner, Space: space}
	id := incarnation{
		Owner:   owner,
		Space:   space,
		Ordinal: n.next[ordinalKey],
	}
	n.next[ordinalKey]++
	n.bindings[key] = normalizedBinding{ID: id, Active: true}
	return id
}

func (n *idNormalizer) retire(id incarnation) error {
	for key, binding := range n.bindings {
		if binding.ID == id {
			if !binding.Active {
				return fmt.Errorf("%s is already retired", id)
			}
			binding.Active = false
			n.bindings[key] = binding
			return nil
		}
	}
	return fmt.Errorf("%s is not bound", id)
}

func (n *idNormalizer) question(caller lifecycleSide, wire uint32) questionRef {
	return questionRef{ID: n.normalize(caller, spaceQuestion, wire)}
}

func (n *idNormalizer) capability(exporter lifecycleSide, wire uint32, promise bool) capRef {
	return capRef{
		Export:  n.normalize(exporter, spaceExport, wire),
		Promise: promise,
	}
}

type actionKind uint8

const (
	actionInvalid actionKind = iota
	actionLocalCall
	actionLocalCancel
	actionLocalRelease
	actionLocalResolve
	actionLocalClose
	actionPeerCall
	actionPeerReturn
	actionPeerFinish
	actionPeerResolve
	actionPeerRelease
	actionPeerDisembargo
	actionPeerAbort
	actionWireCall
	actionWireReturn
	actionWireFinish
	actionWireResolve
	actionWireRelease
	actionWireDisembargo
	actionObserveImportPromise
	actionObserveExportPromise
	actionObserveCompletion
	actionObserveDelivery
	actionObserveCapabilityRelease
	actionObserveConnDone
)

func (a actionKind) valid() bool {
	return a > actionInvalid && a <= actionObserveConnDone
}

func (a actionKind) String() string {
	switch a {
	case actionLocalCall:
		return "LocalCall"
	case actionLocalCancel:
		return "LocalCancel"
	case actionLocalRelease:
		return "LocalRelease"
	case actionLocalResolve:
		return "LocalResolve"
	case actionLocalClose:
		return "LocalClose"
	case actionPeerCall:
		return "PeerCall"
	case actionPeerReturn:
		return "PeerReturn"
	case actionPeerFinish:
		return "PeerFinish"
	case actionPeerResolve:
		return "PeerResolve"
	case actionPeerRelease:
		return "PeerRelease"
	case actionPeerDisembargo:
		return "PeerDisembargo"
	case actionPeerAbort:
		return "PeerAbort"
	case actionWireCall:
		return "WireCall"
	case actionWireReturn:
		return "WireReturn"
	case actionWireFinish:
		return "WireFinish"
	case actionWireResolve:
		return "WireResolve"
	case actionWireRelease:
		return "WireRelease"
	case actionWireDisembargo:
		return "WireDisembargo"
	case actionObserveImportPromise:
		return "ObserveImportPromise"
	case actionObserveExportPromise:
		return "ObserveExportPromise"
	case actionObserveCompletion:
		return "ObserveCompletion"
	case actionObserveDelivery:
		return "ObserveDelivery"
	case actionObserveCapabilityRelease:
		return "ObserveCapabilityRelease"
	case actionObserveConnDone:
		return "ObserveConnDone"
	default:
		return "Invalid"
	}
}

type ruleInput struct {
	Action            actionKind
	Question          questionRef
	Capability        capRef
	Replacement       capRef
	Count             uint32
	ReferencesBefore  uint32
	ResultCaps        uint32
	NoFinishNeeded    bool
	ReleaseResultCaps bool
}

func probeFor(action actionKind) ruleInput {
	return ruleInput{
		Action:   action,
		Question: questionRef{ID: incarnation{Owner: sideGo, Space: spaceQuestion}},
		Capability: capRef{
			Export:  incarnation{Owner: sidePeer, Space: spaceExport},
			Promise: true,
		},
		Replacement: capRef{
			Export: incarnation{Owner: sidePeer, Space: spaceExport, Ordinal: 1},
		},
		Count:             1,
		ReferencesBefore:  1,
		ResultCaps:        0,
		NoFinishNeeded:    true,
		ReleaseResultCaps: true,
	}
}

type eventPattern struct {
	Kind       actionKind
	Label      string
	Question   questionRef
	Capability capRef
	HasQuestion,
	HasCapability bool
}

func eventFor(kind actionKind) eventPattern {
	return eventPattern{Kind: kind}
}

func eventForQuestion(kind actionKind, q questionRef) eventPattern {
	return eventPattern{Kind: kind, Question: q, HasQuestion: true}
}

func eventForCapability(kind actionKind, c capRef) eventPattern {
	return eventPattern{Kind: kind, Capability: c, HasCapability: true}
}

func labeledEvent(kind actionKind, label string) eventPattern {
	return eventPattern{Kind: kind, Label: label}
}

func (p eventPattern) valid() bool {
	return p.Kind.valid()
}

func (p eventPattern) matches(event lifecycleEvent) bool {
	if p.Kind != event.Kind || p.Label != "" && p.Label != event.Label {
		return false
	}
	if p.HasQuestion && p.Question != event.Question {
		return false
	}
	if p.HasCapability && p.Capability != event.Capability {
		return false
	}
	return true
}

func (p eventPattern) String() string {
	var qualifiers []string
	if p.Label != "" {
		qualifiers = append(qualifiers, "label="+p.Label)
	}
	if p.HasQuestion {
		qualifiers = append(qualifiers, "question="+p.Question.String())
	}
	if p.HasCapability {
		qualifiers = append(qualifiers, "capability="+p.Capability.String())
	}
	if len(qualifiers) == 0 {
		return p.Kind.String()
	}
	return p.Kind.String() + "[" + strings.Join(qualifiers, ",") + "]"
}

type cardinality struct {
	Event eventPattern
	Min   uint32
	Max   uint32
}

type orderEdge struct {
	Before eventPattern
	After  eventPattern
}

type referenceBalance struct {
	Capability capRef
	Want       int64
}

type constraintAlternative struct {
	Name       string
	Counts     []cardinality
	Orders     []orderEdge
	References []referenceBalance
}

type constraintAlternatives []constraintAlternative

func validateConstraintAlternatives(alternatives constraintAlternatives) error {
	if len(alternatives) == 0 {
		return errors.New("no constraint alternatives")
	}
	names := make(map[string]struct{}, len(alternatives))
	for i, alternative := range alternatives {
		if alternative.Name == "" {
			return fmt.Errorf("alternative %d has no name", i)
		}
		if _, ok := names[alternative.Name]; ok {
			return fmt.Errorf("duplicate alternative %q", alternative.Name)
		}
		names[alternative.Name] = struct{}{}
		if len(alternative.Counts)+len(alternative.Orders)+len(alternative.References) == 0 {
			return fmt.Errorf("alternative %q has no constraints", alternative.Name)
		}
		for _, count := range alternative.Counts {
			if !count.Event.valid() {
				return fmt.Errorf("alternative %q has invalid cardinality event", alternative.Name)
			}
			if count.Min > count.Max {
				return fmt.Errorf("alternative %q has cardinality min %d greater than max %d", alternative.Name, count.Min, count.Max)
			}
		}
		for _, edge := range alternative.Orders {
			if !edge.Before.valid() || !edge.After.valid() {
				return fmt.Errorf("alternative %q has invalid order edge", alternative.Name)
			}
		}
	}
	return nil
}

type lifecycleEvent struct {
	Kind           actionKind
	Label          string
	Question       questionRef
	Capability     capRef
	ReferenceDelta int64
}

func (event lifecycleEvent) String() string {
	pattern := eventPattern{
		Kind:          event.Kind,
		Label:         event.Label,
		Question:      event.Question,
		Capability:    event.Capability,
		HasQuestion:   event.Question.ID.Space != spaceInvalid,
		HasCapability: event.Capability.Export.Space != spaceInvalid,
	}
	s := pattern.String()
	if event.ReferenceDelta != 0 {
		s += fmt.Sprintf("{ref=%+d}", event.ReferenceDelta)
	}
	return s
}

type lifecycleLedger []lifecycleEvent

func (ledger lifecycleLedger) count(pattern eventPattern) uint32 {
	var count uint32
	for _, event := range ledger {
		if pattern.matches(event) {
			count++
		}
	}
	return count
}

func (ledger lifecycleLedger) indexes(pattern eventPattern) []int {
	var indexes []int
	for i, event := range ledger {
		if pattern.matches(event) {
			indexes = append(indexes, i)
		}
	}
	return indexes
}

func (ledger lifecycleLedger) referenceBalance(capability capRef) int64 {
	var balance int64
	for _, event := range ledger {
		if event.Capability == capability {
			balance += event.ReferenceDelta
		}
	}
	return balance
}

func (ledger lifecycleLedger) check(alternative constraintAlternative) error {
	var failures []string
	for _, count := range alternative.Counts {
		got := ledger.count(count.Event)
		if got < count.Min || got > count.Max {
			failures = append(failures,
				fmt.Sprintf("%s count=%d want=%d..%d", count.Event, got, count.Min, count.Max))
		}
	}
	for _, edge := range alternative.Orders {
		before := ledger.indexes(edge.Before)
		after := ledger.indexes(edge.After)
		if len(before) == 0 || len(after) == 0 {
			failures = append(failures,
				fmt.Sprintf("%s before %s is unobservable", edge.Before, edge.After))
			continue
		}
		if before[len(before)-1] >= after[0] {
			failures = append(failures,
				fmt.Sprintf("%s does not happen before %s", edge.Before, edge.After))
		}
	}
	for _, reference := range alternative.References {
		if got := ledger.referenceBalance(reference.Capability); got != reference.Want {
			failures = append(failures,
				fmt.Sprintf("%s balance=%+d want=%+d", reference.Capability, got, reference.Want))
		}
	}
	if len(failures) > 0 {
		return errors.New(strings.Join(failures, "; "))
	}
	return nil
}

func checkConstraintAlternatives(alternatives constraintAlternatives, ledger lifecycleLedger) error {
	if err := validateConstraintAlternatives(alternatives); err != nil {
		return err
	}
	var failures []string
	for _, alternative := range alternatives {
		if err := ledger.check(alternative); err == nil {
			return nil
		} else {
			failures = append(failures, alternative.Name+": "+err.Error())
		}
	}
	return errors.New(strings.Join(failures, "\n"))
}

func buildQuestionRetirementConstraints(input ruleInput) (constraintAlternatives, error) {
	switch input.Action {
	case actionLocalCall, actionWireCall:
		return constraintAlternatives{{
			Name: "call-bound",
			Counts: []cardinality{{
				Event: eventForQuestion(actionWireCall, input.Question),
				Min:   1,
				Max:   1,
			}},
		}}, nil
	case actionPeerReturn:
		return constraintAlternatives{{
			Name: "returned",
			Counts: []cardinality{
				{Event: eventForQuestion(actionPeerReturn, input.Question), Min: 1, Max: 1},
				{Event: eventForQuestion(actionObserveCompletion, input.Question), Min: 1, Max: 1},
			},
		}}, nil
	case actionWireFinish:
		return constraintAlternatives{{
			Name: "finished",
			Counts: []cardinality{{
				Event: eventForQuestion(actionWireFinish, input.Question),
				Min:   1,
				Max:   1,
			}},
		}}, nil
	default:
		return nil, fmt.Errorf("question retirement does not apply to %s", input.Action)
	}
}

func buildCancelResultConstraints(input ruleInput) (constraintAlternatives, error) {
	switch input.Action {
	case actionLocalCancel:
		return constraintAlternatives{
			{
				Name: "return-wins",
				Counts: []cardinality{
					{Event: labeledEvent(actionObserveCompletion, "success"), Min: 1, Max: 1},
					{Event: labeledEvent(actionWireFinish, "releaseResultCaps=false"), Min: 1, Max: 1},
				},
			},
			{
				Name: "cancel-wins",
				Counts: []cardinality{
					{Event: labeledEvent(actionObserveCompletion, "canceled"), Min: 1, Max: 1},
					{Event: labeledEvent(actionWireFinish, "releaseResultCaps=true"), Min: 1, Max: 1},
				},
				References: []referenceBalance{{
					Capability: input.Capability,
					Want:       0,
				}},
			},
		}, nil
	case actionWireFinish, actionPeerFinish:
		return constraintAlternatives{{
			Name: "implicit-result-release",
			Counts: []cardinality{{
				Event: labeledEvent(input.Action, "releaseResultCaps=true"),
				Min:   1,
				Max:   1,
			}},
			References: []referenceBalance{{
				Capability: input.Capability,
				Want:       0,
			}},
		}}, nil
	case actionPeerReturn, actionWireReturn:
		return constraintAlternatives{{
			Name: "late-result-balanced",
			Counts: []cardinality{{
				Event: eventFor(input.Action),
				Min:   1,
				Max:   1,
			}},
			References: []referenceBalance{{
				Capability: input.Capability,
				Want:       0,
			}},
		}}, nil
	default:
		return nil, fmt.Errorf("result release does not apply to %s", input.Action)
	}
}

func buildNoFinishNeededConstraints(input ruleInput) (constraintAlternatives, error) {
	switch input.Action {
	case actionPeerReturn, actionWireReturn:
		if !input.NoFinishNeeded {
			return nil, errors.New("noFinishNeeded Return does not set noFinishNeeded")
		}
		if input.ResultCaps != 0 {
			return nil, fmt.Errorf("noFinishNeeded Return carries %d result capabilities", input.ResultCaps)
		}
		return constraintAlternatives{
			{
				Name: "retire-without-finish",
				Counts: []cardinality{
					{Event: labeledEvent(input.Action, "noFinishNeeded=true,resultCaps=0"), Min: 1, Max: 1},
					{Event: eventFor(actionWireFinish), Min: 0, Max: 0},
				},
			},
			{
				Name: "compatibility-finish",
				Counts: []cardinality{
					{Event: labeledEvent(input.Action, "noFinishNeeded=true,resultCaps=0"), Min: 1, Max: 1},
					{Event: eventFor(actionWireFinish), Min: 1, Max: 1},
				},
			},
		}, nil
	case actionWireFinish, actionPeerFinish:
		return constraintAlternatives{{
			Name: "compatibility-finish",
			Counts: []cardinality{{
				Event: eventFor(input.Action),
				Min:   1,
				Max:   1,
			}},
		}}, nil
	default:
		return nil, fmt.Errorf("noFinishNeeded does not apply to %s", input.Action)
	}
}

func buildPromiseSingleResolutionConstraints(input ruleInput) (constraintAlternatives, error) {
	switch input.Action {
	case actionObserveImportPromise, actionObserveExportPromise:
		resolve := actionPeerResolve
		if input.Action == actionObserveExportPromise {
			resolve = actionWireResolve
		}
		return constraintAlternatives{
			{
				Name: "resolved-once",
				Counts: []cardinality{{
					Event: eventForCapability(resolve, input.Capability),
					Min:   1,
					Max:   1,
				}},
			},
			{
				Name: "released-before-resolve",
				Counts: []cardinality{{
					Event: eventForCapability(resolve, input.Capability),
					Min:   0,
					Max:   1,
				}},
			},
		}, nil
	case actionLocalResolve, actionPeerResolve, actionWireResolve:
		return constraintAlternatives{{
			Name: "single-resolution",
			Counts: []cardinality{{
				Event: eventForCapability(input.Action, input.Capability),
				Min:   1,
				Max:   1,
			}},
		}}, nil
	case actionLocalRelease, actionPeerRelease, actionWireRelease:
		return constraintAlternatives{{
			Name: "early-release",
			Counts: []cardinality{{
				Event: eventForCapability(input.Action, input.Capability),
				Min:   1,
				Max:   1,
			}},
		}}, nil
	default:
		return nil, fmt.Errorf("single resolution does not apply to %s", input.Action)
	}
}

func buildLateResolveConstraints(input ruleInput) (constraintAlternatives, error) {
	switch input.Action {
	case actionLocalRelease, actionWireRelease:
		return constraintAlternatives{{
			Name: "promise-released",
			Counts: []cardinality{{
				Event: eventForCapability(actionWireRelease, input.Capability),
				Min:   1,
				Max:   1,
			}},
		}}, nil
	case actionPeerResolve:
		return constraintAlternatives{{
			Name: "late-capability-released",
			Counts: []cardinality{
				{Event: eventForCapability(actionPeerResolve, input.Capability), Min: 1, Max: 1},
				{Event: eventForCapability(actionWireRelease, input.Replacement), Min: 1, Max: 1},
				{Event: eventFor(actionWireDisembargo), Min: 0, Max: 0},
			},
			References: []referenceBalance{{
				Capability: input.Replacement,
				Want:       0,
			}},
		}}, nil
	default:
		return nil, fmt.Errorf("late Resolve does not apply to %s", input.Action)
	}
}

func buildReferenceDecrementConstraints(input ruleInput) (constraintAlternatives, error) {
	switch input.Action {
	case actionLocalRelease, actionPeerRelease, actionWireRelease:
		if input.Count > input.ReferencesBefore {
			return nil, fmt.Errorf(
				"Release count %d exceeds reference count %d",
				input.Count,
				input.ReferencesBefore,
			)
		}
		return constraintAlternatives{{
			Name: "reference-decrement",
			Counts: []cardinality{{
				Event: eventForCapability(input.Action, input.Capability),
				Min:   1,
				Max:   1,
			}},
			References: []referenceBalance{{
				Capability: input.Capability,
				Want:       int64(input.ReferencesBefore - input.Count),
			}},
		}}, nil
	default:
		return nil, fmt.Errorf("reference decrement does not apply to %s", input.Action)
	}
}

func buildDisembargoConstraints(input ruleInput) (constraintAlternatives, error) {
	switch input.Action {
	case actionLocalCall, actionWireCall, actionPeerReturn, actionPeerResolve,
		actionWireDisembargo, actionPeerDisembargo, actionObserveDelivery:
		return constraintAlternatives{
			{
				Name: "same-path",
				Counts: []cardinality{
					{Event: eventFor(actionWireDisembargo), Min: 0, Max: 0},
					{Event: labeledEvent(actionObserveDelivery, "pre-resolution"), Min: 1, Max: 1},
					{Event: labeledEvent(actionObserveDelivery, "post-resolution"), Min: 1, Max: 1},
				},
				Orders: []orderEdge{{
					Before: labeledEvent(actionObserveDelivery, "pre-resolution"),
					After:  labeledEvent(actionObserveDelivery, "post-resolution"),
				}},
			},
			{
				Name: "path-shortened",
				Counts: []cardinality{
					{Event: labeledEvent(actionObserveDelivery, "pre-resolution"), Min: 1, Max: 1},
					{Event: labeledEvent(actionObserveDelivery, "post-resolution"), Min: 1, Max: 1},
				},
				Orders: []orderEdge{{
					Before: labeledEvent(actionObserveDelivery, "pre-resolution"),
					After:  labeledEvent(actionObserveDelivery, "post-resolution"),
				}},
			},
		}, nil
	default:
		return nil, fmt.Errorf("Disembargo ordering does not apply to %s", input.Action)
	}
}

func buildDisconnectConstraints(input ruleInput) (constraintAlternatives, error) {
	switch input.Action {
	case actionLocalClose, actionPeerAbort, actionObserveCompletion,
		actionObserveImportPromise, actionObserveCapabilityRelease:
		return constraintAlternatives{{
			Name: "four-tables-lost",
			Counts: []cardinality{
				{Event: eventFor(actionObserveCompletion), Min: 1, Max: ^uint32(0)},
				{Event: eventFor(actionObserveCapabilityRelease), Min: 1, Max: ^uint32(0)},
			},
		}}, nil
	default:
		return nil, fmt.Errorf("disconnect does not apply to %s", input.Action)
	}
}

func goConnectionDrainConstraints() constraintAlternatives {
	return constraintAlternatives{{
		Name: "go-terminal-drain",
		Counts: []cardinality{
			{Event: eventFor(actionObserveCompletion), Min: 1, Max: ^uint32(0)},
			{Event: eventFor(actionObserveCapabilityRelease), Min: 1, Max: ^uint32(0)},
			{Event: eventFor(actionObserveConnDone), Min: 1, Max: 1},
		},
		// Conn shutdown initiates export release before Done, but releasing a
		// capnp.Client may schedule its Shutdown callback asynchronously.  The
		// public callback therefore has no ordering edge with Conn.Done.
		Orders: []orderEdge{{
			Before: eventFor(actionObserveCompletion),
			After:  eventFor(actionObserveConnDone),
		}},
	}}
}

type divergence struct {
	ActionIndex int
	Origin      string
	Reason      string
	Expected    string
	Actual      string
	Replay      string
}

func (d divergence) Error() string {
	return fmt.Sprintf(
		"rpc lifecycle divergence at action %d\norigin: %s\nreason: %s\nexpected: %s\nactual: %s\nreplay:\n%s",
		d.ActionIndex, d.Origin, d.Reason, d.Expected, d.Actual, d.Replay)
}

func stableExpected(alternatives constraintAlternatives) string {
	var names []string
	for _, alternative := range alternatives {
		names = append(names, alternative.Name)
	}
	slices.Sort(names)
	return strings.Join(names, " | ")
}

func stableLedger(ledger lifecycleLedger) string {
	parts := make([]string, len(ledger))
	for i, event := range ledger {
		parts[i] = event.String()
	}
	return strings.Join(parts, " -> ")
}

func stableReplay(events lifecycleLedger) string {
	var b strings.Builder
	b.WriteString("rpc-lifecycle-trace-v1\n")
	for i, event := range events {
		fmt.Fprintf(&b, "%02d %s\n", i, event)
	}
	return strings.TrimSuffix(b.String(), "\n")
}

func checkWithDiagnostic(
	actionIndex int,
	origin string,
	alternatives constraintAlternatives,
	ledger lifecycleLedger,
) error {
	if err := checkConstraintAlternatives(alternatives, ledger); err != nil {
		return divergence{
			ActionIndex: actionIndex,
			Origin:      origin,
			Reason:      err.Error(),
			Expected:    stableExpected(alternatives),
			Actual:      stableLedger(ledger),
			Replay:      stableReplay(ledger),
		}
	}
	return nil
}

type questionState struct {
	Returned       bool
	Finished       bool
	NoFinishNeeded bool
	ResultCaps     uint32
}

type promiseState uint8

const (
	promiseUnresolved promiseState = iota + 1
	promiseResolved
	promiseReleased
)

type promiseRecord struct {
	State      promiseState
	References uint32
}

// lifecycleOracle deliberately models only protocol-visible identities and
// transitions needed by the six named traces. It does not mirror question.go,
// import.go, or export.go implementation states.
type lifecycleOracle struct {
	ids              *idNormalizer
	pendingLocalCall bool
	questions        map[questionRef]questionState
	imports          map[capRef]promiseRecord
	exports          map[capRef]promiseRecord
}

func newLifecycleOracle() *lifecycleOracle {
	return &lifecycleOracle{
		ids:       newIDNormalizer(),
		questions: make(map[questionRef]questionState),
		imports:   make(map[capRef]promiseRecord),
		exports:   make(map[capRef]promiseRecord),
	}
}

func (o *lifecycleOracle) localCall() error {
	if o.pendingLocalCall {
		return errors.New("a local call is already waiting for WireCall")
	}
	o.pendingLocalCall = true
	return nil
}

func (o *lifecycleOracle) wireCall(rawQuestionID uint32) (questionRef, error) {
	if !o.pendingLocalCall {
		return questionRef{}, errors.New("WireCall has no pending LocalCall")
	}
	o.pendingLocalCall = false
	question := o.ids.question(sideGo, rawQuestionID)
	if _, exists := o.questions[question]; exists {
		return questionRef{}, fmt.Errorf("question %s reused before retirement", question)
	}
	o.questions[question] = questionState{}
	return question, nil
}

func (o *lifecycleOracle) peerReturn(
	question questionRef,
	noFinishNeeded bool,
	resultCaps uint32,
) error {
	state, ok := o.questions[question]
	if !ok {
		return fmt.Errorf("Return references unknown question %s", question)
	}
	if state.Returned {
		return fmt.Errorf("question %s returned more than once", question)
	}
	if noFinishNeeded && resultCaps != 0 {
		return fmt.Errorf(
			"question %s has noFinishNeeded with %d result capabilities",
			question,
			resultCaps,
		)
	}
	state.Returned = true
	state.NoFinishNeeded = noFinishNeeded
	state.ResultCaps = resultCaps
	o.questions[question] = state
	return o.maybeRetireQuestion(question)
}

func (o *lifecycleOracle) wireFinish(question questionRef) error {
	state, ok := o.questions[question]
	if !ok {
		return fmt.Errorf("Finish references unknown question %s", question)
	}
	if state.Finished {
		return fmt.Errorf("question %s finished more than once", question)
	}
	state.Finished = true
	o.questions[question] = state
	return o.maybeRetireQuestion(question)
}

func (o *lifecycleOracle) maybeRetireQuestion(question questionRef) error {
	state := o.questions[question]
	if !state.Returned || !state.Finished && !state.NoFinishNeeded {
		return nil
	}
	delete(o.questions, question)
	return o.ids.retire(question.ID)
}

func (o *lifecycleOracle) observeImportedPromise(rawExportID uint32) capRef {
	capability := o.ids.capability(sidePeer, rawExportID, true)
	record, ok := o.imports[capability]
	if !ok {
		record.State = promiseUnresolved
	}
	record.References++
	o.imports[capability] = record
	return capability
}

func (o *lifecycleOracle) observeExportedPromise(rawExportID uint32) capRef {
	capability := o.ids.capability(sideGo, rawExportID, true)
	record, ok := o.exports[capability]
	if !ok {
		record.State = promiseUnresolved
	}
	record.References++
	o.exports[capability] = record
	return capability
}

func (o *lifecycleOracle) localRelease(capability capRef) error {
	return o.localReleaseCount(capability, 1)
}

func (o *lifecycleOracle) localReleaseCount(capability capRef, count uint32) error {
	record, ok := o.imports[capability]
	if !ok {
		return fmt.Errorf("LocalRelease references non-imported %s", capability)
	}
	if record.State == promiseReleased {
		return fmt.Errorf("imported %s released more than once", capability)
	}
	if count > record.References {
		return fmt.Errorf(
			"LocalRelease count %d exceeds %d references to imported %s",
			count,
			record.References,
			capability,
		)
	}
	record.References -= count
	if record.References == 0 {
		record.State = promiseReleased
		if err := o.ids.retire(capability.Export); err != nil {
			return err
		}
	}
	o.imports[capability] = record
	return nil
}

func (o *lifecycleOracle) peerResolve(capability capRef) error {
	record, ok := o.imports[capability]
	if !ok {
		return fmt.Errorf("PeerResolve references non-imported %s", capability)
	}
	if record.State == promiseResolved {
		return fmt.Errorf("imported %s resolved more than once", capability)
	}
	if record.State == promiseUnresolved {
		record.State = promiseResolved
		o.imports[capability] = record
	}
	// A late Resolve for a released import is tolerated. The replacement
	// obligations live in buildLateResolveConstraints.
	return nil
}

func (o *lifecycleOracle) localResolve(capability capRef) error {
	record, ok := o.exports[capability]
	if !ok {
		return fmt.Errorf("LocalResolve references non-exported %s", capability)
	}
	if record.State != promiseUnresolved {
		return fmt.Errorf("exported %s cannot resolve from state %d", capability, record.State)
	}
	record.State = promiseResolved
	o.exports[capability] = record
	return nil
}

func (o *lifecycleOracle) peerRelease(capability capRef) error {
	return o.peerReleaseCount(capability, 1)
}

func (o *lifecycleOracle) peerReleaseCount(capability capRef, count uint32) error {
	record, ok := o.exports[capability]
	if !ok {
		return fmt.Errorf("PeerRelease references non-exported %s", capability)
	}
	if record.State == promiseReleased {
		return fmt.Errorf("exported %s released more than once", capability)
	}
	if count > record.References {
		return fmt.Errorf(
			"PeerRelease count %d exceeds %d references to exported %s",
			count,
			record.References,
			capability,
		)
	}
	record.References -= count
	if record.References == 0 {
		record.State = promiseReleased
		if err := o.ids.retire(capability.Export); err != nil {
			return err
		}
	}
	o.exports[capability] = record
	return nil
}

func TestLifecycleIDNormalizerSeparatesDirectionAndGeneration(t *testing.T) {
	ids := newIDNormalizer()
	goQuestion := ids.question(sideGo, 7)
	peerQuestion := ids.question(sidePeer, 7)
	if goQuestion == peerQuestion {
		t.Fatal("equal raw question IDs from opposite callers normalized to one identity")
	}
	if got, want := goQuestion.String(), "go/question/0@g0"; got != want {
		t.Fatalf("Go question = %q; want %q", got, want)
	}
	if err := ids.retire(goQuestion.ID); err != nil {
		t.Fatal(err)
	}
	reused := ids.question(sideGo, 7)
	if got, want := reused.String(), "go/question/0@g1"; got != want {
		t.Fatalf("reused Go question = %q; want %q", got, want)
	}
	if got := ids.question(sidePeer, 7); got != peerQuestion {
		t.Fatalf("Go retirement changed peer identity from %s to %s", peerQuestion, got)
	}
}

func TestLifecycleCapabilityIdentityBelongsToExporter(t *testing.T) {
	ids := newIDNormalizer()
	senderPromise := ids.capability(sidePeer, 9, true)
	receiverHosted := ids.capability(sideGo, 9, false)
	if senderPromise.Export == receiverHosted.Export {
		t.Fatal("equal raw export IDs owned by opposite exporters normalized to one identity")
	}
	if senderPromise.Export.Owner != sidePeer || receiverHosted.Export.Owner != sideGo {
		t.Fatalf("export ownership = %s and %s; want peer and go", senderPromise, receiverHosted)
	}
}

func TestLifecycleLocalCallBindsOnlyAtWireObservation(t *testing.T) {
	oracle := newLifecycleOracle()
	if err := oracle.localCall(); err != nil {
		t.Fatal(err)
	}
	if len(oracle.questions) != 0 {
		t.Fatal("LocalCall chose a concrete question before WireCall observation")
	}
	question, err := oracle.wireCall(42)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := question.String(), "go/question/0@g0"; got != want {
		t.Fatalf("bound question = %q; want %q", got, want)
	}
	if err := oracle.peerReturn(question, false, 0); err != nil {
		t.Fatal(err)
	}
	if err := oracle.localCall(); err != nil {
		t.Fatal(err)
	}
	if _, err := oracle.wireCall(42); err == nil || !strings.Contains(err.Error(), "reused before retirement") {
		t.Fatalf("early ID reuse error = %v; want retirement diagnostic", err)
	}
	if err := oracle.wireFinish(question); err != nil {
		t.Fatal(err)
	}
	if err := oracle.localCall(); err != nil {
		t.Fatal(err)
	}
	if _, err := oracle.wireCall(42); err != nil {
		t.Fatal("reuse after Return and Finish:", err)
	}
}

func TestLifecyclePromiseDirectionsRejectCrossedActions(t *testing.T) {
	oracle := newLifecycleOracle()
	imported := oracle.observeImportedPromise(4)
	exported := oracle.observeExportedPromise(4)
	if imported.Export == exported.Export {
		t.Fatal("imported and exported promises share an exporter-owned identity")
	}
	if err := oracle.localResolve(imported); err == nil {
		t.Fatal("LocalResolve accepted an imported promise")
	}
	if err := oracle.peerResolve(exported); err == nil {
		t.Fatal("PeerResolve accepted an exported promise")
	}
	if err := oracle.localRelease(exported); err == nil {
		t.Fatal("LocalRelease accepted a peer-owned release")
	}
	if err := oracle.peerRelease(imported); err == nil {
		t.Fatal("PeerRelease accepted a Go-owned release")
	}
}

func TestLifecycleConstraintAlternativesAreCoherent(t *testing.T) {
	input := probeFor(actionLocalCancel)
	alternatives, err := buildCancelResultConstraints(input)
	if err != nil {
		t.Fatal(err)
	}
	frankenstein := lifecycleLedger{
		{Kind: actionObserveCompletion, Label: "success"},
		{Kind: actionWireFinish, Label: "releaseResultCaps=true"},
	}
	err = checkConstraintAlternatives(alternatives, frankenstein)
	if err == nil {
		t.Fatal("mixed observations from different alternatives were accepted")
	}
	if !strings.Contains(err.Error(), "return-wins") || !strings.Contains(err.Error(), "cancel-wins") {
		t.Fatalf("alternative failure = %v; want both complete alternatives", err)
	}
}

func TestLifecyclePartialOrderIgnoresUnrelatedInterleaving(t *testing.T) {
	input := probeFor(actionObserveDelivery)
	alternatives, err := buildDisembargoConstraints(input)
	if err != nil {
		t.Fatal(err)
	}
	ordered := lifecycleLedger{
		{Kind: actionObserveDelivery, Label: "pre-resolution"},
		{Kind: actionObserveCompletion, Label: "unrelated"},
		{Kind: actionObserveDelivery, Label: "post-resolution"},
	}
	if err := checkConstraintAlternatives(alternatives[1:], ordered); err != nil {
		t.Fatal("unrelated event overconstrained the partial order:", err)
	}
	overtaken := lifecycleLedger{
		{Kind: actionObserveDelivery, Label: "post-resolution"},
		{Kind: actionObserveCompletion, Label: "unrelated"},
		{Kind: actionObserveDelivery, Label: "pre-resolution"},
	}
	if err := checkConstraintAlternatives(alternatives[1:], overtaken); err == nil {
		t.Fatal("post-resolution delivery overtook pre-resolution delivery")
	}
}

func TestLifecycleLateResolveBalancesReplacement(t *testing.T) {
	input := probeFor(actionPeerResolve)
	alternatives, err := buildLateResolveConstraints(input)
	if err != nil {
		t.Fatal(err)
	}
	ledger := lifecycleLedger{
		{Kind: actionPeerResolve, Capability: input.Capability},
		{Kind: actionPeerReturn, Label: "unrelated"},
		{Kind: actionObserveImportPromise, Capability: input.Replacement, ReferenceDelta: 1},
		{Kind: actionWireRelease, Capability: input.Replacement, ReferenceDelta: -1},
	}
	if err := checkConstraintAlternatives(alternatives, ledger); err != nil {
		t.Fatal(err)
	}
	ledger = append(ledger, lifecycleEvent{Kind: actionWireDisembargo})
	if err := checkConstraintAlternatives(alternatives, ledger); err == nil {
		t.Fatal("late Resolve accepted an outbound Disembargo")
	}
}

func TestLifecycleGoDrainPolicyIsNotAttributedToSchema(t *testing.T) {
	ledger := lifecycleLedger{
		{Kind: actionObserveCompletion},
		{Kind: actionObserveCapabilityRelease},
		{Kind: actionObserveConnDone},
	}
	if err := checkWithDiagnostic(
		2,
		"go-rpc:connection-terminal-drain",
		goConnectionDrainConstraints(),
		ledger,
	); err != nil {
		t.Fatal(err)
	}
}

func TestLifecycleGoDrainAllowsAsynchronousCapabilityShutdown(t *testing.T) {
	ledger := lifecycleLedger{
		{Kind: actionObserveCompletion},
		{Kind: actionObserveConnDone},
		{Kind: actionObserveCapabilityRelease},
	}
	if err := checkConstraintAlternatives(goConnectionDrainConstraints(), ledger); err != nil {
		t.Fatal(err)
	}
}

func TestLifecycleDiagnosticIsStableAndAlphaRenamed(t *testing.T) {
	ids := newIDNormalizer()
	question := ids.question(sideGo, 900)
	alternatives := constraintAlternatives{{
		Name: "normal-retirement",
		Counts: []cardinality{{
			Event: eventForQuestion(actionWireFinish, question),
			Min:   1,
			Max:   1,
		}},
	}}
	ledger := lifecycleLedger{{
		Kind:     actionWireCall,
		Question: question,
	}}
	err := checkWithDiagnostic(
		1,
		"rpc-v2@0b821519:Call.questionId:question-id-retirement",
		alternatives,
		ledger,
	)
	const want = `rpc lifecycle divergence at action 1
origin: rpc-v2@0b821519:Call.questionId:question-id-retirement
reason: normal-retirement: WireFinish[question=go/question/0@g0] count=0 want=1..1
expected: normal-retirement
actual: WireCall[question=go/question/0@g0]
replay:
rpc-lifecycle-trace-v1
00 WireCall[question=go/question/0@g0]`
	if err == nil || err.Error() != want {
		t.Fatalf("diagnostic:\n%s\nwant:\n%s", err, want)
	}
}

// These contract tests are deliberately wire-free.  The synctest traces feed
// their observations into the same builders, so an accidental weakening of
// the oracle fails here without needing a particular transport schedule.

func TestLifecycleContractQuestionReuseAfterReturnAndFinish(t *testing.T) {
	oracle := newLifecycleOracle()
	if err := oracle.localCall(); err != nil {
		t.Fatal(err)
	}
	first, err := oracle.wireCall(7)
	if err != nil {
		t.Fatal(err)
	}
	if err := oracle.peerReturn(first, false, 0); err != nil {
		t.Fatal(err)
	}
	if err := oracle.wireFinish(first); err != nil {
		t.Fatal(err)
	}
	if err := oracle.localCall(); err != nil {
		t.Fatal(err)
	}
	second, err := oracle.wireCall(7)
	if err != nil {
		t.Fatal(err)
	}
	if first.ID.Generation+1 != second.ID.Generation {
		t.Fatalf("reused question generation = %d; want %d", second.ID.Generation, first.ID.Generation+1)
	}
}

func TestLifecycleContractCancelBalancesLateCapabilityReturn(t *testing.T) {
	input := probeFor(actionLocalCancel)
	alternatives, err := buildCancelResultConstraints(input)
	if err != nil {
		t.Fatal(err)
	}
	ledger := lifecycleLedger{
		{Kind: actionObserveCompletion, Label: "canceled"},
		{Kind: actionPeerReturn, Capability: input.Capability, ReferenceDelta: 1},
		{
			Kind:           actionWireFinish,
			Label:          "releaseResultCaps=true",
			Capability:     input.Capability,
			ReferenceDelta: -1,
		},
	}
	if err := checkConstraintAlternatives(alternatives, ledger); err != nil {
		t.Fatal(err)
	}
}

func TestLifecycleContractNoFinishNeededAllowsQuestionRetirement(t *testing.T) {
	input := probeFor(actionPeerReturn)
	alternatives, err := buildNoFinishNeededConstraints(input)
	if err != nil {
		t.Fatal(err)
	}
	ledger := lifecycleLedger{{
		Kind:  actionPeerReturn,
		Label: "noFinishNeeded=true,resultCaps=0",
	}}
	if err := checkConstraintAlternatives(alternatives, ledger); err != nil {
		t.Fatal(err)
	}

	oracle := newLifecycleOracle()
	if err := oracle.localCall(); err != nil {
		t.Fatal(err)
	}
	first, err := oracle.wireCall(11)
	if err != nil {
		t.Fatal(err)
	}
	if err := oracle.peerReturn(first, true, 0); err != nil {
		t.Fatal(err)
	}
	if err := oracle.localCall(); err != nil {
		t.Fatal(err)
	}
	second, err := oracle.wireCall(11)
	if err != nil {
		t.Fatal(err)
	}
	if first.ID.Generation+1 != second.ID.Generation {
		t.Fatalf("reused question generation = %d; want %d", second.ID.Generation, first.ID.Generation+1)
	}

	invalid := probeFor(actionPeerReturn)
	invalid.ResultCaps = 1
	if _, err := buildNoFinishNeededConstraints(invalid); err == nil {
		t.Fatal("capability-bearing noFinishNeeded Return was accepted")
	}
}

func TestLifecycleContractReleasedPromiseBalancesLateResolve(t *testing.T) {
	input := probeFor(actionPeerResolve)
	alternatives, err := buildLateResolveConstraints(input)
	if err != nil {
		t.Fatal(err)
	}
	ledger := lifecycleLedger{
		{Kind: actionPeerResolve, Capability: input.Capability},
		{Kind: actionObserveImportPromise, Capability: input.Replacement, ReferenceDelta: 1},
		{Kind: actionWireRelease, Capability: input.Replacement, ReferenceDelta: -1},
	}
	if err := checkConstraintAlternatives(alternatives, ledger); err != nil {
		t.Fatal(err)
	}
}

func TestLifecyclePromiseReferenceCountsAndGeneration(t *testing.T) {
	oracle := newLifecycleOracle()
	first := oracle.observeImportedPromise(23)
	repeated := oracle.observeImportedPromise(23)
	if repeated != first {
		t.Fatalf("repeated senderPromise = %s; want %s", repeated, first)
	}
	if err := oracle.peerResolve(first); err != nil {
		t.Fatal(err)
	}
	if err := oracle.localReleaseCount(first, 1); err != nil {
		t.Fatal(err)
	}
	if got := oracle.imports[first].References; got != 1 {
		t.Fatalf("references after partial Release = %d; want 1", got)
	}
	if err := oracle.peerResolve(first); err == nil {
		t.Fatal("second Resolve for one senderPromise incarnation was accepted")
	}
	if err := oracle.localReleaseCount(first, 1); err != nil {
		t.Fatal(err)
	}
	reused := oracle.observeImportedPromise(23)
	if reused.Export.Generation != first.Export.Generation+1 {
		t.Fatalf(
			"reused senderPromise generation = %d; want %d",
			reused.Export.Generation,
			first.Export.Generation+1,
		)
	}
	if err := oracle.peerResolve(reused); err != nil {
		t.Fatal("Resolve for recycled senderPromise:", err)
	}
}

func TestLifecycleReferenceDecrementAllowsPartialRelease(t *testing.T) {
	input := probeFor(actionWireRelease)
	input.ReferencesBefore = 3
	input.Count = 1
	alternatives, err := buildReferenceDecrementConstraints(input)
	if err != nil {
		t.Fatal(err)
	}
	ledger := lifecycleLedger{
		{Kind: actionObserveImportPromise, Capability: input.Capability, ReferenceDelta: 3},
		{Kind: actionWireRelease, Capability: input.Capability, ReferenceDelta: -1},
	}
	if err := checkConstraintAlternatives(alternatives, ledger); err != nil {
		t.Fatal(err)
	}
	input.Count = 4
	if _, err := buildReferenceDecrementConstraints(input); err == nil {
		t.Fatal("Release reference-count underflow was accepted")
	}
}

func TestLifecycleContractDisembargoPreventsOvertaking(t *testing.T) {
	input := probeFor(actionObserveDelivery)
	alternatives, err := buildDisembargoConstraints(input)
	if err != nil {
		t.Fatal(err)
	}
	ordered := lifecycleLedger{
		{Kind: actionObserveDelivery, Label: "pre-resolution"},
		{Kind: actionWireDisembargo},
		{Kind: actionObserveDelivery, Label: "post-resolution"},
	}
	if err := checkConstraintAlternatives(alternatives[1:], ordered); err != nil {
		t.Fatal(err)
	}
}

func TestLifecycleContractDisconnectDrainsBeforeConnDone(t *testing.T) {
	ledger := lifecycleLedger{
		{Kind: actionObserveCompletion},
		{Kind: actionObserveCapabilityRelease},
		{Kind: actionObserveConnDone},
	}
	schemaConstraints, err := buildDisconnectConstraints(probeFor(actionLocalClose))
	if err != nil {
		t.Fatal(err)
	}
	if err := checkConstraintAlternatives(schemaConstraints, ledger); err != nil {
		t.Fatal(err)
	}
	if err := checkConstraintAlternatives(goConnectionDrainConstraints(), ledger); err != nil {
		t.Fatal(err)
	}
}
