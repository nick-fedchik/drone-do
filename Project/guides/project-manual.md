# Project manual

* **WP type:** Guide
* **Version:** 0.3.0
* **Date:** 2026-11-12
* **Status:** Working guide

## 1. Start work

1. Decide whether the need is a research question, experiment, task,
   finding, change request, or risk. Use the corresponding GitHub Issue form.
2. Link the Issue to existing requirements, specifications, source material,
   or other relevant Work Products.
3. If an experiment is involved, agree on the method and pass/fail criteria
   before collecting data. Use `Research/experiments/` for exploratory
   experiments and `Engineering/validation/` to verify requirements.

## 2. Create the result

Choose the main area according to purpose:

- Research questions, original sources, analytical results, and
  knowledge-generating experiments go in `Research/`.
- Conceptual renders and procedural demonstrations go in
  `Research/demonstrations/`; label them as illustrative, not validation.
- System/technical requirements, specifications, designs, code, and
  validation outputs go in `Engineering/`.
- Technical decisions go in `Engineering/decisions/`. Project-scope,
  governance, and resource-priority decisions go in `Project/decisions/`.
- Project rules, templates, registers, and indexes go in `Project/`.

Create a decision or experiment directory when its first Work Product is
ready; avoid empty placeholder directories. Link decisions to affected
requirements, specifications, Issues, and experiment evidence.

Use the appropriate [Work Product template](../templates/README.md), add
the minimum metadata, and link related artifacts rather than duplicating
their contents.

## 3. Preserve evidence

Keep raw observations separate from interpretation. Record equipment and
revision, software, conditions, method, and source revision where known.
Mark unknowns and limitations; do not convert an illustrative demo or a
single trial into a verified performance claim.

For imported documents, preserve source attribution and state whether claims
were rechecked. Migration and formatting are not technical review.

## 4. Work Product status

- `draft`: content is being prepared and has not been reviewed.
- `proposed`: ready for a named owner or decision process to consider; not
  yet the project baseline.
- `accepted`: an authorized owner has adopted this exact version as a
  reference or baseline. Record the approver and basis in the related Issue
  or Decision.
- `verified`: the exact version has passed stated verification criteria,
  with evidence linked. This does not imply general flight safety,
  certification, or validation beyond those criteria.
- `archived`: retained for traceability but superseded or no longer active.

Acceptance and verification are different claims. An accepted design can
remain unverified; a verified test result does not itself authorize a design
or operational baseline.

## 5. Review and close

Review changes for scope, traceability, broken links, and consistency with
current project decisions. Close the WI only after its result is linked or
the reason no WP was produced is recorded. Keep the owner’s approval separate
from an author's statement that work is complete.

## 6. Version history

| Version | Date | Change |
|---|---|---|
| 0.1.0 | 2026-11-12 | Initial project workflow. |
| 0.2.0 | 2026-11-12 | Add experiment/decision placement and clarify WP status meanings. |
| 0.3.0 | 2026-11-12 | Define placement and limitations for concept demonstrations. |
