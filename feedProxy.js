/**
 * Phabricator-Slack integration bot
 *
 * This is simple bot that will forward Phabricator feed events into a Slack
 * channel. Besides this server, you will need to configure Phabricator and
 * Slack independently to get the bot working.
 *
 * See more in the README
 *
 * @author Eduardo Trujillo <ed@chromabits.com>
 */

var express = require('express'),
    bodyParser = require('body-parser'),
    config = require('config'),
    concat = require('concat-stream'),
    Qs = require('qs');
    Canduit = require('canduit'),
    superagent = require('superagent');

var app = express(),
    conduit,
    fetchPhidData,
    setupConduit,
    setupCA,
    setupApp,
    startListening;

// Some Phabricator servers could be using self-signed certificates
// If that is the case, we can give users a choice to disable this check if
// they really want to
setupCA = function () {
    if (config.get('misc.ignore-ca')) {
        process.env['NODE_TLS_REJECT_UNAUTHORIZED'] = '0';
    }
};

// Here we create the instance of the Conduit client for fetching PHID info
// from Phabricator
setupConduit = function (callback) {
    conduit = Canduit(config.get('conduit'), function (err) {
        if (err) {
            throw err;
        }

        console.log('Successfully connected to Phabricator over Conduit');

        callback();
    });
};

// This fetches information about a Phabricator object by using its PHID
fetchPhidData = function (phid, callback) {
    conduit.exec(
        'phid.query',
        {
            "phids": {
                "0": phid
            }
        },
        function (err, response) {
            if (err) {
                callback(err);
            }

            // This call will return a hash of ids and objects so we need to
            // get the first key
            var keys = Object.keys(response);

            callback(null, response[keys[0]]);
        }
    )
};

// Here we just setup a very simple Express application for receiving the HTTP
// hook calls from Phabricator
setupApp = function () {
    app.use(function(req, res, next){
        req.pipe(concat(function(data){
            req.raw = data;
            next();
        }));
    });

    app.use(function(req, res, next) {
        req.body = Qs.parse(req.raw.toString('utf-8'));
        next();
    });

    app.post('/', function(req, res){
        res.send('OK');

        if (req.body.storyData && req.body.storyData.objectPHID) {
            fetchPhidData(req.body.storyData.objectPHID, function (err, data) {
                superagent
                    .post(config.get('slack.url'))
                    .send({
                        "username": config.get('slack.username'),
                        "text":
                        req.body.storyText +
                        " (<" + data.uri + "|More info>)"
                    })
                    .end(function(error, res){

                    });
            });
        }

        console.log(req.raw.toString('utf-8'));
    });
};

// Starts the HTTP server
startListening = function () {
    app.listen(Number(config.get('server.port')));
    console.log(
        'Phabricator-Slack connector server started on port %s',
        config.get('server.port')
    );
};

// Here we actually launch everything we previously defined
setupCA();
setupConduit(function () {
    setupApp();
    startListening();
});
