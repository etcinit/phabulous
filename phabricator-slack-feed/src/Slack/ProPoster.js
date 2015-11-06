'use strict';

let superagent = require('superagent');

/**
 * Class ProPoster
 *
 * Post messages to Slack
 */
class ProPoster
{
    /**
     * Construct an instance of a ProPoster
     *
     * @param Config
     */
    constructor (Config)
    {
        this.config = Config;
    }

    /**
     * Post a message to slack
     *
     * @param username
     * @param message
     * @param callback
     */
    send (username, message, channel, callback)
    {
	var payload = {
            username: username,
            text: message
        };
        if (channel) {
            payload.channel = channel;
        }
        superagent
            .post(this.config.get('slack.url'))
            .send(payload)
            .end(function(error, res){
                if (callback) {
                    callback(error, res);
                }
            });
    }
}

module.exports = ProPoster;
