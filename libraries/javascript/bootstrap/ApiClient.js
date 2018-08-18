const axios = require("axios");
const utils = require("./utils");

class ApiClient {
  constructor(apiGateway) {
    this.client = axios.create({
      baseURL: `${apiGateway}/`,
      validateStatus: status => status >= 200 && status < 300,
      transformResponse: [
        function(data) {
          console.log;
          data = JSON.parse(data);

          // If the data field is empty, there was probably an error so
          // return as-is and handle in the request method below.
          if (!data.data) return data;

          // The response will be enveloped but we no longer care about the
          // extra metadata. Also convert keys to camcel case because JavaScript.
          return utils.toCamelCase(data.data);
        }
      ]
    });
  }

  async request(args) {
    try {
      const response = await this.client.request.apply(this, arguments);
      return response.data;
    } catch (err) {
      console.error("Error making request");
      console.error(args);
      console.error(err);

      // If the backend actuall response, extract useful information from it.
      if (err.response) {
        console.error(err.response);

        // The response should include an array of errors. Turn them into a comma-separated list.
        const errors = err.response.data.errors.reduce((s, e) => `${s}, ${e}`);

        // Add all of the extra messages from the response to the error message
        err.message += `: ${err.response.data.message}: errors: [${errors}]`;
      }

      throw err;
    }
  }

  getDevices(controllerName) {
    return this.request({
      url: "service.registry.device/devices",
      method: "get",
      params: {
        controller_name: controllerName
      }
    });
  }
}

exports = module.exports = ApiClient;
