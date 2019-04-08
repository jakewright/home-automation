const axios = require("axios");
const utils = require("./utils");

class HTTPClient {
  constructor() {
    this.client = axios.create({
      validateStatus: () => true
    });
  }

  setApiGateway(apiGateway) {
    this.client.defaults.baseURL = apiGateway;
  }

  async request(args) {
    args.params = utils.toSnakeCase(args.params);

    // Make the request
    let rsp;
    try {
      rsp = await this.client.request(args);
    } catch (err) {
      if (err.request) {
        console.error("No response received", args.url);
      } else {
        console.log("Axios error: ", err.message);
      }

      throw err;
    }

    // Validate the status
    if (!this.validStatus(rsp.status)) {
      console.error(
        "request failed with status",
        rsp.status,
        args.url,
        rsp.data,
        rsp.headers
      );
      let msg = `request failed with status ${rsp.status} ${rsp.statusText}`;

      // Try to pull out a message if any JSON exists in the body
      try {
        let data = JSON.parse(rsp.data);
        if (data.message) msg += `: ${data.message}`;
      } catch (err) {
        if (rsp.data) msg += `: ${rsp.data}`;
      }

      throw new Error(msg);
    }

    if (!("data" in rsp.data)) {
      console.error("Invalid response", args.url, rsp.status, rsp.data);
      throw new Error(`data not found in response: ${rsp.data}`);
    }

    return utils.toCamelCase(rsp.data.data);
  }

  validStatus(status) {
    return status >= 200 && status < 300;
  }

  get(url, params) {
    return this.request({
      url,
      method: "get",
      params
    });
  }

  patch(url, data) {
    return this.request({
      url,
      method: "patch",
      data
    });
  }
}

const httpClient = new HTTPClient();
exports = module.exports = httpClient;
