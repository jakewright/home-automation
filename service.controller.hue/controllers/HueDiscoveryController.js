const huejay = require("huejay");

class HueDiscoveryController {
  constructor(express, client) {
    this.client = client;

    /* Routes */
    express.get("/hue/discover", this.discoverBridges.bind(this));
    express.post("/hue/bridge/users", this.createUser.bind(this));
    express.get("/hue/bridge/users", this.getAllUsers.bind(this));
    express.get("/hue/bridge/lights", this.getAllLights.bind(this));
  }

  discoverBridges(req, res, next) {
    this.client
      .discover()
      .then(bridges => {
        res.json({ message: `${bridges.length} bridges found`, data: bridges });
      })
      .catch(next);
  }

  createUser(req, res, next) {
    this.client
      .createUser()
      .then(user => {
        res.json({ message: `User ${user.username} created`, data: user });
      })
      .catch(err => {
        if (err instanceof huejay.Error && err.type === 101) {
          res.status(412);
          res.json({ message: "Link button not pressed" });
          return;
        }

        next(err);
      });
  }

  getAllUsers(req, res, next) {
    this.client
      .getAllUsers()
      .then(users => {
        res.json({ message: "Success", data: users });
      })
      .catch(next);
  }

  getAllLights(req, res, next) {
    this.client
      .getAllLights()
      .then(lights => {
        res.json({ message: "Success", data: lights });
      })
      .catch(err);
  }
}

exports = module.exports = HueDiscoveryController;
