<template>
  <div class="rgb-control">
    <input
      ref="input"
      :value="value"
      type="color"
      @click.prevent="focusColorCircle"
    >
  </div>
</template>

<script>
import { isHexColor } from '../../../utils/validators';


export default {
  name: 'RgbControl',

  props: {
    value: {
      type: String,
      required: true,
      validator: isHexColor,
    },

    deviceId: {
      type: String,
      required: true,
    },
  },

  methods: {
    focusColorCircle() {
      this.$router.push({ name: 'rgb', params: { deviceId: this.deviceId } });
    },
  },


};

</script>
