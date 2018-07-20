<template>
  <div class="grid-container full-height">
    <RoomSelector class="grid-item room-selector" />
    <router-view class="main-container grid-item"/>
  </div>
</template>

<script>
    import RoomSelector from './RoomSelector';

    export default {
        components: {RoomSelector}
    }
</script>
