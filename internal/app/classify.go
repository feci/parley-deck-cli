package app

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"

	"parley-deck-cli/internal/track"
)

// runClassify implements `parley classify` — a pure, script-checkable §4.0 track
// classifier (idea track-aware-driver, Slice 2). It maps objective inputs to a
// track and, with --declared, validates a declared track against the computed
// floor (exit 4 on an under-tier so CI can gate). Plain output prints only the
// track so scripts can do `t=$(parley classify ...)`.
func runClassify(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("classify", flag.ContinueOnError)
	fs.SetOutput(stderr)
	var (
		files          = fs.Int("files", 0, "number of files touched")
		loc            = fs.Int("loc", 0, "lines of code changed")
		reversible     = fs.Bool("reversible", false, "the change is fully reversible")
		mechVerifiable = fs.Bool("mechanically-verifiable", false, "mechanically verifiable (lint/type/test)")
		security       = fs.Bool("security", false, "security/auth/secrets/payments/privacy/production-infra surface")
		irreversible   = fs.Bool("irreversible", false, "irreversible/destructive op")
		dataMigration  = fs.Bool("data-migration", false, "data migration")
		protocolChange = fs.Bool("protocol-change", false, "protocol change (see COOPERATION.md §7)")
		autoImplement  = fs.Bool("auto-implement", false, "auto_implement idea")
		strictGate     = fs.Bool("strict-gate", false, "strict_gate idea")
		pipeline       = fs.Bool("pipeline", false, "pipeline or action block")
		apiBreak       = fs.Bool("api-break", false, "public-API break")
		schemaBreak    = fs.Bool("schema-break", false, "persisted-schema break")
		declared       = fs.String("declared", "", "optional declared track to validate (fast|standard|deliberation)")
		jsonOut        = fs.Bool("json", false, "print JSON output")
	)
	if err := fs.Parse(args); err != nil {
		return 2
	}

	in := track.Inputs{
		Files: *files, LOC: *loc, Reversible: *reversible, MechVerifiable: *mechVerifiable,
		ProtocolChange: *protocolChange, Security: *security, Irreversible: *irreversible,
		DataMigration: *dataMigration, StrictGate: *strictGate, AutoImplement: *autoImplement,
		Pipeline: *pipeline, APIBreak: *apiBreak, SchemaBreak: *schemaBreak,
	}
	computed, reason := track.Classify(in)

	valid := true
	var msg string
	if *declared != "" {
		dt, ok := track.Normalize(*declared)
		if !ok {
			fmt.Fprintf(stderr, "unknown --declared track %q (use fast|standard|deliberation)\n", *declared)
			return 2
		}
		if trackRigor(dt) < trackRigor(computed) {
			valid = false
			msg = fmt.Sprintf("declared track %q is under-tiered; the classifier floor is %q (%s)", dt, computed, reason)
		}
	}

	if *jsonOut {
		out := map[string]any{"track": string(computed), "reason": reason}
		if *declared != "" {
			out["declared"] = *declared
			out["valid"] = valid
			if !valid {
				out["message"] = msg
			}
		}
		b, _ := json.MarshalIndent(out, "", "  ")
		fmt.Fprintln(stdout, string(b))
	} else {
		fmt.Fprintln(stdout, string(computed))
		if !valid {
			fmt.Fprintln(stderr, msg)
		}
	}
	if !valid {
		return 4
	}
	return 0
}

// trackRigor orders the tracks by strictness for under-tier validation.
func trackRigor(t track.Track) int {
	switch t {
	case track.Fast:
		return 0
	case track.Deliberation:
		return 2
	default:
		return 1
	}
}
