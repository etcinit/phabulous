'use strict';

/**
 * Class Fetcher
 *
 * Gets information from Conduit API
 */
class Fetcher
{
    /**
     * Construct an instance of a Fetcher
     *
     * @param Conduit
     */
    constructor (Conduit)
    {
        this.conduit = Conduit;
    }

    /**
     * Go and fetch information from Conduit
     *
     * @param phid
     * @param endpoint
     * @param callback
     */
    fetch (phid, endpoint, callback)
    {
        this.conduit.exec(
            endpoint,
            {
                "phids": {
                    "0": phid
                }
            },
            (err, response) => {
                if (err) {
                    callback(err);
                }

                // This call will return a hash of ids and objects so we need to
                // get the first key
                if (!!response) {
                  var keys = Object.keys(response);

                	callback(null, response[keys[0]]);
                }
            }
        )
    }
}

module.exports = Fetcher;
