{{- if .Schema}}
create schema if not exists {{.Schema}};
{{end -}}

create or replace function {{.Function}}()
returns trigger as $func$
begin
   new.updated = now();
   return new;
end;
$func$ language plpgsql;

create table if not exists {{.Table}} (
  id text primary key default ('m_' || encode(gen_random_bytes(16), 'hex')),
  created timestamptz not null default now(),
  updated timestamptz not null default now(),
  queue text not null,
  body bytea not null,
  timeout timestamptz not null default now(),
  received integer not null default 0,
  priority integer not null default 0
);

drop trigger if exists {{.Trigger}} on {{.Table}};

create trigger {{.Trigger}}
before update on {{.Table}}
for each row execute procedure {{.Function}}();

create index if not exists {{.Index}} on {{.Table}} (queue, priority desc, created);
