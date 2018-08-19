const config = require('config');
const express = require('express');
const http = require('http');
const redis = require('redis');
const url = require('url');
const WebSocket = require('ws');

// Connect to the Redis instance
redisClient = redis.createClient({
    host: config.get('redis.host'),
    port: config.get('redis.port'),
});

// Log if something goes wrong with the Redis connection
redisClient.on('error', err => {
    console.error(`Redis error: ${err}`);
});

// Subscribe to all state change events
redisClient.psubscribe("device-state-changed.*");

// Create an express app
const app = express();

// Create a websocket server
const server = http.createServer(app);
const wss = new WebSocket.Server({ server });

// When a message is received as a result of psubscribe
redisClient.on('pmessage', (pattern, channel, message) => {
    console.log(pattern);
    console.log(channel);
    console.log(message);

    message = JSON.parse(message);

    wss.clients.forEach(client => {
      if (client.readyState === WebSocket.OPEN) {
        client.send(JSON.stringify({channel, message}));
      }
    });
});

// Start the server
server.listen(config.get('port'), () => {
  console.log('Listening on port %d', server.address().port);
});
