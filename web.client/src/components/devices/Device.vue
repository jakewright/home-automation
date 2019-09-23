<template>
  <div class="device">
    <h3 class="title">
      {{ header.name }}
    </h3>

    <template v-if="fetchError">
      <p>Failed to fetch device: {{ fetchError }}</p>
    </template>
    <template v-else-if="notFound">
      Device not found
    </template>
    <template v-else-if="!device">
      Loading…
    </template>
    <template v-else>
      <ul>
        <li
          v-for="(property, propertyName) in device.state"
          :key="propertyName"
          :class="propertyName"
        >
          <template v-if="property.type === 'bool'">
            <ToggleControl
              :value="property.value"
              @input="updateProperty(propertyName, $event)"
            />
          </template>

          <template v-else-if="property.type === 'int' && property.interpolation === 'continuous'">
            <SliderControl
              :name="propertyName"
              :pretty-name="property.prettyName"
              :value="property.value"
              :min="property.min"
              :max="property.max"
              @input="updateProperty(propertyName, $event)"
            />
          </template>

          <template v-else-if="property.type === 'int'">
            <NumberControl
              :value="property.value"
              @input="updateProperty(propertyName, $event)"
            />
          </template>

          <template v-else-if="property.options">
            <SelectControl
              :value="property.value"
              :options="property.options"
              @input="updateProperty(propertyName, $event)"
            />
          </template>

          <template v-else-if="property.type === 'rgb'">
            <RgbControl
              :value="property.value"
              :device-id="device.identifier"
              @input="updateProperty(propertyName, $event)"
            />
          </template>

          <template v-else>
            {{ property.value }}
          </template>
        </li>
      </ul>
    </template>
  </div>
</template>

<script>
import DeviceHeader from '../../domain/DeviceHeader';
import ToggleControl from './controls/ToggleControl';
import SliderControl from './controls/SliderControl';
import NumberControl from './controls/NumberControl';
import SelectControl from './controls/SelectControl';
import RgbControl from './controls/RgbControl';

export default {
  name: 'Device',
  components: {
    NumberControl, RgbControl, SelectControl, SliderControl, ToggleControl,
  },

  props: {
    header: {
      type: DeviceHeader,
      required: true,
    },
  },
  data() {
    return {
      notFound: false,
      fetchError: null,
    };
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
      this.fetchError = err.message;
    }

    if (!this.device) {
      this.notFound = true;
    }
  },

  methods: {
    async updateProperty(name, value) {
      try {
        await this.$store.dispatch('updateDeviceProperty', {
          deviceId: this.header.identifier,
          name,
          value,
        });
      } catch (err) {
        console.log(err);
        await this.$store.dispatch('enqueueError', err);
      }
    },
  },
};
</script>
