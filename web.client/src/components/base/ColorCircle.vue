<template>
    <div class="color-circle">
        <input
            :class="$style.hide"
            type="color"
            :value="value"
            @input="$emit('input', $event.target.value)"
            @change="$emit('change', $event.target.value)"
            ref="input">

        <div class="target" ref="target"></div>
    </div>
</template>

<script>
    import { isHexColor } from '../../utils/validators';
    import ColorCircle from 'color-circle';


    export default {
        name: 'ColorCircle',

        props: {
            value: {
                type: String,
                required: true,
                validator: isHexColor,
            },
        },

        mounted: function() {
            new ColorCircle({
                input: this.$refs.input,
                target: this.$refs.target,
                diameterRatio: 0.25,
                thumbDiameter: 30,
            });
        },
    };

</script>

<style module>
    .hide {
        display: none;
    }
</style>
