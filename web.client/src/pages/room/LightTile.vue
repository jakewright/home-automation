<template>
  <Tile
    :icon="['fal', 'lightbulb']"
  >
    <template slot="primary">
      {{ header.name }}
    </template>
  </Tile>
</template>

<script>
  import Tile from "../../components/base/Tile";
  import DeviceHeader from '../../domain/DeviceHeader';

  export default {
    name: "LightTile",

    components: { Tile },

    props: {
      header: {
        type: DeviceHeader,
        required: true,
      },
    },

    computed: {
      device() {
        return this.$store.getters.device(this.header.identifier);
      },
    },

    async created() {
      if (this.device) return;

      try {
        await this.$store.dispatch('fetchDevice', this.header);
      } catch (err) {
        console.error(err);
        this.error = err.message;
      }
    },

  };
</script>
