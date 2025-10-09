-- Migration: Add avatar_url and email_verified columns to users table
-- This supports OAuth profile pictures and email verification tracking

-- Add avatar_url column for storing profile picture URLs (from OAuth or manual upload)
ALTER TABLE users ADD COLUMN avatar_url TEXT DEFAULT '';

-- Add email_verified column to track if user's email has been confirmed
-- OAuth users get this set to TRUE automatically since providers verify emails
-- Local account users will need email verification flow (future enhancement)
ALTER TABLE users ADD COLUMN email_verified BOOLEAN DEFAULT FALSE;

-- Update existing OAuth users to have email_verified = TRUE
-- Since they authenticated via OAuth, their emails are already verified by the provider
UPDATE users SET email_verified = TRUE WHERE provider != '';
