import { rgbToXy, xyToRgb } from './conversions';

export default RgbDecorator = (hueLight) => {
    hueLight.rgb = '#000000';

    hueLight.setState = (state) => {
        huelight.setState(state);
        if ('rgb' in state) hueLight.setRgb(state.rgb);
    }

    hueLight.prepareLight = () => {
        hueLight.prepareLight();

        hueLight.light.xy = rgbToXy(
            parseInt(hueLight.rgb.substring(1, 3), 16),
            parseInt(hueLight.rgb.substring(3, 5), 16),
            parseInt(hueLight.rgb.substring(5, 7), 16),
          );

        return hueLight.light;
    }

    hueLight.applyRemoteState = (light) => {
        const rgb = xyToRgb(light.xy[0], light.xy[1]);
        hueLight.setRgb(`#${rgb.map(x => x.toString(16).padStart(2, '0')).join('')}`);

        hueLight.applyRemoteState(light);
    }

    hueLight.setRgb = (rgb) => {
        if (hueLight.rgb === rgb) return;

        const ok = /^#[0-9A-F]{6}$/i.test(rgb);
        if (!ok) throw new Error(`Invalid hex color '${rgb}'`);

        hueLight.rgb = rgb;
        hueLight.power = true;
    }

    hueLight.toJSON = () => {
        let json = hueLight.toJSON();
        json['rgb'] = hueLight.rgb;
        return json;
    }

    hueLight.getProperties = () => {
        let properties = hueLight.getProperties();
        properties['rgb'] = {type: 'rgb'};
        return properties;
    }
}
