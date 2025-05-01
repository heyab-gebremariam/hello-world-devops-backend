-- Create "tasks" table
CREATE TABLE "tasks" (
  "id" bigserial NOT NULL,
  "title" text NOT NULL,
  "description" text NULL,
  "completed" boolean NULL DEFAULT false,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_tasks_deleted_at" to table: "tasks"
CREATE INDEX "idx_tasks_deleted_at" ON "tasks" ("deleted_at");
