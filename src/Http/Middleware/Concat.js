'use strict';

let concat = require('concat-stream');

/**
 * Class Concat
 *
 * Concat middleware
 */
class Concat
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
        req.pipe(concat(function(data){
            req.raw = data;
            next();
        }));
    }
}

module.exports = Concat;
