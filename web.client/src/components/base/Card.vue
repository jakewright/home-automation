<template>
  <div class="card">
    <FontAwesomeIcon
      class="card__icon"
      :icon="icon"
      size="1x"
    />
    <div class="card__content">
      <span><slot name="primary" /></span>
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