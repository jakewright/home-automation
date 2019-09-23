<template>
  <div class="tile" :class="{ 'tile--active': active, 'tile--clickable': clickable }">
    <FontAwesomeIcon
      v-if="error"
      class="tile__icon"
      :icon="['fal', 'exclamation-triangle']"
      size="1x"
    />
    <FontAwesomeIcon
      v-else-if="loading"
      class="tile__icon"
      :icon="['fal', 'asterisk']"
      spin
      size="1x"
    />
    <FontAwesomeIcon
      v-else
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
      // The Font Awesome icon to display in the top left
      icon: {
        type: Array,
        required: true
      },

      // Affects the tile's background color
      active: {
        type: Boolean,
        required: false
      },

      // Affects the cursor
      clickable: {
        type: Boolean,
        required: false,
      },

      // If true then the icon will be replaced with an exclamation
      // triangle. This takes precedence over loading.
      error: {
        type: Boolean,
        required: false,
      },

      // If true then the icon will be replaced with a loading icon
      loading: {
        type: Boolean,
        required: false
      },
    }
  };
</script>