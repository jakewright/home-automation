<template>
  <div class="tile" :class="{ 'tile--active': checked }">
    <FontAwesomeIcon
      class="tile__icon"
      :icon="icon"
      size="1x"
    />
    <div class="tile__content">
      <span><slot name="primary" /></span>
    </div>
  </div>
</template>

<script>
  export default {
    name: "Tile",

    props: {
      icon: {
        type: Array,
        required: true
      },

      checked: {
        type: Boolean,
        required: false
      }
    }
  };
</script>