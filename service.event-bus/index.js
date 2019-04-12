const express = require("express");
const http = require("http");
const WebSocket = require("ws");

const bootstrap = require("../libraries/javascript/bootstrap");
const config = require("../libraries/javascript/config");
const firehose = require("../libraries/javascript/firehose");

bootstrap("service.event-bus")
  .then(() => {
    // Create an express app
    const app = express();

    // Create a websocket server
    const server = http.createServer(app);
    const wss = new WebSocket.Server({ server });

    firehose.subscribe("device-state-changed.*", (channel, message) => {
      console.log(`Message received on channel '${channel}'\n${message}\n`);

      message = JSON.parse(message);

      wss.clients.forEach(client => {
        if (client.readyState === WebSocket.OPEN) {
          client.send(JSON.stringify({ channel, message }));
        }
      });
    });

    // Start the server
    server.listen(config.get("port", 80), () => {
      console.log("Listening on port %d", server.address().port);
    });
  })
  .catch(err => {
    console.error("Error initialising service", err);
  });
