<template>
  <div>
    <Tile
      :icon="['fal', 'lightbulb']"
      :error="error"
      :loading="!device"
    >
      <template slot="primary">
        {{ header.name }}
      </template>
    </Tile>

    <input
      type="checkbox"
      :active="checked"
    />
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
      }
    },

    computed: {
      device() {
        return this.$store.getters.device(this.header.identifier);
      },

      checked() {
        return this.device ? this.device.state["power"] : false;
      },
    },

    async created() {
      if (this.device) return;

      try {
        await this.$store.dispatch("fetchDevice", this.header);
      } catch (err) {
        console.error(err);
        this.error = true;
      }
    }

  };
</script>
