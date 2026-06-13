ALTER TABLE company_settings ADD COLUMN smtp_host TEXT NOT NULL DEFAULT '';
ALTER TABLE company_settings ADD COLUMN smtp_port INT NOT NULL DEFAULT 587;
ALTER TABLE company_settings ADD COLUMN smtp_user TEXT NOT NULL DEFAULT '';
ALTER TABLE company_settings ADD COLUMN smtp_password TEXT NOT NULL DEFAULT '';
ALTER TABLE company_settings ADD COLUMN smtp_from TEXT NOT NULL DEFAULT '';
