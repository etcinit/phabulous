'use strict';

/**
 * Class IndexController
 *
 * Handles index routes
 */
class IndexController
{
    /**
     * POST /
     *
     * @param req
     * @param res
     */
    static postIndex (req, res)
    {
        res.send('OK');

        if (req.body.storyData && req.body.storyData.objectPHID) {
            var processor = container.make('Processor');

            // Defer processing to a service
            processor.handle(req.body);
        }

        // Some debug
        console.log(req.raw.toString('utf-8'));
    }
}

module.exports = IndexController;
