var express = require('express'),
    bodyParser = require('body-parser'),
    config = require('config'),
    concat = require('concat-stream');

var app = express();

app.use(bodyParser.json());
app.use(function(req, res, next){
    req.pipe(concat(function(data){
        req.raw = data;
        next();
    }));
});

app.post('/', function(req, res){
    res.send('OK');

    console.log(req.raw.toString('utf8'));
});

app.listen(Number(config.get('server.port')));
console.log(
    'Phabricator-Slack connector server started on port %s',
    config.get('server.port')
);
