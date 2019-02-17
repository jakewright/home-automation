const bodyParser = require("body-parser");
const express = require("express");
const config = require("../config");

class Router {
  constructor() {
    this.app = express();

    // JSON body parser
    this.app.use(bodyParser.json());

    // Request logger
    this.app.use((req, res, next) => {
      console.log(
        `${req.method} ${req.originalUrl} ${JSON.stringify(req.body)}`
      );
      next();
    });
  }

  listen() {
    // Add an error handler that returns valid JSON
    this.app.use(function(err, req, res, next) {
      console.error(err.stack);
      res.status(500);
      res.json({ message: err.message });
    });

    const port = config.get("port", 80);
    this.app.listen(port, () => {
      console.log(`Service running on port ${port}`);
    });
  }

  use(path, handler) {
    this.app.use(path, handler);
  }

  get(path, handler) {
    this.app.get(path, handler);
  }

  put(path, handler) {
    this.app.put(path, handler);
  }

  post(path, handler) {
    this.app.post(path, handler);
  }

  patch(path, handler) {
    this.app.patch(path, handler);
  }
}

const router = new Router();
exports = module.exports = router;
