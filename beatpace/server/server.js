// server.js

const express = require('express');
const querystring = require('querystring');
const { generateRandomString } = require('./utils'); // Assume you have a function to generate random string

const app = express();
const port = 8000; // or any port you prefer

const client_id = '66834eefd8464e0ab313121e1b8d4119';
const redirect_uri = 'http://localhost:3000/';

app.get('/login', function(req, res) {
  const state = generateRandomString(16);
  const scope = 'user-library-read user-top-read user-read-recently-played';

  const authorizeUrl = 'https://accounts.spotify.com/authorize?' +
    querystring.stringify({
      response_type: 'code',
      client_id,
      scope,
      redirect_uri,
      state
    });

  res.redirect(authorizeUrl);
});

app.listen(port, () => {
  console.log(`Server is running on http://localhost:${port}`);
});
