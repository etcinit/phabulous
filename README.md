# Phabulous

![Example](http://i.imgur.com/128Gkjw.png)

Phabulous forwards feed events from Phabricator to Slack.

> **Node.js version:** The Javascript version of this project has been replaced
with a rewrite from scratch in Go. The code for the Javascript version is
available at the **legacy** branch, but it won't be actively maintained.

## Features

- Route specific events (Tasks, Revisions, Commits) into specific channels.
- Push all feed events into a single channel (This may flood a channel if your
  organization is big enough).
- Pretty icons ;)

## Requirements

- Phabricator admin access and a certificate
- Slack API token

## Getting started

TODO

## Troubleshooting

### Self-signed certificates

If you are using self-signed certificates for your Phabricator instance, you
can disable checking at your own risk by setting `misc.ignore-ca` to `true` on
your configuration file.
