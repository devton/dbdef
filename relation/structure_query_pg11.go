package relation

const SQLStructPG11 = `
select
	current_database() as dbname,
    n.nspname as schema,
    c.proname as relation,
    'functions' as kind,
	d.description as description,
	(case when c.prokind <> 'a' then pg_get_functiondef(c.oid) else c.prosrc end) as definition,
    'CREATE OR REPLACE FUNCTION '
            || quote_ident(n.nspname) || '.'
            || quote_ident(c.proname) || '('
            || pg_catalog.pg_get_function_identity_arguments(c.oid)
            || ')' AS signature
from pg_catalog.pg_proc c
join pg_catalog.pg_namespace n on n.oid = c.pronamespace
left join pg_description d on d.objoid = c.oid
left join information_schema.views iv on iv.table_name = c.proname and iv.table_schema = n.nspname
where
  n.nspname ~* $1
  and c.probin is null
  and c.prosrc is not null
union
select
  current_database() as dbname,
  n.nspname as schema,
  c.relname as relation,
  (
    case c.relkind::text
    when 'v' then 'views'
    when 'r' then 'tables'
    when 'm' then 'materialized_views'
    when 'f' then 'foreign'
     end
  ) as kind,
  d.description as description,
  (
    case c.relkind::text
    when 'v' then iv.view_definition
    when 'r' then tbdef.definition
    when 'm' then 'materialized_view'
    when 'f' then 'foreign'
     end
  ) as definition,
  '' as signature
from
  pg_catalog.pg_class c
join pg_catalog.pg_namespace n on n.oid = c.relnamespace
left join pg_description d on d.objoid = c.oid
left join information_schema.views iv on iv.table_name = c.relname and iv.table_schema = n.nspname
left join lateral (
    SELECT
      'CREATE TABLE ' || n.nspname || '.' || relname || E'\n(\n' ||
      array_to_string(
        array_agg(
          '    ' || column_name || ' ' ||  type || ' '|| not_null
        )
        , E',\n'
      ) || E'\n);\n\n' ||
      (SELECT array_to_string((array_agg(indexdef)), E'\n\n') FROM pg_indexes WHERE tablename = relname and schemaname = n.nspname) as definition
    from
    (
      SELECT
        c2.relname, a.attname AS column_name,
        pg_catalog.format_type(a.atttypid, a.atttypmod) as type,
        case
          when a.attnotnull
        then 'NOT NULL'
        else 'NULL'
        END as not_null
      FROM pg_class c2,
       pg_attribute a,
       pg_type t
       WHERE c2.oid = c.oid
       and a.attnum > 0
       AND a.attrelid = c2.oid
       AND a.atttypid = t.oid
     ORDER BY a.attnum
    ) as tabledefinition
    group by relname
) as tbdef on c.relkind = 'r'
where
  c.relkind in ('r', 'v', 'm', 'f')
  and n.nspname ~* $1
`
