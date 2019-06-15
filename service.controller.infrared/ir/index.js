const childProcess = require("child_process");

// Create a lock because even though JavaScript is single-threaded, the
// sleep function is asynchronous and we don't want to interleave instructions.
let mutex = false;

const execute = async (instructions) => {
  while (mutex) {
    // Sleep asynchronously so other code can execute
    await sleep(100);
  }

  // JavaScript is single-threaded so we don't need to worry about atomicity
  mutex = true;

  for (let i = 0; i < instructions.length; i++) {
    if (instructions[i] === "wait") {
      await sleep(instructions[++i]);
      continue;
    }

    send(instructions[i], instructions[++i]);
  }

  mutex = false;
};

const send = (device, key) => {
  try {
    childProcess.execSync(`irsend SEND_ONCE ${device} ${key}`);
  } catch (err) {
    console.error(`Process exited with status '${err.status}'`);
    console.error(err.message);
    console.error(`stderr: ${err.stderr.toString()}`);
    console.error(`stdout: ${err.stdout.toString()}`);

    throw new Error(`Failed to call irsend: ${err.message}`);
  }
};

const sleep = ms => {
  return new Promise(resolve => setTimeout(resolve, ms));
};

exports = module.exports = { execute };
