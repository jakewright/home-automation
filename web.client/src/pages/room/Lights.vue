<template>
  <div>
    Lights

    <div class="card-grid">
      <Card
        v-for="header in lights"
        :key="header.identifier"
        :icon="['far', 'lightbulb']"
      >
        <template slot="primary">
          {{ header.name }}
        </template>
      </Card>

    </div>
  </div>

</template>

<script>
  import Card from "../../components/base/Card";

  export default {
    name: "Lights",
    components: { Card },
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
    }
  };
</script>