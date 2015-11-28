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

## Guides

- [Getting Started](http://phabricator.chromabits.com/w/phabulous/start/):
A guide on how to setup Phabulous for the first time.
- [Upgrade Notes](http://phabricator.chromabits.com/w/phabulous/upgrade/):
Instructions on how to upgrade to newer versions of Phabulous.
- [Help & Troubleshooting](http://phabricator.chromabits.com/w/phabulous/faq/):
Tips and answers to common problems.
- [Wiki](http://phabricator.chromabits.com/w/phabulous/): More articles and
information about Phabulous.

## Requirements

- Phabricator admin access and a certificate
- Slack API token

## Compiling from source

```
go get github.com/etcinit/phabulous

// or, for cross-compiling:

go get github.com/mitchellh/gox
git clone git@github.com:etcinit/phabulous.git
cd phabulous
make
```
