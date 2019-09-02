<template>
  <div id="home">
    <h1>Hello</h1>
    <h2>Favourites</h2>
    <p>Coming soon...</p>

    <h2>Rooms</h2>
    <RoomSelector />

    <h2>Scenes</h2>
    <p>Coming soon...</p>
  </div>
</template>

<script>
import RoomSelector from "./RoomSelector";
export default {
  components: { RoomSelector }
};
</script>