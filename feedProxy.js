var express = require('express'),
    bodyParser = require('body-parser');

var app = express.createServer();

app.use(bodyParser.json());

app.post('/feed', function(req, res){
    res.send('OK');

    console.log(req.body);
});

app.listen();
console.log('Express server started on port %s', app.address().port);
