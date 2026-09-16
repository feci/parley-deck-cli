#!/usr/bin/env python3
"""Counterbalanced allocation and honest paired summaries for the pilot.

Allocation does not freeze a study by itself. Missing measurements remain null;
only observed terminal failures receive the preregistered zero quality score.
"""
import hashlib
import itertools
import json
import math
import random
import statistics
from collections import Counter
from decimal import Decimal, InvalidOperation, localcontext

ARMS = ("solo", "duo", "full")
FAILURES = {"provider-failure", "timeout", "process-failure", "protocol-block", "missing-output"}
STATES = FAILURES | {"not-run", "running", "completed", "ungraded"}


def allocation(tasks, roster, seed=20260905):
    ids = [task["id"] for task in tasks]
    if len(ids) != 12 or len(set(ids)) != 12:
        raise ValueError("Exactly twelve unique task IDs are required")
    if len(roster) != 6 or len(set(roster)) != 6:
        raise ValueError("The full arm requires six distinct frozen roster IDs")
    rng = random.Random(seed)
    people = list(roster)
    rng.shuffle(people)
    order = list(ids)
    rng.shuffle(order)
    solos = people * 2
    rng.shuffle(solos)
    # Two directed six-cycles give every participant four pair appearances and
    # two duo drafter assignments without repeating an unordered pair.
    pairs = [(people[i], people[(i + step) % 6]) for step in (1, 2) for i in range(6)]
    rng.shuffle(pairs)
    treatments = list(itertools.permutations(ARMS)) * 2
    rng.shuffle(treatments)
    rows = []
    for index, task_id in enumerate(order):
        rotated = people[index % 6:] + people[:index % 6]
        duet = list(pairs[index])
        assignments = {
            "solo": {"participants": [solos[index]], "drafter": solos[index], "critic": solos[index]},
            "duo": {"participants": duet, "drafter": duet[0], "critic": duet[1]},
            "full": {"participants": rotated, "drafter": rotated[0], "critic": rotated[1]},
        }
        rows.append({"task": task_id, "execution_index": index,
                     "arms": [{"arm": arm, **assignments[arm]} for arm in treatments[index]]})
    return {"schema": 1, "status": "draft-allocation-not-preregistered", "seed": seed,
            "roster": list(roster), "tasks": rows,
            "limits": "Balanced overall, not a fully crossed model-by-category factorial design."}


def validate_cells(tasks, cells):
    planned = {(task["id"], arm) for task in tasks for arm in ARMS}
    seen, invocations = set(), set()
    for cell in cells:
        key = (cell.get("task"), cell.get("arm"))
        if key not in planned or key in seen:
            raise ValueError("Unplanned or duplicate task/arm cell")
        seen.add(key)
        if cell.get("status") not in STATES:
            raise ValueError("Unknown cell status")
        score = cell.get("quality_score")
        if score is not None and (isinstance(score, bool) or not isinstance(score, (int, float)) or
                                  not math.isfinite(score) or not 0 <= score <= 100):
            raise ValueError("Quality score must be finite and between 0 and 100")
        for identity in cell.get("invocation_ids", []):
            if not isinstance(identity, str) or not identity or identity in invocations:
                raise ValueError("Invocation ID is missing, duplicated or assigned to two arms")
            invocations.add(identity)
        if cell["status"] == "not-run" and (cell.get("invocation_ids") or score is not None):
            raise ValueError("Not-run cell cannot contain treatment invocations or a score")
        cost = cell.get("cost_usd")
        if cost is not None:
            if not isinstance(cost, str):
                raise ValueError("Cost must be a decimal string or null")
            try:
                amount = Decimal(cost)
            except InvalidOperation as error:
                raise ValueError("Cost is not a decimal") from error
            if not amount.is_finite() or amount < 0:
                raise ValueError("Cost must be finite and nonnegative")
        elapsed = cell.get("elapsed_seconds")
        if elapsed is not None and (isinstance(elapsed, bool) or not isinstance(elapsed, (int, float)) or
                                    not math.isfinite(elapsed) or elapsed < 0):
            raise ValueError("Elapsed time must be finite and nonnegative")
    return planned


def quality(cell):
    if not cell or cell["status"] in ("not-run", "running", "ungraded"):
        return None
    if cell["status"] in FAILURES:
        return 0
    return cell.get("quality_score")


def percentile(values, proportion):
    ordered = sorted(values)
    position = (len(ordered) - 1) * proportion
    low = int(position)
    high = min(low + 1, len(ordered) - 1)
    return ordered[low] + (ordered[high] - ordered[low]) * (position - low)


def paired_delta(tasks, by_key, arm, reference="solo", seed=20260905, samples=10000):
    pairs = []
    missing = []
    for task in tasks:
        left, right = quality(by_key.get((task["id"], arm))), quality(by_key.get((task["id"], reference)))
        if left is None or right is None:
            missing.append(task["id"])
        else:
            pairs.append({"task": task["id"], "delta": left-right})
    deltas = [pair["delta"] for pair in pairs]
    result = {"arm": arm, "reference": reference, "pairs": pairs, "n": len(pairs),
              "missing_tasks": missing, "mean_delta": None, "median_delta": None,
              "bootstrap_95_percentile": None, "leave_one_out_mean_range": None}
    if not deltas:
        return result
    result.update(mean_delta=statistics.mean(deltas), median_delta=statistics.median(deltas))
    if len(deltas) > 1:
        rng = random.Random(seed)
        means = [statistics.mean(rng.choices(deltas, k=len(deltas))) for _ in range(samples)]
        result["bootstrap_95_percentile"] = [percentile(means, .025), percentile(means, .975)]
        leave_one_out = [statistics.mean(deltas[:i] + deltas[i+1:]) for i in range(len(deltas))]
        result["leave_one_out_mean_range"] = [min(leave_one_out), max(leave_one_out)]
    return result


def summarize(tasks, cells, seed=20260905, samples=10000):
    planned = validate_cells(tasks, cells)
    by_key = {(cell["task"], cell["arm"]): cell for cell in cells}
    arm_results = []
    for arm in ARMS:
        rows = [by_key.get((task["id"], arm), {"task": task["id"], "arm": arm, "status": "not-run"}) for task in tasks]
        attempted = [row for row in rows if row["status"] != "not-run"]
        terminal = [row for row in attempted if row["status"] != "running"]
        completed = [row for row in rows if row["status"] == "completed"]
        scores = [quality(row) for row in rows if quality(row) is not None]
        costs = [row["cost_usd"] for row in terminal if row.get("cost_usd") is not None]
        elapsed = [row["elapsed_seconds"] for row in terminal if row.get("elapsed_seconds") is not None]
        completion_unknown = sum(row["status"] in ("not-run", "running", "ungraded") for row in rows)
        with localcontext() as context:
            context.prec = 80
            subtotal = format(sum((Decimal(value) for value in costs), Decimal(0)), "f") if costs else None
        arm_results.append({"arm": arm, "planned": len(rows), "attempted": len(attempted),
                            "terminal": len(terminal), "completed": len(completed),
                            "statuses": dict(Counter(row["status"] for row in rows)),
                            "quality_n": len(scores), "quality_mean": statistics.mean(scores) if scores else None,
                            "known_cost_subtotal_usd": subtotal, "cost_known_n": len(costs),
                            "cost_unknown_n": len(terminal)-len(costs), "cost_basis": "reported, not invoiced",
                            "median_elapsed_seconds": statistics.median(elapsed) if elapsed else None,
                            "elapsed_n": len(elapsed),
                            "completion_rate_attempted": len(completed)/len(attempted) if attempted else None,
                            "planned_completion_bounds": [len(completed)/len(rows),
                                (len(completed)+completion_unknown)/len(rows)] if rows else None})
    return {"schema": 1, "planned_cells": len(planned), "recorded_cells": len(cells),
            "arms": arm_results, "paired": [paired_delta(tasks, by_key, arm, seed=seed, samples=samples) for arm in ("duo", "full")],
            "uncertainty": {"method": "paired task-level percentile bootstrap", "seed": seed, "samples": samples,
                            "caveats": "Small curated sample; no population or causal historical claim. Missing grades are not zero. Terminal failed arms score zero by design; not-run cells do not."}}


def blind_identifier(task, arm, salt):
    if not salt or len(salt) < 24:
        raise ValueError("Use an independently generated private salt, not a public seed")
    digest = hashlib.sha256((salt + "\0" + task + "\0" + arm).encode()).hexdigest()
    return "candidate-" + digest[:16]


def canonical_hash(value):
    return hashlib.sha256(json.dumps(value, sort_keys=True, separators=(",", ":"), ensure_ascii=True).encode()).hexdigest()
