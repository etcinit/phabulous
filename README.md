# Phabricator Slack Bridge

This small server will forward news feed items from Phabricator into a Slack
channel

## Requirements

- Node.js
- Phabricator admin access and a certificate
- Slack's incoming webhook integration

## Setup instructions

1. First, clone the repository and run `npm install`
2. While dependencies install, head over to your Slack admin panel and create a
**"Incoming Webhook"** integration. Give it a nice icon 
(like [this one](http://blogs.gnome.org/aklapper/files/2014/05/Phab_logo.png)).
After that you should take note of the **Webhook URL**
3. Go to your Phabricator's settings panel and get your Conduit certificate:
https://such.phabricator/settings/panel/apitokens/
4. Create a `production.json` inside the config directory by copying the
contents of `default.json`. Place your Phabricator username, url, and
certificate under the `user`, `api`, `cert` keys respectively under the
`conduit` key. In the `slack` key, set a username for the bot and a the url
of the webhook.
5. Start the server: `NODE_ENV=production node feedProxy` and check for errors
6. If everything is working fine, add an entry to your Phabricator's config to
call the webhook. The URL should be the address and port of the server you are
running the bot in. Example for bot running in the same server as Phabricator:

```js
    //...
    "feed.http-hooks": [
        "http://localhost:8085"
    ]
    //...
```

## Misc

### Self-signed certificates

If you are using self-signed certificates for your Phabricator instance, you
can disable checking at your own risk by setting `misc.ignore-ca` to `true`
