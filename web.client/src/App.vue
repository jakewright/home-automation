<template>
  <div id="app">
    <RouterView />
    <ErrorList />
  </div>
</template>

<script>

import ErrorList from './components/errors/ErrorList';

export default {
  name: 'App',
  components: { ErrorList },
};
</script>

<style lang="scss">
    @import '@variables';
    @import '@design';
</style>
