<template>
  <div>
    <Tile
      :icon="['fal', 'lightbulb']"
      :error="error"
      :loading="loading"
      :active="checked"
      :clickable="!loading"
      @click.native="toggle"
    >
      <template slot="primary">
        {{ header.name }}
      </template>
    </Tile>
  </div>
</template>

<script>
  import Tile from "../../components/base/Tile";
  import DeviceHeader from "../../domain/DeviceHeader";

  export default {
    name: "LightTile",

    components: { Tile },

    props: {
      header: {
        type: DeviceHeader,
        required: true
      }
    },

    data: function () {
      return {
        error: false,
        loading: true,
      }
    },

    computed: {
      device() {
        return this.$store.getters.device(this.header.identifier);
      },

      checked() {
        return this.device ? this.device.state["power"]["value"] : false;
      },
    },

    async created() {
      if (!this.device) {
        try {
          await this.$store.dispatch("fetchDevice", this.header);
        } catch (err) {
          console.error(err);
          this.error = true;
        }
      }

      this.loading = false;
    },

    methods: {
      async toggle() {
        // Return early if the device hasn't been loaded yet
        if (!this.device) return;

        // Return early if already busy doing something
        if (this.loading) return;

        // Enable the loading icon if the request takes more than 200ms
        const loading = setTimeout(() => {
          this.loading = true;
        }, 200);

        try {
          await this.$store.dispatch("updateDeviceProperty", {
            deviceId: this.header.identifier,
            name: "power",
            value: !this.checked,
          });
        } catch (err) {
          console.error(err);
          await this.$store.dispatch('enqueueError', err);
        }

        // Cancel the timeout and set loading to false
        clearTimeout(loading);
        this.loading = false;
      },
    },

  };
</script>
