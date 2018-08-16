import _ from 'lodash';

const HUE_MIN = 0;
const HUE_MATCH = 65536;
const SAT_MIN = 0;
const SAT_MAX = 254;

export default ColorDecorator = (hueLight) => {
    hueLight.color = { hue: 0, saturation: 0 };

    hueLight.setColor = ({ hue, saturation }) => {
        if (_.isEqual(hueLight.color, { hue, saturation })) return;

        if (!hue) throw new Error('Must set hue');
        if (!saturation) throw new Error('Must set saturation');

        if (hue < HUE_MIN || hue > HUE_MAX) throw new Error(`Invalid hue '${hue}'`);
        if (saturation < SAT_MIN || saturation > SAT_MAX) throw new Error(`Invalid saturation '${sat}'`);

        hueLight.color = { hue, saturation };
        hueLight.power = true;
    }

    hueLight.setState = (state) => {
        huelight.setState(state);
        if ('color' in state) hueLight.setColor(state.color);
    }

    hueLight.prepareLight = () => {
        hueLight.prepareLight();

        hueLight.light.hue = hueLight.color.hue;
        hueLight.light.saturation = hueLight.color.saturation;

        return hueLight.light;
    }

    hueLight.applyRemoteState = (light) => {
        hueLight.setColor({
            hue: light.hue,
            saturation: light.saturation,
        });

        hueLight.applyRemoteState(light);
    }

    hueLight.toJSON = () => {
        let json = hueLight.toJSON();
        json['color'] = hueLight.color;
        return json;
    }

    hueLight.getProperties = () => {
        let properties = hueLight.getProperties();
        properties['color'] = {type: 'color'};
        return properties;
    }
}
