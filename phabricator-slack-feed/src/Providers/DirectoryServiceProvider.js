'use strict';

let ServiceProvider = use('Chromabits/Container/ServiceProvider');

let listings = use('Listings');
/**
 * Class DirectoryServiceProvider
 *
 * Providers lookups for Slack usernames
 */
class DirectoryServiceProvider extends ServiceProvider
{
    /**
     * Register services
     *
     * @param app
     */
    register (app)
    {
        app.instance('Directory', listings);
    }
}

module.exports = DirectoryServiceProvider;
