-- Add new schema named "public"
CREATE SCHEMA IF NOT EXISTS "public";
-- Set comment to schema: "public"
COMMENT ON SCHEMA "public" IS 'standard public schema';
-- Create enum type "priority_level"
CREATE TYPE "public"."priority_level" AS ENUM ('low', 'med', 'high');
-- Create "schema_migrations" table
CREATE TABLE "public"."schema_migrations" (
  "version" bigint NOT NULL,
  "dirty" boolean NOT NULL,
  PRIMARY KEY ("version")
);
-- Create "users" table
CREATE TABLE "public"."users" (
  "id" serial NOT NULL,
  "email" character varying(255) NOT NULL,
  "name" character varying(255) NOT NULL,
  "password_hash" character varying(255) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  CONSTRAINT "users_email_key" UNIQUE ("email")
);
-- Create "refresh_tokens" table
CREATE TABLE "public"."refresh_tokens" (
  "id" serial NOT NULL,
  "user_id" integer NOT NULL,
  "token_hash" character varying(255) NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "device_info" character varying(255) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  CONSTRAINT "refresh_tokens_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_refresh_tokens_expires_at" to table: "refresh_tokens"
CREATE INDEX "idx_refresh_tokens_expires_at" ON "public"."refresh_tokens" ("expires_at");
-- Create index "idx_refresh_tokens_token_hash" to table: "refresh_tokens"
CREATE INDEX "idx_refresh_tokens_token_hash" ON "public"."refresh_tokens" ("token_hash");
-- Create index "idx_refresh_tokens_user_id" to table: "refresh_tokens"
CREATE INDEX "idx_refresh_tokens_user_id" ON "public"."refresh_tokens" ("user_id");
-- Create "todos" table
CREATE TABLE "public"."todos" (
  "id" serial NOT NULL,
  "user_id" integer NOT NULL,
  "title" character varying(255) NOT NULL,
  "description" text NULL,
  "completed" boolean NOT NULL DEFAULT false,
  "priority" "public"."priority_level" NOT NULL DEFAULT 'med',
  "due_date" timestamptz NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "todos_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_todos_completed" to table: "todos"
CREATE INDEX "idx_todos_completed" ON "public"."todos" ("completed");
-- Create index "idx_todos_deleted_at" to table: "todos"
CREATE INDEX "idx_todos_deleted_at" ON "public"."todos" ("deleted_at");
-- Create index "idx_todos_due_date" to table: "todos"
CREATE INDEX "idx_todos_due_date" ON "public"."todos" ("due_date");
-- Create index "idx_todos_priority" to table: "todos"
CREATE INDEX "idx_todos_priority" ON "public"."todos" ("priority");
-- Create index "idx_todos_user_id" to table: "todos"
CREATE INDEX "idx_todos_user_id" ON "public"."todos" ("user_id");
