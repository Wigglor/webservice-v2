ALTER TABLE organizations
ADD COLUMN subscription_id VARCHAR(255),
ADD COLUMN plan_type VARCHAR(50), -- e.g., 'free', 'standard', 'enterprise'
ADD COLUMN subscription_status VARCHAR(50), -- e.g., 'active', 'canceled'
ADD COLUMN next_billing_date TIMESTAMPTZ;