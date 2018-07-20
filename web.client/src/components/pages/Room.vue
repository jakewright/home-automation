<template>

  <BaseLayout
    :done="false"
    class="room grid-container"
  >
    <template slot="heading">{{ name }}</template>
    <template slot="content">


      <div v-if="notFound">
        Not found
      </div>
      <template v-else>


        <div v-if="!room && !fetchError">Loading…</div>

        <div v-else-if="fetchError">
          <p>Failed to fetch room: {{ fetchError }}</p>
        </div>

        <div v-else-if="room">

          <div 
            v-for="deviceHeader in room.deviceHeaders" 
            :key="deviceHeader.identifier">
            <Component 
              :is="getComponentForDeviceType(deviceHeader.type)" 
              :device-header="deviceHeader" />
          </div>

        </div>
      </template>

    </template>
  </BaseLayout>


</template>

<script>
    import Device from '../devices/Device';
    import BaseLayout from '../layouts/BaseLayout';

    export default {
        name: 'Room',

        components: {
            Device, BaseLayout
        },

        data () {
            return {
                notFound: false,
                fetchError: null,
            };
        },

        computed: {
            room () {
                return this.$store.getters.room(this.$route.params.roomId);
            },

            name () {
                return this.room ? this.room.name : "";
            }
        },

        async created () {
            // If the store already has the data, return early.
            if (this.room) return;

            // Otherwise fetch the remote data
            try {
                await this.$store.dispatch('fetchRoom', this.$route.params.roomId);
            } catch (err) {
                console.error(err);
                this.fetchError = err.message;
            }

            // If the store is still empty, assume the room was not found.
            if (!this.room) {
                this.notFound = true;
            }
        },

        methods: {
            getComponentForDeviceType (type) {
                switch (type) {
                    default:
                        return 'Device';
                }
            },
        },
    };
</script>
