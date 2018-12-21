<template>
  <div class="color-circle">
    <input
      ref="input"
      :class="$style.hide"
      :value="value"
      type="color"
      @input="$emit('input', $event.target.value)"
      @change="$emit('change', $event.target.value)"
    >

    <div
      ref="target"
      class="target"
    />
  </div>
</template>

<script>
import ColorCircle from 'color-circle';
import { isHexColor } from '../../utils/validators';


export default {
  name: 'ColorCircle',

  props: {
    value: {
      type: String,
      required: true,
      validator: isHexColor,
    },
  },

  mounted() {
    new ColorCircle({
      input: this.$refs.input,
      target: this.$refs.target,
      diameterRatio: 0.25,
      thumbDiameter: 30,
    });
  },
};

</script>

<style module>
    .hide {
        display: none;
    }
</style>
