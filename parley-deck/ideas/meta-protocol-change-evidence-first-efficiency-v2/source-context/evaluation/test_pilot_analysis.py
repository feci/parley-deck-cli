import collections
import itertools
import unittest

from pilot_analysis import allocation, blind_identifier, canonical_hash, summarize, validate_cells

TASKS = [{"id": "T%02d" % index} for index in range(12)]
ROSTER = ["participant-%d" % index for index in range(6)]


class PilotAnalysisTests(unittest.TestCase):
    def test_allocation_is_balanced_and_reproducible(self):
        result = allocation(TASKS, ROSTER)
        self.assertEqual(result, allocation(TASKS, ROSTER))
        self.assertNotEqual(result, allocation(TASKS, ROSTER, seed=9))
        orders, solos, duos, drafters, pairs, positions = [collections.Counter() for _ in range(6)]
        for task in result["tasks"]:
            orders[tuple(row["arm"] for row in task["arms"])] += 1
            for row in task["arms"]:
                if row["arm"] == "solo": solos.update(row["participants"])
                if row["arm"] == "duo":
                    duos.update(row["participants"])
                    drafters[row["drafter"]] += 1
                    pairs[tuple(sorted(row["participants"]))] += 1
                if row["arm"] == "full":
                    self.assertEqual(set(row["participants"]), set(ROSTER))
                    positions.update((person, index) for index, person in enumerate(row["participants"]))
        self.assertEqual(set(orders), set(itertools.permutations(("solo", "duo", "full"))))
        self.assertEqual(set(orders.values()), {2})
        self.assertEqual(set(solos.values()), {2})
        self.assertEqual(set(duos.values()), {4})
        self.assertEqual(set(drafters.values()), {2})
        self.assertEqual(len(pairs), 12)
        self.assertEqual(set(pairs.values()), {1})
        self.assertEqual(set(positions.values()), {2})

    def test_full_roster_cannot_shrink(self):
        with self.assertRaises(ValueError): allocation(TASKS, ROSTER[:4])
        with self.assertRaises(ValueError): allocation(TASKS, ROSTER[:5] + [ROSTER[0]])

    def test_empty_results_are_unknown_not_free_or_failed(self):
        result = summarize(TASKS, [])
        self.assertEqual(result["planned_cells"], 36)
        for arm in result["arms"]:
            self.assertEqual(arm["attempted"], 0)
            self.assertIsNone(arm["quality_mean"])
            self.assertIsNone(arm["known_cost_subtotal_usd"])
            self.assertEqual(arm["planned_completion_bounds"], [0, 1])
        self.assertEqual(result["paired"][0]["n"], 0)

    def test_paired_not_unpaired_and_failures_retained(self):
        cells = [
            {"task": "T00", "arm": "solo", "status": "completed", "quality_score": 60, "cost_usd": "0.1"},
            {"task": "T00", "arm": "duo", "status": "completed", "quality_score": 80, "cost_usd": "0.2"},
            {"task": "T01", "arm": "solo", "status": "timeout", "cost_usd": None},
            {"task": "T01", "arm": "duo", "status": "completed", "quality_score": 90},
            {"task": "T02", "arm": "solo", "status": "completed", "quality_score": 100},
        ]
        result = summarize(TASKS, cells, samples=100)
        self.assertEqual(result["paired"][0]["n"], 2)
        self.assertEqual(result["paired"][0]["mean_delta"], 55)
        self.assertEqual(result["arms"][0]["quality_mean"], 160/3)
        self.assertEqual(result["arms"][0]["cost_unknown_n"], 2)
        self.assertEqual(result["arms"][0]["known_cost_subtotal_usd"], "0.1")
        self.assertEqual(result, summarize(TASKS, cells, samples=100))

    def test_ungraded_is_not_imputed(self):
        cells = [{"task": "T00", "arm": "solo", "status": "ungraded", "quality_score": None}]
        result = summarize(TASKS, cells)
        self.assertEqual(result["arms"][0]["quality_n"], 0)
        self.assertEqual(result["arms"][0]["attempted"], 1)
        self.assertEqual(result["arms"][0]["planned_completion_bounds"], [0, 1])

    def test_duplicate_events_and_bad_cells_rejected(self):
        base = {"task": "T00", "arm": "solo", "status": "completed", "quality_score": 80, "invocation_ids": ["unique-id"]}
        with self.assertRaises(ValueError): validate_cells(TASKS, [base, base])
        second = {**base, "arm": "duo"}
        with self.assertRaises(ValueError): validate_cells(TASKS, [base, second])
        for update in [{"quality_score": True}, {"quality_score": float("nan")}, {"cost_usd": "NaN"}, {"cost_usd": "bad"},
                       {"elapsed_seconds": -1}, {"status": "not-run"}, {"arm": "full-minus-one"}]:
            with self.subTest(update=update), self.assertRaises(ValueError):
                validate_cells(TASKS, [{**base, **update}])

    def test_blind_ids_do_not_embed_treatment(self):
        salt = "private-fixture-salt-not-for-live-use"
        identities = [blind_identifier("T00", arm, salt) for arm in ("solo", "duo", "full")]
        self.assertEqual(len(set(identities)), 3)
        self.assertTrue(all(value.startswith("candidate-") for value in identities))
        with self.assertRaises(ValueError): blind_identifier("T00", "solo", "public")

    def test_canonical_hash_stable_across_key_order(self):
        self.assertEqual(canonical_hash({"x": 1, "y": 2}), canonical_hash({"y": 2, "x": 1}))


if __name__ == "__main__":
    unittest.main()
