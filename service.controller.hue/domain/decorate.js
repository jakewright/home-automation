const decorate = (hueLight, decorator) => {
  const transform = hueLight.transform;
  hueLight.transform = state => {
    let t = transform.call(hueLight, state);
    return decorator.transform.call(hueLight, state, t);
  };

  const applyRemoteState = hueLight.applyRemoteState;
  hueLight.applyRemoteState = state => {
    applyRemoteState.call(hueLight, state);
    decorator.applyRemoteState.call(hueLight, state);
  };

  const getProperties = hueLight.getProperties;
  hueLight.getProperties = () => {
    let properties = getProperties.call(hueLight);
    return decorator.getProperties.call(hueLight, properties);
  };
};

exports = module.exports = decorate;
