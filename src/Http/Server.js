'use strict';

let Concat = use('Http/Middleware/Concat'),
    Qs = use('Http/Middleware/Qs'),
    IndexController = use('Http/Controllers/IndexController');

let express = require('express');

/**
 * Class Server
 *
 * Main application server
 */
class Server
{
    /**
     * Construct an instance of a Server
     *
     * @param Config
     */
    constructor (Config)
    {
        this.config = Config;

        this.exp = express();
    }

    /**
     * Setup the server
     */
    setup ()
    {
        this.exp.use(Concat.handle);
        this.exp.use(Qs.handle);

        this.exp.post('/', IndexController.postIndex);
    }

    /**
     * Begin listening
     */
    listen ()
    {
        this.exp.listen(Number(this.config.get('server.port')));

        console.log(
            'Phabricator-Slack connector server started on port %s',
            this.config.get('server.port')
        );
    }
}

module.exports = Server;
