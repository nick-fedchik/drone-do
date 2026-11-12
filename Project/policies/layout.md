# Repository layout policy

* **WP type:** Policy
* **Version:** 0.4.0
* **Date:** 2026-11-12
* **Status:** Working agreement

## Purpose

Keep research evidence, engineering artifacts, and project administration
separate while allowing links between them. A folder indicates the primary
work area; it does not determine the Work Product type.

## Top-level areas

| Path | Owns |
|---|---|
| `Research/` | Research questions, source material, knowledge-generating experiments, analysis reports, and concept demonstrations |
| `Engineering/` | Requirements, specifications, technical decisions, design, software, hardware descriptions, and verification evidence |
| `Project/` | Project rules, templates, registers, plans, and project-scope or governance decisions |
| `.github/` | GitHub-specific Issue and pull request forms |

Place an artifact where its primary purpose belongs and link to it from
other relevant indexes. Avoid duplicate authoritative copies. A research
result that becomes an engineering baseline should be linked and explicitly
reviewed before adoption.

## Experiments and decisions

- Record experiments intended to answer open research questions under
  `Research/experiments/<topic>/<experiment-id>/`. Keep the method, conditions,
  observations, and interpretation traceable; do not treat an exploratory
  result as verification of a requirement.
- Record verification against requirements under
  `Engineering/validation/`, linking the exact requirement and evidence.
- Put conceptual renders and procedural demonstrations under
  `Research/demonstrations/`. Label them as illustrative; they are not
  evidence of physical performance or validation.
- Put technical decisions (architecture, interfaces, protocols, or component
  choices) under `Engineering/decisions/`. Put project-scope, governance, and
  resource-priority decisions under `Project/decisions/`.
- Place each independently built software module or component in its own
  folder under `Engineering/software/`, regardless of language. Keep its
  source, tests, and build manifests together; language-specific module roots
  belong inside that folder.
- Create these directories when the first corresponding Work Product exists;
  do not add empty placeholder trees.

## Naming and metadata

- Use descriptive, stable, lowercase kebab-case names for new files when
  practical. Preserve an imported source filename when changing it would
  weaken traceability.
- Markdown Work Products identify their type, version, date, and status near
  the title. Record origin and migration state for imported material.
- Do not imply approval, implementation, or validation merely by placing a
  document in `Engineering/`.
- Keep original sources and raw experiment data unchanged; describe them in
  an adjacent index or provenance record.
- Keep media, generated output, credentials, and large raw data out of
  documentation-only migrations unless they are explicitly in scope.

## Version history

| Version | Date | Change |
|---|---|---|
| 0.1.0 | 2026-11-12 | Initial three-area structure. |
| 0.2.0 | 2026-11-12 | Clarify experiment, verification, and decision-record locations. |
| 0.3.0 | 2026-11-12 | Specify separate, language-independent software module folders. |
| 0.4.0 | 2026-11-12 | Add a home and evidence boundary for concept demonstrations. |
