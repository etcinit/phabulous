'use strict';

let Concat = use('Http/Middleware/Concat'),
    Qs = use('Http/Middleware/Qs');

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

        let indexController = container.make('Http/Controllers/IndexController');
        this.exp.post('/', indexController.postIndex);
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
