ALTER TABLE polls ADD  CONSTRAINT polls_expires_at_check
CHECK (
  expires_at >= now() + interval '1 hour'
);
ALTER TABLE polls ADD CONSTRAINT polls_options_length_check CHECK (cardinality(options) >= 2)