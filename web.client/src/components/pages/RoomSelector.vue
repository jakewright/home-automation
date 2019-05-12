<template>
  <BaseLayout
    class="room-selector"
  >
    <template slot="heading">
      Rooms
    </template>
    <template slot="content">
      <ul
        v-if="rooms.length > 0"
        class="rooms"
      >
        <li
          v-for="room in rooms"
          :key="room.identifier"
          class="room"
        >
          <RouterLink :to="{ name: 'room', params: { roomId: room.identifier } }">
            {{ room.name }}
          </RouterLink>
        </li>
      </ul>

      <div v-if="fetchError">
        <p>Failed to fetch rooms: {{ fetchError }}</p>
      </div>
    </template>
  </BaseLayout>
</template>

<script>
import { mapGetters } from 'vuex';
import BaseLayout from '../layouts/BaseLayout';

export default {
  name: 'RoomSelector',

  components: {
    BaseLayout,
  },

  data() {
    return {
      loading: true,
      fetchError: null,
    };
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
  },
};
</script>
