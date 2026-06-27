alter table "public"."recipes" alter column "published_at" drop not null;
alter table "public"."recipes" add column "published_at" timestamptz;
