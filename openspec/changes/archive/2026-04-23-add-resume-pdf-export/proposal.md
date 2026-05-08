## Why

The resume version detail page currently exports raw text as a `.txt` file, which is not suitable for formal resume delivery. Adding direct PDF export lets users download a shareable resume from the system without manually reformatting content in another tool.

## What Changes

- Replace the current text-file export flow on the resume version detail page with PDF export.
- Render the current resume version into a print-friendly layout before generating the file.
- Show export progress and clear failure feedback so users know whether the PDF was created successfully.
- Keep the exported file name aligned with the current resume title and version number.

## Capabilities

### New Capabilities
- `resume-pdf-export`: Export a viewed resume version as a downloadable PDF with readable formatting and predictable file naming.

### Modified Capabilities

## Impact

- Affected UI: `src/views/resume-version/view.vue`
- Likely affected frontend support code: shared export/render utilities and dependency management in `package.json`
- New client-side PDF generation dependency is expected
- No backend or mock API contract changes are required for the initial version
