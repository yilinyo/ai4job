# resume-pdf-export Specification

## Purpose
Define how users export resume versions as readable PDF files from the resume version detail flow.

## Requirements

### Requirement: User can export the viewed resume as a PDF
The system SHALL allow a user on the resume version detail page to download the currently viewed resume version as a PDF file.

#### Scenario: Export current resume version
- **WHEN** the user opens a resume version detail page and clicks the export action
- **THEN** the system downloads a PDF generated from the currently loaded resume version

### Requirement: Exported PDF uses a resume-focused layout
The system SHALL generate the PDF from a layout that excludes management UI controls and includes resume content in a readable document format.

#### Scenario: Exclude page controls from PDF
- **WHEN** the system generates a PDF for a resume version
- **THEN** action buttons, management tags, and navigation controls are not included in the exported document

#### Scenario: Preserve readable resume content
- **WHEN** the system generates a PDF for resume content containing line breaks
- **THEN** the exported document preserves the content order and visible line separation in the output

### Requirement: Export interaction communicates status clearly
The system SHALL provide user feedback during and after PDF generation.

#### Scenario: Prevent duplicate export during generation
- **WHEN** a PDF export is already in progress
- **THEN** the export action is disabled or ignored until the current generation finishes

#### Scenario: Report export success
- **WHEN** the PDF file is generated successfully
- **THEN** the system informs the user that the export completed successfully

#### Scenario: Report export failure
- **WHEN** PDF generation fails for any reason
- **THEN** the system informs the user that the export failed and does not report success

### Requirement: Exported PDF file naming is predictable
The system SHALL name the exported PDF file using the current resume title and version number so the downloaded file is identifiable.

#### Scenario: Name file from title and version
- **WHEN** the user exports a resume version titled "Frontend Resume" with version "v1.2.0"
- **THEN** the downloaded file name includes "Frontend Resume" and "v1.2.0" and uses the `.pdf` extension
