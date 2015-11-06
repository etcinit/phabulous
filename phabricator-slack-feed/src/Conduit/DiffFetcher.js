'use strict';

/**
 * Class DiffFetcher
 *
 * Gets information about diffs
 */
class DiffFetcher
{
    /**
     * Construct an instance of a DiffFetcher
     *
     * @param Conduit
     */
    constructor (Conduit)
    {
        this.conduit = Conduit;
    }

    /**
     * Go and fetch information about a diff from Conduit
     *
     * @param phid
     * @param callback
     */
    fetch (phid, callback)
    {
        this.conduit.exec(
            'differential.query',
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

module.exports = DiffFetcher;
