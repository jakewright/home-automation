<template>
  <SideColumn
    :loading="!device"
    done
    icon="lightbulb"
    class="color-picker"
  >
    <template slot="heading">
      Colour
    </template>
    <template slot="content">
      <ColorCircle
        :value="value"
        @input="handleInput"
        @change="handleChange"
      />
    </template>
  </SideColumn>
</template>

<script>
import _ from 'lodash';
import ColorCircle from '../base/ColorCircle';
import SideColumn from '../layouts/SideColumn';
import { isHexColor } from '../../utils/validators';

export default {
  name: 'ColorPicker',

  components: { SideColumn, ColorCircle },

  props: {
    initialValue: {
      type: String,
      required: false,
      default: '#000000',
      validator: isHexColor,
    },
  },

  // created () {
  //     this.value = this.initialValue;
  // },

  data() {
    return {
      throttled: _.throttle(async () => await this.updateRgb(), 200),
    };
  },

  computed: {
    device() {
      return this.$store.getters.device(this.$route.params.deviceId);
    },

    value() {
      return this.device ? this.device.state.rgb.value : this.initialValue;
    },
  },

  methods: {
    handleInput(value) {
      this._value = value;
      this.throttled();
    },

    handleChange() {
      this.throttled.flush();
    },

    async updateRgb() {
      try {
        await this.$store.dispatch('updateDeviceProperty', {
          deviceId: this.$route.params.deviceId,
          name: 'rgb',
          value: this._value,
        });
      } catch (err) {
        console.log(err);
        await this.$store.dispatch('enqueueError', err);
      }
    },
  },
};
</script>
