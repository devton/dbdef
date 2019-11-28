package main

var structureQuery = `
select
    n.nspname as schema,
    c.proname as relation,
    'function' as kind,
    c.prosrc as definition,
    'CREATE OR REPLACE FUNCTION '
            || quote_ident(n.nspname) || '.'
            || quote_ident(c.proname) || '('
            || pg_catalog.pg_get_function_identity_arguments(c.oid)
            || ')' AS signature
from pg_catalog.pg_proc c
join pg_catalog.pg_namespace n on n.oid = c.pronamespace
left join information_schema.views iv on iv.table_name = c.proname and iv.table_schema = n.nspname
where
  n.nspname in ('public')
  and c.probin is null
  and c.prosrc is not null
union
select
  n.nspname as schema,
  c.relname as relation,
  (
    case c.relkind::text
    when 'v' then 'view'
    when 'r' then 'table'
    when 'm' then 'materialized_view'
    when 'f' then 'foreign'
     end
  )as kind,
  iv.view_definition as definition,
  '' as signature
from
  pg_catalog.pg_class c
join pg_catalog.pg_namespace n on n.oid = c.relnamespace
left join information_schema.views iv on iv.table_name = c.relname and iv.table_schema = n.nspname
where
  c.relkind in ('r', 'v', 'm', 'f')
  and n.nspname in ('public');
`
