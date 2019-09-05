<template>
  <div>
    <h1>Not Found</h1>
    <h2>
      <router-link :to="{ name: 'home' }">
        <FontAwesomeIcon :icon="['far', 'home']"/> Home
      </router-link>
    </h2>
  </div>
</template>

<script>
  export default {
    name: "NotFound"
  };
</script>