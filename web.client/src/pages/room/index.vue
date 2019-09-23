<template>
  <Base
    v-if="room"
    :back="back"
    :error="error"
    id="room"
  >
    <template slot="heading">
      {{ name }}
    </template>
    <template slot="content">

      <div class="tabs">
        <h2><a href="#lights" :class="{'tab--active': tab === '#lights'}">Lights</a></h2>
        <h2><a href="#devices" :class="{'tab--active': tab === '#devices'}">Devices</a></h2>
      </div>

      <Lights v-if="tab === '#lights'" :headers="room.deviceHeaders" />

<!--      <div class="devices">-->

<!--      </div>-->

<!--      <div-->
<!--        v-for="deviceHeader in room.deviceHeaders"-->
<!--        :key="deviceHeader.identifier"-->
<!--      >-->
<!--        <Component-->
<!--          :is="getComponentForDeviceType(deviceHeader.type)"-->
<!--          :device-header="deviceHeader"-->
<!--        />-->
<!--      </div>-->
    </template>

  </Base>
  <NotFound v-else></NotFound>
</template>

<script>
  import Base from "../../templates/Base";
  import Device from "../../components/devices/Device";
  import Store from "../../store";
  import NotFound from "../404/index";
  import Lights from "./Lights";

  export default {
    name: "Room",

    components: {
      Lights,
      NotFound,
      Base, Device
    },

    props: {
      // The page the back button should go to
      back: {
        type: String,
        required: false,
        default: "home"
      }
    },

    data() {
      return {
        error: null
      };
    },

    computed: {
      room() {
        return this.$store.getters.room(this.$route.params.roomId);
      },

      name() {
        return this.room ? this.room.name : "";
      },

      tab() {
        return this.$route.hash || "#lights";
      },
    },

    async beforeRouteEnter(to, from, next) {
      // If the store already has the data, resolve root.
      if (Store.getters.room(to.params.roomId)) next();

      // Otherwise fetch the remote data
      try {
        await Store.dispatch("fetchRoom", to.params.roomId);
      } catch (err) {
        console.error(err);
        next(vm => {
          vm.error = err.toString();
        });
      }

      next();
    },

    // methods: {
    //   getComponentForDeviceType(type) {
    //     switch (type) {
    //       default:
    //         return "Device";
    //     }
    //   }
    // }
  };
</script>
