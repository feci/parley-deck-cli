package telemetry

// Shared runner/preflight grammar. A bare 503 is a whole line, not a number
// embedded in a byte count, source location or diagnostic sentence.
const ProviderUnavailablePattern = `(?im)overloaded|overloaded_error|serverOverloaded|server[_ ]error|internalServerError|^\s*503(?:\s+service unavailable)?\s*$|\b(?:http(?:/[0-9.]+)?|status(?:code)?|api error|response(?: status)?|error code)\s*[:=]?\s*5[0-9][0-9]\b|\b5[0-9][0-9]\s+(?:unavailable|service unavailable|internal server error)\b`
