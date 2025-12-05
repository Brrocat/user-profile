-- Create user_profile table
CREATE TABLE IF NOT EXISTS user_profiles
(
    id              UUID PRIMARY KEY         DEFAULT gen_random_uuid(),
    user_id         UUID         NOT NULL UNIQUE,
    first_name      VARCHAR(100) NOT NULL,
    last_name       VARCHAR(100) NOT NULL,
    phone           VARCHAR(20),
    date_of_birth   DATE,
    avatar_url      TEXT,
    address         TEXT,
    city            VARCHAR(100),
    country         VARCHAR(100),
    postal_code     VARCHAR(20),
    driving_licence VARCHAR(50),
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create index on user_id for faster lookup
CREATE INDEX IF NOT EXISTS idx_user_profile_user_id ON user_profiles(user_id);

-- Create index on email-related fields for search
CREATE INDEX IF NOT EXISTS idx_user_profile_name ON user_profiles(first_name, last_name);

-- Create updated_at trigger
CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.update_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE OR REPLACE TRIGGER update_user_profiles_updated_at
    BEFORE UPDATE ON user_profiles
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
