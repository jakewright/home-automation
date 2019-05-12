<template>
  <form class="slider">
    <label
      class="property-name-label"
      :for="name"
    >
      {{ prettyName ? prettyName : name }}
    </label>

    <output>
      {{ currentValue }}
    </output>

    <input
      v-model="currentValue"
      :min="min"
      :max="max"
      :name="name"

      type="range"

      @change="$emit('input', $event.target.value)"
    >
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

    prettyName: {
      type: String,
      required: false,
      default: '',
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
