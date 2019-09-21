<template>
  <label :class="$style.switch">
    <input
      :checked="value"
      type="checkbox"
      @change="$emit('input', $event.target.checked)"
    >
    <span :class="[$style.slider, $style.round]" />
  </label>
</template>

<script>
export default {
  name: 'ToggleControl',

  props: {
    value: {
      type: Boolean,
      required: true,
    },
  },
};
</script>

<style module lang="scss">
    @import '@variables';

    .switch {
        position: relative;
        display: inline-block;
        width: 60px;
        height: 34px;
    }

    .switch input {display:none;}

    .slider {
        position: absolute;
        cursor: pointer;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background-color: $color-gray-dark;
        -webkit-transition: .4s;
        transition: .4s;
    }

    .slider:before {
        position: absolute;
        content: "";
        height: 26px;
        width: 26px;
        left: 4px;
        bottom: 4px;
        background-color: white;
        -webkit-transition: .4s;
        transition: .4s;
    }

    input:checked + .slider {
        background-color: $color-green;
    }

    input:focus + .slider {
        box-shadow: 0 0 1px $color-green;
    }

    input:checked + .slider:before {
        -webkit-transform: translateX(26px);
        -ms-transform: translateX(26px);
        transform: translateX(26px);
    }

    .slider.round {
        border-radius: 34px;
    }

    .slider.round:before {
        border-radius: 50%;
    }

</style>
