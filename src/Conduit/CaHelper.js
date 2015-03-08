'use strict';

/**
 * Class CaHelper
 *
 * Allow self-signed CA certs if required
 */
class CaHelper
{
    /**
     * Construct an instance of a CaHelper
     *
     * @param Config
     */
    constructor (Config)
    {
        this.config = Config;
    }

    /**
     * Configure
     */
    setup ()
    {
        if (this.config.get('misc.ignore-ca')) {
            process.env['NODE_TLS_REJECT_UNAUTHORIZED'] = '0';
        }
    }
}

module.exports = CaHelper;
