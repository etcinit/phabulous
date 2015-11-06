'use strict';

/**
 * Class UserFetcher
 *
 * Gets information about users
 */
class UserFetcher
{
    /**
     * Construct an instance of a UserFetcher
     *
     * @param Conduit
     */
    constructor (Conduit)
    {
        this.conduit = Conduit;
    }

    /**
     * Go and fetch information about a user from Conduit
     *
     * @param phid
     * @param callback
     */
    fetch (phid, callback)
    {
        this.conduit.exec(
            'user.query',
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
