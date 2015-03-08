'use strict';

let superagent = require('superagent');

/**
 * Class Poster
 *
 * Post messages to Slack
 */
class Poster
{
    /**
     * Construct an instance of a Poster
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
    send (username, message, callback)
    {
        superagent
            .post(this.config.get('slack.url'))
            .send({
                username: username,
                text: message
            })
            .end(function(error, res){
                if (callback) {
                    callback(error, res);
                }
            });
    }
}

module.exports = Poster;
