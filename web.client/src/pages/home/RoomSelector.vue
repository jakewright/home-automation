<template>
  <div v-if="fetchError">
    {{ fetchError }}
  </div>
  <div v-else class="room-selector">
    <div class="scroller">
      <RouterLink
        v-for="room in rooms"
        :key="room.identifier"
        :to="{ name: 'room', params: {roomId: room.identifier } }"
        class="room"
      >
        <div class="box">
          <div class="text-container">{{ room.name }}</div>
        </div>
      </RouterLink>
    </div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex';

export default {
  name: 'RoomSelector',

  data() {
    return {
      loading: true,
      fetchError: null,
    }
  },

  computed: mapGetters({
    rooms: 'allRooms',
  }),

  async created() {
    try {
      await this.$store.dispatch('fetchRooms');
    } catch (err) {
      console.error(err);
      this.fetchError = err.message;
    }

    this.loading = false;
  }
}
</script>