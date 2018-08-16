const CT_MIN = 2000;
const CT_MAX = 6535;

export default ColorTempDecorator = (hueLight) => {
    hueLight.colorTemp = 2700;

    hueLight.setColorTemp = (temp) => {
        if (hueLight.colorTemp === temp) return;

        // Despite what the docs say, huejay sometimes returns numbers higher than 6500
        if (temp < CT_MIN || temp > CT_MAX) {
            throw new Error(`Invalid colour temperature '${temp}'`);
        }

        hueLight.colorTemp = temp;
        hueLight.power = true;
    }

    hueLight.setState = (state) => {
        huelight.setState(state);
        if ('colorTemp' in state) hueLight.setColorTemp(state.colorTemp);
    }

    hueLight.prepareLight = () => {
        hueLight.prepareLight();

        // Convert from Kelvin to Mirek (Huejay wants a valye between 153 and 500)
        hueLight.light.colorTemp = Math.floor(1000000 / hueLight.colorTemp)

        return hueLight.light;
    }

    hueLight.applyRemoteState = (light) => {
        hueLight.setColorTemp(Math.ceil(1000000 / light.colorTemp))
        hueLight.applyRemoteState(light);
    }

    hueLight.getProperties = () => {
        let properties = hueLight.getProperties();
        properties['colorTemp'] = {type: 'int', min: CT_MIN, max: CT_MAX, interpolation: 'continuous'};
        return properties;
    }

    hueLight.toJSON = () => {
        let json = hueLight.toJSON();
        json['colorTemp'] = hueLight.colorTemp;
        return json;
    }
}
