CREATE SCHEMA IF NOT EXISTS public;

-- Create the Users table
CREATE TABLE IF NOT EXISTS public.users (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  username text NOT NULL UNIQUE,
  email text NOT NULL UNIQUE,
  avatar text,
  blocked boolean,
  created_at timestamp with time zone DEFAULT current_timestamp,
  updated_at timestamp with time zone DEFAULT current_timestamp
);

-- Create the MagicLinks table
CREATE TABLE IF NOT EXISTS public.magic_links (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  code text NOT NULL,
  user_id uuid NOT NULL REFERENCES public.users(id) ON DELETE RESTRICT ON UPDATE CASCADE,
  cookie text NOT NULL UNIQUE,
  valid_until timestamp with time zone NOT NULL,
  used boolean NOT NULL DEFAULT false
);

-- Create the Sessions table
CREATE TABLE IF NOT EXISTS public.sessions (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid NOT NULL REFERENCES public.users(id) ON DELETE RESTRICT ON UPDATE CASCADE,
  valid_until timestamp with time zone NOT NULL,
  user_agent text,
  ip text,
  latitude double precision,
  longitude double precision,
  city text,
  region text,
  country text,
  isp text,
  created_at timestamp with time zone DEFAULT current_timestamp,
  updated_at timestamp with time zone DEFAULT current_timestamp
);
