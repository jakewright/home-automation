<template>
    <div class="rgb-control">
        <input type="color" :value="value" @click.prevent="focusColorCircle" ref="input">
    </div>
</template>

<script>
    import { isHexColor } from '../../../utils/validators';

    
    export default {
        name: 'RgbControl',

        props: {
            value: {
                type: String,
                required: true,
                validator: isHexColor,
            },

            deviceId: {
                type: String,
                required: true,
            }
        },

        methods: {
            focusColorCircle() {
                this.$router.push({ name: 'rgb', params: { deviceId: this.deviceId }})
            },
        },


    };

</script>
