<template>
  <div class="card">
    <FontAwesomeIcon
      class="card__icon"
      :icon="icon"
      size="2x"
    />
    <div class="card__content">
      <slot name="primary" />
    </div>
  </div>
</template>

<script>
  export default {
    name: "Card",

    props: {
      icon: {
        type: Array,
        required: true
      }
    }
  };
</script>