'use strict';

/**
 * Class FeedProcessor
 *
 * Takes Phabricator feed items and sends them to Slack
 */
class FeedProcessor
{
    /**
     * Construct an instance of a FeedProcessor
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
     * Process a feed item
     *
     * @param reqBody
     */
    handle (reqBody)
    {
        // Now, we fetch additional information from Conduit
        this.fetcher.fetch(reqBody.storyData.objectPHID, (err, data) => {
            // If we have a URI, add it to the message
            if (!!data && data.uri) {
                this.poster.send(
                    this.config.get('slack.username'),
                    reqBody.storyText + " (<" + data.uri + "|More info>)"
                );
            } else {
                this.poster.send(
                    this.config.get('slack.username'),
                    reqBody.storyText
                );
            }
        });
    }
}

module.exports = FeedProcessor;
