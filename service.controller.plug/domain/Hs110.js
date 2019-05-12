const Hs100 = require("./Hs100");

class Hs110 extends Hs100 {
  constructor(config) {
    super(config);

    this.state.watts = {
      watts: { type: "float", immutable: true }
    };

    // Generate powerset of device identifiers
    const ps = getPowerset(Object.keys(config.attributes.devices));

    // Array of objects {combination: [id_1, ..., id_n], watts: x}
    this.powerMap = [];
    for (const combination of ps) {
      const watts = combination.reduce(
        (sum, identifier) => sum + config.attributes.devices[identifier],
        0
      );
      this.powerMap.push({ combination, watts });
    }
  }

  applyState(state) {
    super.applyState(state);

    let closest = null;
    let combination = null;

    for (const map of this.powerMap) {
      const d = Math.abs(map.watts - state.watts);

      if (closest === null || d < closest) {
        closest = d;
        combination = map.combination;
      }
    }

    this.combination = combination;
  }
}

/**
 * Return the powerset of the given array.
 * E.g., given [1, 2, 3], will return [[], [1], [2], [1, 2], [3], [1, 3], [2, 3], [1, 2, 3]].
 *
 * @param {array} arr
 * @return {array} powerset
 */
const getPowerset = arr => {
  let ps = [[]];
  for (let i = 0; i < arr.length; i++) {
    for (let j = 0, len = ps.length; j < len; j++) {
      ps.push(ps[j].concat(arr[i]));
    }
  }
  return ps;
};

exports = module.exports = Hs110;
