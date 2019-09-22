<template>
  <div class="tile" :class="{ 'tile--active': active }">
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
      icon: {
        type: Array,
        required: true
      },

      active: {
        type: Boolean,
        required: false
      },

      error: {
        type: Boolean,
        required: false,
      },

      loading: {
        type: Boolean,
        required: false
      },
    }
  };
</script>