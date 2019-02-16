const huejay = require("huejay");
const hueClient = require("../api/hueClient");

const discover = async (req, res, next) => {
  try {
    const bridges = await hueClient.discover();
    res.json({ data: bridges });
  } catch (err) {
    next(err);
  }
};

const createUser = async (req, res, next) => {
  try {
    const user = await hueClient.createUser();
    res.json({ data: user });
  } catch (err) {
    if (err instanceof huejay.Error && err.type === 101) {
      res.status(412);
      res.json({ message: "Link button not pressed" });
      return;
    }

    next(err);
  }
};

const getAllUsers = async (req, res, next) => {
  try {
    const users = await hueClient.getAllUsers();
    res.json({ data: users });
  } catch (err) {
    next(err);
  }
};

const getAllLights = async (req, res, next) => {
  try {
    const lights = await hueClient.getAllLights();
    res.json({ data: lights });
  } catch (err) {
    next(err);
  }
};


exports = module.exports = { discover, createUser, getAllUsers, getAllLights };
