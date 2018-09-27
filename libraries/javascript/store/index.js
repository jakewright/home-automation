const EventEmitter = require("events");

class Store extends EventEmitter {
  constructor({ state, getters, actions, mutations }) {
    super();

    this.state = state;
    this.getters = getters;
    this.mutations = mutations;
    this.actions = actions;

    this.updateCache();
  }

  updateCache() {
    this.cache = {};
    for (let key in this.state) {
      this.cache[key] = JSON.stringify(this.state[key]);
    }
  }

  get(getter, payload) {
    return this.getters[getter](
      {
        state: this.state,
        get: this.get.bind(this)
      },
      payload
    );
  }

  commit(mutation, payload) {
    const result = this.mutations[mutation](
      {
        state: this.state,
        get: this.get.bind(this),
        commit: this.commit.bind(this)
      },
      payload
    );

    for (let key in this.state) {
      if (!(key in this.cache)) {
        super.emit("key-added", key);
      } else if (this.cache[key] !== JSON.stringify(this.state[key])) {
        super.emit("key-changed", key);
      }
    }

    for (let key in this.cache) {
      if (!(key in this.state)) {
        super.emit("key-deleted", key);
      }
    }

    this.updateCache();

    return result;
  }

  dispatch(action, payload) {
    return this.actions[action](
      {
        state: this.state,
        get: this.get.bind(this),
        commit: this.commit.bind(this),
        dispatch: this.dispatch.bind(this)
      },
      payload
    );
  }
}

exports = module.exports = Store;
