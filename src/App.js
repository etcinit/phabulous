'use strict';

/**
 * Class App
 *
 * Entrypoint of the application
 */
class App
{
    /**
     * Construct an instance of an App
     *
     * @param Config
     * @param Conduit_CanduitFactory
     * @param Conduit_CaHelper
     */
    constructor (Config, Conduit_CanduitFactory, Conduit_CaHelper)
    {
        this.config = Config;
        this.canduit = Conduit_CanduitFactory;
        this.cahelper = Conduit_CaHelper;
    }

    /**
     * Run the application
     */
    main ()
    {
        // Setup CA config if needed
        this.cahelper.setup();

        let Processor = use('FeedProcessor');

        // Begin setting up the app
        this.canduit.make((err, conduit) => {
            container.instance('Conduit', conduit);

            // `true` indicates that Processor is shared (aka a singleton)
            container.bind('FeedProcessor', Processor, true);

            let server = container.make('Http/Server');

            server.setup();
            server.listen();
        });
    }
}

module.exports = App;
