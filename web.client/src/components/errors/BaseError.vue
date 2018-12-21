<template>
  <div class="notification is-danger">
    <button
      class="delete"
      @click="$emit('close')"
    />
    <slot />
  </div>
</template>
