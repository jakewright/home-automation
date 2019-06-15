<template>
  <div class="rgb-control">
    <div
      class="color"
      :style="{ backgroundColor: value }"
      @click="focusColorCircle"
    />
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

<style scoped>
  .color {
    margin: 10px 0;
    width: 60px;
    height: 20px;
    border: 2px solid #FFF;
    border-radius: 5px;
    cursor: pointer;
  }
</style>
