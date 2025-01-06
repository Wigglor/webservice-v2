ALTER TABLE users
  ADD CONSTRAINT users_email_uq UNIQUE (email);

ALTER TABLE users
  ADD CONSTRAINT users_subid_uq UNIQUE (sub_id);
