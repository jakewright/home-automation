const { Device } = require("../../libraries/javascript/device");

const irName = "ONKYO_HT_R380";

const inputOptions = [
  { value: "BD_DVD", text: "MacBook" },
  // { value: "VCR_DVR", text: "VCR/DVD" },
  { value: "CBL_SAT", text: "Roku" },
  { value: "GAME", text: "Other" },
  // { value: "AUX", text: "Auxiliary" },
  // { value: "TUNER", text: "Radio" },
  // { value: "TV_CD", text: "TV/CD" },
  // { value: "PORT", text: "PORT" }
];

const keys = {
  power: "KEY_POWER",
  BD_DVD: "KEY_DVD",
  VCR_DVR: "KEY_VCR",
  CBL_SAT: "KEY_SAT",
  GAME: "BTN_GAMEPAD",
  AUX: "KEY_AUX",
  TUNER: "KEY_TUNER",
  TV_CD: "KEY_TV",
  PORT: "KEY_TV2",
  volumeUp: "KEY_VOLUMEUP",
  volumeDown: "KEY_VOLUMEDOWN",
  mute: "KEY_MUTE"
};

class OnkyoHTR380 extends Device {
  constructor(config) {
    super(config);

    this.state = {
      power: { type: "bool" }
    };
  }

  // Perform validation and sanitization of the input
  conform(original) {
    let c = {};

    if ("power" in original) {
      c.power = Boolean(original.power);
    }

    if ("volume" in original) {
      if (!Number.isInteger(original.volume)) {
        throw new Error(`Volume must be an integer`);
      }

      // Volume is relative because we don't know what it's currently set to
      // but check that the number isn't ridiculous.
      if (Math.abs(original.volume) > 10) {
        throw new Error("Relative volume change must be [-10, 10]");
      }

      c.volume = original.volume;
    }

    if ("input" in original) {
      // Validate the input name
      if (!inputOptions.find(option => option.value === original.input)) {
        let options = inputOptions.map(option => option.value);
        throw new Error(
          `Invalid input '${original.input}'; must be one of: ${options}`
        );
      }

      // Make sure the key map is defined
      if (!(original.input in keys)) {
        throw new Error(`Key not defined for '${original.input}'`);
      }

      c.input = original.input;
    }

    if ("mute" in original && original.mute === true) {
      c.mute = true;
    }

    return c;
  }

  generateInstructions(state) {
    let instructions = [];

    if ("power" in state && state.power !== this.state.power.value) {
      instructions.push(irName, keys.power);

      // Give the AV receiver plenty of time to get its affairs in order
      instructions.push("wait", 5000);
    }

    if ("volume" in state) {
      const key = state.volume > 0 ? keys.volumeUp : keys.volumeDown;

      // Send the key n + 1 times because the key needs to be
      // pressed to activate the volume control before it changes
      for (let i = 0; i <= Math.abs(state.volume); i++) {
        instructions.push(irName, key);
        instructions.push("wait", 200);
      }

      // Wait until the volume control goes away
      // before attempting anything else
      instructions.push("wait", 3000);
    }

    if ("input" in state) {
      instructions.push(irName, keys[state.input]);
      instructions.push("wait", 3000);
    }

    if ("mute" in state) {
      instructions.push(irName, keys.mute);
      instructions.push("wait", 3000);
    }

    return instructions;
  }
}

exports = module.exports = OnkyoHTR380;