<template>
  <div class="base-layout">
    <div class="grid-item">
      <div class="nav-bar">
        <div class="nav-bar__element">
          <button
            v-if="back"
            @click="goBack"
          >
            Back
          </button>
        </div>
        <div class="nav-bar__element nav-bar__element--right">
          <button
            v-if="done"
            @click="goBack"
          >
            Done
          </button>
        </div>
      </div>

      <div v-if="error">
        {{ error }}
      </div>

      <div
        v-else-if="loading"
        class="loading"
      >
        <FontAwesomeIcon
          class="loading__icon"
          icon="spinner-third"
          size="4x"
          spin
        />
        <h1 class="loading__heading">
          Loading
        </h1>
      </div>

      <template v-else>
        <div
          :class="{ 'header--icon': icon }"
          class="header"
        >
          <FontAwesomeIcon
            :icon="icon"
            class="header__icon"
            size="4x"
          />
          <h1 class="header__heading">
            <slot name="heading" />
          </h1>
        </div>

        <div class="content">
          <slot name="content" />
        </div>
      </template>
    </div>

    <Transition :name="childTransition">
      <RouterView class="extra-details grid-item" />
    </Transition>
  </div>
</template>

<script>
export default {
  name: 'PageColumn',

  props: {
    // Whether to show a "done" button in the nav bar
    done: {
      type: Boolean,
      required: false,
      default: false,
    },

    // Whether to show a "back" button in the nav bar
    back: {
      type: Boolean,
      required: false,
      default: false,
    },

    // The Font Awesome icon name to show with the heading
    icon: {
      type: String,
      required: false,
      default: '',
    },

    // The name of the CSS transition to apply to a child router view
    childTransition: {
      type: String,
      required: false,
      default: 'slide-up',
    },

    // Shows a loading screen while true
    loading: {
      type: Boolean,
      required: false,
      default: false,
    },

    // An error to show instead of the content
    error: {
      type: String,
      required: false,
    },
  },

  methods: {
    goBack() {
      window.history.length > 1
        ? this.$router.go(-1)
        : this.$router.push('/');
    },
  },
};
</script>
