DROP INDEX IF EXISTS idx_invoices_invoice_number_no_company;
DROP INDEX IF EXISTS idx_invoices_company_invoice_number;
ALTER TABLE company_settings DROP COLUMN IF EXISTS next_invoice_number;
ALTER TABLE invoices DROP COLUMN IF EXISTS invoice_number;
