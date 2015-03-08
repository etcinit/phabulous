'use strict';

/**
 * Class IndexController
 *
 * Handles index routes
 */
class IndexController
{
    /**
     * Construct an instance of an IndexController
     *
     * @param Config
     * @param Slack_Poster
     * @param Conduit_PhidFetcher
     */
    constructor (Config, Slack_Poster, Conduit_PhidFetcher)
    {
        this.config = Config;
        this.poster = Slack_Poster;
        this.fetcher = Conduit_PhidFetcher;
    }

    /**
     * POST /
     *
     * @param req
     * @param res
     */
    postIndex (req, res)
    {
        res.send('OK');

        if (req.body.storyData && req.body.storyData.objectPHID) {
            // Now, we fetch additional information from Conduit
            this.fetcher.go(req.body.storyData.objectPHID, (err, data) => {
                // If we have a URI, add it to the message
                if (data.uri) {
                    this.poster.send(
                        this.config.get('slack.username'),
                        req.body.storyText + " (<" + data.uri + "|More info>)"
                    );
                } else {
                    this.poster.send(
                        this.config.get('slack.username'),
                        req.body.storyText
                    );
                }
            });
        }

        // Some debug
        console.log(req.raw.toString('utf-8'));
    }
}

module.exports = IndexController;
