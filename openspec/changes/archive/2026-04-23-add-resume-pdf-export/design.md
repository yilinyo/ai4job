## Context

The current resume detail page loads a resume version from Vuex and displays its content inside a simple HTML container. Export is implemented entirely in the browser by turning `currentVersion.content` into a plain-text Blob and downloading a `.txt` file. The requested change is to let users export that same viewed resume as a PDF without introducing a backend rendering service.

This project is a Vue 2 + TypeScript + Element UI application with mock APIs. Because resume content is already rendered on the client, the lowest-friction implementation is to build the PDF from a dedicated DOM fragment in the browser and keep the first version client-only.

## Goals / Non-Goals

**Goals:**
- Allow users to export the currently viewed resume version as a PDF from the existing detail page.
- Produce a PDF from a controlled resume layout rather than from raw text content.
- Keep the interaction responsive with loading state, success feedback, and failure feedback.
- Preserve a stable file name based on the resume title and version.

**Non-Goals:**
- Building a full WYSIWYG print editor or template system.
- Adding server-side PDF rendering or storage.
- Supporting batch export or exporting multiple resume versions at once.
- Perfectly reproducing arbitrary rich text or Markdown edge cases in the first release.

## Decisions

Use a client-side HTML-to-PDF library.
Rationale: the app already renders resume content in the browser, and a client-only flow avoids backend API changes. A library such as `html2pdf.js` can convert a dedicated DOM container into a downloadable PDF with explicit margins, page size, and file naming.
Alternatives considered:
- `window.print()`: simple, but it depends on browser print dialogs and does not provide a predictable direct-download flow.
- Server-side PDF generation: higher fidelity options exist, but it would add new backend endpoints and deployment complexity that this project does not currently have.

Generate the PDF from a dedicated export layout instead of the on-screen card directly.
Rationale: the current page contains headers, buttons, tags, and management metadata that should not appear in the final resume. A dedicated hidden or off-screen export container lets us include only resume-relevant fields and control page spacing.
Alternatives considered:
- Capture the visible card as-is: faster to wire up, but it would include UI chrome and produce inconsistent output.

Treat resume content as formatted plain text with line preservation in the initial version.
Rationale: current rendering is newline-to-`<br/>` conversion, so the exported PDF should match that behavior first. This keeps implementation aligned with existing content storage and avoids broad editor/rendering changes.
Alternatives considered:
- Introduce a Markdown parser before export: useful later, but it changes rendering semantics and would expand the scope beyond PDF export itself.

Disable repeat export attempts while generation is in progress.
Rationale: PDF generation can take noticeable time in the browser. A guarded loading state prevents duplicate downloads, accidental double clicks, and confusing success messages.
Alternatives considered:
- Leave the button active: simpler, but risks concurrent generation and inconsistent UX.

## Risks / Trade-offs

- [Large resumes may generate slowly in the browser] -> Mitigation: show loading state, keep export DOM minimal, and avoid rendering non-resume UI content.
- [HTML-to-PDF output may differ slightly from on-screen styling] -> Mitigation: use a dedicated export stylesheet with constrained typography, spacing, and page margins.
- [Automatic page breaks may split sections awkwardly] -> Mitigation: add page-break rules for major blocks and validate with long sample resumes.
- [Special characters or unsupported CSS may render inconsistently across browsers] -> Mitigation: keep the export layout simple and test in the primary supported browser before release.

## Migration Plan

No data migration is required. Ship the new frontend dependency and replace the existing text export action with PDF export on the resume version view page. If export quality is unacceptable after release, rollback can revert the UI handler to the previous `.txt` export implementation because there are no persisted data changes.

## Open Questions

- Should the initial PDF include version metadata such as creation time and active status, or only the resume title and body content?
- Do we want to keep a secondary text export option later, or is PDF the only supported download format going forward?
