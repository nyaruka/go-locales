v1.5.0 (2026-08-31)
-------------------------
 * Reject malformed FDCC input previously accepted silently - NUL chars, invalid UTF-8, bad escapes and mismatched END directives
 * Report an error rather than overflowing the stack when locale categories copy from each other in a cycle
 * Return copies from Query and Codes so callers can't modify the shared database
 * Correct LC_XLITERATE which named LC_MESSAGES and so returned data for the wrong category
 * Load locale data on first use rather than at init
 * Update supported Go versions

v1.4.0 (2025-04-25)
-------------------------
 * Update locale data
 * Update dependencies and go version

v1.3.0 (2024-05-31)
-------------------------
 * Update metadata from latest glibc
 * Test on multiple go versions
 * Use std lib errors

