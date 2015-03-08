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
            var processor = container.make('FeedProcessor');

            // Defer processing to a service
            processor.handle(req.body);
        }

        // Some debug
        console.log(req.raw.toString('utf-8'));
    }
}

module.exports = IndexController;
