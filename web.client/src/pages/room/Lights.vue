<template>
  <div>
    Lights

    <div class="tile-grid">
      <LightTile
        v-for="header in lights"
        :key="header.identifier"
        :header="header"
        />

    </div>
  </div>

</template>

<script>
  import Tile from "../../components/base/Tile";
  import LightTile from "./LightTile";

  export default {
    name: "Lights",
    components: { LightTile, Tile },

    computed: {
      lights() {
        return this.headers.filter(header => header.kind === "lamp");
      }
    },

    props: {
      headers: {
        type: Array,
        required: true
      }
    },

    methods: {
      async updateProperty(deviceId, name, value) {
        try {
          await this.$store.dispatch('updateDeviceProperty', {
            deviceId,
            name,
            value,
          })
        } catch (err) {
          console.error(err);
          await this.$store.dispatch('enqueueError', err)
        }
      }
    }
  };
</script>