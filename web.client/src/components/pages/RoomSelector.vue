<template>
  <div 
    id="room-selector" 
    class="room-selector">
    <h2>Rooms</h2>

    <div v-if="loading">Loading…</div>

    <ul 
      v-if="rooms.length > 0" 
      class="rooms">
      <li 
        v-for="room in rooms" 
        :key="room.identifier" 
        class="room">
        <router-link :to="{ name: 'room', params: { roomId: room.identifier } }" >{{ room.name }}</router-link>
      </li>
    </ul>

    <div v-if="fetchError">
      <p>Failed to fetch rooms: {{ fetchError }}</p>
    </div>
  </div>
</template>

<script>
    import { mapGetters } from 'vuex';

    export default {
        name: 'RoomSelector',
        data () {
            return {
                loading: true,
                fetchError: null,
            };
        },

        computed: mapGetters({
            rooms: 'allRooms',
        }),

        async created () {
            try {
                await this.$store.dispatch('fetchRooms');
            } catch (err) {
                console.error(err);
                this.fetchError = err.message;
            }

            this.loading = false;
        },
    };
</script>
