<template>
  <div>
    <div class="tile-grid">
      <LightTile
        v-for="header in lights"
        :key="header.identifier"
        :header="header"
      />
    </div>

    <Component
      v-for="header in lights"
      :key="header.identifier"
      :is="getComponentForDeviceType(header.type)"
      :header="header"
    />
  </div>
</template>

<script>
  import Tile from "../../components/base/Tile";
  import LightTile from "./LightTile";
  import Device from "../../components/devices/Device";

  export default {
    name: "Lights",
    components: { LightTile, Tile, Device },

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
      getComponentForDeviceType(type) {
        switch (type) {
          default:
            return "Device";
        }
      }
    }
  };
</script>