## 1. Export foundation

- [x] 1.1 Add and register the client-side PDF export dependency required for browser-based resume download
- [x] 1.2 Define a dedicated resume PDF export layout and style rules for readable page output

## 2. Resume detail page integration

- [x] 2.1 Replace the current text export handler in `src/views/resume-version/view.vue` with PDF generation from the viewed resume data
- [x] 2.2 Add export loading guards, success feedback, and failure handling to prevent duplicate actions and surface errors

## 3. Validation

- [x] 3.1 Verify exported file naming uses resume title plus version with a `.pdf` extension
- [x] 3.2 Validate short and long resume samples to confirm layout readability and acceptable page breaks
