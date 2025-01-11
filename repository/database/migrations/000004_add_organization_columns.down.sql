ALTER TABLE organizations
  DROP COLUMN IF EXISTS subscription_id,
  DROP COLUMN IF EXISTS plan_type,
  DROP COLUMN IF EXISTS subscription_status,
  DROP COLUMN IF EXISTS next_billing_date;
