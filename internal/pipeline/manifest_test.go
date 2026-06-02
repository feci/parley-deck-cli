package pipeline

import "testing"

const validLinearYAML = `
schema_version: 1
idea_slug: demo-pipeline
transport: local-dir
participants: [codex, claude, hermes]
blocks:
  - id: business-spec
    kind: deliberation
    stage: business-spec
    role_lens: product-analyst
    output_artifact: BUSINESS_SPEC.md
  - id: technical-spec
    kind: deliberation
    stage: technical-spec
    role_lens: architect
    input_artifacts: [BUSINESS_SPEC.md]
    output_artifact: TECHNICAL_SPEC.md
  - id: impl-design
    kind: deliberation
    stage: implementation-design
    role_lens: tech-lead
    output_artifact: IMPLEMENTATION_DESIGN.md
`

func TestParseValidLinearManifest(t *testing.T) {
	m, err := Parse([]byte(validLinearYAML))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(m.Blocks) != 3 {
		t.Fatalf("blocks = %d, want 3", len(m.Blocks))
	}
	if m.Autonomy != AutonomySupervised {
		t.Fatalf("autonomy = %q, want supervised (default)", m.Autonomy)
	}
	if m.Blocks[0].Kind != KindDeliberation {
		t.Fatalf("block[0].kind = %q", m.Blocks[0].Kind)
	}
}

func TestParseLinearManifestWithMatchingEdges(t *testing.T) {
	y := validLinearYAML + `
edges:
  - {from: business-spec, to: technical-spec}
  - {from: technical-spec, to: impl-design}
`
	if _, err := Parse([]byte(y)); err != nil {
		t.Fatalf("matching linear edges should be valid: %v", err)
	}
}

func TestRejectNonLinearEdges(t *testing.T) {
	// A branch: business-spec fans out to both later blocks.
	y := validLinearYAML + `
edges:
  - {from: business-spec, to: technical-spec}
  - {from: business-spec, to: impl-design}
`
	if _, err := Parse([]byte(y)); err == nil {
		t.Fatal("expected non-linear (branching) edges to be rejected")
	}
}

func TestRejectWrongEdgeCount(t *testing.T) {
	y := validLinearYAML + `
edges:
  - {from: business-spec, to: technical-spec}
`
	if _, err := Parse([]byte(y)); err == nil {
		t.Fatal("expected wrong edge count to be rejected")
	}
}

func TestRejectUnknownKind(t *testing.T) {
	y := `
schema_version: 1
idea_slug: x
transport: local-dir
participants: [a, b]
blocks:
  - {id: only, kind: bogus}
`
	if _, err := Parse([]byte(y)); err == nil {
		t.Fatal("expected unknown block kind to be rejected")
	}
}

func TestRejectDuplicateBlockID(t *testing.T) {
	y := `
schema_version: 1
idea_slug: x
transport: local-dir
participants: [a, b]
blocks:
  - {id: dup, kind: deliberation}
  - {id: dup, kind: deliberation}
`
	if _, err := Parse([]byte(y)); err == nil {
		t.Fatal("expected duplicate block id to be rejected")
	}
}

func TestRejectBadTransport(t *testing.T) {
	y := `
schema_version: 1
idea_slug: x
transport: carrier-pigeon
participants: [a, b]
blocks:
  - {id: b1, kind: deliberation}
`
	if _, err := Parse([]byte(y)); err == nil {
		t.Fatal("expected invalid transport to be rejected")
	}
}

func TestRejectBadSchemaVersion(t *testing.T) {
	y := `
schema_version: 2
idea_slug: x
transport: local-dir
participants: [a, b]
blocks:
  - {id: b1, kind: deliberation}
`
	if _, err := Parse([]byte(y)); err == nil {
		t.Fatal("expected unsupported schema_version to be rejected")
	}
}

func TestRejectProductionMisspelledRisk(t *testing.T) {
	y := `
schema_version: 1
idea_slug: x
transport: local-dir
participants: [a, b]
blocks:
  - {id: deploy, kind: action, risk: prod}
`
	if _, err := Parse([]byte(y)); err == nil {
		t.Fatal("expected invalid risk to be rejected")
	}
}
