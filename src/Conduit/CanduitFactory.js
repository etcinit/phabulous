'use strict';

let canduit = require('canduit');

/**
 * Class Canduit
 *
 * Provides an interface for communicating with Conduit
 */
class Canduit
{
    /**
     * Create an instance of a Canduit
     *
     * @param Config
     */
    constructor (Config)
    {
        this.config = Config;
    }

    /**
     * Make a connection to Conduit
     *
     * @param callback
     *
     * @returns {*|exports}
     */
    make (callback)
    {
        return canduit(this.config.get('conduit'), function (err, canduit) {
            if (err) {
                throw err;
            }

            if (canduit) {
                console.log(
                    'Successfully connected to Phabricator over Conduit'
                );

                if (callback) {
                    callback(err, canduit);
                }
            } else {
                throw new Error('Something went wrong connecting to Conduit');
            }
        });
    }
}

module.exports = Canduit;
