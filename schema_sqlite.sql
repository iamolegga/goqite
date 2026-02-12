create table if not exists {{.Table}} (
  id text primary key default ('m_' || lower(hex(randomblob(16)))),
  created text not null default (strftime('%Y-%m-%dT%H:%M:%fZ')),
  updated text not null default (strftime('%Y-%m-%dT%H:%M:%fZ')),
  queue text not null,
  body blob not null,
  timeout text not null default (strftime('%Y-%m-%dT%H:%M:%fZ')),
  received integer not null default 0,
  priority integer not null default 0
) strict;

create trigger if not exists {{.Trigger}} after update on {{.Table}} begin
  update {{.Table}} set updated = strftime('%Y-%m-%dT%H:%M:%fZ') where id = old.id;
end;

create index if not exists {{.Index}} on {{.Table}} (queue, priority desc, created);
