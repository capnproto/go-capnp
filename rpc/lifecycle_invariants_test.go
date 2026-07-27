package rpc_test

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"testing"
)

type schemaIdentity struct {
	GitBlob string
	SHA256  string
}

var rpcSchemaSource = struct {
	Canonical schemaIdentity
	Local     schemaIdentity
}{
	Canonical: schemaIdentity{
		GitBlob: "0b821519acf9b221801639dfaf8c7adc25a26b08",
		SHA256:  "2ecc3049d4f7f2d48a3a368dbb9ef4b97b31c1365996d615bd19c267983a1931",
	},
	Local: schemaIdentity{
		GitBlob: "bc8ad6e0d84cce2b3bb159ba855ddbd8368ad5e0",
		SHA256:  "3a2646562b0b8d5f421e1281b7283b39fb8cf7f27a44fa9a6679f6011b2a2da7",
	},
}

type invariantStatus uint8

const (
	invariantStatusInvalid invariantStatus = iota
	invariantSupported
	invariantDeferred
	invariantLocalSchemaDivergent
)

func (s invariantStatus) String() string {
	switch s {
	case invariantSupported:
		return "supported"
	case invariantDeferred:
		return "deferred"
	case invariantLocalSchemaDivergent:
		return "local-schema-divergent"
	default:
		return "invalid"
	}
}

type clauseSourceKind uint8

const (
	clauseSourceInvalid clauseSourceKind = iota
	clauseDeclaration
	clauseSection
)

type reviewedClause struct {
	SourceKind    clauseSourceKind
	Declaration   string
	Excerpt       string
	ExcerptDigest string
}

type constraintBuilder func(ruleInput) (constraintAlternatives, error)

type invariantRule struct {
	Action actionKind
	Probe  ruleInput
	Build  constraintBuilder
}

func newInvariantRule(action actionKind, build constraintBuilder) invariantRule {
	return invariantRule{
		Action: action,
		Probe:  probeFor(action),
		Build:  build,
	}
}

type lifecycleInvariant struct {
	ID            string
	CanonicalBlob string
	Declaration   string
	Slug          string
	Paraphrase    string
	Status        invariantStatus
	Clauses       []reviewedClause
	Rules         []invariantRule

	// A divergent invariant's canonical clauses must be absent from the local
	// schema, while this exact gap-spanning anchor must remain present.
	LocalAbsenceAnchor       string
	LocalAbsenceAnchorDigest string
}

func newLifecycleInvariant(
	declaration string,
	slug string,
	paraphrase string,
	status invariantStatus,
	clauses []reviewedClause,
	rules ...invariantRule,
) lifecycleInvariant {
	return lifecycleInvariant{
		ID:            invariantID(declaration, slug),
		CanonicalBlob: rpcSchemaSource.Canonical.GitBlob,
		Declaration:   declaration,
		Slug:          slug,
		Paraphrase:    paraphrase,
		Status:        status,
		Clauses:       clauses,
		Rules:         rules,
	}
}

func invariantID(declaration, slug string) string {
	return "rpc-v2@" + rpcSchemaSource.Canonical.GitBlob[:8] + ":" + declaration + ":" + slug
}

func clause(declaration, excerpt, digest string) reviewedClause {
	return reviewedClause{
		SourceKind:    clauseDeclaration,
		Declaration:   declaration,
		Excerpt:       excerpt,
		ExcerptDigest: digest,
	}
}

func sectionClause(path, excerpt, digest string) reviewedClause {
	return reviewedClause{
		SourceKind:    clauseSection,
		Declaration:   path,
		Excerpt:       excerpt,
		ExcerptDigest: digest,
	}
}

func lifecycleInvariants() []lifecycleInvariant {
	invariants := []lifecycleInvariant{
		newLifecycleInvariant(
			"Call.questionId",
			"question-id-retirement",
			"A question ID remains occupied until Return and Finish retire it.",
			invariantSupported,
			[]reviewedClause{clause(
				"Call.questionId",
				`A question ID can be reused once both:
- A matching Return has been received from the callee.
- A matching Finish has been sent from the caller.`,
				"2fc81182310ac7fbcf4ad85762066111eb75550d5557df41df5e8650dd527ad0",
			)},
			newInvariantRule(actionLocalCall, buildQuestionRetirementConstraints),
			newInvariantRule(actionWireCall, buildQuestionRetirementConstraints),
			newInvariantRule(actionPeerReturn, buildQuestionRetirementConstraints),
			newInvariantRule(actionWireFinish, buildQuestionRetirementConstraints),
		),
		newLifecycleInvariant(
			"Finish.releaseResultCaps",
			"result-cap-release",
			"Cancel-owned Finish releases every capability in a later result.",
			invariantSupported,
			[]reviewedClause{clause(
				"Finish.releaseResultCaps",
				"If true, all capabilities that were in the results should be considered released.\n"+
					"The sender must not send separate `Release` messages for them.",
				"44dce5d2dd5142f2d2e0f8bac5624c169b67c6ea136a3e91993853a9d0fa9433",
			)},
			newInvariantRule(actionLocalCancel, buildCancelResultConstraints),
			newInvariantRule(actionWireFinish, buildCancelResultConstraints),
			newInvariantRule(actionPeerReturn, buildCancelResultConstraints),
		),
		newLifecycleInvariant(
			"Return.noFinishNeeded",
			"question-retirement",
			"A capability-free noFinishNeeded Return retires the question; a compatibility Finish is harmless.",
			invariantSupported,
			[]reviewedClause{clause(
				"Return.noFinishNeeded",
				"If true, the sender does not need the receiver to send a `Finish` message; its answer table\n"+
					"entry has already been cleaned up. This implies that the results do not contain any\n"+
					"capabilities, since the `Finish` message would normally release those capabilities from\n"+
					"promise pipelining responsibility. The caller may still send a `Finish` message if it wants,\n"+
					"which will be silently ignored by the callee.",
				"e0a3a15a1a573e3d2043432ae38121fe1fcb2e86869cba7c48ebcb2f75c01c1a",
			)},
			newInvariantRule(actionPeerReturn, buildNoFinishNeededConstraints),
			newInvariantRule(actionWireFinish, buildNoFinishNeededConstraints),
			newInvariantRule(actionPeerFinish, buildNoFinishNeededConstraints),
		),
		newLifecycleInvariant(
			"CapDescriptor.senderPromise",
			"single-resolution",
			"Each senderPromise incarnation receives at most one Resolve.",
			invariantSupported,
			[]reviewedClause{
				clause(
					"CapDescriptor.senderPromise",
					"A promise that the sender will resolve later. The sender will send exactly one Resolve\n"+
						"message at a future point in time to replace this promise. Note that even if the same\n"+
						"`senderPromise` is received multiple times, only one `Resolve` is sent to cover all of\n"+
						"them. If `senderPromise` is released before the `Resolve` is sent, the sender (of this\n"+
						"`CapDescriptor`) may choose not to send the `Resolve` at all.",
					"f9e40d053b672f3a85e3da394d13a1b879bc49bf8c28ac1dedfe0c09235288c1",
				),
				clause(
					"Resolve.promiseId",
					"promiseId @0 :ExportId;\n"+
						"The ID of the promise to be resolved.\n\n"+
						"Unlike all other instances of `ExportId` sent from the exporter, the `Resolve` message does\n"+
						"_not_ increase the reference count of `promiseId`. In fact, it is expected that the receiver\n"+
						"will release the export soon after receiving `Resolve`, and the sender will not send this\n"+
						"`ExportId` again until it has been released and recycled.\n\n"+
						"When an export ID sent over the wire (e.g. in a `CapDescriptor`) is indicated to be a promise,\n"+
						"this indicates that the sender will follow up at some point with a `Resolve` message. If the\n"+
						"same `promiseId` is sent again before `Resolve`, still only one `Resolve` is sent. If the\n"+
						"same ID is sent again later _after_ a `Resolve`, it can only be because the export's\n"+
						"reference count hit zero in the meantime and the ID was re-assigned to a new export, therefore\n"+
						"this later promise does _not_ correspond to the earlier `Resolve`.",
					"184087dca8b5720cca48c314fcc1ae86288acaef0bfc44a2d3143c4f9548ba4c",
				),
				clause(
					"ExportId.referenceCounting",
					"ExportId/ImportIds are subject to reference counting. Whenever an `ExportId` is sent over the\n"+
						"wire (from the exporter to the importer), the export's reference count is incremented (unless\n"+
						"otherwise specified). The reference count is later decremented by a `Release` message. Since\n"+
						"the `Release` message can specify an arbitrary number by which to reduce the reference count, the\n"+
						"importer should usually batch reference decrements and only send a `Release` when it believes the\n"+
						"reference count has hit zero. Of course, it is possible that a new reference to the export is\n"+
						"in-flight at the time that the `Release` message is sent, so it is necessary for the exporter to\n"+
						"keep track of the reference count on its end as well to avoid race conditions.",
					"d12924dc161e6c484bed6d74651d47106e74d7a9ca6e644552a50bb2a2158239",
				),
			},
			newInvariantRule(actionPeerResolve, buildPromiseSingleResolutionConstraints),
			newInvariantRule(actionWireResolve, buildPromiseSingleResolutionConstraints),
			newInvariantRule(actionLocalResolve, buildPromiseSingleResolutionConstraints),
		),
		newLifecycleInvariant(
			"Resolve.promiseId",
			"late-resolution-release",
			"A late Resolve for a released imported promise releases its replacement without reviving the promise.",
			invariantSupported,
			[]reviewedClause{
				clause(
					"Resolve.promiseId",
					"If a promise ID's reference count reaches zero before a `Resolve` is sent, the `Resolve`\n"+
						"message may or may not still be sent (the `Resolve` may have already been in-flight when\n"+
						"`Release` was sent, but if the `Release` is received before `Resolve` then there is no longer\n"+
						"any reason to send a `Resolve`). Thus a `Resolve` may be received for a promise of which\n"+
						"the receiver has no knowledge, because it already released it earlier. In this case, the\n"+
						"receiver should simply release the capability to which the promise resolved.",
					"bc9d0aa20433b598b0bab8e75b7858856416baf9b3c3e73b04a59b3c850a2e25",
				),
				clause(
					"Resolve.cap",
					"cap @1 :CapDescriptor;\n"+
						"The object to which the promise resolved.\n\n"+
						"The sender promises that from this point forth, until `promiseId` is released, it shall\n"+
						"simply forward all messages to the capability designated by `cap`. This is true even if\n"+
						"`cap` itself happens to designate another promise, and that other promise later resolves --\n"+
						"messages sent to `promiseId` shall still go to that other promise, not to its resolution.\n"+
						"This is important in the case that the receiver of the `Resolve` ends up sending a\n"+
						"`Disembargo` message towards `promiseId` in order to control message ordering -- that\n"+
						"`Disembargo` really needs to reflect back to exactly the object designated by `cap` even\n"+
						"if that object is itself a promise.",
					"10258fb3803ddde67b236d440fe77dd1e1677a0287e92413cf452385ed6d39f0",
				),
				clause(
					"CapDescriptor.referenceCounting",
					"Keep in mind that `ExportIds` in a `CapDescriptor` are subject to reference counting. See the\n"+
						"description of `ExportId`.",
					"59d2d6daad803f57b474677a4dd7bba1568777774f2d12b7843de56f36ffb964",
				),
				clause(
					"ExportId.referenceCounting",
					"ExportId/ImportIds are subject to reference counting. Whenever an `ExportId` is sent over the\n"+
						"wire (from the exporter to the importer), the export's reference count is incremented (unless\n"+
						"otherwise specified). The reference count is later decremented by a `Release` message. Since\n"+
						"the `Release` message can specify an arbitrary number by which to reduce the reference count, the\n"+
						"importer should usually batch reference decrements and only send a `Release` when it believes the\n"+
						"reference count has hit zero. Of course, it is possible that a new reference to the export is\n"+
						"in-flight at the time that the `Release` message is sent, so it is necessary for the exporter to\n"+
						"keep track of the reference count on its end as well to avoid race conditions.",
					"d12924dc161e6c484bed6d74651d47106e74d7a9ca6e644552a50bb2a2158239",
				),
			},
			newInvariantRule(actionLocalRelease, buildLateResolveConstraints),
			newInvariantRule(actionWireRelease, buildLateResolveConstraints),
			newInvariantRule(actionPeerResolve, buildLateResolveConstraints),
		),
		newLifecycleInvariant(
			"Release.referenceCount",
			"reference-decrement",
			"Release decrements the exporter-owned reference count and retires only at zero.",
			invariantSupported,
			[]reviewedClause{
				clause(
					"Release.referenceCount",
					"referenceCount @1 :UInt32;\n"+
						"The amount by which to decrement the reference count. The export is only actually released\n"+
						"when the reference count reaches zero.",
					"7aad0df169e059d72aa358590cca189da7f0c2bc095774b65e792c3dc3f301c9",
				),
				clause(
					"ExportId.referenceCounting",
					"ExportId/ImportIds are subject to reference counting. Whenever an `ExportId` is sent over the\n"+
						"wire (from the exporter to the importer), the export's reference count is incremented (unless\n"+
						"otherwise specified). The reference count is later decremented by a `Release` message. Since\n"+
						"the `Release` message can specify an arbitrary number by which to reduce the reference count, the\n"+
						"importer should usually batch reference decrements and only send a `Release` when it believes the\n"+
						"reference count has hit zero. Of course, it is possible that a new reference to the export is\n"+
						"in-flight at the time that the `Release` message is sent, so it is necessary for the exporter to\n"+
						"keep track of the reference count on its end as well to avoid race conditions.",
					"d12924dc161e6c484bed6d74651d47106e74d7a9ca6e644552a50bb2a2158239",
				),
				clause(
					"ImportId.lifetime",
					"An `ImportId` remains valid in importer -> exporter messages until the importer has sent\n"+
						"`Release` messages that (it believes) have reduced the reference count to zero.",
					"97c851b796f1f066979f5c522eae91d1f1bac9f86efa8474edfece5ecccebe3c",
				),
			},
			newInvariantRule(actionLocalRelease, buildReferenceDecrementConstraints),
			newInvariantRule(actionPeerRelease, buildReferenceDecrementConstraints),
			newInvariantRule(actionWireRelease, buildReferenceDecrementConstraints),
		),
		newLifecycleInvariant(
			"Disembargo",
			"promise-resolution-ordering",
			"Disembargo preserves the required delivery edge across a promise path change.",
			invariantSupported,
			[]reviewedClause{
				clause(
					"Disembargo",
					"Message sent to indicate that an embargo on a recently-resolved promise may now be lifted.\n\n"+
						"Embargos are used to enforce E-order in the presence of promise resolution. That is, if an\n"+
						"application makes two calls foo() and bar() on the same capability reference, in that order,\n"+
						"the calls should be delivered in the order in which they were made. But if foo() is called\n"+
						"on a promise, and that promise happens to resolve before bar() is called, then the two calls\n"+
						"may travel different paths over the network, and thus could arrive in the wrong order. In\n"+
						"this case, the call to `bar()` must be embargoed, and a `Disembargo` message must be sent along\n"+
						"the same path as `foo()` to ensure that the `Disembargo` arrives after `foo()`. Once the\n"+
						"`Disembargo` arrives, `bar()` can then be delivered.",
					"d6af96339c82cba73124d8fd95824f24af02c682154e46a102b93787cf468f69",
				),
				sectionClause(
					"Protocol.eOrder",
					"Unless otherwise specified, messages must be delivered to the receiving application in the same\n"+
						"order in which they were initiated by the sending application. The goal is to support \"E-Order\",\n"+
						"which states that two calls made on the same reference must be delivered in the order which they\n"+
						"were made:\n"+
						"http://erights.org/elib/concurrency/partial-order.html",
					"76b68be8a32fb4413e43adb735795019693c6842f3c2ad7a7c78dd41071912cf",
				),
				sectionClause(
					"Disembargo.samePathAlternative",
					"Note that in the case where Carol actually lives in Vat B (i.e., the same vat that the promise\n"+
						"already pointed at), no embargo is needed, because the pipelined calls are delivered over the\n"+
						"same path as the later direct calls.\n\n"+
						"Keep in mind that promise resolution happens both in the form of Resolve messages as well as\n"+
						"Return messages (which resolve PromisedAnswers). Embargos apply in both cases.",
					"a1151b6b770d9b51b9fa3b84e0490ddc5dd67a9c45bd411792eaff48d360055d",
				),
				clause(
					"Disembargo.context.senderLoopback/receiverLoopback",
					"senderLoopback @1 :EmbargoId;\n"+
						"The sender is requesting a disembargo on a promise that is known to resolve back to a\n"+
						"capability hosted by the sender. As soon as the receiver has echoed back all pipelined calls\n"+
						"on this promise, it will deliver the Disembargo back to the sender with `receiverLoopback`\n"+
						"set to the same value as `senderLoopback`. This value is chosen by the sender, and since\n"+
						"it is also consumed be the sender, the sender can use whatever strategy it wants to make sure\n"+
						"the value is unambiguous.\n\n"+
						"The receiver must verify that the target capability actually resolves back to the sender's\n"+
						"vat. Otherwise, the sender has committed a protocol error and should be disconnected.\n\n"+
						"receiverLoopback @2 :EmbargoId;\n"+
						"The receiver previously sent a `senderLoopback` Disembargo towards a promise resolving to\n"+
						"this capability, and that Disembargo is now being echoed back.",
					"9c8d0f3d3211386d47809376f661f0433203388480ce8f480360548c8f5c2755",
				),
			},
			newInvariantRule(actionWireCall, buildDisembargoConstraints),
			newInvariantRule(actionPeerResolve, buildDisembargoConstraints),
			newInvariantRule(actionWireDisembargo, buildDisembargoConstraints),
			newInvariantRule(actionPeerDisembargo, buildDisembargoConstraints),
			newInvariantRule(actionObserveDelivery, buildDisembargoConstraints),
		),
		newLifecycleInvariant(
			"FourTables.connectionLoss",
			"terminal-drain",
			"Disconnect fails questions, breaks imports, and implicitly releases exports and answers.",
			invariantSupported,
			[]reviewedClause{
				sectionClause(
					"FourTables.connectionLoss",
					"When a Cap'n Proto connection is lost, everything on the four tables is lost. All questions are\n"+
						"canceled and throw exceptions. All imports become broken (all future calls to them throw\n"+
						"exceptions). All exports and answers are implicitly released. The only things not lost are\n"+
						"persistent capabilities (`SturdyRef`s). The application must plan for this and should respond by\n"+
						"establishing a new connection and restoring from these persistent capabilities.",
					"58cfecf459448a64a40b153147d98095844622698f502e5fac091af24a363902",
				),
				clause(
					"ExportId.connectionLoss",
					"When a connection is lost, all exports are implicitly released. It is not possible to restore\n"+
						"a connection state after disconnect (although a transport layer could implement a concept of\n"+
						"persistent connections if it is transparent to the RPC layer).",
					"669cc81e9893efaa6927f21d4e2d321342c32e04c48063339a9002a8224e4782",
				),
				clause(
					"Message.abort",
					"abort @1 :Exception;\n"+
						"Sent when a connection is being aborted due to an unrecoverable error. This could be e.g.\n"+
						"because the sender received an invalid or nonsensical message or because the sender had an\n"+
						"internal error. The sender will shut down the outgoing half of the connection after `abort`\n"+
						"and will completely close the connection shortly thereafter (it's up to the sender how much\n"+
						"of a time buffer they want to offer for the client to receive the `abort` before the\n"+
						"connection is reset).",
					"472e6cd2d128015cdb97b5e482f8b8588d467e4486e1c83aeea913969cb4f789",
				),
			},
			newInvariantRule(actionLocalClose, buildDisconnectConstraints),
			newInvariantRule(actionPeerAbort, buildDisconnectConstraints),
		),
		newLifecycleInvariant(
			"Return.releaseParamCaps",
			"param-cap-release",
			"Return releases each parameter capability reference when releaseParamCaps is true.",
			invariantDeferred,
			[]reviewedClause{clause(
				"Return.releaseParamCaps",
				"If true, all capabilities that were in the params should be considered released. The sender\n"+
					"must not send separate `Release` messages for them. Level 0 implementations in particular\n"+
					"should always set this true. This defaults true because if level 0 implementations forget to\n"+
					"set it they'll never notice (just silently leak caps), but if level >=1 implementations forget\n"+
					"to set it to false they'll quickly get errors.\n\n"+
					"The receiver should act as if the sender had sent a release message with count=1 for each\n"+
					"CapDescriptor in the original Call message.",
				"1210e77a0dcff7ed80de75074c3e5042e2a91a5148f54ab2ce02d3c832ee74b4",
			)},
		),
	}

	divergent := newLifecycleInvariant(
		"CapDescriptor.senderPromise",
		"resolved-target-shortening",
		"A reused senderPromise must not name an export whose Resolve was already sent.",
		invariantLocalSchemaDivergent,
		[]reviewedClause{clause(
			"CapDescriptor.senderPromise",
			"`senderPromise` must not refer to an export for which `Resolve` was already sent in the past.\n"+
				"The reason for this is that the receiver may no longer have a record of that past `Resolve`.\n"+
				"In fact, it's extremely likely that when the receiver received the `Resolve`, it promptly\n"+
				"released the promise from its import table, since it updated all references to the new taget.\n"+
				"This means that if the old promise's export ID were sent again here, the receiver wouldn't\n"+
				"know anything about it and wouldn't know that it's already resolved.\n\n"+
				"In order to avoid the above case, if the sender detects that promise A is resolving to\n"+
				"promise B, but promise B has already previously resolved to C, then the sender should simply\n"+
				"resolve promise A directly to C.",
			"56484603507c1ff849b6c1f0eb8eedc2f5178054baa506e7f64b19f19bd04fea",
		)},
	)
	divergent.LocalAbsenceAnchor = "If `senderPromise` is released before the `Resolve` is sent, the sender (of this\n" +
		"`CapDescriptor`) may choose not to send the `Resolve` at all.\n\n" +
		"receiverHosted @3 :ImportId;"
	divergent.LocalAbsenceAnchorDigest = "c7987ccde4230d26f02c520aaf810e99244479c31767eb4897519426e71a1368"
	return append(invariants, divergent)
}

var invariantSlugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

func validateLifecycleInvariants(invariants []lifecycleInvariant, schemaText string) error {
	normalizedSchema := normalizeSchemaText(schemaText)
	seen := make(map[string]struct{}, len(invariants))
	for i, inv := range invariants {
		if inv.ID == "" {
			return fmt.Errorf("invariant %d: empty ID", i)
		}
		if _, ok := seen[inv.ID]; ok {
			return fmt.Errorf("invariant %d: duplicate ID %q", i, inv.ID)
		}
		seen[inv.ID] = struct{}{}
		if inv.ID != invariantID(inv.Declaration, inv.Slug) {
			return fmt.Errorf("invariant %q: ID does not match declaration and slug", inv.ID)
		}
		if inv.CanonicalBlob != rpcSchemaSource.Canonical.GitBlob {
			return fmt.Errorf("invariant %q: canonical blob %q; want %q",
				inv.ID, inv.CanonicalBlob, rpcSchemaSource.Canonical.GitBlob)
		}
		if !invariantSlugPattern.MatchString(inv.Slug) {
			return fmt.Errorf("invariant %q: invalid slug %q", inv.ID, inv.Slug)
		}
		if inv.Paraphrase == "" {
			return fmt.Errorf("invariant %q: empty paraphrase", inv.ID)
		}
		if len(inv.Clauses) == 0 {
			return fmt.Errorf("invariant %q: no reviewed clauses", inv.ID)
		}
		for j, clause := range inv.Clauses {
			if clause.SourceKind != clauseDeclaration && clause.SourceKind != clauseSection {
				return fmt.Errorf("invariant %q clause %d: invalid source kind %d", inv.ID, j, clause.SourceKind)
			}
			if clause.Declaration == "" || clause.Excerpt == "" {
				return fmt.Errorf("invariant %q clause %d: empty declaration or excerpt", inv.ID, j)
			}
			if err := validateDigest(clause.ExcerptDigest); err != nil {
				return fmt.Errorf("invariant %q clause %q: %v", inv.ID, clause.Declaration, err)
			}
			if got := digestNormalizedText(clause.Excerpt); got != clause.ExcerptDigest {
				return fmt.Errorf("invariant %q clause %q: excerpt digest = %s; want %s",
					inv.ID, clause.Declaration, got, clause.ExcerptDigest)
			}
			excerptPresent := strings.Contains(normalizedSchema, normalizeSchemaText(clause.Excerpt))
			if inv.Status == invariantLocalSchemaDivergent && excerptPresent {
				return fmt.Errorf("invariant %q clause %q: divergent excerpt unexpectedly present in local schema",
					inv.ID, clause.Declaration)
			}
			if inv.Status != invariantLocalSchemaDivergent && !excerptPresent {
				return fmt.Errorf("invariant %q clause %q: excerpt not found in local schema",
					inv.ID, clause.Declaration)
			}
		}

		switch inv.Status {
		case invariantSupported:
			if len(inv.Rules) == 0 {
				return fmt.Errorf("invariant %q: supported invariant has no executable rules", inv.ID)
			}
			if inv.LocalAbsenceAnchor != "" || inv.LocalAbsenceAnchorDigest != "" {
				return fmt.Errorf("invariant %q: supported invariant has a local-absence anchor", inv.ID)
			}
			if err := validateInvariantRules(inv); err != nil {
				return err
			}
		case invariantDeferred:
			if len(inv.Rules) != 0 {
				return fmt.Errorf("invariant %q: deferred invariant has executable rules", inv.ID)
			}
			if inv.LocalAbsenceAnchor != "" || inv.LocalAbsenceAnchorDigest != "" {
				return fmt.Errorf("invariant %q: deferred invariant has a local-absence anchor", inv.ID)
			}
		case invariantLocalSchemaDivergent:
			if len(inv.Rules) != 0 {
				return fmt.Errorf("invariant %q: local-schema-divergent invariant has executable rules", inv.ID)
			}
			if inv.LocalAbsenceAnchor == "" {
				return fmt.Errorf("invariant %q: local-schema-divergent invariant has no absence anchor", inv.ID)
			}
			if err := validateDigest(inv.LocalAbsenceAnchorDigest); err != nil {
				return fmt.Errorf("invariant %q local-absence anchor: %v", inv.ID, err)
			}
			if got := digestNormalizedText(inv.LocalAbsenceAnchor); got != inv.LocalAbsenceAnchorDigest {
				return fmt.Errorf("invariant %q: local-absence anchor digest = %s; want %s",
					inv.ID, got, inv.LocalAbsenceAnchorDigest)
			}
			if !strings.Contains(normalizedSchema, normalizeSchemaText(inv.LocalAbsenceAnchor)) {
				return fmt.Errorf("invariant %q: local-absence anchor not found in local schema", inv.ID)
			}
		default:
			return fmt.Errorf("invariant %q: invalid status %d", inv.ID, inv.Status)
		}
	}
	return nil
}

func validateInvariantRules(inv lifecycleInvariant) error {
	rules := slices.Clone(inv.Rules)
	slices.SortFunc(rules, func(a, b invariantRule) int {
		return int(a.Action) - int(b.Action)
	})
	for i, rule := range rules {
		if !rule.Action.valid() {
			return fmt.Errorf("invariant %q: invalid rule action %d", inv.ID, rule.Action)
		}
		if i > 0 && rules[i-1].Action == rule.Action {
			return fmt.Errorf("invariant %q: duplicate rule action %d", inv.ID, rule.Action)
		}
		if rule.Build == nil {
			return fmt.Errorf("invariant %q action %d: nil constraint builder", inv.ID, rule.Action)
		}
		alternatives, err := rule.Build(rule.Probe)
		if err != nil {
			return fmt.Errorf("invariant %q action %d: build constraints: %v", inv.ID, rule.Action, err)
		}
		if err := validateConstraintAlternatives(alternatives); err != nil {
			return fmt.Errorf("invariant %q action %d: invalid constraints: %v", inv.ID, rule.Action, err)
		}
	}
	return nil
}

func validateDigest(digest string) error {
	decoded, err := hex.DecodeString(digest)
	if err != nil || len(decoded) != sha256.Size {
		return fmt.Errorf("invalid SHA-256 digest %q", digest)
	}
	return nil
}

func normalizeSchemaText(s string) string {
	lines := strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "#"))
		}
		lines[i] = line
	}
	return collapseASCIIWhitespace(strings.Join(lines, "\n"))
}

func collapseASCIIWhitespace(s string) string {
	var b strings.Builder
	space := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		isSpace := c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f' || c == '\v'
		if isSpace {
			space = b.Len() > 0
			continue
		}
		if space {
			b.WriteByte(' ')
			space = false
		}
		b.WriteByte(c)
	}
	return b.String()
}

func digestNormalizedText(s string) string {
	sum := sha256.Sum256([]byte("capnp-rpc-clause-v1\x00" + normalizeSchemaText(s)))
	return hex.EncodeToString(sum[:])
}

func gitBlobID(data []byte) string {
	header := "blob " + strconv.Itoa(len(data)) + "\x00"
	sum := sha1.Sum(append([]byte(header), data...))
	return hex.EncodeToString(sum[:])
}

func readLocalRPCSchema(t *testing.T) []byte {
	t.Helper()
	_, sourceFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("locate lifecycle invariant test source")
	}
	path := filepath.Join(filepath.Dir(sourceFile), "..", "std", "capnp", "rpc.capnp")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read local RPC schema: %v", err)
	}
	return data
}

func TestLifecycleInvariantRegistry(t *testing.T) {
	schema := readLocalRPCSchema(t)
	sum := sha256.Sum256(schema)
	if got := hex.EncodeToString(sum[:]); got != rpcSchemaSource.Local.SHA256 {
		t.Fatalf("local rpc.capnp SHA-256 = %s; want %s", got, rpcSchemaSource.Local.SHA256)
	}
	if got := gitBlobID(schema); got != rpcSchemaSource.Local.GitBlob {
		t.Fatalf("local rpc.capnp git blob = %s; want %s", got, rpcSchemaSource.Local.GitBlob)
	}
	if rpcSchemaSource.Canonical == rpcSchemaSource.Local {
		t.Fatal("canonical and local schema identities are not distinct")
	}
	if got := len(lifecycleInvariants()); got != 10 {
		t.Fatalf("registry length = %d; want 10", got)
	}
	if err := validateLifecycleInvariants(lifecycleInvariants(), string(schema)); err != nil {
		t.Fatal(err)
	}
}

func TestLifecycleInvariantSchemaPathIsCWDIndependent(t *testing.T) {
	want := readLocalRPCSchema(t)
	t.Chdir(t.TempDir())
	if got := readLocalRPCSchema(t); string(got) != string(want) {
		t.Fatal("local RPC schema changed after working-directory change")
	}
}

func TestLifecycleInvariantRegistryDiagnostics(t *testing.T) {
	schema := string(readLocalRPCSchema(t))
	tests := []struct {
		name string
		edit func([]lifecycleInvariant)
		want string
	}{
		{
			name: "duplicate ID",
			edit: func(invariants []lifecycleInvariant) {
				invariants[1].ID = invariants[0].ID
			},
			want: "duplicate ID",
		},
		{
			name: "bad clause digest",
			edit: func(invariants []lifecycleInvariant) {
				invariants[0].Clauses[0].ExcerptDigest = "not-a-digest"
			},
			want: "invalid SHA-256 digest",
		},
		{
			name: "invalid action",
			edit: func(invariants []lifecycleInvariant) {
				invariants[0].Rules[0].Action = actionKind(255)
			},
			want: "invalid rule action",
		},
		{
			name: "supported without rules",
			edit: func(invariants []lifecycleInvariant) {
				invariants[0].Rules = nil
			},
			want: "supported invariant has no executable rules",
		},
		{
			name: "deferred with rule",
			edit: func(invariants []lifecycleInvariant) {
				invariants[len(invariants)-2].Rules = []invariantRule{
					newInvariantRule(actionLocalCall, buildQuestionRetirementConstraints),
				}
			},
			want: "deferred invariant has executable rules",
		},
		{
			name: "supported excerpt drift",
			edit: func(invariants []lifecycleInvariant) {
				invariants[0].Clauses[0].Excerpt = "not present in the schema"
				invariants[0].Clauses[0].ExcerptDigest = digestNormalizedText(invariants[0].Clauses[0].Excerpt)
			},
			want: "excerpt not found in local schema",
		},
		{
			name: "divergence loses local anchor",
			edit: func(invariants []lifecycleInvariant) {
				last := len(invariants) - 1
				invariants[last].LocalAbsenceAnchor = "not present in the schema"
				invariants[last].LocalAbsenceAnchorDigest = digestNormalizedText(invariants[last].LocalAbsenceAnchor)
			},
			want: "local-absence anchor not found",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			invariants := lifecycleInvariants()
			test.edit(invariants)
			err := validateLifecycleInvariants(invariants, schema)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("validateLifecycleInvariants() error = %v; want substring %q", err, test.want)
			}
		})
	}
}
