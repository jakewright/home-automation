<template>
  <div :class="$style.errorList">
    <TransitionGroup name="error-item">
      <BaseError
        v-for="err in errors"
        :key="err.id"
        @close="removeError(err.id)"
      >
        {{ err.message }}
      </BaseError>
    </TransitionGroup>
  </div>
</template>

<script>
import { mapGetters } from 'vuex';
import BaseError from './BaseError';

export default {
  name: 'ErrorList',

  components: {
    BaseError,
  },

  computed: mapGetters({
    errors: 'allErrors',
  }),

  methods: {
    removeError(id) {
      this.$store.dispatch('removeError', id);
    },
  },

};
</script>

<style module lang="scss">
    .errorList {
        position: fixed;
        left: 20px;
        right: 20px;
        bottom: 20px;
    }
</style>
<style>
    .error-item-enter-active, .error-item-leave-active {
        transition: all 1s;
    }
    .error-item-enter, .error-item-leave-to /* .list-leave-active below version 2.1.8 */ {
        opacity: 0;
    }

    .error-item-move {
        transition: transform 0.5s;
    }
</style>
