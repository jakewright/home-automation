<template>
  <div>
    <div
      v-for="header in lights"
      :key="header.identifier"
    >
      <Component
        :is="getComponentForDeviceType(deviceHeader.type)"
        :device-header="deviceHeader"
      />
    </div>
  </div>
</template>

<script>
  export default {
    name: "Lights",

    computed: {
      lights() {
        return this.headers.filter(header => header.type == "light")
      },
    },

    props: {
      headers: {
        type: Array,
        required: true,
      }
    }
  };
</script>