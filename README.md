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
2. Download the latest release for your OS from the [releases page](https://github.com/etcinit/phabulous/releases) and
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
      "http://localhost:8086/v1/feed/receive"
  ]
  //...
```

## Help & Troubleshooting

### Environment

It is possible to override any variable in the configuration file through
environment variables:

```
export SLACK_TOKEN=mytoken
```

is the same as setting:

```yaml
slack:
  token: mytoken
```

See https://github.com/jacobstr/confer for more information.

### Events are not showing up

Make sure that the `feed.http-hooks` setting on your Phabricator instance is
setup correctly, and that the server can communicate with the Phabulous API.
An easy way to test this is using `curl` from the server hosting Phabricator:

```sh
curl http://localhost:8086
```

The command above should return something like:

```json
{
  "messages": ["Welcome to the Phabulous API"],
  "status": "success",
  "version": "1.0.0"
}
```

### Self-signed certificates

If you are using self-signed certificates for your Phabricator instance, you
can disable checking at your own risk by setting `misc.ignore-ca` to `true` on
your configuration file.

### OMG, the feed is flooding everything

The `channels.feed` setting tells Phabulous where to post about every single
feed event from Phabricator. This might get too noisy if you have a constant
stream of events. To disable the feed channel, just set it to an empty string:

```yaml
channels:
  feed: ''
```

### Routing events

Phabulous supports routing events concerning Revisions, Tasks and Commits to
specific channels, such as Project's or Repo's channel:

```yaml
channels:
  feed: '#phabricator'
  repositories:
    CALLSIGN: '#channel'
    OTHERCALLSIGN: '#otherchannel'
  projects:
    10: '#anotherchannel'
```

Specifying repository-channel mappings will cause Revision and Commit events
for that repository to be sent to the provided channel. The same applies for
project-channels mappings and Task events.

Project IDs can be found in the URL of a project.

## Compiling from source

```
go get github.com/etcinit/phabulous

// or, for cross-compiling:

go get github.com/mitchellh/gox
git clone git@github.com:etcinit/phabulous.git
cd phabulous
make
```

## Roadmap

- Improve error handling in general.
- Add support for various commands, such as looking up objects or creating
memes using macros.
- Add support for etcd or some database configuration backend, so the server
does not need to be restarted in order to update its configuration.
- Windows support?
