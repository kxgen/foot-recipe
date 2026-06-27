alter table "public"."recipes" add column "status" varchar(20) not null default 'draft';
alter table "public"."recipes" add constraint "recipes_status_check" check (status IN ('draft', 'published', 'archived'));
