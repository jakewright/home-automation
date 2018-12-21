<template>
  <form class="slider">
    <input
      v-model="currentValue"
      :min="min"
      :max="max"
      :name="name"

      type="range"

      @change="$emit('input', $event.target.value)"
    >

    <output>
      {{ currentValue }}
    </output>
  </form>
</template>

<script>
export default {
  name: 'SliderControl',

  props: {
    value: {
      type: Number,
      required: true,
    },

    name: {
      type: String,
      required: false,
      default: 'slider',
    },

    min: {
      type: Number,
      required: false,
      default: 0,
    },

    max: {
      type: Number,
      required: false,
      default: 100,
    },
  },

  data() {
    return {
      currentValue: this.value,
    };
  },

  watch: {
    value(newValue) {
      this.currentValue = newValue;
    },
  },
};
</script>
