package agents

import "strings"

// Launch-argv placeholders resolved from the Spec's effective configuration.
//
// Before this existed, a built-in spec carried its model in TWO places: the `Model`
// field AND a literal baked into `HeadlessArgs`. Config layers set the field and never
// rewrote the args (`applyOverride` in internal/config/runtime.go), and the runner
// substituted only {root}/{prompt} — so `model = "..."` in any config layer changed the
// DISPLAYED value while the process kept launching the built-in literal. The only way to
// actually re-pin a model was to override the whole `headless_args` vector, which is how
// hermes silently lost `--yolo`.
//
// One value, one place: the args carry a placeholder and the resolved field fills it.
const (
	ModelPlaceholder  = "{model}"
	EffortPlaceholder = "{effort}"
)

// LaunchArgsStatus reports placeholders that could not be bound. Empty means everything
// resolved. These strings are STATUS codes in the roster contract.
type LaunchArgsStatus struct {
	ModelUnbound  bool
	EffortUnbound bool
}

// Codes renders the status as roster STATUS codes, sorted for stable output.
func (s LaunchArgsStatus) Codes() []string {
	var out []string
	if s.ModelUnbound {
		out = append(out, "model-unbound")
	}
	if s.EffortUnbound {
		out = append(out, "effort-unknown")
	}
	return out
}

// bindable reports whether a configured value is a real value rather than an absent one.
// `cli-default` means "we could not discover it" — passing it as a flag value would send
// the literal string "cli-default" to the vendor CLI.
func bindable(v string) bool {
	v = strings.TrimSpace(v)
	return v != "" && v != CLIDefault
}

// ResolveLaunchArgs substitutes {model} and {effort} in headlessArgs from the spec's
// effective values.
//
// When a placeholder cannot be bound, the placeholder AND the flag that introduces it are
// BOTH dropped, because a value-taking flag left without its value makes the CLI abort with
// a parse error and write no output — a silent launch failure. Dropping the pair lets the
// vendor CLI fall back to its own configured default, which is the honest outcome, and the
// caller reports `model-unbound` / `effort-unknown` rather than printing a declaration the
// process never received.
func (s Spec) ResolveLaunchArgs() ([]string, LaunchArgsStatus) {
	var status LaunchArgsStatus
	out := make([]string, 0, len(s.HeadlessArgs))
	for _, arg := range s.HeadlessArgs {
		switch arg {
		case ModelPlaceholder:
			if bindable(s.Model) {
				out = append(out, s.Model)
				continue
			}
			status.ModelUnbound = true
			out = dropIntroducingFlag(out)
		case EffortPlaceholder:
			if bindable(s.Reasoning) {
				out = append(out, s.Reasoning)
				continue
			}
			status.EffortUnbound = true
			out = dropIntroducingFlag(out)
		default:
			out = append(out, arg)
		}
	}
	return out, status
}

// dropIntroducingFlag removes the trailing token when it is the flag that would have taken
// the unbound placeholder as its value. A non-flag trailing token is left alone: it is a
// positional argument that happens to precede the placeholder, not its flag.
func dropIntroducingFlag(args []string) []string {
	if n := len(args); n > 0 && strings.HasPrefix(args[n-1], "-") {
		return args[:n-1]
	}
	return args
}

// EffectiveModel returns the model the launch will actually pass, and whether it passes one
// at all. A spec whose args carry no {model} placeholder and no literal model flag launches
// without a model argument: the vendor CLI reads its own config, and the roster must NOT
// claim the declared value as effective.
func (s Spec) EffectiveModel() (string, bool) {
	args, status := s.ResolveLaunchArgs()
	if status.ModelUnbound {
		return "", false
	}
	for i, a := range args {
		if (a == "--model" || a == "-m") && i+1 < len(args) {
			return args[i+1], true
		}
	}
	return "", false
}

// EffectiveEffort mirrors EffectiveModel for the reasoning/effort axis.
func (s Spec) EffectiveEffort() (string, bool) {
	args, status := s.ResolveLaunchArgs()
	if status.EffortUnbound {
		return "", false
	}
	for i, a := range args {
		if (a == "--effort" || a == "--reasoning" || a == "--variant") && i+1 < len(args) {
			return args[i+1], true
		}
	}
	return "", false
}
