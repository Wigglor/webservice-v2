ALTER TABLE user_organizations
ADD COLUMN updated_at TIMESTAMPTZ DEFAULT now();