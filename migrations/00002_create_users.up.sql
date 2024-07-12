CREATE SCHEMA IF NOT EXISTS auth;

CREATE EXTENSION IF NOT EXISTS citext;

CREATE OR REPLACE FUNCTION auth.set_current_timestamp_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$;

CREATE TABLE IF NOT EXISTS auth.users (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    created_at timestamp DEFAULT now() NOT NULL,
    updated_at timestamp DEFAULT now() NOT NULL,
    name varchar(50) NOT NULL,
    email varchar(40) NOT NULL,
    phone varchar(25),
    password varchar(90),
    profile_picture varchar(255),
    email_confirmed boolean DEFAULT false,
    email_confirmed_at timestamp,
    email_confirmation_token varchar(255),
    email_confirmation_token_sent_at timestamp DEFAULT now(),
    phone_confirmed boolean DEFAULT false,
    phone_confirmed_at timestamp,
    phone_confirmation_token_sent_at timestamp DEFAULT now(),
    phone_confirmation_token varchar(255),
    recovery_token varchar(255),
    recovery_token_sent_at timestamp,
    email_change_token varchar(255),
    email_change_token_sent_at timestamp,
    phone_change_token varchar(255),
    phone_change_token_sent_at timestamp,
    last_login_at timestamp,
    incorrect_login_attempts integer DEFAULT 0,
    last_incorrect_login_attempt_at timestamp,
    CONSTRAINT check_email_confirmation CHECK (((email_confirmed = false) OR (email_confirmed_at IS NOT NULL))),
    CONSTRAINT check_phone_confirmation CHECK (((phone_confirmed = false) OR (phone_confirmed_at IS NOT NULL))),
    CONSTRAINT check_phone_or_email_are_provided CHECK (((email IS NOT NULL) OR (phone IS NOT NULL))),
    CONSTRAINT users_email_key UNIQUE (email),
    CONSTRAINT users_phone_key UNIQUE (phone),
    CONSTRAINT users_pkey PRIMARY KEY (id)
);

CREATE OR REPLACE TRIGGER set_auth_users_updated_at BEFORE UPDATE ON auth.users FOR EACH ROW EXECUTE FUNCTION auth.set_current_timestamp_updated_at();
