CREATE SCHEMA main AUTHORIZATION dbuser;
CREATE SCHEMA api AUTHORIZATION dbuser;

CREATE TABLE main."user" (
	id serial NOT NULL,
	"name" varchar NOT NULL
);

ALTER TABLE main."user" OWNER TO dbuser;
GRANT ALL ON TABLE main."user" TO dbuser;

CREATE TABLE main."comment" (
	id serial NOT NULL,
	msg text NOT NULL,
	user_id int4 NOT NULL,
	create_dt timestamp NOT NULL
);

ALTER TABLE main."comment" OWNER TO dbuser;
GRANT ALL ON TABLE main."comment" TO dbuser;

INSERT INTO main."user" ("name")
select md5(random()::text) 
from generate_series(1, 100);

insert into main."comment" (msg, user_id, create_dt)
select 
  md5(random()::text) msg,
  generate_series(1, 100) user_id,
  generate_series('2019-11-06 00:00'::timestamp,
				  '2019-11-10 03:00', '1 hours') create_dt;

create or replace function api.get_comment_by_id(_id integer)
 RETURNS json
 LANGUAGE sql
 SECURITY DEFINER
AS $function$
 select row_to_json(t)
 from (
	select 
	 id,
	 msg,
	 user_id,
	 create_dt
	from main."comment"
	where id = _id
    ) t;
$function$
;

create or replace function api.get_comments()
 RETURNS json
 LANGUAGE sql
 SECURITY DEFINER
AS $function$
select array_to_json(array_agg(row_to_json(t)))
 from (
    select 
	 id,
	 msg,
	 user_id,
	 create_dt
	from main."comment"
    ) t;
$function$
;

create or replace function api.update_comment(_id integer, _msg text)
 RETURNS integer
 LANGUAGE sql
 SECURITY DEFINER
AS $function$
 UPDATE main."comment"
 SET msg = _msg,
	 create_dt = now()
 where id = _id
returning id;
$function$
;

create or replace function api.create_comment(_user_id integer, _msg text)
 RETURNS integer
 LANGUAGE sql
 SECURITY DEFINER
AS $function$
 INSERT INTO main."comment"
 (msg, user_id, create_dt)
 VALUES(_msg, _user_id, now())
 returning id;
$function$
;

create or replace function api.delete_comment(_id integer)
 RETURNS integer
 LANGUAGE sql
 SECURITY DEFINER
AS $function$
 DELETE FROM main."comment"
 WHERE id = _id
 returning id;
$function$
;

create or replace function api.get_users()
 RETURNS json
 LANGUAGE sql
 SECURITY DEFINER
AS $function$
select array_to_json(array_agg(row_to_json(t)))
 from (
    select 
	 id,
	 name
	from main."user"
    ) t;
$function$
;

CREATE OR REPLACE FUNCTION api.create_user(_user_name varchar)
 RETURNS integer
 LANGUAGE sql
 SECURITY DEFINER
AS $function$
 INSERT INTO main."user"
 ("name")
 VALUES(_user_name)
 returning id;
$function$
;

CREATE OR REPLACE FUNCTION api.update_user(_id integer, _user_name text)
 RETURNS integer
 LANGUAGE sql
 SECURITY DEFINER
AS $function$
 UPDATE main."user"
 SET "name" = _user_name
 where id = _id
returning id;
$function$
;

CREATE OR REPLACE FUNCTION api.delete_user(_uid integer)
 RETURNS integer
 LANGUAGE sql
 SECURITY DEFINER
AS $function$
 DELETE FROM main."user"
 WHERE id = _uid
 returning id;
$function$
;