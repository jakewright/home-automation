<template>
  <div class="base-layout">
    <div class="header">
      <div class="header__element">
        <h1 class="header__heading">
          <router-link v-if="back" :to="{ name: back }">
            <FontAwesomeIcon :icon="['far', 'arrow-left']"/>
          </router-link>
          <slot name="heading"/>
        </h1>
      </div>

    </div>

    <div v-if="error">
      {{ error }}
    </div>

    <div v-else class="content">
      <slot name="content"/>
    </div>

    <Transition :name="childTransition">
      <RouterView class=""/>
    </Transition>

  </div>
</template>

<script>
  export default {
    name: "Base",

    props: {
      // The page the back button should link to
      back: {
        type: String,
        required: false,
        default: false
      },

      // The name of the CSS transition to apply to a child router view
      childTransition: {
        type: String,
        required: false,
        default: "slide-up"
      },

      // An error to show instead of the content
      error: {
        type: String,
        required: false
      }
    },

    methods: {
      goBack() {
        window.history.length > 1
          ? this.$router.go(-1)
          : this.$router.push("/");
      }
    }
  };
</script>
