import unittest

from pilot_analysis import canonical_hash
from pilot_grading import compare_grades, prepare_blind_bundle, require_disjoint_graders, strip_frontmatter, validate_grade, validate_grade_batch


def grade(identity="candidate-abcd", value=3):
    return {"candidate_id": identity, "criteria": {key: value for key in ("correctness", "completeness", "verification", "tradeoffs", "actionability")},
            "passes": True, "critical_findings": [], "supported_findings": [], "unsupported_findings": [],
            "missed_seeded_findings": [], "evidence": [{"locator": "answer paragraph 1", "claim": "Concrete verification is specified"}]}


class PilotGradingTests(unittest.TestCase):
    def test_blinding_preserves_substance_and_hashes(self):
        task = {"id": "R01", "title": "Review", "task": "Inspect the transaction", "category": "review", "author": "codex-1"}
        answer = "---\nagent: claude-1\nmodel: pinned-model\n---\n# Actual answer\nI disagree with codex-1 because the row is not locked.\n"
        bundle, mapping = prepare_blind_bundle(task, [{"arm": "full", "answer": answer,
                                                      "model": "pinned-model", "cost_usd": "12"}],
                                              "a-private-test-salt-longer-than-24", "judge-a", 42,
                                              identity_terms=["codex-1", "claude-1"])
        public = bundle["candidates"][0]
        self.assertNotIn("agent:", public["answer"])
        self.assertIn("I disagree with codex-1", public["answer"])
        self.assertNotIn("arm", public)
        self.assertNotIn("model", public)
        self.assertNotIn("author", bundle["task"])
        self.assertEqual(mapping["candidates"][0]["possible_identity_leaks"], ["codex-1"])
        self.assertEqual(mapping["bundle_sha256"], canonical_hash(bundle))
        self.assertNotEqual(mapping["candidates"][0]["original_sha256"], mapping["candidates"][0]["blinded_body_sha256"])

    def test_invalid_frontmatter_not_silently_deleted(self):
        with self.assertRaises(ValueError): strip_frontmatter("---\nagent: bad\nanswer")
        with self.assertRaises(ValueError): strip_frontmatter("---\n- list\n---\nanswer")
        self.assertEqual(strip_frontmatter("plain answer"), ("plain answer", None))

    def test_grader_cannot_be_task_or_candidate_author(self):
        authors = [{"id": "writer", "requested_model": "model-a", "reported_model": "actual-a"}]
        graders = [{"id": "judge-a", "requested_model": "model-b"}, {"id": "judge-b", "requested_model": "model-c"}]
        self.assertTrue(require_disjoint_graders(graders, authors)["configured_model_disjoint"])
        for replacement in [{"id": "writer", "requested_model": "model-b"},
                            {"id": "judge-a", "requested_model": "actual-a"},
                            {"id": "judge-a", "requested_model": "model-a"}]:
            with self.assertRaises(ValueError): require_disjoint_graders([replacement, graders[1]], authors)
        with self.assertRaises(ValueError): require_disjoint_graders(graders[:1], authors)
        with self.assertRaises(ValueError):
            require_disjoint_graders([{**row, "reported_model": "same-underlying"} for row in graders], authors)

    def test_agreement_scored_only_when_complete(self):
        result = compare_grades(grade(), grade())
        self.assertFalse(result["needs_adjudication"])
        self.assertEqual(result["quality_score"], 75)
        broken = grade()
        del broken["criteria"]["verification"]
        with self.assertRaises(ValueError): compare_grades(grade(), broken)

    def test_contradictions_are_not_averaged(self):
        first = grade()
        for update in [{"passes": False}, {"critical_findings": ["C1"]},
                       {"criteria": grade(value=1)["criteria"]}]:
            result = compare_grades(first, {**grade(), **update})
            self.assertTrue(result["needs_adjudication"])
            self.assertIsNone(result["quality_score"])
            self.assertIsNone(result["passes"])
        first["supported_findings"] = ["F1"]
        second = grade()
        second["unsupported_findings"] = ["F1"]
        self.assertIn("factual-finding-conflict", compare_grades(first, second)["reasons"])

    def test_executable_failure_and_harness_error_differ(self):
        result = compare_grades(grade(), grade(), {"status": "failed"})
        self.assertIn("grader-conflicts-with-executable-result", result["reasons"])
        result = compare_grades(grade(), grade(), {"status": "harness-error"})
        self.assertEqual(result["reasons"], ["acceptance-harness-unavailable"])

    def test_batch_completeness_and_input_hash(self):
        bundle = {"candidates": [{"candidate_id": "a"}, {"candidate_id": "b"}]}
        digest = canonical_hash(bundle)
        self.assertEqual(validate_grade_batch(bundle, [grade("a"), grade("b")], digest), [75, 75])
        with self.assertRaises(ValueError): validate_grade_batch(bundle, [grade("a")], digest)
        with self.assertRaises(ValueError): validate_grade_batch(bundle, [grade("a"), grade("a")], digest)
        with self.assertRaises(ValueError): validate_grade_batch({**bundle, "changed": True}, [grade("a"), grade("b")], digest)

    def test_invalid_scores_and_evidence_rejected(self):
        for value in [True, 2.5, -1, 5, "3"]:
            value_grade = grade()
            value_grade["criteria"]["correctness"] = value
            with self.assertRaises(ValueError): validate_grade(value_grade, "candidate-abcd")
        value_grade = grade()
        value_grade["evidence"] = []
        with self.assertRaises(ValueError): validate_grade(value_grade, "candidate-abcd")
        with self.assertRaises(ValueError): validate_grade({**grade(), "candidate_id": None}, None)


if __name__ == "__main__":
    unittest.main()
