-- 2. New organizations table
CREATE TABLE organizations (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL UNIQUE,
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now()
);

-- 3. New user_organizations join table
--    This links many users to many organizations.
CREATE TABLE user_organizations (
  user_id INT NOT NULL,
  organization_id INT NOT NULL,

  -- Optional: store a "role" or membership type (e.g., "admin", "member", etc.)
  role VARCHAR(50),

  created_at TIMESTAMPTZ DEFAULT now(),

  -- Primary key ensures a user can't join the same org more than once.
  PRIMARY KEY (user_id, organization_id),

  -- Foreign key constraints
  FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
  FOREIGN KEY (organization_id) REFERENCES organizations (id) ON DELETE CASCADE
);