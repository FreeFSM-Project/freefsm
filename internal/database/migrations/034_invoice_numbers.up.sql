ALTER TABLE invoices ADD COLUMN IF NOT EXISTS invoice_number BIGINT;

UPDATE invoices
SET invoice_number = id
WHERE invoice_number IS NULL;

ALTER TABLE invoices ALTER COLUMN invoice_number SET NOT NULL;

ALTER TABLE company_settings ADD COLUMN IF NOT EXISTS next_invoice_number BIGINT NOT NULL DEFAULT 1;

UPDATE company_settings
SET next_invoice_number = GREATEST(
    next_invoice_number,
    COALESCE((SELECT MAX(invoice_number) + 1 FROM invoices WHERE company_settings.company_id IS NOT NULL AND invoices.company_id = company_settings.company_id), 1),
    COALESCE((SELECT MAX(invoice_number) + 1 FROM invoices WHERE company_settings.company_id IS NULL AND invoices.company_id IS NULL), 1)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_invoices_company_invoice_number
    ON invoices(company_id, invoice_number)
    WHERE company_id IS NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_invoices_invoice_number_no_company
    ON invoices(invoice_number)
    WHERE company_id IS NULL;
