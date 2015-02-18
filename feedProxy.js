var express = require('express'),
    bodyParser = require('body-parser'),
    config = require('config'),
    concat = require('concat-stream'),
    superagent = require('superagent');

var app = express();

app.use(bodyParser.urlencoded({ extended: false }));
app.use(bodyParser.json());
app.use(function(req, res, next){
    req.pipe(concat(function(data){
        req.raw = data;
        next();
    }));
});

app.post('/', function(req, res){
    res.send('OK');

    superagent
        .post(config.get('slack.url'))
        .send({
            "username": config.get('slack.username'),
            "text": req.body.storyText
        })
        .end(function(error, res){

        });

    console.log(req.body);
});

app.listen(Number(config.get('server.port')));
console.log(
    'Phabricator-Slack connector server started on port %s',
    config.get('server.port')
);
