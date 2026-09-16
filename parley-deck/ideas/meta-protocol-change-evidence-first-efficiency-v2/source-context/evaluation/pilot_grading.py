#!/usr/bin/env python3
"""Prepare blinded inputs and reject incomplete or contradictory pilot grades."""
import hashlib
import random
import re

import yaml

from pilot_analysis import blind_identifier, canonical_hash

WEIGHTS = {"correctness": 30, "completeness": 25, "verification": 20, "tradeoffs": 15, "actionability": 10}


def strip_frontmatter(text):
    lines = text.splitlines(keepends=True)
    if not lines or lines[0].strip() != "---":
        return text, None
    end = next((index for index in range(1, len(lines)) if lines[index].strip() == "---"), None)
    if end is None:
        raise ValueError("Unclosed candidate frontmatter")
    metadata = yaml.safe_load("".join(lines[1:end]))
    if metadata is not None and not isinstance(metadata, dict):
        raise ValueError("Candidate frontmatter must be a mapping")
    return "".join(lines[end+1:]), metadata


def require_disjoint_graders(graders, authors):
    if len(graders) < 2:
        raise ValueError("Two actual independent grading configurations are required")
    ids, models = set(), set()
    author_ids = {row["id"] for row in authors}
    author_models = {row["requested_model"] for row in authors}
    author_models.update(row.get("reported_model") for row in authors if row.get("reported_model"))
    for grader in graders:
        identity, model = grader.get("id"), grader.get("requested_model")
        resolved = {model, grader.get("reported_model")} - {None, ""}
        if not identity or not model or identity in ids or resolved & models:
            raise ValueError("Grader identity or configured model is missing or duplicated")
        if identity in author_ids or model in author_models or grader.get("reported_model") in author_models:
            raise ValueError("Task or candidate author cannot serve as a blind grader")
        ids.add(identity)
        models.update(resolved)
    return {"configured_model_disjoint": True,
            "resolved_identity_coverage": sum(bool(row.get("reported_model")) for row in graders),
            "limit": "Configuration disjointness is checked; hidden aliases and shared provider ancestry are not independently proved."}


def prepare_blind_bundle(task, candidates, salt, grader_id, seed, identity_terms=()):
    public, mapping = [], []
    seen = set()
    for candidate in candidates:
        arm = candidate["arm"]
        if arm not in ("solo", "duo", "full") or arm in seen:
            raise ValueError("Duplicate or unknown arm")
        seen.add(arm)
        original = candidate["answer"]
        body, metadata = strip_frontmatter(original)
        identity = blind_identifier(task["id"], arm, salt)
        code = candidate.get("source_code")
        combined = body + ("\n" + code if code else "")
        leaks = [term for term in identity_terms if re.search(r"(?<![A-Za-z0-9])" + re.escape(term) + r"(?![A-Za-z0-9])", combined, re.I)]
        acceptance = candidate.get("acceptance")
        if acceptance:
            acceptance = {key: acceptance.get(key) for key in ("status", "executed", "skipped", "cases")}
        public.append({"candidate_id": identity, "answer": body, "source_code": code,
                       "acceptance": acceptance})
        mapping.append({"candidate_id": identity, "task": task["id"], "arm": arm,
                        "original_sha256": hashlib.sha256(original.encode()).hexdigest(),
                        "blinded_body_sha256": hashlib.sha256(body.encode()).hexdigest(),
                        "source_sha256": hashlib.sha256(code.encode()).hexdigest() if code is not None else None,
                        "removed_metadata": metadata, "possible_identity_leaks": leaks})
    rng = random.Random(str(seed) + ":" + grader_id)
    rng.shuffle(public)
    public_task = {key: task[key] for key in ("title", "category", "task", "contract") if key in task}
    bundle = {"schema": 1, "task": public_task, "candidates": public,
              "weights": WEIGHTS, "scale": "Use integer scores 0 through 4 for each criterion.",
              "blindness_limit": "Formatting or answer content may reveal identity; do not use it to score substance."}
    return bundle, {"schema": 1, "grader_id": grader_id, "bundle_sha256": canonical_hash(bundle), "candidates": mapping}


def validate_grade(grade, candidate_id):
    if not isinstance(candidate_id, str) or not candidate_id or grade.get("candidate_id") != candidate_id:
        raise ValueError("Wrong or missing blind candidate ID")
    criteria = grade.get("criteria")
    if not isinstance(criteria, dict) or set(criteria) != set(WEIGHTS):
        raise ValueError("Every rubric criterion is required, with no extra criteria")
    for value in criteria.values():
        if type(value) is not int or not 0 <= value <= 4:
            raise ValueError("Criterion scores must be integers 0 through 4")
    if type(grade.get("passes")) is not bool:
        raise ValueError("An explicit pass/fail judgement is required")
    for key in ("critical_findings", "supported_findings", "unsupported_findings", "missed_seeded_findings"):
        value = grade.get(key)
        if not isinstance(value, list) or any(not isinstance(item, str) or not item.strip() for item in value):
            raise ValueError("Finding categories must be explicit string lists")
    evidence = grade.get("evidence")
    if not isinstance(evidence, list) or not evidence:
        raise ValueError("Grade needs concrete evidence, not only scores")
    for row in evidence:
        if not isinstance(row, dict) or not all(isinstance(row.get(key), str) and row[key].strip() for key in ("locator", "claim")):
            raise ValueError("Evidence requires a locator and claim")
    return sum(criteria[key] * WEIGHTS[key] / 4 for key in WEIGHTS)


def compare_grades(first, second, acceptance=None):
    identity = first.get("candidate_id")
    scores = [validate_grade(first, identity), validate_grade(second, identity)]
    reasons = []
    if first["passes"] != second["passes"]:
        reasons.append("pass-fail-disagreement")
    if set(first["critical_findings"]) != set(second["critical_findings"]):
        reasons.append("critical-invariant-disagreement")
    if (set(first["supported_findings"]) & set(second["unsupported_findings"]) or
            set(second["supported_findings"]) & set(first["unsupported_findings"])):
        reasons.append("factual-finding-conflict")
    if abs(scores[0]-scores[1]) > 15:
        reasons.append("score-spread-over-15")
    if acceptance and acceptance["status"] in ("harness-error", "harness-timeout"):
        reasons.append("acceptance-harness-unavailable")
    elif acceptance and acceptance["status"] != "passed" and (first["passes"] or second["passes"]):
        reasons.append("grader-conflicts-with-executable-result")
    return {"candidate_id": identity, "scores": scores, "needs_adjudication": bool(reasons),
            "reasons": reasons, "quality_score": None if reasons else sum(scores)/2,
            "passes": None if reasons else first["passes"]}


def validate_grade_batch(bundle, grades, expected_hash):
    if canonical_hash(bundle) != expected_hash:
        raise ValueError("Grade inputs changed after dispatch")
    expected = {candidate["candidate_id"] for candidate in bundle["candidates"]}
    actual = [grade.get("candidate_id") for grade in grades]
    if len(actual) != len(set(actual)) or set(actual) != expected:
        raise ValueError("Missing, duplicated or unexpected candidate grade")
    return [validate_grade(grade, grade["candidate_id"]) for grade in grades]
