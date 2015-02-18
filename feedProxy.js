var express = require('express'),
    bodyParser = require('body-parser'),
    config = require('config');

var app = express();

app.use(bodyParser.json());

app.post('/feed', function(req, res){
    res.send('OK');

    console.log(req.body);
});

app.listen(Number(config.get('server.port')));
console.log(
    'Phabricator-Slack connector server started on port %s',
    config.get('server.port')
);
