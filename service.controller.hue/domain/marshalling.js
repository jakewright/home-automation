import HueLight from "./HueLight";

const apiToHueLight = (rsp, huejayClient) => {
  return new HueLight({
    identifier: rsp.identifier,
    name: rsp.name,
    type: rsp.type,
    controllerName: rsp.controllerName,
    hueId: rsp.hueId,
    huejayClient: huejayClient
  });
};

export { apiToHueLight };
