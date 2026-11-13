# FAQ

## Why "bosla"?

Bosla (بصلة) means *boat* — the tool carries your project's services
together, version by version, without drifting apart.

## Why store state in PostgreSQL instead of a local file?

So that a team shares one source of truth: projects, services, and
releases live in a small, inspectable schema you can query with `psql`,
back up, and migrate. The CLI itself is stateless.

## What happens when I delete a project?

Everything under it is removed by foreign-key cascades: services, releases,
service versions, and metadata. `project delete` therefore requires the
`--yes-i-am-sure` flag. Service deletion does **not** have a confirmation
flag — be careful.

## Can a service version appear in multiple releases?

Yes. `service_versions` rows are unique per service, not per release. The
`project_version_services` join table simply points each release at the
service versions it ships, so the same service version can be referenced by
many releases.

## What do `--status`, `--git-hash`, and `--description` do on a release?

They describe the **service versions attached to the release**, and they are
only applied to service versions that the publish command creates. If a
service version already exists, it is reused unchanged.

## Is search a database query?

No. `project search` and `service search` load the full list (of the
project, for services) and match names in-process against a Go regular
expression. Fine for the scale bosla targets; it is not full-text search.

## Does bosla talk to my repositories?

No. Links (project links, service repo URLs) are stored as plain strings
for reference. Bosla does not clone, query, or tag anything remotely — the
`--git-hash` on publish is just metadata you pass in.

## Where do I report problems?

Open an issue on [github.com/felukka/bosla](https://github.com/felukka/bosla).
