import huejay from 'huejay';

export default class HueDiscoveryController {
  constructor(express, client) {
    this.client = client;

    /* Routes */
    express.get('/hue/discover', this.discoverBridges.bind(this));
    express.post('/hue/bridge/users', this.createUser.bind(this));
    express.get('/hue/bridge/users', this.getAllUsers.bind(this));
    express.get('/hue/bridge/lights', this.getAllLights.bind(this));
  }

  async discoverBridges(req, res) {
    const bridges = await client.discover();
    res.json({message: `${bridges.length} bridges found`, data: bridges});
  }

  async createUser(req, res) {
    try {
      const user = await this.client.createUser();
      res.json({message: `User ${user.username} created`, data: user});
    } catch(err) {
      if (err instanceof huejay.Error && err.type === 101) {
        res.status(412);
        res.json({message: 'Link button not pressed'});
        return;
      }

      throw err;
    }
  }

  async getAllUsers(req, res) {
    const users = await this.client.getAllUsers();
    res.json({message: 'Success', data: users});
  }

  async getAllLights(req, res) {
    const lights = await this.client.getAllLights();
    res.json({message: 'Success', data: lights});
  }
}
