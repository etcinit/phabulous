![Phabulous](http://i.imgur.com/0ezr6XZ.png)

![Example](http://i.imgur.com/Uv4nVJa.png)

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

1. Create a directory for the bot: `mkdir phabulous`.
2. Download the latest stable release for your OS on the releases page and
save it there.
3. Create a configuration file `config/main.yml` using the one on this
repository as a template.
4. Get a Slack Web API token from https://api.slack.com/web. Write it down on
the config file (`slack.token` key).
5. Go to your Phabricator's settings panel and get your Conduit certificate: https://such.phabricator.wow/settings/panel/apitokens/, and place it on the
configuration file as well (`conduit.cert`). You will also need to specify the
URL for your Phabricator instance (`conduit.api`).
6. Start the server: `./phabulous server`.
7. If everything is working fine, add an entry to your Phabricator's config to
call Phabulous' event webhook. The URL should be the address and port of the
server you are running the bot in. Example for bot running in the same server
as Phabricator:

```json
  //...
  "feed.http-hooks": [
      "http://localhost:8086"
  ]
  //...
```

## Troubleshooting

### Self-signed certificates

If you are using self-signed certificates for your Phabricator instance, you
can disable checking at your own risk by setting `misc.ignore-ca` to `true` on
your configuration file.
