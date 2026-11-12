# Work Product policy

* **WP type:** Policy
* **Version:** 0.3.0
* **Date:** 2026-11-12
* **Status:** Working agreement

## Definition

A Work Product (WP) is an artifact with a defined purpose: a question,
specification, design, report, data set, source file, or other project
result. It is not the GitHub Issue used to track work on that artifact.

Choose the primary type by the artifact's purpose, not by its extension or
filename. One document may contain several related requirements; not every
requirement needs its own file.

## Project vocabulary

| WP type | Use |
|---|---|
| `Problem` | Research question, context, hypotheses, and scope |
| `Source` | Original or external material retained for traceability |
| `Protocol` | Planned method and conditions for a test or experiment |
| `Report` | Findings and limitations from completed work; subtypes include Analysis, Experiment, Review, and Assessment |
| `Demonstration` | Illustrative render, simulation, or procedural demo; not by itself evidence that a physical system works |
| `Data` | Raw measurements, logs, images, or other captured evidence |
| `Requirement` | A verifiable need or constraint |
| `Specification` | Organized requirements, interfaces, behavior, or constraints |
| `Design` | Proposed or implemented structure and component interactions |
| `Decision` | Choice, rationale, alternatives, and consequences |
| `Solution` | A bounded set of components that addresses a problem |
| `Source Code` | Executable source, tests, and tools |
| `Configuration` | Machine-readable build, deployment, or runtime settings |
| `Hardware` | Description and identity of a physical prototype or device |
| `Plan` | Intended stages and activities |
| `Guide` | Navigation or instructions for readers and contributors |
| `Policy` | Project rule or working agreement |
| `Register` | Maintained list of entities, status, provenance, or relationships |
| `Template` | Reusable form for creating another WP or WI |

## Minimum document metadata

For a Markdown WP, use the applicable fields from this set near the title:

```markdown
* **WP type:** [type / subtype]
* **Version:** [version]
* **Date:** YYYY-MM-DD
* **Status:** [draft / proposed / accepted / verified / archived]
* **Origin:** [source or "new"]
```

Not every status applies to every artifact. State the actual status and
identify the basis for claims of acceptance or verification. For imported
material, add its source revision and say whether content was technically
revalidated. Never fabricate missing authorship, dates, measurements, or
approval.

Use the status meanings in the [Project manual](../guides/project-manual.md).
In particular, `accepted` records an authorized decision about a specific
version; `verified` records evidence against stated criteria. Neither status
implies a broader safety claim.

## Evidence and source handling

- Link claims to primary sources, experiment data, or the relevant Issue.
- Keep raw data and original source documents unchanged; use a sidecar or
  register for metadata when needed.
- Distinguish observation, inference, hypothesis, and decision.
- A report, demo, or successful single trial is not by itself validation of
  the system.
- Mark unverified interfaces and performance figures explicitly.

Use the [Work Product templates](../templates/work-products/README.md) when
they fit. See the [register](../registers/work-products.md) for the current
inventory.

## Version history

| Version | Date | Change |
|---|---|---|
| 0.1.0 | 2026-11-12 | Initial WP types and metadata rules. |
| 0.2.0 | 2026-11-12 | Link status meanings to the project workflow. |
| 0.3.0 | 2026-11-12 | Add an illustrative Demonstration Work Product type. |
