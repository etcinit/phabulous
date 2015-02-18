var express = require('express'),
    bodyParser = require('body-parser'),
    config = require('config'),
    concat = require('concat-stream'),
    Qs = require('qs');
    Canduit = require('canduit'),
    superagent = require('superagent');

var app = express(),
    conduit,
    fetchPhidData;

if (config.get('misc.ignore-ca')) {
    process.env['NODE_TLS_REJECT_UNAUTHORIZED'] = '0';
}

conduit = Canduit(config.get('conduit'), function (err) {
    if (err) {
        throw err;
    }

    console.log('Successfully connected to Phabricator over Conduit');
});

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

            var keys = Object.keys(response);

            callback(null, response[keys[0]]);
        }
    )
};

//app.use(bodyParser.urlencoded({ extended: true }));
//app.use(bodyParser.json());
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

app.listen(Number(config.get('server.port')));
console.log(
    'Phabricator-Slack connector server started on port %s',
    config.get('server.port')
);
