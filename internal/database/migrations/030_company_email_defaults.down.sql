ALTER TABLE company_settings
    DROP COLUMN IF EXISTS estimate_email_body,
    DROP COLUMN IF EXISTS estimate_email_subject,
    DROP COLUMN IF EXISTS invoice_email_body,
    DROP COLUMN IF EXISTS invoice_email_subject;
