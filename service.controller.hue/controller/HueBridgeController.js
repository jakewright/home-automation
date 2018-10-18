const huejay = require("huejay");

class HueBridgeController {
  constructor(express, hueClient) {
    this.hueClient = hueClient;

    /* Routes */
    express.get("/hue/discover", this.discoverBridges.bind(this));
    express.post("/hue/bridge/users", this.createUser.bind(this));
    express.get("/hue/bridge/users", this.getAllUsers.bind(this));
    express.get("/hue/bridge/lights", this.getAllLights.bind(this));
  }

  async discoverBridges(req, res, next) {
    try {
      const bridges = await this.hueClient.discover();
      res.json({ data: bridges });
    } catch (err) {
      next(err);
    }
  }

  async createUser(req, res, next) {
    try {
      const user = await this.hueClient.createUser();
      res.json({ data: user });
    } catch (err) {
      if (err instanceof huejay.Error && err.type === 101) {
        res.status(412);
        res.json({ message: "Link button not pressed" });
        return;
      }

      next(err);
    }
  }

  async getAllUsers(req, res, next) {
    try {
      const users = await this.hueClient.getAllUsers();
      res.json({ data: users });
    } catch (err) {
      next(err);
    }
  }

  async getAllLights(req, res, next) {
    try {
      const lights = await this.hueClient.getAllLights();
      res.json({ data: lights });
    } catch (err) {
      next(err);
    }
  }
}

exports = module.exports = HueBridgeController;
