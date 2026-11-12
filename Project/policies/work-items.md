# Work Item policy

* **WP type:** Policy
* **Version:** 0.1.0
* **Date:** 2026-11-12
* **Status:** Working agreement

## Definition

A Work Item (WI) is a unit of work or an open question tracked in GitHub
Issues. It describes the goal, scope, and expected result. The resulting
specification, code, data, or report is a separate Work Product.

## Work Item types

| Type | Use | Suggested label |
|---|---|---|
| `Research` | Investigate a question or evidence base | `type:research` |
| `Experiment` | Test a stated hypothesis using a recorded method | `type:experiment` |
| `Task` | Produce or change a concrete artifact | `type:task` |
| `Finding` | Record an observed defect, surprising result, or limitation | `type:finding` |
| `Change Request` | Request a decision to change an accepted baseline or scope | `type:change-request` |
| `Risk` | Track a possible adverse event and response | `type:risk` |

The matching GitHub forms are in `.github/ISSUE_TEMPLATE/`. These types are
project conventions, not GitHub or standards-defined classifications.

## Good Work Item boundaries

- State one principal outcome and what is out of scope.
- For an experiment, state a falsifiable hypothesis, method, and success and
  failure criteria before running it.
- For a task, state observable acceptance criteria and link the relevant WP.
- Use a Change Request when the owner must decide whether to change an
  accepted requirement, interface, architecture, or scope. Routine edits
  within an existing agreement can be Tasks.
- Create parent/child Issues only when the child has an independently useful
  outcome. Small steps can remain a checklist.

## State and closure

GitHub Issue state records whether the WI is open or closed. Labels such as
`blocked`, `accepted`, and priority labels are separate project metadata;
do not confuse them with WP approval or experiment success.

Close a WI with a link to the resulting WP or explain why no artifact was
needed. A negative or inconclusive result can complete an Experiment or
Research WI when its method and evidence are recorded.

## Version history

| Version | Date | Change |
|---|---|---|
| 0.1.0 | 2026-11-12 | Initial WI classes and lifecycle. |
