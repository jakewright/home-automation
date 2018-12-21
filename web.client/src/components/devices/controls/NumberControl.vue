<template>
  <div>
    <input
      :min="min"
      :max="max"
      :value="value"
      type="number"
      @change="$emit('input', $event.target.value)"
    >
  </div>
</template>

<script>
export default {
  name: 'NumberControl',

  props: {
    value: {
      type: Number,
      required: true,
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
};
</script>
