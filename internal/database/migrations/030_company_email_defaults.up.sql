ALTER TABLE company_settings
    ADD COLUMN IF NOT EXISTS invoice_email_subject TEXT DEFAULT 'Invoice {invoice_number} from {business_name}',
    ADD COLUMN IF NOT EXISTS invoice_email_body TEXT DEFAULT 'Hello {customer_name},

Please find invoice {invoice_number} attached.

Thank you,
{business_name}',
    ADD COLUMN IF NOT EXISTS estimate_email_subject TEXT DEFAULT 'Estimate {estimate_number} from {business_name}',
    ADD COLUMN IF NOT EXISTS estimate_email_body TEXT DEFAULT 'Hello {customer_name},

Please find estimate {estimate_number} attached.

Thank you,
{business_name}';
