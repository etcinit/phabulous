'use strict';

let qs = require('qs');

/**
 * Class Qs
 *
 * Query string middleware
 */
class Qs
{
    /**
     * Handle a request
     *
     * @param req
     * @param res
     * @param next
     */
    static handle (req, res, next)
    {
        req.body = qs.parse(req.raw.toString('utf-8'));
        next();
    }
}

module.exports = Qs;
