/**
 * Enclosure Bootstrap Script
 *
 * This is the starting point of an Enclosure application.
 * Here we setup anything we need before Enclosure starts doing its thing.
 */

// We want to use ES6 so we enable Babel's transpiler which will automagically
// transform ES6 files into ES5-level code
require('babel/register');

// Next, we tell Enclosure to start the application.
require('enclosure').boot();
