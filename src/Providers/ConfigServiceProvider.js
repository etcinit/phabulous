'use strict';

let ServiceProvider = use('Chromabits/Container/ServiceProvider');

let config = require('config');

/**
 * Class ConfigServiceProvider
 *
 * Provides configuration
 */
class ConfigServiceProvider extends ServiceProvider
{
    /**
     * Register services
     *
     * @param app
     */
    register (app)
    {
        app.instance('Config', config);
    }
}

module.exports = ConfigServiceProvider;
