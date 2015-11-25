'use strict';

var domain = require('domain').create();

/**
 * Class FeedProcessor
 *
 * Takes Phabricator feed items and sends them to Slack
 */
class FeedProcessor
{
  /**
   * Construct an instance of a FeedProcessor
   *
   * @param Config
   * @param Slack_Poster
   * @param Conduit
   * @param Conduit_Fetcher
   */
  constructor (Config, Slack_Poster, Conduit_CanduitFactory, Conduit_Fetcher)
  {
    this.config = Config;
    this.username = this.config.get('slack.username');
    this.directory = this.config.get('users');
    this.poster = Slack_Poster;
    this.canduit = Conduit_CanduitFactory;
    this.fetcher = Conduit_Fetcher;

    // Binds are necessary because functions assigned to objects become "methods"
    // and `this` becomes the containing object. Who knew?
    this.handlers = {
      "differential.query": this.diffHandler.bind(this),
      "phid.query": this.PHIDHandler.bind(this)
    };

    domain.on('error', (err) => {
      console.error("Domain caught an error:");
      console.error(err);
      console.error("Error message: " + err.message);
      this.poster.send(this.username,
        '<http://giphy.com/gifs/help-the-lion-king-cz314BBYiCkiA|I have died!> "' + 
        err.message + '". See logs: `~/slackbot.log`.', (err) => {
        if (err) {
          console.error("Domain threw an error:");
          console.error(err);
          console.error("Error message: " + err.message);
        }
        process.exit(1);
      });
    });

    console.log('FeedProcessor instance initialized.');
  }

  /**
   * Construct the callback for all Conduit API calls
   * May attempt revovery one time via recursion
   *
   * @param phid
   * @param endpoint
   * @param attemptRecover
   * @param body
   * @param customCB
   */
  genericHandlerFactory (phid, endpoint, attemptRecover, body, customCB) {
    return (err, data) => {
      if (err) {
         
        // Currently we only attempt to recover from timed-out sessions
        if (err.code === 'ERR-INVALID-SESSION') {
          if (attemptRecover) {
            this.canduit.make((err, conduit) => {
              console.log('Tried to make a new Conduit connection because session was invalid');
              if (err) {
                console.error('Conduit connection creation failed');
                throw new Error('Conduit `' + endpoint + '` request failed. Recovery unsuccessful: ' + err.message);
              }
              else {
                console.log('Conduit reconnection succeeded');
                container.instance('Conduit', conduit);
                this.fetcher = container.make('Conduit/Fetcher');
                console.log('Created a new fetcher');
                this.fetcher.fetch(phid, endpoint,
                  this.genericHandlerFactory(phid, endpoint, false, body, customCB)
                );
              }
            });
          }
          else {
            throw new Error('Conduit `' + endpoint + '` request failed. Recovery unsuccessful: ' + err.message);
          }
        }
        else {
          throw new Error('Conduit `' + endpoint + '` request failed. Recovery not attempted: ' + err.message);
        }
      }
      else {
        // Default callback if none supplied
        if (typeof customCB === 'undefined') {
          this.handlers[endpoint](data, body);
        }
        else {
          customCB(data, body);
        }
      }
    };
  }

  /**
   * Process a response from `differential.query`
   *
   * @param data
   * @param body
   */
  diffHandler (data, body)
  {
    if (!body) {
      throw new Error('diffHandler() received no request body.');
    }
    else if (!data || !data.authorPHID || !data.reviewers || !data.ccs) {
      // Swap to PHID query
      let phid = body.storyData.objectPHID;
      let endpoint = 'phid.query';
      this.fetcher.fetch(phid, endpoint,
        this.genericHandlerFactory(phid, endpoint, true, body)
      );
    }
    else {
      let authorPHID = data.authorPHID;
      let reviewers = data.reviewers;
      let subscribers = data.ccs;
      let message = body.storyText;
      let query = 'user.query';
      
      // If we have a URI, add it to the message
      if (data.uri) {
        message += " (<" + data.uri + "|More info>)";
      }

      let authorFetch = () => {
        if (!authorPHID) {
          reviewerFetch(0);
        }
        else {
          let phid = authorPHID
          this.fetcher.fetch(phid, query,
            this.genericHandlerFactory(phid, query, true, null, (data) => {
              if (!data || !data.userName) {
                console.error('Incomplete data for PHID ' + phid + ': ' + data);
                throw new Error('Phabricator Conduit `' + query + '` returned incomplete data');
              }
              else {
                message += ' Author: ';
                message += '<' + this.directory[data.userName] + '>';
                reviewerFetch(0);
              }
            })
          );
        }
      }

      let reviewerFetch = (i) => {
        if (i === reviewers.length) {
          if (reviewers.length !== 0) {
            message += ']';
          }
          subscriberFetch(0);
        }
        else {
          let phid = reviewers[i];
          this.fetcher.fetch(phid, query,
            this.genericHandlerFactory(phid, query, true, null, (data) => {
              if (!data || !data.userName) {
                console.error('Incomplete data for PHID ' + phid + ': ' + data);
                throw new Error('Phabricator Conduit `' + query + '` returned incomplete data');
              }
              else {
                if (i === 0) {
                  message += ' Reviewers: [';
                }
                message += '<' + this.directory[data.userName] + '>';
                if (i < reviewers.length - 1) {
                  message += ', ';
                }
                reviewerFetch(i+1);
              }
            })
          );
        }
      }

      let subscriberFetch = (i) => {
        if (i == subscribers.length) {
          if (subscribers.length !== 0) {
            message += ']';
          }
          console.log(message);  // Debug logging
          this.poster.send(this.username, message);
        }
        else {
          let phid = subscribers[i];
          this.fetcher.fetch(phid, query,
            this.genericHandlerFactory(phid, query, true, null, (data) => {
              if (!data || !data.userName) {
                console.error('Incomplete data for PHID ' + phid + ': ' + data);
                throw new Error('Phabricator Conduit `' + query + '` returned incomplete data');
              }
              else {
                if (i === 0) {
                  message += ' Subscribers: [';
                }
                message += '<' + this.directory[data.userName] + '>';
                if (i < subscribers.length - 1) {
                  message += ', ';
                }
                subscriberFetch(i+1);
              }
            })
          );
        }
      }

      authorFetch();

    }
  }

  /**
   * Process a response from `phid.query`
   *
   * @param data
   * @param body
   */
  PHIDHandler (data, body)
  {
    if (!body) {
      throw new Error('PHIDHandler() received no request body.');
    }

    // Empty data payload received
    else if (!data) {
      throw new Error('Conduit `phid.query` returned no data');
    }
    else {
      let message = body.storyText;
      if (data.uri) {
        message += " (<" + data.uri + "|More info>)";
      }
      this.poster.send(this.username, message);
    }
  }

  /**
   * Process a feed item
   *
   * @param body
   */
  handle (body)
  {
    // Bind handler for all errors in scope
    domain.run(() => {
      if (!body || !body.storyData || !body.storyData.objectPHID || !body.storyText) {
        throw new Error('Phabricator webhook returned unexpected or incomplete data');
      }
      else {
        let phid = body.storyData.objectPHID;
        let query = 'differential.query';
        this.fetcher.fetch(phid, query,
          this.genericHandlerFactory(phid, query, true, body)
        );
      }
    });
  }
}

module.exports = FeedProcessor;
